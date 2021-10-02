// +build !windows
package filetracker

import (
	"context"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/St0rmPetrel/chain-log-cat/utils"
	"github.com/fsnotify/fsnotify"
)

type tracker struct {
	root   string
	patern string
	age    time.Duration
}

func New(root, patern string, age time.Duration) *tracker {
	return &tracker{
		root:   root,
		patern: patern,
		age:    age,
	}
}

func addInBuffer(buffer map[time.Duration][]byte,
	data []byte, t time.Duration) {
	b, exist := buffer[t]
	if !exist {
		buffer[t] = data
		return
	}
	buffer[t] = append(b, data...)
	return
}

func (t *tracker) TrackChanges() ([]byte, error) {
	buffer := t.trackChanges()
	keys := make([]time.Duration, 0, len(buffer))
	for k, _ := range buffer {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	ret := []byte{}
	for _, k := range keys {
		ret = append(ret, buffer[k]...)
	}
	return ret, nil
}

func (t *tracker) trackChanges() map[time.Duration][]byte {
	exists_files := must(t.findActualFiles())
	buffer := make(map[time.Duration][]byte)

	interruptSignal := make(chan os.Signal)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	startTime := time.Now()
	var mu sync.Mutex
	for _, file := range exists_files {
		go func(file string) {
			b, t, err := trackFileChanges(ctx, file)
			if err != nil {
				return
			}
			mu.Lock()
			addInBuffer(buffer, b, t.Sub(startTime))
			mu.Unlock()
		}(file)
	}
	<-interruptSignal
	cancel()

	new_files := utils.Except(must(t.findActualFiles()), exists_files)
	for _, file := range new_files {
		b, t, err := readNewFile(file)
		if err != nil {
			continue
		}
		addInBuffer(buffer, b, t.Sub(startTime))
	}
	return buffer
}

func must(data []string, err error) []string {
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return data
}

func (t *tracker) findActualFiles() (files []string, err error) {
	files = make([]string, 0)

	now := time.Now()
	validFile := regexp.MustCompile(t.patern)
	err = filepath.Walk(t.root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && validFile.MatchString(path) &&
				(now.Sub(info.ModTime()) <= t.age) {
				files = append(files, path)
			}
			return nil
		})

	return files, err
}

// Wait while in file something changed
// for example in file sombody write something
// or somebody just touch a file
// TODO change functio for return err
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

func readNewFile(filename string) ([]byte, time.Time, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, time.Now(), err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, time.Now(), err
	}
	t := stat.ModTime()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, time.Now(), err
	}
	return buf, t, nil
}

func trackFileChanges(ctx context.Context,
	filename string) ([]byte, time.Time, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, time.Now(), err
	}
	defer file.Close()

	// Set offset of the next Read to the end of a file
	stat, err := file.Stat()
	if err != nil {
		return nil, time.Now(), err
	}
	filesize := stat.Size()
	file.Seek(filesize, os.SEEK_SET)

	waitFileChange(ctx, filename)
	t := time.Now()

	<-ctx.Done()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return buf, t, nil
}
