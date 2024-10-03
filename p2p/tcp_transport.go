package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/luqxus/dstore/crypto"
)

// TCPPeer represents the remote node over TCP established connection
type TCPPeer struct {

	// underlying connection of the peer
	// which in this case is a TCP connection
	net.Conn

	// if we dial and retrieve a conn => outbound == true
	// if we accept and retrieve a conn => outbound => false
	outbound bool

	wg *sync.WaitGroup

	PublicKey crypto.PublicKey
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

func (p *TCPPeer) SetPublicKey(b []byte) {
	p.PublicKey = b
}

// Addr implements transport interface
func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC, 1024),
	}
}

// consume implements the transport interface, which will return read only channel
// for reading the incoming messages received from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Printf("TCP transport listening on addr : %s\n", t.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			log.Printf("tcp connection closed")
			return
		}

		if err != nil {
			fmt.Printf("TCP accept error : %s\n", err.Error())
		}

		fmt.Printf("new incoming connection : %+v\n", conn)

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err.Error())
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	fmt.Printf("peer public key : %s\n", peer.PublicKey.String())

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	for {

		// Read loop
		rpc := RPC{}
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {
			// log.Printf("tcp error: %s", "connection closed")
			return
		}

		rpc.From = conn.RemoteAddr().String()
		if rpc.Stream {
			// fmt.Println("Waiting")
			peer.wg.Add(1)
			fmt.Println("incoming stream. waiting...")
			peer.wg.Wait()
			// fmt.Println("Done Waiting")
			fmt.Println("stream closed. resuming read")
		}
		t.rpcch <- rpc
	}

}

// Close implements the transport interface.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implements the Transport interface.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}
