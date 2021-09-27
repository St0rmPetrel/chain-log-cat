package main

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/St0rmPetrel/chain-log-cat/utils"
	"github.com/fsnotify/fsnotify"
)

func main() {
	//var logs map[time.Time][]byte = make(map[time.Time][]byte)
	exists_files, err := FindFreshLogFiles("..", 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	new_files, err := FindFreshLogFiles("..", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	new_files = utils.Except(new_files, exists_files)
	for _, file := range exists_files {
		fmt.Println("|" + file + "|")
	}
}

func trackChanges(ctx context.Context, filename string) (time.Time, []byte) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Set offset of the next Read to the end of a file
	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	filesize := stat.Size()
	file.Seek(filesize, os.SEEK_SET)

	waitFileChange(ctx, filename)
	t := time.Now()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return t, buf
}

// Wait while in file something changed
// for example in file sombody write something
// or somebody just touch a file
func waitFileChange(ctx context.Context, filename string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}
	select {
	case _, ok := <-watcher.Events:
		if !ok {
			log.Fatal(err)
		}
		return
	case <-watcher.Errors:
		log.Fatal(err)
	case <-ctx.Done():
		return
	}
}

// Add above package
// //go:build !windows
// // +build !windows

func FindFreshLogFiles(root string,
	age time.Duration) (files []string, err error) {
	files = make([]string, 0)

	now := time.Now()
	err = filepath.Walk(root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && isLog(path) &&
				(now.Sub(info.ModTime()) <= age) {
				files = append(files, path)
			}
			return nil
		})

	return files, err
}

// Hard conde :(
func isLog(path string) bool {
	date := "20[0-9][0-9]-[0-1][0-9]-[0-3][0-9]"
	ret, err := filepath.
		Match(
			"/opt/syslog/*/*/"+date+"/*-*-*h.log",
			path,
		)
	if err != nil {
		//return false
		return ret
	}
	//return ret
	return true
}
