package willie_test

import (
	"net/url"
	"testing"

	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_Sessions(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	res := w.Request("/sessions/get").Get()
	r.NotContains(res.Body.String(), "mark")
	w.Request("/sessions/set").Post(User{Name: "mark"})
	res = w.Request("/sessions/get").Get()
	r.Contains(res.Body.String(), "mark")
}

func Test_Request_URL_Params(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	req := w.Request("/foo?a=%s&b=%s", "A", "B")
	r.Equal("/foo?a=A&b=B", req.URL)
}

func Test_Request_Copies_Headers(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())
	w.Headers["foo"] = "bar"

	req := w.Request("/")
	r.Equal("bar", req.Headers["foo"])
}

func Test_Request_Headers_Dont_Overwrite_App_Headers(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())
	w.Headers["foo"] = "bar"

	req := w.Request("/")
	req.Headers["foo"] = "baz"
	r.Equal("baz", req.Headers["foo"])
	r.Equal("bar", w.Headers["foo"])
}

func Test_Get(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	req := w.Request("/get")
	r.Equal("/get", req.URL)

	res := req.Get()
	r.Equal(201, res.Code)
	r.Contains(res.Body.String(), "METHOD:GET")
	r.Contains(res.Body.String(), "Hello from Get!")
}

func Test_Delete(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	req := w.Request("/delete")
	r.Equal("/delete", req.URL)

	res := req.Delete()
	r.Contains(res.Body.String(), "METHOD:DELETE")
	r.Contains(res.Body.String(), "Goodbye")
}

func Test_Post_Struct(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	req := w.Request("/post")
	res := req.Post(User{Name: "Mark"})
	r.Contains(res.Body.String(), "METHOD:POST")
	r.Contains(res.Body.String(), "NAME:Mark")
}

func Test_Post_Struct_Pointer(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	req := w.Request("/post")
	res := req.Post(&User{Name: "Mark"})
	r.Contains(res.Body.String(), "METHOD:POST")
	r.Contains(res.Body.String(), "NAME:Mark")
}

func Test_Post_Values(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	req := w.Request("/post")
	vals := url.Values{}
	vals.Add("Name", "Mark")
	res := req.Post(vals)
	r.Contains(res.Body.String(), "METHOD:POST")
	r.Contains(res.Body.String(), "NAME:Mark")
}

func Test_Put(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	req := w.Request("/put")
	res := req.Put(User{Name: "Mark"})
	r.Contains(res.Body.String(), "METHOD:PUT")
	r.Contains(res.Body.String(), "NAME:Mark")
}

func Test_Put_Struct_Pointer(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	req := w.Request("/put")
	res := req.Put(&User{Name: "Mark"})
	r.Contains(res.Body.String(), "METHOD:PUT")
	r.Contains(res.Body.String(), "NAME:Mark")
}

func Test_Put_Values(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	req := w.Request("/put")
	vals := url.Values{}
	vals.Add("Name", "Mark")
	res := req.Put(vals)
	r.Contains(res.Body.String(), "METHOD:PUT")
	r.Contains(res.Body.String(), "NAME:Mark")
}

func Test_Put_Struct(t *testing.T) {
	r := require.New(t)
	w := willie.New(App())

	req := w.Request("/put")
	vals := struct {
		Name  string
		Email string
	}{"Antonio", "ap@ap.com"}

	res := req.Put(vals)

	r.Contains(res.Body.String(), "METHOD:PUT")
	r.Contains(res.Body.String(), "EMAIL:ap@ap.com")
}
