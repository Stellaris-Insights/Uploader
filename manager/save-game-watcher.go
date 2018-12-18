// Copyright © 2018 C45tr0 <william.the.developer+stellaris@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

type SaveGameWatcher struct {
	watcher FSNotify
	uploader Uploader
}

func NewSaveGameWatcher(fsn FSNotify, uploader Uploader) SaveGameWatcher {
	return SaveGameWatcher {
		watcher: fsn,
		uploader: uploader,
	}
}

func (w SaveGameWatcher) Start(userdataDir string) {
	var err error

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer w.watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-w.watcher.Events():
				w.processEvent(event)

			// watch for errors
			case err := <-w.watcher.Errors():
				fmt.Println("ERROR: ", err)
			}
		}
	}()

	// fsnotify doesn't support recrusive folder watching yet...
	// https://github.com/fsnotify/fsnotify/issues/18
	// So we need to register every subfolder for watching

	// Get all files in save games folder
	files, err := ioutil.ReadDir(path.Join(userdataDir, "save games"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dirs := []string{}

	// Limit to just directories so we can watch them
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}

	// Watch current existing directories in save games folder
	for _, d := range dirs {
		if err := w.watcher.Add(path.Join(userdataDir, "save games", d)); err != nil {
			fmt.Println(err)
		}
	}

	// Watch parent directory for new folder
	if err := w.watcher.Add(path.Join(userdataDir, "save games")); err != nil {
		fmt.Println(err)
	}

	<-done
}

func (w SaveGameWatcher) processEvent(event fsnotify.Event) {
	fmt.Printf("EVENT! %#v\n", event)
	if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
		if fi, err := os.Stat(event.Name); err == nil && fi.IsDir() {
			fmt.Println("dir")
			if err := w.watcher.Add(event.Name); err != nil {
				fmt.Println(err)
			}
			return
		}

		fmt.Println("create|write")
		uploaded, err := w.uploader.Upload(event.Name)

		if err != nil {
			fmt.Println(err)
		}

		if !uploaded {
			fmt.Println("Failed to upload file")
		}
	}
}
