package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/rbxb/smeckle/model"
	"github.com/rbxb/smeckle/walker"
	"github.com/rbxb/smeckle/wavefront"
)

var source string
var ex string
var threads int

func init() {
	flag.StringVar(&source, "source", "./source", "The source directory or file.")
	flag.StringVar(&ex, "ex", "./ex", "The save directory.")
	flag.IntVar(&threads, "threads", 8, "The number of threads to run on.")
}

func main() {
	flag.Parse()
	ex = filepath.Join(ex, "models")
	walker.Walk(source, threads, convertFile)
}

func convertFile(path string, info os.FileInfo) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	if !isModel(f) {
		return
	}
	f.Seek(0, 0)
	b := make([]byte, info.Size())
	if _, err := f.Read(b); err != nil {
		panic(err)
	}
	f.Close()
	r := bytes.NewReader(b)
	m := model.ConvertModel(r)
	name := filepath.Join(ex, m.Name+".obj")
	os.MkdirAll(filepath.Dir(name), os.ModePerm)
	f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	wavefront.Encode(m, buf)
	f.Truncate(int64(buf.Len()))
	if _, err := buf.WriteTo(f); err != nil {
		panic(err)
	}
}

func isModel(r io.ReadSeeker) bool {
	b := make([]byte, 39)
	r.Read(b)
	if string(b[:4]) == "RSC0" && string(b[38]) == "/" {
		return true
	}
	return false
}
