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

type Server struct {
	Router   *mux.Router
	listener net.Listener
}

func (s *Server) Serve() {
	log.Printf("listening on %s", s.listener.Addr())
	http.Serve(s.listener, s.Router)
}

func Listener() (net.Listener, error) {
	if isSock, _ := config.GroupBool("server", "unixSock"); isSock {
		f := config.Required.GroupString("server", "unixSockFile")
		return listeners.NewSOCK(f, os.ModePerm)
	}

	p := config.Required.GroupInt("server", "port")
	h, ok := config.GroupString("server", "host")
	if !ok {
		h = ""
	}
	return listeners.NewHTTP(h, p)
}

func New(endpoint string) *Server {
	listener, err := Listener()
	if err != nil {
		panic(err)
	}

	// Create the subrouter, so all additional paths are children of `endpoint`
	router := mux.NewRouter().PathPrefix(endpoint).Subrouter()

	// Add status & ping routes
	ping := router.PathPrefix("/ping")
	status := router.PathPrefix("/status")
	ping.HandlerFunc(handlers.Ping)
	status.Handler(handlers.SignedFunc(handlers.Status))

	return &Server{router, listener}
}
