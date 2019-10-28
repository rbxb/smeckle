package main

import (
	"flag"
	"os"
	"bytes"
	"sync"
	"path/filepath"
	"github.com/rbxb/smeckle/rsc0"
	"github.com/rbxb/smeckle/wavefront"
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
	ex = filepath.Join(ex, "models")
	info, err := os.Stat(source)
	if err != nil {
		panic(err)
	}
	if info.IsDir() {
		limit = make(chan byte, threads)
		wait = sync.WaitGroup{}
		if err := filepath.Walk(source, walker); err != nil {
			panic(err)
		}
		wait.Wait()
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
	limit <- 0
	wait.Add(1)
	go func(){
		if err := convertFile(path, info); err != nil {
			panic(err)
		}
		<- limit
		wait.Done()
	}()
	return nil
}

func convertFile(path string, info os.FileInfo) error {
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