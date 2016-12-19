package main

import "net/http"

func main() {
	server := NewServer()
	server.AddHandler("/view/", viewHandler)
	server.Serve()
}

func viewHandler(w http.ResponseWriter, r *http.Request, path string) {
	w.Write([]byte("Jo hollo"))
}

type Server struct {
	Handlers map[string]func(http.ResponseWriter, *http.Request, string)
}

func NewServer() *Server {
	server := new(Server)
	server.Handlers = make(map[string]func(http.ResponseWriter, *http.Request, string))
	return server
}

func (server *Server) AddHandler(path string, handler func(http.ResponseWriter, *http.Request, string)) {
	server.Handlers[path] = handler
}

func (server *Server) ServeOnPort(port string) {
	for path, handler := range server.Handlers {
		http.HandleFunc(path, makeHandler(path, handler))
	}
	http.ListenAndServe(port, nil)
}

func (server *Server) Serve() {
	server.ServeOnPort(":8080")
}

func makeHandler(path string, handler func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, r.URL.Path[len(path):])
	}
}
