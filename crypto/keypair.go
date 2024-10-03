package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
)

// private key type
type PrivateKey struct {
	key *ecdsa.PrivateKey // key
}

// sign, signs data and []byte data
// return Signature pointer | error
func (k PrivateKey) Sign(data []byte) (*Signature, error) {

	// sign data
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		// on signing error
		return nil, err
	}

	// return signature | error
	return &Signature{
		S: s,
		R: r,
	}, nil
}

// public key derived from private key
type PublicKey []byte

// signature type

type Signature struct {
	S *big.Int
	R *big.Int
}

// encode signature to sign
func (sig *Signature) String() string {
	b := append(sig.S.Bytes(), sig.R.Bytes()...)

	return hex.EncodeToString(b)
}

// venrify signature
func (sig *Signature) Verify(pubKey PublicKey, data []byte) bool {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), pubKey)

	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return ecdsa.Verify(key, data, sig.R, sig.S)
}

// private generate private key from rand.Reader
func newPrivateKeyFromReader(r io.Reader) PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), r)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

// Public function Generate private key
func GeneratePrivateKey() PrivateKey {
	// call and return new private from reader
	// which generates a new privateKey using the elliptic curve
	return newPrivateKeyFromReader(rand.Reader)
}

// derive public key from private key
// returns public key
func (k PrivateKey) PublicKey() PublicKey {
	return elliptic.MarshalCompressed(
		k.key.PublicKey,
		k.key.PublicKey.X,
		k.key.PublicKey.Y,
	)
}

// encode public key to string
// returns string
func (k PublicKey) String() string {
	return hex.EncodeToString(k)
}

// derive address from public key
// returns address
func (k PublicKey) Address() Address {
	h := sha256.Sum256(k)

	// gets address from public key []byte
	return AddressFromBytes(h[len(h)-20:])
}

// address type with byte size of 20
// type uint8
type Address [20]uint8

// get address from bytes
// returns address [20]uint20
func AddressFromBytes(b []byte) Address {
	// check if the provided []byte size is valid
	// (equals 20)
	if len(b) != 20 {
		// if invalid byte size panic
		msg := fmt.Sprintf("invalid bytes length (%d) should be (20)", len(b))
		panic(msg)
	}

	// new uint8 slice of size 20
	var value [20]uint8
	for i := 0; i < 20; i++ {
		// parse byte to uint8
		value[i] = b[i]
	}

	// return address
	return Address(value)
}

// encode address []uint8 to string
// returns string
func (a Address) String() string {
	return hex.EncodeToString(a.Slice())
}

func (a Address) Slice() []byte {
	b := make([]byte, 20)

	for i := 0; i < 20; i++ {
		b[i] = a[i]
	}

	return b
}
