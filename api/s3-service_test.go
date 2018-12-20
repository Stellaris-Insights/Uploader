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

package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stellaris-insights/uploader/api"
	"github.com/stellaris-insights/uploader/testutils"
)

func TestGetSignedUploadUrl(t *testing.T) {
	signedURL := "https://postbin.com"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		testutils.Equals(t, req.URL.String(), "/v1/sessions/abcd/uploads/s3-signed-url")

		r, err := json.Marshal(
			api.GetSignedUploadURLResponse{
				SignedURL: signedURL,
			},
		)
		testutils.Ok(t, err)

		_, err = rw.Write(r)
		testutils.Ok(t, err)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	s := api.NewS3ApiService(server.Client(), server.URL)
	url, err := s.GetSignedUploadSaveGameURL("abcd", "1234")

	testutils.Ok(t, err)
	testutils.Equals(t, signedURL, url)
}
