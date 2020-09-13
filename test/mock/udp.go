package mock

import (
	"net"
	"time"
)

// UDP mock implementation
type UDP struct {
	ReadFromFun func(p []byte) (n int, addr net.Addr, err error)
}

// ReadFrom mock implementation
func (U UDP) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	if U.ReadFromFun != nil {
		return U.ReadFromFun(p)
	}
	return 0, nil, nil
}

// WriteTo mock implementation
func (U UDP) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	panic("not implemented")
}

// Close mock implementation
func (U UDP) Close() error {
	panic("not implemented")
}

// LocalAddr mock implementation
func (U UDP) LocalAddr() net.Addr {
	panic("not implemented")
}

// SetDeadline mock implementation
func (U UDP) SetDeadline(t time.Time) error {
	panic("not implemented")
}

// SetReadDeadline mock implementation
func (U UDP) SetReadDeadline(t time.Time) error {
	panic("not implemented")
}

// SetWriteDeadline mock implementation
func (U UDP) SetWriteDeadline(t time.Time) error {
	panic("not implemented")
}



