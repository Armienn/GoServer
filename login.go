package goserver

import "net/http"

func loginGetHandler(w http.ResponseWriter, r *http.Request, info Info) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<form method=\"post\" action=\"/login\"><input type=\"text\" name=\"user\"><input type=\"password\" name=\"password\"><input type=\"submit\" name=\"submit\"></form>"))
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
	if pass, ok := server.Users[user]; ok && pass == password {
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
