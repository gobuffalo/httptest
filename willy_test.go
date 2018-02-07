package willie_test

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
)

var Store sessions.Store = sessions.NewCookieStore([]byte("something-very-secret"))

type User struct {
	Name string `form:"name"`
}

func App() http.Handler {
	p := pat.New()
	p.Get("/get", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(201)
		fmt.Fprintln(res, "METHOD:"+req.Method)
		fmt.Fprint(res, "Hello from Get!")
	})
	p.Delete("/delete", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(201)
		fmt.Fprintln(res, "METHOD:"+req.Method)
		fmt.Fprint(res, "Goodbye")
	})
	p.Post("/post", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, "METHOD:"+req.Method)
		renderForm(res, req)
	})
	p.Put("/put", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, "METHOD:"+req.Method)
		renderForm(res, req)
	})
	p.Post("/sessions/set", func(res http.ResponseWriter, req *http.Request) {
		applyFormToSession(res, req)
	})
	p.Get("/sessions/get", func(res http.ResponseWriter, req *http.Request) {
		renderSession(res, req)
	})
	return p
}

func renderForm(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	for k := range req.PostForm {
		fmt.Fprintln(res, strings.ToUpper(k)+":"+req.PostFormValue(k))
	}
}

func applyFormToSession(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	for k := range req.PostForm {
		sess, _ := Store.Get(req, "my-session")
		sess.Values[k] = req.PostFormValue(k)
		sess.Save(req, res)
	}
}

func renderSession(res http.ResponseWriter, req *http.Request) {
	sess, _ := Store.Get(req, "my-session")
	for k := range sess.Values {
		fmt.Fprintln(res, strings.ToUpper(k.(string))+":"+sess.Values[k].(string))
	}
}
