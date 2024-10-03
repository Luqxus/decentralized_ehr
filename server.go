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

// file server options
type FileServerOpts struct {

	// root local network storage folder
	StorageRoot string

	// path transform function
	PathTransformFunc PathTransformFunc

	// transport
	Transport p2p.Transport

	// nodes | remote networks to connect to on starting server
	BootstrapNodes []string
}

// file server
type FileServer struct {
	// server options
	FileServerOpts

	// storage
	store *Store

	// peers map lock
	peerLock sync.Mutex

	// connected remote peers map
	peers map[string]p2p.Peer

	// quit channel
	quitch chan struct{}
}

// Message carries payload and is sent over the wire
type Message struct {

	// payload can of any type
	Payload any
}

// MessageStoreFile tells the receiver that a file is
// being transmitted for storage
type MessageStoreFile struct {

	// file path
	Key string

	// file size
	Size int64
}

// MessageGetFile tells the receiver to check and send file with Key
type MessageGetFile struct {

	// file path
	Key string
}

func (s *FileServer) stream(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	multiWriter := io.MultiWriter(peers...)
	return gob.NewEncoder(multiWriter).Encode(msg)
}

// broadcast, broadcasts file to all connected peers
// over the wire | returns [error]
func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)

	// encode msg for transmission
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	// loop over all connected peers
	for _, peer := range s.peers {
		// send file to each peer
		peer.Send([]byte{p2p.IncomingMessage})
		if err := peer.Send(buf.Bytes()); err != nil {
			// on send error return error
			return err
		}
	}

	// return nil on success
	return nil
}

// Get reads check and reads file from local network
// if file not found check file over connected peers remote network
func (s *FileServer) Get(key string) (int64, io.Reader, error) {
	// check if file exists in local network
	ok := s.store.Has(key)
	if ok {

		// if file found, read file
		fmt.Println("serving file from local disk")
		return s.store.Read(key)
	}

	fmt.Println("file not found locally, searching on network...")

	// if file is not found. prepare message of type MessageGetFile
	msg := Message{
		Payload: MessageGetFile{
			Key: key,
		},
	}

	// broadcast message over wire to request file from connected peerss
	if err := s.broadcast(&msg); err != nil {
		return 0, nil, err
	}

	time.Sleep(time.Millisecond * 500)

	// loop through all connected peers
	for _, peer := range s.peers {
		var fileSize int64

		// read file size from peer
		binary.Read(peer, binary.LittleEndian, &fileSize)

		// read file from peer and write to local network
		n, err := s.store.Write(key, io.LimitReader(peer, fileSize))
		if err != nil {
			return 0, nil, err
		}
		fmt.Printf("received (%d) bytes from peer", n)
		// fmt.Println(fileBuffer.String())

		// close read stream
		peer.CloseStream()

	}

	// return file size (int64) | file reader (io.Reader) | error (error)
	return s.store.Read(key)
}

// store file to local network and broadcast file over wire
// to all connected peers
// returns error
func (s *FileServer) Store(key string, r io.Reader) error {

	fileBuffer := new(bytes.Buffer)
	tee := io.TeeReader(r, fileBuffer)

	// write file to local network
	size, err := s.store.Write(key, tee)
	if err != nil {
		return err
	}

	// prepare message of type MessageStoreFile
	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,  // file path
			Size: size, // file size
		},
	}

	// broadcast MessageStoreFile
	// tells remote peers to store incoming file
	err = s.broadcast(&msg)
	if err != nil {
		return nil
	}

	time.Sleep(time.Millisecond * 3)

	// broadcast file to all remote peers
	for _, peer := range s.peers {
		// send file stream type
		peer.Send([]byte{p2p.IncomingStream})

		// send stream
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
}

// implements OnPeer transport interface
// applies logic on connected peer
// return error
func (s *FileServer) OnPeer(peer p2p.Peer) error {
	// lock peers map
	s.peerLock.Lock()

	// on function ends | unlock peers map
	defer s.peerLock.Unlock()

	log.Printf("connection with remote %s", peer.RemoteAddr())

	// add connected peer to peers map
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

// start server and listen
// returns error
func (s *FileServer) Start() error {

	// start transport and listen
	err := s.Transport.ListenAndAccept()
	if err != nil {
		// on error starting transport
		return err
	}

	// check if node to bootstrap is not zero
	if len(s.BootstrapNodes) != 0 {

		// if not zero |  bootstrap nodes
		s.bootstrapNetwork()
	}

	// start read loop
	s.loop()

	// return nil
	return nil
}

// read loop
// reads and handles messages from peers
func (s *FileServer) loop() {

	defer func() {
		log.Println("file server stopped")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			// new message
			var msg Message

			// decode message to messageb
			err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg)
			if err != nil {
				log.Println(err)
			}

			// handle message
			err = s.handleMessage(rpc.From, &msg)
			if err != nil {
				// on error handling message
				log.Println(err)
			}

		case <-s.quitch:
			// on quitch triggered
			return
		}
	}
}

// handle message from peer
func (s *FileServer) handleMessage(from string, msg *Message) error {

	// check message type
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		// on message type is MessageStoreFile
		return s.handleMessageStoreFile(from, v)

	case MessageGetFile:
		// on message typoe is MessageGetFile
		return s.handleMessageGetFile(from, v)
	}
	return nil
}

// handle MessageGetFile message from peer
// checks for file in local network
// if file found then writes file to peer
// return error
func (s *FileServer) handleMessageGetFile(from string, msg MessageGetFile) error {
	// check if file in local network storage
	if !s.store.Has(msg.Key) {
		// if file not found
		fmt.Printf("file (%s) is does not exist on disk\n", msg.Key)
		return fmt.Errorf("file (%s) is does not exist on disk", msg.Key)
	}

	fmt.Println("serving file over the network")

	// read file from local network storage
	n, r, err := s.store.Read(msg.Key)
	if err != nil {
		return err
	}

	// check if peer if in peers map
	peer, ok := s.peers[from]
	if !ok {
		// if peer on in peers return error
		return fmt.Errorf("peer (%s) not in peers list", from)
	}

	// send transmission type to peer
	peer.Send([]byte{p2p.IncomingStream})

	// write file size to peer
	binary.Write(peer, binary.LittleEndian, n)

	// write file to peer
	_, err = io.Copy(peer, r)
	if err != nil {
		// on error writing file
		return err
	}

	fmt.Printf("written (%d) bytes to peer\n", n)

	return nil
}

// handle MessageStoreFile message from peer
// return error
func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	// check if the sender peer is in peers map
	peer, ok := s.peers[from]
	if !ok {
		// if not in peers return error
		return fmt.Errorf("peer (%s) not found in peers", from)
	}

	// write file to local network storage
	n, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		// on error writing file
		return err
	}

	log.Printf("written (%d) bytes to disk.\n", n)

	// close stream on done
	peer.CloseStream()

	// return nil
	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
