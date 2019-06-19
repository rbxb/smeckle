package rsc0

import (
	"io"
	"bytes"
	"encoding/binary"
	"github.com/rbxb/smeckle/model"
)

const (
	VERTEX_CHUNK_USAGE_POSITION = 		0x00
	VERTEX_CHUNK_USAGE_NORMAL = 		0x01
	VERTEX_CHUNK_USAGE_UV = 			0x05

	VERTEX_CHUNK_STORED_AS_FLOAT32 = 	0x04
	VERTEX_CHUNK_STORED_AS_UINT8 = 		0x05
)

type datatype struct {
	usage, unkown, storedAs byte
	valueCount int
	offset int64
}

func Decode(r io.ReadSeeker) (* model.Model, error) {
	m := model.Model{}

	r.Seek(36,1)
	//Skip the first 36 bytes.
	//I don't know what these bytes are for.

	var materialCount uint16
	if err := binary.Read(r, binary.LittleEndian, &materialCount); err != nil {
		panic(err)
	}
	//Read the number of materials used by the model.

	materials := make([]string, materialCount)
	for i := range materials {
		b := make([]byte, 1)
		buf := bytes.NewBuffer(nil)
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
	//Read the material names.

	m.Name = materials[0]

	var objectCount uint16
	if err := binary.Read(r, binary.LittleEndian, &objectCount); err != nil {
		panic(err)
	}
	//Read the number of objects in the file.

	r.Seek(2,1)
	//Skip 2 bytes.
	//I don't know what these bytes are for.
	//Maybe texture count?

	m.Objects = make([]model.Object, objectCount)
	for i := range m.Objects {
		r.Seek(4,1)
		//Skip 4 bytes.
		//I don't know what these bytes are for.

		var b uint32
		if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
			panic(err)
		}
		//Read the objectIndexCount for this object.

		m.Objects[i].Faces = make([]int, int(b))
		//Make the face buffer in the object

		r.Seek(63,1)
		//Skip 63 bytes.
		//I don't know what these bytes are for.
	}
	//Read out the number of face indices in each object.

	r.Seek(28,1)
	//Skip 28 bytes.
	//I don't know what these bytes are for.

	var vertCount int
	var vertCountUint32 uint32
	if err := binary.Read(r, binary.LittleEndian, &vertCountUint32); err != nil {
		panic(err)
	}
	vertCount = int(vertCountUint32)
	//Read the number of vertices.

	m.Verts = 	make([]float32, vertCount * 3)
	m.UVs = 	make([]float32, vertCount * 2)
	m.Norms = 	make([]float32, vertCount * 3)
	//Make the model buffers.

	r.Seek(6,1)
	//Skip 6 bytes.
	//I don't know what these bytes are for.

	var chunkSize uint32
	if err := binary.Read(r, binary.LittleEndian, &chunkSize); err != nil {
		panic(err)
	}
	//Read the vertex chunk size.

	var datatypeCount uint8
	if err := binary.Read(r, binary.LittleEndian, &datatypeCount); err != nil {
		panic(err)
	}
	//Read the number of datatypes per chunk.

	datatypes := make([]datatype, int(datatypeCount))
	for i := range datatypes {
		b := make([]byte, 4)
		if _, err := r.Read(b); err != nil {
			panic(err)
		}
		datatypes[i].usage = b[0]
		datatypes[i].unkown = b[1]
		datatypes[i].storedAs = b[2]
		datatypes[i].valueCount = int(uint8(b[3]))
		var offset uint32
		if err := binary.Read(r, binary.LittleEndian, &offset); err != nil {
			panic(err)
		}
		datatypes[i].offset = int64(offset)
	}
	//Populate the datatypes.

	for i := 0; i < vertCount; i++ {
		chunkPos, _ := r.Seek(0,1)
		for _, dt := range datatypes {
			r.Seek(chunkPos + dt.offset,0)
			switch dt.usage {
			case VERTEX_CHUNK_USAGE_POSITION:
				switch dt.storedAs {
				case VERTEX_CHUNK_STORED_AS_FLOAT32:
					if err := binary.Read(r, binary.LittleEndian, m.Verts[i*3+0:i*3+3]); err != nil {
						panic(err)
					}
				}
			case VERTEX_CHUNK_USAGE_NORMAL:
				switch dt.storedAs {
				case VERTEX_CHUNK_STORED_AS_FLOAT32:
					if err := binary.Read(r, binary.LittleEndian, m.Norms[i*3+0:i*3+3]); err != nil {
						panic(err)
					}
				}
			case VERTEX_CHUNK_USAGE_UV:
				switch dt.storedAs {
				case VERTEX_CHUNK_STORED_AS_FLOAT32:
					if err := binary.Read(r, binary.LittleEndian, m.UVs[i*2+0:i*2+2]); err != nil {
						panic(err)
					}
					m.UVs[i*2+1] = 1 - m.UVs[i*2+1]
				case VERTEX_CHUNK_STORED_AS_UINT8:
					b := make([]uint8, 2)
					if err := binary.Read(r, binary.LittleEndian, b); err != nil {
						panic(err)
					}
					m.UVs[i*2+0] = float32(b[0]) / 255
					m.UVs[i*2+1] = 1 - float32(b[1]) / 255
				}
			}
		}
		r.Seek(chunkPos + int64(chunkSize),0)
	}
	//Read the vertex chunks.

	for i, object := range m.Objects {
		for f := range object.Faces {
			if vertCount <= 255 {
				var b uint8
				if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
					panic(err)
				}
				m.Objects[i].Faces[f] = int(b)
			} else if vertCount <= 65535 {
				var b uint16
				if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
					panic(err)
				}
				m.Objects[i].Faces[f] = int(b)
			} else {
				var b uint32
				if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
					panic(err)
				}
				m.Objects[i].Faces[f] = int(b)
			}
		}
	}
	//Read the face indices.

	return &m, nil
}