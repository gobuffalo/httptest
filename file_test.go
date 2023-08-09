package httptest

import (
	"fmt"
	"github.com/gobuffalo/httptest/testassets"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FileUpload(t *testing.T) {
	r := require.New(t)
	w := New(App())

	foo := func(filename, expectedType string) {
		f := struct {
			Name string
		}{"Foo"}

		rr, err := testassets.FS.Open(filename)
		r.NoError(err)
		wf := File{
			ParamName: "MyFile",
			FileName:  filename,
			Reader:    rr,
		}
		res, err := w.HTML("/up").MultiPartPost(f, wf)
		r.NoError(err)
		r.Equal(200, res.Code)
		r.Equal(fmt.Sprintf("Foo\n%s\n%s\n", filename, expectedType), res.Body.String())
	}

	foo("test.jpg", "image/jpeg")
	foo("test.png", "image/png")
	foo("test.pdf", "application/pdf")
	foo("embed.go", "text/plain; charset=utf-8")
	foo("random.bin", "application/octet-stream")
}
