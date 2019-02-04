package httptest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/require"
)

type mux struct {
	routes map[string]map[string]http.HandlerFunc
}

func (m mux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if len(m.routes) == 0 {
		m.routes = map[string]map[string]http.HandlerFunc{}
	}
	verb := req.Method
	vm, ok := m.routes[verb]
	if !ok {
		res.WriteHeader(500)
		fmt.Fprintf(res, "couldn't find map for %s", verb)
		return
	}
	if h, ok := vm[req.URL.Path]; ok {
		h(res, req)
		return
	}
	res.WriteHeader(500)
	fmt.Fprintf(res, "couldn't find map for %s", req.URL.Path)
}

func (m *mux) Handle(verb string, route string, h http.HandlerFunc) {
	if len(m.routes) == 0 {
		m.routes = map[string]map[string]http.HandlerFunc{}
	}
	vm, ok := m.routes[verb]
	if !ok {
		vm = map[string]http.HandlerFunc{}
		m.routes[verb] = vm
	}

	vm[route] = h

}

var Store sessions.Store = sessions.NewCookieStore([]byte("something-very-secret"))

type User struct {
	Name string `form:"name" xml:"name"`
}

func App() http.Handler {
	p := &mux{}
	p.Handle("GET", "/get", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(201)
		fmt.Fprintln(res, "METHOD:"+req.Method)
		fmt.Fprint(res, "Hello from Get!")
	})
	p.Handle("DELETE", "/delete", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(201)
		fmt.Fprintln(res, "METHOD:"+req.Method)
		fmt.Fprint(res, "Goodbye")
	})
	p.Handle("POST", "/post", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, "METHOD:"+req.Method)
		fmt.Fprint(res, "NAME:"+req.PostFormValue("name"))
	})
	p.Handle("PUT", "/put", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, "METHOD:"+req.Method)
		fmt.Fprint(res, "NAME:"+req.PostFormValue("name"))
	})
	p.Handle("POST", "/sessions/set", func(res http.ResponseWriter, req *http.Request) {
		sess, _ := Store.Get(req, "my-session")
		sess.Values["name"] = req.PostFormValue("name")
		sess.Save(req, res)
	})
	p.Handle("GET", "/sessions/get", func(res http.ResponseWriter, req *http.Request) {
		sess, _ := Store.Get(req, "my-session")
		if sess.Values["name"] != nil {
			fmt.Fprint(res, "NAME:"+sess.Values["name"].(string))
		}
	})
	p.Handle("POST", "/up", func(res http.ResponseWriter, req *http.Request) {
		if err := req.ParseMultipartForm(5 * 1024); err != nil {
			res.WriteHeader(500)
			fmt.Fprint(res, err.Error())
		}
		_, h, err := req.FormFile("MyFile")
		if err != nil {
			res.WriteHeader(500)
			fmt.Fprint(res, err.Error())
		}
		fmt.Fprintln(res, req.FormValue("Name"))
		fmt.Fprintln(res, h.Filename)
	})
	return p
}

func Test_Sessions(t *testing.T) {
	r := require.New(t)
	w := New(App())

	res := w.HTML("/sessions/get").Get()
	r.NotContains(res.Body.String(), "mark")
	w.HTML("/sessions/set").Post(User{Name: "mark"})
	res = w.HTML("/sessions/get").Get()
	r.Contains(res.Body.String(), "mark")
}

func Test_Request_URL_Params(t *testing.T) {
	r := require.New(t)
	w := New(App())

	req := w.HTML("/foo?a=%s&b=%s", "A", "B")
	r.Equal("/foo?a=A&b=B", req.URL)
}

func Test_Request_Copies_Headers(t *testing.T) {
	r := require.New(t)
	w := New(App())
	w.Headers["foo"] = "bar"

	req := w.HTML("/")
	r.Equal("bar", req.Headers["foo"])
}
