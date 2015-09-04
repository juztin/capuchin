// Copyright 2015 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package capuchin

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"code.minty.io/capuchin/handlers"
	"code.minty.io/marbles/listeners"
	"github.com/gorilla/mux"
)

const (
	DEFAULT_PORT = 9000
	DEFAULT_HOST = "127.0.0.1"
)

// A wrapper around gorilla.mux and net.listener.
// The listener is created from settings within `config.json`.
// The endpoints `ping` and `status` are automatically added to the mux.
type Server struct {
	endpoint string
	Router   *mux.Router
	listener net.Listener
}

// Invokes http.Serve
func (s *Server) Serve() {
	log.Printf("listening at %s on %s", s.endpoint, s.listener.Addr())
	http.Serve(s.listener, s.Router)
}

// Adds defaults routes, ping and status.
func addRoutes(r *mux.Router) {
	ping := handlers.RecoveryFunc(handlers.Ping)
	time := handlers.RecoveryFunc(handlers.Time)
	status := handlers.RecoveryFunc(handlers.Status)
	r.Handle("/status/", status).Methods("GET")
	r.Handle("/ping/", ping).Methods("GET")
	r.Handle("/time/", time).Methods("GET")
}

// Returns a new gorilla Router for the given endpoint, with ping and status routes added.
func newRouter(endpoint string) *mux.Router {
	r := Router(endpoint)
	r.StrictSlash(true)
	addRoutes(r)
	return r
}

// Return either a UNIX socket, or a TCP net.Listener based on `config.json`.
func Listener() (net.Listener, error) {
	h := os.Getenv("SERVICE_HOST")
	p := os.Getenv("SERVICE_PORT")
	f := os.Getenv("SERVICE_SOCK_FILE")
	if f != "" {
		return listeners.NewSOCK(f, os.ModePerm)
	}
	if h == "" {
		h = DEFAULT_HOST
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		port = DEFAULT_PORT
	}
	return listeners.NewHTTP(h, port)
}

func Router(endpoint string) *mux.Router {
	router := mux.NewRouter()
	if endpoint != "" {
		// Create the subrouter, so all additional paths are children of `endpoint`
		router = router.PathPrefix(endpoint).Subrouter()
	}
	return router
}

// Returns a new Server, handling routes at the given endpoint.
func New(endpoint string) *Server {
	listener, err := Listener()
	if err != nil {
		panic(err)
	}
	router := newRouter(endpoint)
	return &Server{endpoint, router, listener}
}

// Same as New, using the given net.Listener.
func NewWithListener(endpoint string, listener net.Listener) *Server {
	router := newRouter(endpoint)
	return &Server{endpoint, router, listener}
}
