package httptest

import (
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type xBody struct {
	Method  string `xml:"method"`
	Name    string `xml:"name"`
	Message string `xml:"message"`
}

func XMLApp() http.Handler {
	p := &mux{}
	p.Handle("GET", "/get", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(201)
		xml.NewEncoder(res).Encode(xBody{
			Method:  req.Method,
			Message: "Hello from Get!",
		})
	})
	p.Handle("DELETE", "/delete", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(201)
		xml.NewEncoder(res).Encode(xBody{
			Method:  req.Method,
			Message: "Goodbye",
		})
	})
	p.Handle("POST", "/post", func(res http.ResponseWriter, req *http.Request) {
		jb := xBody{}
		xml.NewDecoder(req.Body).Decode(&jb)
		jb.Method = req.Method
		xml.NewEncoder(res).Encode(jb)
	})
	p.Handle("PUT", "/put", func(res http.ResponseWriter, req *http.Request) {
		jb := xBody{}
		xml.NewDecoder(req.Body).Decode(&jb)
		jb.Method = req.Method
		xml.NewEncoder(res).Encode(jb)
	})
	p.Handle("PATCH", "/patch", func(res http.ResponseWriter, req *http.Request) {
		jb := xBody{}
		xml.NewDecoder(req.Body).Decode(&jb)
		jb.Method = req.Method
		xml.NewEncoder(res).Encode(jb)
	})
	p.Handle("POST", "/sessions/set", func(res http.ResponseWriter, req *http.Request) {
		sess, _ := Store.Get(req, "my-session")
		sess.Values["name"] = req.PostFormValue("name")
		sess.Save(req, res)
	})
	p.Handle("GET", "/sessions/get", func(res http.ResponseWriter, req *http.Request) {
		sess, _ := Store.Get(req, "my-session")
		if sess.Values["name"] != nil {
			xml.NewEncoder(res).Encode(xBody{
				Method: req.Method,
				Name:   sess.Values["name"].(string),
			})
		}
	})
	return p
}

func Test_XML_Headers_Dont_Overwrite_App_Headers(t *testing.T) {
	r := require.New(t)
	w := New(XMLApp())
	w.Headers["foo"] = "bar"

	req := w.XML("/")
	req.Headers["foo"] = "baz"
	r.Equal("baz", req.Headers["foo"])
	r.Equal("bar", w.Headers["foo"])
}

func Test_XML_Get(t *testing.T) {
	r := require.New(t)
	w := New(XMLApp())

	req := w.XML("/get")
	r.Equal("/get", req.URL)

	res := req.Get()
	r.Equal(201, res.Code)
	jb := &xBody{}
	res.Bind(jb)
	r.Equal("GET", jb.Method)
	r.Equal("Hello from Get!", jb.Message)
}

func Test_XML_Delete(t *testing.T) {
	r := require.New(t)
	w := New(XMLApp())

	req := w.XML("/delete")
	r.Equal("/delete", req.URL)

	res := req.Delete()
	jb := &xBody{}
	res.Bind(jb)
	r.Equal("DELETE", jb.Method)
	r.Equal("Goodbye", jb.Message)
}

func Test_XML_Post_Struct(t *testing.T) {
	r := require.New(t)
	w := New(XMLApp())

	req := w.XML("/post")
	res := req.Post(User{Name: "Mark"})

	jb := &xBody{}
	res.Bind(jb)
	r.Equal("POST", jb.Method)
	r.Equal("Mark", jb.Name)
}

func Test_XML_Post_Struct_Pointer(t *testing.T) {
	r := require.New(t)
	w := New(XMLApp())

	req := w.XML("/post")
	res := req.Post(&User{Name: "Mark"})

	jb := &xBody{}
	res.Bind(jb)
	r.Equal("POST", jb.Method)
	r.Equal("Mark", jb.Name)
}

func Test_XML_Put(t *testing.T) {
	r := require.New(t)
	w := New(XMLApp())

	req := w.XML("/put")
	res := req.Put(User{Name: "Mark"})

	jb := &xBody{}
	res.Bind(jb)
	r.Equal("PUT", jb.Method)
	r.Equal("Mark", jb.Name)
}

func Test_XML_Put_Struct_Pointer(t *testing.T) {
	r := require.New(t)
	w := New(XMLApp())

	req := w.XML("/put")
	res := req.Put(&User{Name: "Mark"})

	jb := &xBody{}
	res.Bind(jb)
	r.Equal("PUT", jb.Method)
	r.Equal("Mark", jb.Name)
}

func Test_XML_Patch(t *testing.T) {
	r := require.New(t)
	w := New(XMLApp())

	req := w.XML("/patch")
	res := req.Patch(User{Name: "Mark"})

	jb := &xBody{}
	res.Bind(jb)
	r.Equal("PATCH", jb.Method)
	r.Equal("Mark", jb.Name)
}

func Test_XML_Patch_Struct_Pointer(t *testing.T) {
	r := require.New(t)
	w := New(XMLApp())

	req := w.XML("/patch")
	res := req.Patch(&User{Name: "Mark"})

	jb := &xBody{}
	res.Bind(jb)
	r.Equal("PATCH", jb.Method)
	r.Equal("Mark", jb.Name)
}
