package raft_udp_transport

import (
	"errors"
	"io"
	"net"
	"time"

	"github.com/hashicorp/raft"
)

var (
	errNotAdvertisable = errors.New("local bind address is not advertisable")
	errNotUDP          = errors.New("local address is not a UDP address")
)

type UDPStreamLayer struct {
	advertise net.Addr
	listener  net.Listener
}

// NewUDPTransport returns a NetworkTransport that is built on top of
// a UDP streaming transport layer.
func NewUDPTransport(
	bindAddr string,
	advertise net.Addr,
	maxPool int,
	timeout time.Duration,
	logOutput io.Writer,
) (*raft.NetworkTransport, error) {
	return newUDPTransport(bindAddr, advertise, func(stream raft.StreamLayer) *raft.NetworkTransport {
		return raft.NewNetworkTransport(stream, maxPool, timeout, logOutput)
	})
}

func newUDPTransport(bindAddr string,
	advertise net.Addr,
	transportCreator func(stream raft.StreamLayer) *raft.NetworkTransport) (*raft.NetworkTransport, error) {
	// Try to bind
	list, err := net.Listen("udp", bindAddr)
	if err != nil {
		return nil, err
	}

	// Create stream
	stream := &UDPStreamLayer{
		advertise: advertise,
		listener:  list.(net.Listener),
	}

	// Verify that we have a usable advertise address
	addr, ok := stream.Addr().(*net.TCPAddr)
	if !ok {
		list.Close()
		return nil, errNotUDP
	}
	if addr.IP.IsUnspecified() {
		list.Close()
		return nil, errNotAdvertisable
	}

	// Create the network transport
	trans := transportCreator(stream)
	return trans, nil
}

// Dial implements the StreamLayer interface.
func (t *UDPStreamLayer) Dial(address raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("udp", string(address), timeout)
}

// Accept implements the net.Listener interface.
func (t *UDPStreamLayer) Accept() (c net.Conn, err error) {
	return t.listener.Accept()
}

// Close implements the net.Listener interface.
func (t *UDPStreamLayer) Close() (err error) {
	return t.listener.Close()
}

// Addr implements the net.Listener interface.
func (t *UDPStreamLayer) Addr() net.Addr {
	// Use an advertise addr if provided
	if t.advertise != nil {
		return t.advertise
	}
	return t.listener.Addr()
}
