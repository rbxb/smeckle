// Package model is used to decode 3D models from Super Evil Megacorp's format.
package model

import (
	"bytes"
	"encoding/binary"
	"io"
)

// VertexChunkUsageX are bytes that define what the data in a vertex chunk is used for.
const (
	VertexChunkUsagePosition = 0x00
	VertexChunkUsageNormal   = 0x01
	VertexChunkUsageUV       = 0x05
)

// VertexChunkStoredAsX are bytes that define how a number in a vertex chunk is encoded.
const (
	VertexChunkStoredAsFloat32 = 0x04
	VertexChunkStoredAsUint8   = 0x05
)

// An Object stores part of a 3D model as an array of face indices.
type Object struct {
	Faces []int
}

// A Model stores a 3D model.
type Model struct {
	Name              string
	Verts, UVs, Norms []float32
	Objects           []Object
}

// A vectorHeader contains values needed to read a vector.
// Each vertex in the 3D model may have multiple vectors associated with it.
// Each vector has it's own usage, e.g. a vector may define the position of a vertex.
type vectorHeader struct {
	usage      byte
	storedAs   byte
	valueCount int
	offset     int64
}

// ConvertModel reads model data that is encoded in SEMC's format from a reader and puts it in a model object.
func ConvertModel(r io.ReadSeeker) *Model {
	model := &Model{}

	// Skip the first 36 bytes.
	// I don't know what these bytes are for.
	r.Seek(36, 1)

	// Read the number of materials used by the model.
	// Materials are in separate files which reference the textures and shaders used by a rendered 3D model.
	var materialCount uint16
	if err := binary.Read(r, binary.LittleEndian, &materialCount); err != nil {
		panic(err)
	}
	materials := make([]string, materialCount)

	// Read the names of the materials.
	for i := range materials {
		b := make([]byte, 1)
		buf := bytes.NewBuffer(nil)
		// The names are separated by zero bytes.
		// Read characters until a zero byte is reached.
		for {
			_, err := r.Read(b)
			if err != nil {
				panic(err)
			}
			if b[0] == 0x00 {
				break
			}
			buf.Write(b)
		}
		materials[i] = string(buf.Bytes())
	}

	// Set the name of the model to the name of the first material that is referenced.
	model.Name = materials[0]

	// Read the number of objects in the file.
	// Each object is it's own 3D model made of a collection of faces.
	var objectCount uint16
	if err := binary.Read(r, binary.LittleEndian, &objectCount); err != nil {
		panic(err)
	}
	model.Objects = make([]Object, objectCount)

	// Skip 2 bytes.
	// I don't know what these bytes are for (Maybe texture count).
	r.Seek(2, 1)

	// Read the number of face indices that each object has.
	// Each face has three indices.
	// Each index corresponds to a vertex.
	for i := range model.Objects {
		// Skip 4 bytes.
		// I don't know what these bytes are for.
		r.Seek(4, 1)

		var faces uint32
		if err := binary.Read(r, binary.LittleEndian, &faces); err != nil {
			panic(err)
		}

		// Create the face index buffer.
		// Every three indices represents a face.
		model.Objects[i].Faces = make([]int, int(faces))

		// Skip 63 bytes.
		// I don't know what these bytes are for.
		r.Seek(63, 1)
	}

	// Skip 28 bytes.
	// I don't know what these bytes are for.
	r.Seek(28, 1)

	// Read the number of vertices in the model.
	var vertexCount uint32
	if err := binary.Read(r, binary.LittleEndian, &vertexCount); err != nil {
		panic(err)
	}
	model.Verts = make([]float32, int(vertexCount)*3)
	model.UVs = make([]float32, int(vertexCount)*2)
	model.Norms = make([]float32, int(vertexCount)*3)

	// Skip 6 bytes.
	// I don't know what these bytes are for.
	r.Seek(6, 1)

	// Read the size of each vertex.
	// A vertex is made of multiple vectors.
	var vertexSize uint32
	if err := binary.Read(r, binary.LittleEndian, &vertexSize); err != nil {
		panic(err)
	}

	// Read the number of vectors used in each vertex.
	var vectorCount uint8
	if err := binary.Read(r, binary.LittleEndian, &vectorCount); err != nil {
		panic(err)
	}
	vectorHeaders := make([]vectorHeader, int(vectorCount))

	// Fill in the values for each vector header.
	for i := range vectorHeaders {
		b := make([]byte, 4)
		if _, err := r.Read(b); err != nil {
			panic(err)
		}
		vectorHeaders[i].usage = b[0]
		// I don't know what b[1] is used for.
		vectorHeaders[i].storedAs = b[2]
		vectorHeaders[i].valueCount = int(uint8(b[3]))
		var offset uint32
		if err := binary.Read(r, binary.LittleEndian, &offset); err != nil {
			panic(err)
		}
		vectorHeaders[i].offset = int64(offset)
	}

	// Read the vertices.
	for i := 0; i < int(vertexCount); i++ {
		vertexOffset, _ := r.Seek(0, 1)
		for _, vector := range vectorHeaders {
			r.Seek(vertexOffset+vector.offset, 0)
			switch vector.usage {
			// This vector defines the position of the vertex.
			case VertexChunkUsagePosition:
				switch vector.storedAs {
				case VertexChunkStoredAsFloat32:
					if err := binary.Read(r, binary.LittleEndian, model.Verts[i*3+0:i*3+3]); err != nil {
						panic(err)
					}
				}
			// This vector defines the direction of the normal of the vertex.
			case VertexChunkUsageNormal:
				switch vector.storedAs {
				case VertexChunkStoredAsFloat32:
					if err := binary.Read(r, binary.LittleEndian, model.Norms[i*3+0:i*3+3]); err != nil {
						panic(err)
					}
				}
			// This vector defines the UV position of the vertex.
			case VertexChunkUsageUV:
				switch vector.storedAs {
				case VertexChunkStoredAsFloat32:
					if err := binary.Read(r, binary.LittleEndian, model.UVs[i*2+0:i*2+2]); err != nil {
						panic(err)
					}
					model.UVs[i*2+1] = 1 - model.UVs[i*2+1]
				case VertexChunkStoredAsUint8:
					b := make([]uint8, 2)
					if err := binary.Read(r, binary.LittleEndian, b); err != nil {
						panic(err)
					}
					model.UVs[i*2+0] = float32(b[0]) / 255
					model.UVs[i*2+1] = 1 - float32(b[1])/255
				}
			}
		}

		// Skip to the next vertex.
		r.Seek(vertexOffset+int64(vertexSize), 0)
	}

	// Read the face indices of each object.
	// Each face index is a reference to a vertex.
	for i, object := range model.Objects {
		for face := range object.Faces {
			// The number of bytes used to store the index is dependent on the number of vertices.
			if vertexCount <= 255 {
				var b uint8
				if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
					panic(err)
				}
				model.Objects[i].Faces[face] = int(b)
			} else if vertexCount <= 65535 {
				var b uint16
				if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
					panic(err)
				}
				model.Objects[i].Faces[face] = int(b)
			} else {
				var b uint32
				if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
					panic(err)
				}
				model.Objects[i].Faces[face] = int(b)
			}
		}
	}

	return model
}
