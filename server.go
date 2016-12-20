package goserver

import (
	"fmt"
	"net/http"
)

type Server struct {
	Handlers []struct {
		path    string
		handler func(*Server, http.ResponseWriter, *http.Request, string, Session, interface{})
	}
	SessionManager *SessionManager
	RequireLogin   bool
	Users          map[string]string
}

func NewServer(requireLogin bool) *Server {
	server := new(Server)
	server.Handlers = make([]struct {
		path    string
		handler func(*Server, http.ResponseWriter, *http.Request, string, Session, interface{})
	}, 0)
	server.Users = make(map[string]string)
	server.SessionManager = NewSessionManager("sessionid", 3600)
	server.RequireLogin = requireLogin
	if server.RequireLogin {
		server.AddHandler("/login", loginHandler)
		server.AddHandler("/logout", logoutHandler)
	}
	return server
}

func (server *Server) AddHandler(path string, handler func(*Server, http.ResponseWriter, *http.Request, string, Session, interface{})) {
	for i, existingHandler := range server.Handlers {
		if existingHandler.path == path {
			server.Handlers = append(server.Handlers[:i], server.Handlers[i+1:]...)
			break
		}
	}
	server.Handlers = append(server.Handlers, struct {
		path    string
		handler func(*Server, http.ResponseWriter, *http.Request, string, Session, interface{})
	}{
		path,
		handler,
	})
}

func (server *Server) ServeOnPort(port string) {
	for _, handler := range server.Handlers {
		http.HandleFunc(handler.path, server.makeHandler(handler.path, handler.handler))
	}
	http.ListenAndServe(port, nil)
}

func (server *Server) Serve() {
	server.ServeOnPort(":8080")
}

func (server *Server) makeHandler(path string, handler func(*Server, http.ResponseWriter, *http.Request, string, Session, interface{})) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := server.SessionManager.SessionStart(w, r)
		user, ok := session.Get("user")
		if !ok && server.RequireLogin && path != "login" && path != "/login" && path != "/login/" {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		} else {
			handler(server, w, r, r.URL.Path[len(path):], server.SessionManager.SessionStart(w, r), user)
		}
	}
}

func loginHandler(server *Server, w http.ResponseWriter, r *http.Request, path string, session Session, user interface{}) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<form method=\"post\" action=\"/login\"><input type=\"text\" name=\"user\"><input type=\"password\" name=\"password\"><input type=\"submit\" name=\"submit\"></form>")
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	users, _ := r.Form["user"]
	passwords, _ := r.Form["password"]
	if len(users) > 0 && len(passwords) > 0 && server.Login(users[0], passwords[0], session) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<form method=\"post\" action=\"/login\"><input type=\"text\" name=\"user\"><input type=\"password\" name=\"password\"><input type=\"submit\" name=\"submit\"></form>")
}

func (server *Server) Login(user string, password string, session Session) bool {
	if server.Users[user] == password {
		session.Set("user", user)
		return true
	}
	return false
}

func logoutHandler(server *Server, w http.ResponseWriter, r *http.Request, path string, session Session, user interface{}) {
	server.Logout(session)
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func (server *Server) Logout(session Session) {
	session.Delete("user")
}

func (server *Server) AddUser(user string, password string) {
	server.Users[user] = password
}
