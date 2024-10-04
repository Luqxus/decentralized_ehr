package p2p

type HandshakeFunc func(Peer, func(Peer) error) error

func NOPHandshakeFunc(peer Peer, fn func(Peer) error) error {
	return nil
}

func DefaultHandshakeFunc(peer Peer, verifyFunc func(Peer) error) error {
	return verifyFunc(peer)

}

// func verify() bool {
// 	time.Sleep(time.Millisecond * 600)
// 	return false
// }
