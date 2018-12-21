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
	"path/filepath"
	"strings"
	"time"

	"github.com/stellaris-insights/uploader/api"
)

// S3Uploader is a struct that describes the uploader used to upload to S3.
type S3Uploader struct {
	service             api.S3ApiServicer
	uploadSessionID     string
	uploadSessionSecret string
	basePath            string
	lastUpload          time.Time
}

// NewS3Uploader creates a new instance of S3Uploader.
func NewS3Uploader(
	service api.S3ApiServicer,
	uploadSessionID string,
	uploadSessionSecret string,
	basePath string,
) S3Uploader {
	return S3Uploader{
		service,
		uploadSessionID,
		uploadSessionSecret,
		basePath,
		time.Now(),
	}
}

// Upload will upload the given file to the Stellaris Insights API.
func (u S3Uploader) Upload(file string) (bool, error) {
	if time.Since(u.lastUpload).Minutes() <= 5 {
		return false, nil
	}

	u.lastUpload = time.Now()

	fmt.Printf("Uploading file: %#v\n", file)

	url, err := u.service.GetSignedUploadSaveGameURL(u.uploadSessionID, u.uploadSessionSecret)

	if err != nil {
		return false, err
	}

	absFile, err := filepath.Abs(file)
	if err != nil {
		return false, err
	}

	if !strings.HasPrefix(file, u.basePath) {
		return false, fmt.Errorf("wrong file path: %s is not within the current basepath of %s", file, u.basePath)
	}

	fmt.Printf("Upload url: %#v\n", url)

	/* #nosec */
	fileReader, err := os.Open(absFile)
	if err != nil {
		return false, err
	}

	err = u.service.UploadSaveGame(url, fileReader)
	if err != nil {
		return false, err
	}

	return true, nil
}
