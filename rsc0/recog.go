package rsc0

import (
	"io"
)

func Recog(r io.ReadSeeker) bool {
	b := make([]byte, 39)
	r.Read(b)
	if string(b[:4]) == "RSC0" && string(b[38]) == "/" {
		return true
	}
	return false
}