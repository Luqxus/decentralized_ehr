package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(r io.Reader, rpc *RPC) error
}

type GOBdecoder struct{}

func (dec GOBdecoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {
	buf := make([]byte, 1024)

	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	rpc.Payload = buf[:n]

	return nil
}