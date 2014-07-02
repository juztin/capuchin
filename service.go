package service

import (
	"log"
	"net"
	"net/http"
	"os"

	"bitbucket.org/stampinup/service/handlers"
	"code.minty.io/config"
	"code.minty.io/marbles/listeners"
	"github.com/gorilla/mux"
)

type Mux struct {
	*mux.Router
}

type Service struct {
	Router   *Mux
	listener net.Listener
}

func (m *Mux) HandleSigned(path string, handler http.Handler) {
	m.Handle(path, handlers.Signed(handler))
}

func (m *Mux) HandleSignedFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	m.Handle(path, handlers.Signed(http.HandlerFunc(f)))
}

func (s *Service) Serve() {
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

//func New() (*Service, error) {
func New() *Service {
	listener, err := Listener()
	if err != nil {
		panic(err)
		//return nil, err
	}

	return &Service{
		&Mux{mux.NewRouter()},
		listener,
	} //, err
}
