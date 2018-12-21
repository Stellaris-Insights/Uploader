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

// Package api is a package to access api endpoints for Stellaris Insights.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// S3ApiServicer is an interface for managing
// Stellars Insights S3 api endpoints
type S3ApiServicer interface {
	GetSignedUploadSaveGameURL(string, string) (string, error)
	UploadSaveGame(string, io.Reader) error
}

// S3ApiService is a struct that describes the base data for
// Stellars Insights S3 api endpoints
type S3ApiService struct {
	client  *http.Client
	baseURL string
}

// NewS3ApiService creates a new instance of S3ApiService
func NewS3ApiService(client *http.Client, baseURL string) S3ApiService {
	return S3ApiService{
		client,
		baseURL,
	}
}

// GetSignedUploadURLRequest is a struct for the request body
// to get a signed upload URL for file uploads
type GetSignedUploadURLRequest struct {
	UploadSessionSecret string
}

// GetSignedUploadURLResponse is a struct for the response body
// to get a signed upload URL for file uploads
type GetSignedUploadURLResponse struct {
	SignedURL string
}

// GetSignedUploadSaveGameURL gets a signed upload url
// to upload a save game to S3
func (s3as S3ApiService) GetSignedUploadSaveGameURL(
	uploadSessionID string,
	uploadSessionSecret string) (string, error) {
	// Probally want to do something different
	message := GetSignedUploadURLRequest{
		uploadSessionSecret,
	}

	b, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	resp, err := s3as.client.Post(
		s3as.baseURL+"/v1/sessions/"+uploadSessionID+"/uploads/s3-signed-url",
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return "", err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	var result GetSignedUploadURLResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.SignedURL, nil
}

// UploadSaveGame uploads a file to a signed S3 url
func (s3as S3ApiService) UploadSaveGame(signedURL string, file io.Reader) error {
	req, err := http.NewRequest("PUT", signedURL, file)
	if err != nil {
		return err
	}

	resp, err := s3as.client.Do(req)
	if err != nil {
		return err
	}

	fmt.Printf("%#v", resp)

	return nil
}
