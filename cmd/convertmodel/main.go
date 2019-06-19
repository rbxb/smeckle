package main

import (
	"flag"
	"os"
	"bytes"
	"path/filepath"
	"github.com/rbxb/smeckle/rsc0"
	"github.com/rbxb/smeckle/wavefront"
)

var source string
var ex string

func init() {
	flag.StringVar(&source, "source", "./source", "The source directory or file. (./source)")
	flag.StringVar(&ex, "ex", "./ex", "The save directory. (./ex)")
}

func main() {
	flag.Parse()
	ex = filepath.Join(ex, "models")
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
	if !rsc0.Recog(f) {
		return nil
	}
	f.Seek(0,0)
	b := make([]byte, info.Size())
	if _, err := f.Read(b); err != nil {
		return err
	}
	f.Close()
	r := bytes.NewReader(b)
	m, err := rsc0.Decode(r)
	if err != nil {
		return err
	}
	name := filepath.Join(ex, m.Name + ".obj")
	os.MkdirAll(filepath.Dir(name), os.ModePerm)
	f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	if err := wavefront.Encode(m, buf); err != nil {
		return err
	}
	f.Truncate(int64(buf.Len()))
	if _, err := buf.WriteTo(f); err != nil {
		return err
	}
	return nil
}