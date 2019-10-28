package main

import (
	"flag"
	"path/filepath"
	"os"
	"crypto/sha256"
	"io/ioutil"
)

var dirA string
var dirB string
var destination string
var count int
var sums [][32]byte

func main() {
	flag.Parse()
	if err := filepath.Walk(dirA, preWalker); err != nil {
		panic(err)
	}
	sums = make([][32]byte, count)
	count = 0
	if err := filepath.Walk(dirA, firstWalker); err != nil {
		panic(err)
	}
	if err := filepath.Walk(dirB, secondWalker); err != nil {
		panic(err)
	}

}

func init() {
	flag.StringVar(&dirA, "a", "./a", "The first directory that is walked.")
	flag.StringVar(&dirB, "b", "./b", "The second directory that is walked.")
	flag.StringVar(&destination, "diff", "./diff", "The save directory.")
}

func preWalker(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	count++
	return nil
}

func firstWalker(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(b)
	sums[count] = sum
	count++
	return nil
}

func secondWalker(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(b)
	for _, s := range sums {
		if compareSlice(sum, s) {
			return nil
		}
	}
	name := filepath.Join(destination, info.Name())
	os.MkdirAll(filepath.Dir(name), os.ModePerm)
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		return err
	}
	if err := f.Truncate(int64(len(b))); err != nil {
		return err
	}
	if _, err := f.WriteAt(b, 0); err != nil {
		return err
	}
	return nil
}

func compareSlice(a, b [32]byte) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}