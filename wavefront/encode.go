package wavefront

import (
	"io"
	"github.com/rbxb/smeckle/model"
	"strconv"
)

func f32toA(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

func Encode(m * model.Model, w io.Writer) error {
	for i := 0; i < len(m.Verts) / 3; i++ {
		w.Write([]byte("v " + 
			f32toA(m.Verts[i*3+0]) + " " + 
			f32toA(m.Verts[i*3+1]) + " " + 
			f32toA(m.Verts[i*3+2]) + " \r\n"))
	}

	w.Write([]byte("\r\n"))

	for i := 0; i < len(m.Norms) / 3; i++ {
		w.Write([]byte("vn " + 
			f32toA(m.Norms[i*3+0]) + " " + 
			f32toA(m.Norms[i*3+1]) + " " + 
			f32toA(m.Norms[i*3+2]) + " \r\n"))
	}

	w.Write([]byte("\r\n"))

	for i := 0; i < len(m.UVs) / 2; i++ {
		w.Write([]byte("vt " + 
			f32toA(m.UVs[i*2+0]) + " " + 
			f32toA(m.UVs[i*2+1]) + " \r\n"))
	}

	w.Write([]byte("\r\n"))

	for i, object := range m.Objects {
		w.Write([]byte("g " + strconv.Itoa(i) + " \r\n"))
		for f := 0; f < len(object.Faces) / 3; f++ {
			v := strconv.Itoa(object.Faces[f*3+0] + 1)
			w.Write([]byte("f " + v + "/" + v + "/" + v))
			v = strconv.Itoa(object.Faces[f*3+1] + 1)
			w.Write([]byte(" " + v + "/" + v + "/" + v))
			v = strconv.Itoa(object.Faces[f*3+2] + 1)
			w.Write([]byte(" " + v + "/" + v + "/" + v + " \r\n"))
		}
	}

	return nil
}