package webp

import (
	"io"
)

func Recog(r io.ReadSeeker) bool {
	b := make([]byte, 45)
	r.Read(b)
	if string(b[32:36]) == "RIFF" && string(b[40:44]) == "WEBP" {
		return true
	}
	return false
}