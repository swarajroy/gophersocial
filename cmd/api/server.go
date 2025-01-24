package api

import "net/http"

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	mux := http.NewServeMux()
	return &Server{mux: mux}
}

func (s *Server) RegisterRoutes() {

}
