package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/luqxus/dstore/crypto"
)

type ContractOpts struct {
	ContractAddress string
	Provider        string
	Keystore        *crypto.Keystore
}

type Contract interface {
	VerifyNode(address common.Address, ip string) (bool, error)
	IsAdded(ip string) (bool, error)
	AddNode(ip string) error
	GetPublicKey() (*ecdsa.PublicKey, error)
}

type EthContract struct {
	client   *ethclient.Client
	keystore *crypto.Keystore
	verifier *Verifier
}

func NewEthContract(opts ContractOpts) (*EthContract, error) {

	client, err := ethclient.Dial(opts.Provider)
	if err != nil {
		return nil, err
	}

	v, err := NewVerifier(common.HexToAddress(opts.ContractAddress), client)
	if err != nil {
		return nil, err
	}

	return &EthContract{
		verifier: v,
		client:   client,
		keystore: opts.Keystore,
	}, nil
}

func (c *EthContract) GetPublicKey() (*ecdsa.PublicKey, error) {
	return c.keystore.GetPublicKey()
}

func (c *EthContract) AddNode(ip string) error {

	nonce, err := c.getNonce()
	if err != nil {
		return err
	}

	gasPrice, err := c.suggestedGasPrice()
	if err != nil {
		return err
	}

	addFunc := func(auth *bind.TransactOpts, ip string) (*types.Transaction, error) {
		return c.verifier.Add(auth, ip)
	}

	tx, err := c.keystore.SignAddNodeTx(nonce, gasPrice, ip, addFunc)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Add Node tx : %+v\n", tx)

	return err
}

func (c *EthContract) VerifyNode(address common.Address, ip string) (bool, error) {

	verifyFunc := func(opts *bind.CallOpts, address common.Address, ip string) (bool, error) {
		return c.verifier.Verify(opts, address, ip)
	}

	ok, err := c.keystore.VerifyNode(address, ip, verifyFunc)
	if err != nil {
		return false, err
	}

	// fmt.Printf("%+v\n", tx)

	return ok, nil
}

// FIXME: @luqxus this function is here for testing purposes
func (c *EthContract) IsAdded(ip string) (bool, error) {

	checkFunc := func(opts *bind.CallOpts, ip string) (bool, error) {
		return c.verifier.IsAdded(opts, ip)
	}

	tx, err := c.keystore.CheckAdded(ip, checkFunc)
	if err != nil {
		return false, err
	}

	fmt.Printf("Is Added : %+v\n", tx)

	return false, nil
}

func (c *EthContract) getNonce() (uint64, error) {
	return c.client.PendingNonceAt(context.Background(), c.keystore.Address())
}

func (c *EthContract) suggestedGasPrice() (*big.Int, error) {
	return c.client.SuggestGasPrice(context.Background())
}
