package handlers

import "net/http"

var pong = []byte("pong")

// HTTP handler for /ping.
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write(pong)
}
