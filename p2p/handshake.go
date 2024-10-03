package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/luqxus/dstore/crypto"
)

type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error {
	return nil
}

func DefaultHandshakeFunc(peer Peer) error {
	// FIXME: @luqxus do not generate private key here
	privKey := crypto.GeneratePrivateKey()
	pubKey := privKey.PublicKey()

	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(pubKey)
	if err != nil {
		return err
	}

	if err := peer.Send(buf.Bytes()); err != nil {
		return err
	}

	timeout := 40 * time.Second
	done := make(chan bool, 1)

	go func() {
		pubKey, err := receivePeerPublicKey(peer)
		if err != nil {
			done <- false
			log.Printf("handshake error : %s\n", err.Error())
			return
		}

		if !verify() {
			done <- false
			log.Printf("handshake error :  invalid public key")
			return
		}

		peer.SetPublicKey(pubKey)
		done <- true

	}()

	select {
	case success := <-done:
		if !success {
			return fmt.Errorf("handshake error : failed to verify peer's public key")
		}

	case <-time.After(timeout):
		return fmt.Errorf("handshake error: timeout")
	}

	log.Println("handshake successful")
	return nil
}

func verify() bool {
	time.Sleep(time.Millisecond * 600)
	return false
}

func receivePeerPublicKey(peer Peer) ([]byte, error) {
	b := make([]byte, 33)
	_, err := peer.Read(b)
	if err != nil {
		return nil, err
	}

	return b, err
}
