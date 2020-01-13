package main

import (
	"bytes"
	"flag"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/rbxb/smeckle/texture"
	"github.com/rbxb/smeckle/walker"
)

var source string
var ex string
var threads int

var limit chan byte
var wait sync.WaitGroup

func init() {
	flag.StringVar(&source, "source", "./source", "The source directory or file.")
	flag.StringVar(&ex, "ex", "./ex", "The save directory.")
	flag.IntVar(&threads, "threads", 8, "The number of threads to run on.")
}

func main() {
	flag.Parse()
	ex = filepath.Join(ex, "textures")
	walker.Walk(source, threads, convertFile)
}

func convertFile(path string, info os.FileInfo) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	if !isTexture(f) {
		return
	}
	f.Seek(0, 0)
	b := make([]byte, info.Size())
	if _, err := f.Read(b); err != nil {
		panic(err)
	}
	f.Close()
	r := bytes.NewReader(b)
	diffuse, specular := texture.ConvertTexture(r)
	diffuseName := filepath.Join(ex, info.Name()+".png")
	specularName := filepath.Join(ex, "specular", info.Name()+"_spec.png")
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, diffuse); err != nil {
		panic(err)
	}
	if err := writeFile(diffuseName, buf); err != nil {
		panic(err)
	}
	buf = bytes.NewBuffer(nil)
	if err := png.Encode(buf, specular); err != nil {
		panic(err)
	}
	if err := writeFile(specularName, buf); err != nil {
		panic(err)
	}
}

func isTexture(r io.ReadSeeker) bool {
	b := make([]byte, 45)
	r.Read(b)
	if string(b[32:36]) == "RIFF" && string(b[40:44]) == "WEBP" {
		return true
	}
	return false
}

func writeFile(name string, buf *bytes.Buffer) error {
	os.MkdirAll(filepath.Dir(name), os.ModePerm)
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Truncate(int64(buf.Len()))
	_, err = buf.WriteTo(f)
	return err
}
