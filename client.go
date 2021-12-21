package teletxt

import "net"

func Dial(laddr, raddr *net.TCPAddr) (*Conn, error) {
	conn, err := net.DialTCP("tcp", laddr, raddr)
	if err != nil {
		return nil, err
	}
	return NewConn(conn), nil
}
