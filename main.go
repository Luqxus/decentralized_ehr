package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/luqxus/dstore/p2p"
)

func OnPeer(p p2p.Peer) error {
	p.Close()
	// fmt.Println("doing some logic with the peer outside of TCPTransport")
	return nil
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tr,
		BootstrapNodes:    nodes,
	}

	server := NewFileServer(fileServerOpts)

	tr.OnPeer = server.OnPeer

	return server
}

func main() {
	server1 := makeServer(":3000", "")
	server2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(server1.Start())
	}()

	time.Sleep(2 * time.Second)
	go func() {
		log.Fatal(server2.Start())
	}()

	time.Sleep(2 * time.Second)

	// data := bytes.NewReader([]byte("oisajdop"))

	// server2.StoreData("private key", data)

	r, err := server2.Get("privatekey")
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(b))
	select {}
}