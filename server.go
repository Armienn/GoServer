package goserver

import (
	"net/http"
)

type Server struct {
	Handlers       map[string]func(http.ResponseWriter, *http.Request, string, Session, interface{})
	SessionManager *SessionManager
	RequireLogin   bool
}

func NewServer(requireLogin bool) *Server {
	server := new(Server)
	server.Handlers = make(map[string]func(http.ResponseWriter, *http.Request, string, Session, interface{}))
	server.SessionManager = NewSessionManager("sessionid", 3600)
	server.RequireLogin = requireLogin
	return server
}

func (server *Server) AddHandler(path string, handler func(http.ResponseWriter, *http.Request, string, Session, interface{})) {
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

func (server *Server) makeHandler(path string, handler func(http.ResponseWriter, *http.Request, string, Session, interface{})) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := server.SessionManager.SessionStart(w, r)
		user, ok := session.Get("user")
		if !ok && server.RequireLogin {
			if r.Method == "POST" {
				handleLogin(w, r, session)
			} else {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			}
		} else {
			handler(w, r, r.URL.Path[len(path):], server.SessionManager.SessionStart(w, r), user)
		}
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request, session Session) {
	if r.URL.Path != "/login" {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	users, _ := r.Form["user"]
	passwords, _ := r.Form["password"]
	if len(users) > 0 && len(passwords) > 0 && userExists(users[0], passwords[0]) {
		session.Set("user", users[0])
		http.Redirect(w, r, "/", http.StatusAccepted)
		return
	}
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func userExists(user string, password string) bool {
	return user == password
}
