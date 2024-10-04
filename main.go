package main

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/luqxus/dstore/contract"
	"github.com/luqxus/dstore/crypto"
	"github.com/luqxus/dstore/p2p"
)

func OnPeer(p p2p.Peer) error {
	p.Close()
	// fmt.Println("doing some logic with the peer outside of TCPTransport")
	return nil
}

func makeServer(listenAddr string, walletPath string, tmp string, nodes ...string) *FileServer {

	keystoreOpts := crypto.KeystoreOpts{
		ChainID: big.NewInt(534351),

		// FIXME: wallet path passed through function argument for testing purposes
		WalletPath: walletPath,

		// FIXME: @luqxus remove password here, this is for testing purposes
		Password: "1234567890",

		TmpPath: tmp,
	}

	ks, err := crypto.NewKeystore(keystoreOpts)
	if err != nil {
		log.Fatal(err)
	}

	contractAddr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddr == "" {
		log.Fatal("CONTRACT_ADDRESS not found")
	}

	provider := os.Getenv("ALCHEMY_PROVIDER")
	if provider == "" {
		log.Fatal("ALCHEMY_PROVIDER not found")
	}

	contractOpts := contract.ContractOpts{
		ContractAddress: contractAddr,
		Provider:        provider,
		Keystore:        ks,
	}

	contract, err := contract.NewEthContract(contractOpts)
	if err != nil {
		log.Fatal(err)
	}

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.DefaultHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
		Contract:      contract,
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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	server1 := makeServer(
		":3000",
		"./wallet/UTC--2024-10-04T15-37-19.251086713Z--a56aa73c2a4178a2ecc36125459df6ef4a346c98",
		"./tmp1",
		"",
	)
	server2 := makeServer(
		":4000",
		"./wallet2/UTC--2024-10-04T18-38-28.600393821Z--c4b7384139c3bf033da5e15566803ad894f4c3c7",
		"./tmp2",
		":3000",
	)

	go func() {
		// start server 1
		log.Fatal(server1.Start())
	}()

	time.Sleep(2 * time.Second)

	go func() {
		// start server 2
		log.Fatal(server2.Start())
	}()

	time.Sleep(10 * time.Second)

	_, r, err := server2.Get("privatekey")
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
