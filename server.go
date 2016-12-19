package main

import "net/http"

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8080", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Jo hollo"))
}
