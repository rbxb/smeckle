package walker

import (
	"os"
	"path/filepath"
	"sync"
)

// A WalkCallbackFunc is called on every file when walking a directory.
type WalkCallbackFunc func(string, os.FileInfo)

// Walk walks a directory, creating new goroutines to handle each file.
// The parameter threads inicated the maximum number of goroutines.
func Walk(path string, threads int, callback WalkCallbackFunc) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		walker, wait := createWalker(threads, callback)
		if err := filepath.Walk(path, walker); err != nil {
			return err
		}
		wait.Wait()
	} else {
		callback(path, info)
	}
	return nil
}

func createWalker(threads int, callback WalkCallbackFunc) (filepath.WalkFunc, *sync.WaitGroup) {
	limit := make(chan byte, threads)
	wait := sync.WaitGroup{}
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		limit <- 0
		wait.Add(1)
		go func() {
			callback(path, info)
			<-limit
			wait.Done()
		}()
		return nil
	}, &wait
}
