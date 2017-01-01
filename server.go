package goserver

import "net/http"

type Server struct {
	Handlers       []HandlerInfo
	SessionManager *SessionManager
	RequireLogin   bool
	Users          map[string]string
}

type HandlerInfo struct {
	Path           string
	GetHandler     func(http.ResponseWriter, *http.Request, Info)
	PostHandler    func(http.ResponseWriter, *http.Request, Info)
	GetRestricted  bool
	PostRestricted bool
}

type Info struct {
	Server  *Server
	Session Session
	Path    string
}

func NewServer(requireLogin bool) *Server {
	server := new(Server)
	server.Handlers = make([]HandlerInfo, 0)
	server.Users = make(map[string]string)
	server.SessionManager = NewSessionManager("sessionid", 3600)
	server.RequireLogin = requireLogin
	if server.RequireLogin {
		server.AddGetHandler("/login", loginGetHandler, false)
		server.AddPostHandler("/login", loginPostHandler, false)
		server.AddHandler("/logout", logoutHandler)
	}
	return server
}

func (server *Server) AddHandler(path string, handler func(http.ResponseWriter, *http.Request, Info)) {
	server.AddGetHandler(path, handler, true)
}

func (server *Server) AddHandlerFrom(handlerInfo HandlerInfo) {
	for i, existingHandler := range server.Handlers {
		if existingHandler.Path == handlerInfo.Path {
			server.Handlers = append(server.Handlers[:i], server.Handlers[i+1:]...)
			break
		}
	}
	server.Handlers = append(server.Handlers, handlerInfo)
}

func (server *Server) AddGetHandler(path string, handler func(http.ResponseWriter, *http.Request, Info), restricted bool) {
	for i, existingHandler := range server.Handlers {
		if existingHandler.Path == path {
			server.Handlers[i].GetHandler = handler
			server.Handlers[i].GetRestricted = restricted
			return
		}
	}
	server.Handlers = append(server.Handlers, HandlerInfo{path, handler, nil, restricted, restricted})
}

func (server *Server) AddPostHandler(path string, handler func(http.ResponseWriter, *http.Request, Info), restricted bool) {
	for i, existingHandler := range server.Handlers {
		if existingHandler.Path == path {
			server.Handlers[i].PostHandler = handler
			server.Handlers[i].PostRestricted = restricted
			return
		}
	}
	server.Handlers = append(server.Handlers, HandlerInfo{path, nil, handler, restricted, restricted})
}

func (server *Server) ServeOnPort(port string) {
	for _, handler := range server.Handlers {
		http.HandleFunc(handler.Path, server.makeHandler(handler))
	}
	http.ListenAndServe(port, nil)
}

func (server *Server) Serve() {
	server.ServeOnPort(":8080")
}

func (server *Server) makeHandler(handlerInfo HandlerInfo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler, restricted := getRelevantHandler(handlerInfo, r)
		session := server.SessionManager.SessionStart(w, r)
		_, ok := session.Get("user")
		if !ok && server.RequireLogin && restricted {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		info := Info{server, session, r.URL.Path[len(handlerInfo.Path):]}
		if handler != nil {
			handler(w, r, info)
		}
	}
}

func getRelevantHandler(handlerInfo HandlerInfo, r *http.Request) (handler func(http.ResponseWriter, *http.Request, Info), restricted bool) {
	handler = handlerInfo.GetHandler
	restricted = handlerInfo.GetRestricted
	if r.Method == "POST" {
		if handlerInfo.PostHandler != nil {
			handler = handlerInfo.PostHandler
		}
		restricted = handlerInfo.PostRestricted
	}
	return
}
