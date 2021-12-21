package teletxt

import "net"

type Server struct {
	listener net.Listener
}

func NewServer(l net.Listener) *Server {
	return &Server{l}
}
func (s *Server) Accept() (*Conn, error) {
	netConn, err := s.listener.Accept()
	return NewConn(netConn), err
}
func (s *Server) Close() error {
	return s.Close()
}

func (s *Server) Addr() net.Addr {
	return s.Addr()
}
