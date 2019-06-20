package main

import (
	"flag"
	"os"
	"bytes"
	"path/filepath"
	"github.com/rbxb/smeckle/webp"
	"github.com/rbxb/smeckle/png"
)

var source string
var ex string

func init() {
	flag.StringVar(&source, "source", "./source", "The source directory or file. (./source)")
	flag.StringVar(&ex, "ex", "./ex", "The save directory. (./ex)")
}

func main() {
	flag.Parse()
	ex = filepath.Join(ex, "textures")
	info, err := os.Stat(source)
	if err != nil {
		panic(err)
	}
	if info.IsDir() {
		if err := filepath.Walk(source, convertFile); err != nil {
			panic(err)
		}
	} else {
		if err := convertFile(source, info, nil); err != nil {
			panic(err)
		}
	}
}

func writeFile(name string, buf * bytes.Buffer) error {
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

func convertFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	defer f.Close()
	if err != nil {
		return err
	}
	if !webp.Recog(f) {
		return nil
	}
	f.Seek(0,0)
	b := make([]byte, info.Size())
	if _, err := f.Read(b); err != nil {
		return err
	}
	f.Close()
	r := bytes.NewReader(b)
	diffuse, specular, err := webp.Decode(r)
	if err != nil {
		return err
	}
	diffuseName := filepath.Join(ex, info.Name() + ".png")
	specularName := filepath.Join(ex, "specular", info.Name() + "_spec.png")
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(diffuse, buf); err != nil {
		return err
	}
	if err := writeFile(diffuseName, buf); err != nil {
		return err
	}
	buf = bytes.NewBuffer(nil)
	if err := png.Encode(specular, buf); err != nil {
		return err
	}
	if err := writeFile(specularName, buf); err != nil {
		return err
	}
	return nil
}