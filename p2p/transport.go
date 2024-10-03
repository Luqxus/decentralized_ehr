package p2p

import "net"

// Peer interface represents the remote node.
type Peer interface {
	net.Conn
	Send([]byte) error
	SetPublicKey(b []byte)
	CloseStream()
}

// Transport handles communication between nodes in the network.
// Can be [TCP, UDP, Websockets]
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
