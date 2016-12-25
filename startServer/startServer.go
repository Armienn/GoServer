package main

import (
	"html/template"
	"net/http"

	"github.com/Armienn/GoServer"
)

func main() {
	server := goserver.NewServer(false)
	server.AddHandler("/", viewHandler)
	server.Serve()
}

func viewHandler(w http.ResponseWriter, r *http.Request, info goserver.Info) {
	data := struct{ Count int }{0}
	value, ok := info.Session.Get("musle")
	if ok {
		data.Count = value.(int)
	}
	data.Count++
	info.Session.Set("musle", data.Count)
	temp, err := template.ParseFiles("test.html")
	if err != nil {
		w.Write([]byte("Fejl: " + err.Error()))
	} else {
		temp.Execute(w, data)
	}
}
