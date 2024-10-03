package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/luqxus/dstore/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts
	store *Store

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	quitch chan struct{}
}

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

func (s *FileServer) stream(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	multiWriter := io.MultiWriter(peers...)
	return gob.NewEncoder(multiWriter).Encode(msg)
}

func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		peer.Send([]byte{p2p.IncomingMessage})
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

type MessageGetFile struct {
	Key string
}

func (s *FileServer) Get(key string) (int64, io.Reader, error) {
	ok := s.store.Has(key)
	if ok {
		fmt.Println("serving file from local disk")
		return s.store.Read(key)
	}

	fmt.Println("file not found locally, searching on network...")

	msg := Message{
		Payload: MessageGetFile{
			Key: key,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return 0, nil, err
	}

	time.Sleep(time.Millisecond * 500)

	for _, peer := range s.peers {
		var fileSize int64
		binary.Read(peer, binary.LittleEndian, &fileSize)

		n, err := s.store.Write(key, io.LimitReader(peer, fileSize))
		if err != nil {
			return 0, nil, err
		}
		fmt.Printf("received (%d) bytes from peer", n)
		// fmt.Println(fileBuffer.String())

		peer.CloseStream()

	}

	return s.store.Read(key)
}

func (s *FileServer) Store(key string, r io.Reader) error {

	fileBuffer := new(bytes.Buffer)
	tee := io.TeeReader(r, fileBuffer)

	size, err := s.store.Write(key, tee)
	if err != nil {
		return err
	}

	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}

	err = s.broadcast(&msg)
	if err != nil {
		return nil
	}

	time.Sleep(time.Millisecond * 3)

	for _, peer := range s.peers {
		peer.Send([]byte{p2p.IncomingStream})
		n, err := io.Copy(peer, fileBuffer)
		if err != nil {
			return err
		}

		fmt.Printf("wrote (%d) bytes to peer\n", n)

		// if err := peer.Send(payload); err != nil {
		// 	return err
		// }
	}

	return nil

	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r, buf)

	// err := s.store.Write(key, tee)
	// if err != nil {
	// 	return err
	// }

	// _, err = io.Copy(buf, r)
	// if err != nil {
	// 	return nil
	// }

	// p := &DataMessage{
	// 	key:  key,
	// 	Data: buf.Bytes(),
	// }

	// return s.broadcast()
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()

	defer s.peerLock.Unlock()

	log.Printf("connection with remote %s", peer.RemoteAddr())

	s.peers[peer.RemoteAddr().String()] = peer

	return nil
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			err := s.Transport.Dial(addr)
			if err != nil {
				log.Printf("dial error %s", err.Error())
			}
		}(addr)
	}

	return nil

}

func (s *FileServer) Start() error {
	err := s.Transport.ListenAndAccept()
	if err != nil {
		return err
	}

	if len(s.BootstrapNodes) != 0 {
		s.bootstrapNetwork()
	}

	s.loop()

	return nil
}

func (s *FileServer) loop() {

	defer func() {
		log.Println("file server stopped")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			// fmt.Println(msg)
			var msg Message
			// fmt.Println(msg.Payload)
			err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg)
			if err != nil {
				log.Println(err)
			}

			err = s.handleMessage(rpc.From, &msg)
			if err != nil {
				log.Println(err)
			}

		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, v)

	case MessageGetFile:
		return s.handleMessageGetFile(from, v)
	}
	return nil
}

func (s *FileServer) handleMessageGetFile(from string, msg MessageGetFile) error {
	if !s.store.Has(msg.Key) {
		fmt.Printf("file (%s) is does not exist on disk\n", msg.Key)
		return fmt.Errorf("file (%s) is does not exist on disk", msg.Key)
	}

	fmt.Println("serving file over the network")

	n, r, err := s.store.Read(msg.Key)
	if err != nil {
		return err
	}

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) not in peers list", from)
	}

	peer.Send([]byte{p2p.IncomingStream})

	binary.Write(peer, binary.LittleEndian, n)

	_, err = io.Copy(peer, r)
	if err != nil {
		return err
	}

	fmt.Printf("written (%d) bytes to peer\n", n)

	return nil
}

func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) not found in peers", from)
	}

	n, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk.\n", n)

	peer.CloseStream()

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
