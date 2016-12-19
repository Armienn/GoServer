package main

import (
	"net/http"
	"strconv"

	"github.com/Armienn/GoServer/session"
)

func main() {
	server := NewServer()
	server.AddHandler("/view/", viewHandler)
	server.Serve()
}

func viewHandler(w http.ResponseWriter, r *http.Request, path string, session session.Session) {
	count := 0
	value, ok := session.Get("musle")
	if ok {
		count = value.(int)
	}
	count += 1
	session.Set("musle", count)
	w.Write([]byte("Jo hollo" + strconv.Itoa(count)))
}

type Server struct {
	Handlers       map[string]func(http.ResponseWriter, *http.Request, string, session.Session)
	SessionManager *session.SessionManager
}

func NewServer() *Server {
	server := new(Server)
	server.Handlers = make(map[string]func(http.ResponseWriter, *http.Request, string, session.Session))
	server.SessionManager = session.NewSessionManager("sessionid", 3600)
	return server
}

func (server *Server) AddHandler(path string, handler func(http.ResponseWriter, *http.Request, string, session.Session)) {
	server.Handlers[path] = handler
}

func (server *Server) ServeOnPort(port string) {
	for path, handler := range server.Handlers {
		http.HandleFunc(path, server.makeHandler(path, handler))
	}
	http.ListenAndServe(port, nil)
}

func (server *Server) Serve() {
	server.ServeOnPort(":8080")
}

func (server *Server) makeHandler(path string, handler func(http.ResponseWriter, *http.Request, string, session.Session)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, r.URL.Path[len(path):], server.SessionManager.SessionStart(w, r))
	}
}
