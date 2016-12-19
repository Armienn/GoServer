package main

import (
	"net/http"
	"strconv"

	"github.com/Armienn/GoServer"
)

func main() {
	server := goserver.NewServer(true)
	server.AddHandler("/view/", viewHandler)
	server.Serve()
}

func viewHandler(w http.ResponseWriter, r *http.Request, path string, session goserver.Session, user interface{}) {
	count := 0
	value, ok := session.Get("musle")
	if ok {
		count = value.(int)
	}
	count++
	session.Set("musle", count)
	w.Write([]byte("Jo hollo" + strconv.Itoa(count)))
}
