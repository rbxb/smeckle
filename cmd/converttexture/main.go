package main

import (
	"flag"
	"os"
	"bytes"
	"path/filepath"
	"github.com/rbxb/smeckle/webp"
	"github.com/rbxb/smeckle/png"
	"github.com/rbxb/workload"
)

var source string
var ex string
var distr * workload.Distr

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
		distr = workload.NewDistr(4,8,callback)
		if err := filepath.Walk(source, walker); err != nil {
			panic(err)
		}
	} else {
		if err := convertFile(source, info); err != nil {
			panic(err)
		}
	}
}

func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	var a interface{} = &action{path,info}
	distr.Add(&a)
	return nil
}

type action struct {
	path string
	info os.FileInfo
}

func callback(a * interface{}) {
	ac := (*a).(*action)
	if err := convertFile(ac.path, ac.info); err != nil {
		panic(err)
	}
}

func convertFile(path string, info os.FileInfo) error {
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
	image, err := webp.Decode(r)
	if err != nil {
		return err
	}
	name := filepath.Join(ex, info.Name() + ".png")
	os.MkdirAll(filepath.Dir(name), os.ModePerm)
	f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(image, buf); err != nil {
		return err
	}
	f.Truncate(int64(buf.Len()))
	if _, err := buf.WriteTo(f); err != nil {
		return err
	}
	return nil
}