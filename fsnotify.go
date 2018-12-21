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

// Package uploader is used to handle save game processing
// and upload them to Stellaris Insights.
package uploader

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

// FSNotifier is an interface used for file system watchers.
type FSNotifier interface {
	Close() error
	Add(string) error
	Remove(string) error

	Events() chan fsnotify.Event
	Errors() chan error
}

// FSNotify is a struct that is used to wrap the fsnotify libraries Watcher.
// This is primarily to abstract away the use for testing, but needed to
// make the channels functions.
type FSNotify struct {
	watcher *fsnotify.Watcher
}

// NewFSNotifyWrapper creates a new instance of FSNotify.
func NewFSNotifyWrapper() FSNotify {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return FSNotify{
		watcher: watcher,
	}
}

// Close the current fsnotify.Watcher.
func (fsn FSNotify) Close() error {
	return fsn.watcher.Close()
}

// Add a new path to watch to fsnotify.Watcher.
func (fsn FSNotify) Add(path string) error {
	return fsn.watcher.Add(path)
}

// Remove a path that is being watched by fsnotify.Watcher.
func (fsn FSNotify) Remove(path string) error {
	return fsn.watcher.Remove(path)
}

// Events contains the channel listening for file system changes.
func (fsn FSNotify) Events() chan fsnotify.Event {
	return fsn.watcher.Events
}

// Errors contains the channel listening for errors.
func (fsn FSNotify) Errors() chan error {
	return fsn.watcher.Errors
}
