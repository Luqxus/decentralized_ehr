package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"sync"

	"github.com/luqxus/dstore/contract"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
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

	PublicKey ecdsa.PublicKey
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
	Contract      contract.Contract
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

func (p *TCPPeer) SetPublicKey(key ecdsa.PublicKey) {
	p.PublicKey = key
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

	// FIXME: @luqxus adding node to smart contract here for testing purposes
	// exists, err := t.Contract.IsAdded(t.ListenAddr)
	// if err != nil {
	// 	return err
	// }

	// if !exists {
	// 	if err := t.Contract.AddNode(t.ListenAddr); err != nil {
	// 		return err
	// 	}
	// }

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

	fmt.Printf("New Connected Peer : %+v\n", peer)
	if err = t.HandshakeFunc(peer, t.DefaultHandshakeFunc); err != nil {
		return
	}

	fmt.Printf("peer public key : %+v\n", peer.PublicKey)

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

func (t *TCPTransport) DefaultHandshakeFunc(peer Peer) error {
	// FIXME: @luqxus do not generate private key here

	myKey, err := t.Contract.GetPublicKey()
	if err != nil {
		return err
	}

	fmt.Printf("Our Public Key : %+v\n", toSerializablePubKey(myKey))

	buf := new(bytes.Buffer)
	err = gob.NewEncoder(buf).Encode(&myKey)
	if err != nil {
		return err
	}

	if err := peer.Send(buf.Bytes()); err != nil {
		return err
	}

	// timeout := 40 * time.Second
	// done := make(chan bool, 1)
	// var keyCh chan ecdsa.PublicKey
	// go func() {

	b, err := receivePeerPublicKey(peer)
	if err != nil {
		// done <- false
		log.Printf("handshake error : %s\n", err.Error())
		return err
	}

	var serializablePubKey SerializablePubKey

	err = gob.NewDecoder(bytes.NewBuffer(b)).Decode(&serializablePubKey)
	if err != nil {
		// done <- false
		log.Printf("handshake error. error decoding: %s", err.Error())
		return err
	}

	fmt.Printf("skdjklasjdlkjaskldklas %+v\n", serializablePubKey)

	pubKey, err := fromSerializablePubKey(serializablePubKey)
	if err != nil {
		log.Printf("handshake error. error decoding: %s", err.Error())
		return err
	}

	fmt.Printf("Received Public Key : %+v\n", pubKey)

	// keyCh <- pubKey
	fmt.Printf("Address Verified : %s\n", peer.LocalAddr().Network())
	ok, err := t.Contract.VerifyNode(crypto.PubkeyToAddress(*pubKey), peer.LocalAddr().Network())
	if err != nil {
		// done <- false
		log.Printf("handshake error :  invalid public key")
		return err
	}

	fmt.Printf("Decoded Address : %s\n", crypto.PubkeyToAddress(*pubKey).Hex())

	if !ok {
		return fmt.Errorf("handshake error : node not validated")
	}

	// done <- ok

	// }()

	// select {
	// case success := <-done:
	// 	if !success {
	// 		<-keyCh
	// 		return fmt.Errorf("handshake error : failed to verify peer's public key")
	// 	}
	peer.SetPublicKey(*pubKey)
	// 	return nil

	// case <-time.After(timeout):
	// 	return fmt.Errorf("handshake error: timeout")
	// }

	return nil

}

func receivePeerPublicKey(peer Peer) ([]byte, error) {
	b := make([]byte, 512)
	n, err := peer.Read(b)
	if err != nil {
		return nil, err
	}

	fmt.Printf("number kdsjfkdsfkldslfksdf : %d\n", n)

	return b[:n], err
}

func init() {
	gob.Register(ecdsa.PublicKey{})
	gob.Register(secp256k1.BitCurve{})
}

type SerializablePubKey struct {
	CurveName string // Name of the curve
	X, Y      *big.Int
}

// toSerializablePubKey converts ecdsa.PublicKey to SerializablePubKey
func toSerializablePubKey(pubKey *ecdsa.PublicKey) SerializablePubKey {
	return SerializablePubKey{
		CurveName: "secp256k1", // Use the name for the specific curve
		X:         pubKey.X,
		Y:         pubKey.Y,
	}
}

// fromSerializablePubKey converts SerializablePubKey back to ecdsa.PublicKey
func fromSerializablePubKey(s SerializablePubKey) (*ecdsa.PublicKey, error) {
	// var curve elliptic.Curve
	// if s.CurveName == "secp256k1" {
	// 	curve = secp256k1.S256() // Assign the secp256k1 curve
	// } else {
	// 	return nil, log.Output(0, "Unsupported curve")
	// }

	return &ecdsa.PublicKey{
		Curve: secp256k1.S256(),
		X:     s.X,
		Y:     s.Y,
	}, nil
}
