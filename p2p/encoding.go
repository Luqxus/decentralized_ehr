package p2p

import (
	"encoding/gob"
	"io"
)

// encoder interface
type Decoder interface {

	// encode reader data to rpc
	// returns error
	Decode(r io.Reader, rpc *RPC) error
}

// GOB decoder for testing purposes
type GOBdecoder struct{}

// GOB decoder Decode method for testing purposes
func (dec GOBdecoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}

// default decoder implemets Decoder interface
type DefaultDecoder struct{}

// default implementation of the Decode interface method
// returns  error
func (dec DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {
	// peek buffer for checking incoming transmission type
	peekBuf := make([]byte, 1)

	// read reader data into peekBuf
	_, err := r.Read(peekBuf)
	if err != nil {
		return err
	}

	// check if incoming transmission is a stream
	stream := peekBuf[0] == IncomingStream
	if stream {
		// if tranmission is stream
		// set rpc.Stream to true
		rpc.Stream = true

		// return nil
		return nil
	}

	buf := make([]byte, 1024)

	// read transmission reader data into buf
	n, err := r.Read(buf)
	if err != nil {
		// on error reading return
		return err
	}

	// set Payload to read buf
	rpc.Payload = buf[:n]

	// return nil on success
	return nil
}
