// Copyright 2014 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package capuchin

import (
	"log"
	"net"
	"net/http"
	"os"

	"code.minty.io/capuchin/handlers"
	"code.minty.io/config"
	"code.minty.io/marbles/listeners"
	"github.com/gorilla/mux"
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

// Return either a UNIX socket, or a TCP net.Listener based on `config.json`.
func Listener() (net.Listener, error) {
	if isSock, _ := config.GroupBool("server", "unixSock"); isSock {
		f := config.RequiredGroupString("server", "unixSockFile")
		return listeners.NewSOCK(f, os.ModePerm)
	}

	p := config.RequiredGroupInt("server", "port")
	h, ok := config.GroupString("server", "host")
	if !ok {
		h = ""
	}
	return listeners.NewHTTP(h, p)
}

func Router(endpoint string) *mux.Router {
	router := mux.NewRouter()
	if endpoint != "" {
		// Create the subrouter, so all additional paths are children of `endpoint`
		router = router.PathPrefix(endpoint).Subrouter()
	}
	return router
}

func NewWithListener(endpoint string, listener net.Listener) *Server {
	router := Router(endpoint)
	ping := router.PathPrefix("/ping/")
	status := router.PathPrefix("/status/")
	ping.Handler(handlers.RecoveryFunc(handlers.Ping))
	status.Handler(handlers.Recovery(handlers.SignedFunc(handlers.Status)))
	return &Server{endpoint, router, listener}
}

// Returns a new Server.
func New(endpoint string) *Server {
	listener, err := Listener()
	if err != nil {
		panic(err)
	}

	// Create the router for the endpoint.
	router := Router(endpoint)
	// Add status & ping routes
	ping := router.PathPrefix("/ping/")
	status := router.PathPrefix("/status/")
	ping.Handler(handlers.RecoveryFunc(handlers.Ping))
	status.Handler(handlers.Recovery(handlers.SignedFunc(handlers.Status)))
	return &Server{endpoint, router, listener}
}
