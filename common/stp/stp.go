package stp

import (
	"fmt"
	"net"
)

type Delegate interface {
}

type CTL int32

const (
	CTL_DAT CTL = iota
	CTL_ACK
)

type STP struct {
	conn *net.UDPConn
	wbuf []byte
}

func NewSTP(conv uint32, conn *net.UDPConn) *STP {
	stp := &STP{
		conn: conn,
		wbuf: make([]byte, 65535),
	}

	return stp
}

func (s *STP) Write(b []byte) (int, error) {
	return 0, fmt.Errorf("closed")
}

func (s *STP) Close() {

	s.conn.Close()
}
