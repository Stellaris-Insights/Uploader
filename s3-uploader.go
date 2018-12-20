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

package uploader

import (
	"fmt"
	"os"
	"time"

	"github.com/stellaris-insights/uploader/api"
)

type S3Uploader struct {
	service api.S3ApiServicer
	lastUpload time.Time
	uploadSessionId string
	uploadSessionSecret string
}

func NewS3Uploader(service api.S3ApiServicer, uploadSessionId string, uploadSessionSecret string) S3Uploader {
	return S3Uploader {
		service,
		time.Now(),
		uploadSessionId,
		uploadSessionSecret,
	}
}

func (u S3Uploader) Upload(file string) (bool, error) {
	if time.Since(u.lastUpload).Minutes() <= 5 {
		return false, nil
	}

	u.lastUpload = time.Now()

	fmt.Printf("Uploading file: %#v\n", file)

	url, err := u.service.GetSignedUploadSaveGameURL(u.uploadSessionId, u.uploadSessionSecret)

	if err != nil {
		return false, err
	}

	fmt.Printf("Upload url: %#v\n", url)

	fileReader, err := os.Open(file)
	if err != nil {
		return false, err
	}

	err = u.service.UploadSaveGame(url, fileReader)
	if err != nil {
		return false, err
	}

	return true, nil
}
