package httptest

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

type File struct {
	io.Reader
	ParamName string
	FileName  string
}

func (r *Request) MultiPartPost(body interface{}, files ...File) (*Response, error) {
	req, err := newMultipart(r.URL, "POST", body, files...)
	if err != nil {
		return nil, err
	}
	return r.Perform(req), nil
}

func (r *Request) MultiPartPut(body interface{}, files ...File) (*Response, error) {
	req, err := newMultipart(r.URL, "PUT", body, files...)
	if err != nil {
		return nil, err
	}
	return r.Perform(req), nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// this helper method was inspired by this blog post by Matt Aimonetti:
// https://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
func newMultipart(url string, method string, body interface{}, files ...File) (*http.Request, error) {
	bb := &bytes.Buffer{}
	writer := multipart.NewWriter(bb)
	defer writer.Close()
	for _, f := range files {
		fBuffer, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				escapeQuotes(f.ParamName), escapeQuotes(f.FileName)))
		h.Set("Content-Type", http.DetectContentType(fBuffer))
		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, err
		}
		fReader := bytes.NewReader(fBuffer)
		_, err = io.Copy(part, fReader)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range toURLValues(body) {
		for _, vv := range v {
			err := writer.WriteField(k, vv)
			if err != nil {
				return nil, err
			}
		}
	}

	req, err := http.NewRequest(method, url, bb)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
