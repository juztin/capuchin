package handlers

import "net/http"

var pong = []byte("pong")

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write(pong)
}
