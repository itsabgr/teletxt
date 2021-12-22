package teletxt

import (
	"bufio"
	"bytes"
	"io"
)

type Conn struct {
	conn io.ReadWriter
	buff *bufio.Reader
}

func (c *Conn) Stream() io.ReadWriter {
	return c.conn
}
func (c *Conn) WriteKV(k, v []byte) (int, error) {
	b := bytes.NewBuffer(nil)
	b.Write(k)
	b.Write([]byte{':', ' '})
	b.Write(v)
	b.WriteByte('\n')
	return c.conn.Write(b.Bytes())
}
func (c *Conn) WriteValue(v []byte) (int, error) {
	b := bytes.NewBuffer(nil)
	b.Write(v)
	b.Write([]byte{'\r', '\n'})
	return c.conn.Write(b.Bytes())
}
func (c *Conn) WriteEmptyLine() (int, error) {
	return c.conn.Write([]byte{'\n', '\n'})
}
func ParseKeyPair(b []byte) (k []byte, v []byte) {
	if b == nil {
		return nil, nil
	}
	parts := bytes.SplitN(b, []byte{':'}, 2)
	if len(parts) == 2 {
		v = bytes.TrimSpace(parts[1])
	}
	k = bytes.TrimSpace(parts[0])
	if len(v) == 0 {
		v = nil
	}
	if len(k) == 0 {
		k = nil
	}
	return k, v
}
func (c *Conn) ReadKV() ([]byte, []byte, error) {
	line, err := c.ReadLine()
	k, v := ParseKeyPair(line)
	return k, v, err
}

func (c *Conn) ReadByte() (byte, error) {
	return c.buff.ReadByte()
}
func (c *Conn) Discard(n int) (int, error) {
	return c.buff.Discard(n)
}
func (c *Conn) UnreadByte() error {
	return c.buff.UnreadByte()
}

func (c *Conn) DiscardUntilEmptyLine() (err error) {
	for {
		line, err := c.ReadLine()
		if err != nil {
			return err
		}
		if line == nil || len(line) == 0 {
			return nil
		}
	}
}

func (c *Conn) ReadLine() (line []byte, err error) {
	for {
		chunk, isPrefix, err := c.buff.ReadLine()
		if chunk != nil {
			line = append(line, chunk...)
		}
		if err != nil {
			return line, err
		}
		if !isPrefix {
			break
		}
	}
	return line, err
}

func NewConn(conn io.ReadWriter) *Conn {
	return &Conn{
		conn: conn,
		buff: bufio.NewReader(conn),
	}
}
