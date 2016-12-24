package goserver

import (
	"fmt"
	"net/http"
)

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
	AllowAnonymous bool
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
		server.AddHandlerFrom(HandlerInfo{"/login", loginGetHandler, loginPostHandler, true})
		server.AddHandler("/logout", logoutHandler)
	}
	return server
}

func (server *Server) AddHandler(path string, handler func(http.ResponseWriter, *http.Request, Info)) {
	server.AddHandlerFrom(HandlerInfo{path, handler, nil, false})
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
		session := server.SessionManager.SessionStart(w, r)
		_, ok := session.Get("user")
		if !ok && server.RequireLogin && !handlerInfo.AllowAnonymous {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		info := Info{server, session, r.URL.Path[len(handlerInfo.Path):]}
		if r.Method == "POST" && handlerInfo.PostHandler != nil {
			handlerInfo.PostHandler(w, r, info)
		} else if handlerInfo.GetHandler != nil {
			handlerInfo.GetHandler(w, r, info)
		}
	}
}

func loginGetHandler(w http.ResponseWriter, r *http.Request, info Info) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<form method=\"post\" action=\"/login\"><input type=\"text\" name=\"user\"><input type=\"password\" name=\"password\"><input type=\"submit\" name=\"submit\"></form>")
}

func loginPostHandler(w http.ResponseWriter, r *http.Request, info Info) {
	err := r.ParseForm()
	if err != nil {
		loginGetHandler(w, r, info)
		return
	}
	users, _ := r.Form["user"]
	passwords, _ := r.Form["password"]
	if len(users) > 0 && len(passwords) > 0 && info.Server.Login(users[0], passwords[0], info.Session) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	loginGetHandler(w, r, info)
}

func (server *Server) Login(user string, password string, session Session) bool {
	if server.Users[user] == password {
		session.Set("user", user)
		return true
	}
	return false
}

func logoutHandler(w http.ResponseWriter, r *http.Request, info Info) {
	info.Server.Logout(info.Session)
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func (server *Server) Logout(session Session) {
	session.Delete("user")
}

func (server *Server) AddUser(user string, password string) {
	server.Users[user] = password
}

func (info *Info) User() string {
	user, ok := info.Session.Get("user")
	if !ok {
		return ""
	}
	return user.(string)
}
