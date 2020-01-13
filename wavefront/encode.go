// Package wavefront is used to encode a 3D model into the wavefront (.obj) format.
package wavefront

import (
	"io"
	"strconv"

	"github.com/rbxb/smeckle/model"
)

func f32toA(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

// Encode encodes a 3D model in wavefront format.
func Encode(model *model.Model, w io.Writer) {
	for i := 0; i < len(model.Verts)/3; i++ {
		w.Write([]byte("v " +
			f32toA(model.Verts[i*3+0]) + " " +
			f32toA(model.Verts[i*3+1]) + " " +
			f32toA(model.Verts[i*3+2]) + " \r\n"))
	}

	w.Write([]byte("\r\n"))

	for i := 0; i < len(model.Norms)/3; i++ {
		w.Write([]byte("vn " +
			f32toA(model.Norms[i*3+0]) + " " +
			f32toA(model.Norms[i*3+1]) + " " +
			f32toA(model.Norms[i*3+2]) + " \r\n"))
	}

	w.Write([]byte("\r\n"))

	for i := 0; i < len(model.UVs)/2; i++ {
		w.Write([]byte("vt " +
			f32toA(model.UVs[i*2+0]) + " " +
			f32toA(model.UVs[i*2+1]) + " \r\n"))
	}

	w.Write([]byte("\r\n"))

	for i, object := range model.Objects {
		w.Write([]byte("g " + strconv.Itoa(i) + " \r\n"))
		for f := 0; f < len(object.Faces)/3; f++ {
			v := strconv.Itoa(object.Faces[f*3+0] + 1)
			w.Write([]byte("f " + v + "/" + v + "/" + v))
			v = strconv.Itoa(object.Faces[f*3+1] + 1)
			w.Write([]byte(" " + v + "/" + v + "/" + v))
			v = strconv.Itoa(object.Faces[f*3+2] + 1)
			w.Write([]byte(" " + v + "/" + v + "/" + v + " \r\n"))
		}
	}
}
