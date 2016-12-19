package main

import (
	"net/http"
	"strconv"

	"github.com/Armienn/GoServer"
)

func main() {
	server := goserver.NewServer()
	server.AddHandler("/view/", viewHandler)
	server.Serve()
}

func viewHandler(w http.ResponseWriter, r *http.Request, path string, session goserver.Session) {
	count := 0
	value, ok := session.Get("musle")
	if ok {
		count = value.(int)
	}
	count += 1
	session.Set("musle", count)
	w.Write([]byte("Jo hollo" + strconv.Itoa(count)))
}
