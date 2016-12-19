package goserver

import (
	"net/http"
)

type Server struct {
	Handlers       map[string]func(http.ResponseWriter, *http.Request, string, Session)
	SessionManager *SessionManager
}

func NewServer() *Server {
	server := new(Server)
	server.Handlers = make(map[string]func(http.ResponseWriter, *http.Request, string, Session))
	server.SessionManager = NewSessionManager("sessionid", 3600)
	return server
}

func (server *Server) AddHandler(path string, handler func(http.ResponseWriter, *http.Request, string, Session)) {
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

func (server *Server) makeHandler(path string, handler func(http.ResponseWriter, *http.Request, string, Session)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, r.URL.Path[len(path):], server.SessionManager.SessionStart(w, r))
	}
}
