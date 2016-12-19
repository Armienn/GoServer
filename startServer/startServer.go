package main

import (
	"html/template"
	"net/http"

	"io/ioutil"

	"github.com/Armienn/GoServer"
)

func main() {
	server := goserver.NewServer(false)
	server.AddHandler("/js/", jsHandler)
	server.AddHandler("/", viewHandler)
	server.Serve()
}

func viewHandler(server *goserver.Server, w http.ResponseWriter, r *http.Request, path string, session goserver.Session, user interface{}) {
	data := struct{ Count int }{0}
	value, ok := session.Get("musle")
	if ok {
		data.Count = value.(int)
	}
	data.Count++
	session.Set("musle", data.Count)
	temp, err := template.ParseFiles("test.html")
	if err != nil {
		w.Write([]byte("Fejl: " + err.Error()))
	} else {
		temp.Execute(w, data)
	}
	//w.Write([]byte("Jo hollo" + strconv.Itoa(count)))
}

func jsHandler(server *goserver.Server, w http.ResponseWriter, r *http.Request, path string, session goserver.Session, user interface{}) {
	file, _ := ioutil.ReadFile(path)
	w.Write(file)
}
