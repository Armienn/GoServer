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

func viewHandler(w http.ResponseWriter, r *http.Request, info goserver.Info) {
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

func jsHandler(w http.ResponseWriter, r *http.Request, info goserver.Info) {
	file, _ := ioutil.ReadFile(path)
	w.Write(file)
}
