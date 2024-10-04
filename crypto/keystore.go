package crypto

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type KeystoreOpts struct {
	ChainID    *big.Int
	WalletPath string
	TmpPath    string
	// FIXME: @luqxus remove password from here
	Password string
}

type Keystore struct {
	ks *keystore.KeyStore
	KeystoreOpts
}

func NewKeystore(opts KeystoreOpts) (*Keystore, error) {
	// ks := keystore.NewKeyStore("./wallet2", keystore.StandardScryptN, keystore.StandardScryptP)
	ks := importKs(opts.WalletPath, opts.TmpPath, opts.Password)
	// _, err := ks.NewAccount(opts.Password)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return &Keystore{
		ks:           ks,
		KeystoreOpts: opts,
	}, nil
}

func (ks *Keystore) GetPublicKey() (*ecdsa.PublicKey, error) {

	jsonBytes, err := os.ReadFile(ks.WalletPath)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(jsonBytes, ks.Password)
	if err != nil {
		return nil, err
	}

	return &key.PrivateKey.PublicKey, nil

}

func (ks *Keystore) Address() common.Address {
	ks.ks.Unlock(ks.ks.Accounts()[0], ks.Password)
	addr := ks.ks.Accounts()[0].Address
	log.Printf("Addr Hex : %s\n", addr.Hex())
	return addr
}

func (ks *Keystore) SignAddNodeTx(nonce uint64,
	gasPrice *big.Int, value string,
	fn func(*bind.TransactOpts, string) (*types.Transaction, error),
) (*types.Transaction, error) {
	auth, err := ks.newTransactor(nonce, gasPrice)
	if err != nil {
		return nil, err
	}

	return fn(auth, value)
}

func (ks *Keystore) VerifyNode(
	address common.Address,
	ip string,
	fn func(*bind.CallOpts, common.Address, string) (bool, error),
) (bool, error) {

	opts := &bind.CallOpts{
		From: ks.Address(),
	}

	return fn(opts, address, ip)
}

func (ks *Keystore) CheckAdded(
	value string,
	fn func(*bind.CallOpts, string) (bool, error),
) (bool, error) {

	opts := &bind.CallOpts{
		From: ks.Address(),
	}

	return fn(opts, value)
}

func (ks *Keystore) newTransactor(nonce uint64, gasPrice *big.Int) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyStoreTransactorWithChainID(ks.ks, ks.ks.Accounts()[0], ks.ChainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	return auth, nil
}

func importKs(walletPath string, tmp string, password string) *keystore.KeyStore {
	ks := keystore.NewKeyStore(tmp, keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := os.ReadFile(walletPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	return ks
}
