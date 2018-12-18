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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Uploader interface {
	Upload(file string) (bool, error)
}

type S3Uploader struct {
	lastUpload time.Time
	uploadSessionId string
	uploadSessionSecret string
}

type GetSignedUrlForUploadRequest struct {
	uploadSessionId string
	uploadSessionSecret string
}

type GetSignedUrlForUploadResponse struct {
	signedUrl string
}

func NewS3Uploader(uploadSessionId string, uploadSessionSecret string) S3Uploader {
	return S3Uploader {
		uploadSessionId: uploadSessionId,
		uploadSessionSecret: uploadSessionSecret,
	}
}

func (u S3Uploader) Upload(file string) (bool, error) {
	if time.Since(u.lastUpload).Minutes() <= 5 {
		return false, nil
	}

	u.lastUpload = time.Now()

	fmt.Printf("Uploading file: %#v\n", file)

	url, err := u.getUploadUrl()

	if err != nil {
		return false, err
	}

	fmt.Printf("Upload url: %#v\n", url)

	fileReader, err := os.Open(file)
	if err != nil {
		return false, err
	}

	err = u.uploadFile(url, fileReader)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (u S3Uploader) getUploadUrl() (string, error) {
	// Probally want to do something different
	message := GetSignedUrlForUploadRequest{
		u.uploadSessionId,
		u.uploadSessionSecret,
	}

	b, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Post(
		"https://api.stellarisinsights.com/v1/signed-url", // maybe this should include the uploadSessionId in the url?
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result GetSignedUrlForUploadResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	fmt.Printf("%#v", result)

	return result.signedUrl, nil
}

func (u S3Uploader) uploadFile(signedUrl string, file io.Reader) (error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("PUT", signedUrl, file)
    if err != nil {
        return err
	}

	resp, err := client.Do(req)
    if err != nil {
        return err
	}
	
	fmt.Printf("%#v", resp)

	return nil
}