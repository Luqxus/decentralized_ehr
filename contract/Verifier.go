// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// VerifierMetaData contains all meta data concerning the Verifier contract.
var VerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_ip\",\"type\":\"string\"}],\"name\":\"add\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"}],\"name\":\"isAdded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_ip\",\"type\":\"string\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506108c78061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c80632e3d616e14610043578063b0c8f9dc14610073578063e8fc92731461008f575b5f5ffd5b61005d6004803603810190610058919061032f565b6100bf565b60405161006a9190610394565b60405180910390f35b61008d6004803603810190610088919061032f565b610135565b005b6100a960048036038101906100a49190610407565b610240565b6040516100b69190610394565b60405180910390f35b5f82826040516100d09291906104a0565b60405180910390205f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2060010160405161012491906105a7565b604051809103902014905092915050565b60405180604001604052803373ffffffffffffffffffffffffffffffffffffffff16815260200183838080601f0160208091040260200160405190810160405280939291908181526020018383808284375f81840152601f19601f820116905080830192505050505050508152505f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f820151815f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506020820151816001019081610238919061079d565b509050505050565b5f82826040516102519291906104a0565b60405180910390205f5f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206001016040516102a591906105a7565b6040518091039020036102bb57600190506102bf565b5f90505b9392505050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f8401126102ef576102ee6102ce565b5b8235905067ffffffffffffffff81111561030c5761030b6102d2565b5b602083019150836001820283011115610328576103276102d6565b5b9250929050565b5f5f60208385031215610345576103446102c6565b5b5f83013567ffffffffffffffff811115610362576103616102ca565b5b61036e858286016102da565b92509250509250929050565b5f8115159050919050565b61038e8161037a565b82525050565b5f6020820190506103a75f830184610385565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6103d6826103ad565b9050919050565b6103e6816103cc565b81146103f0575f5ffd5b50565b5f81359050610401816103dd565b92915050565b5f5f5f6040848603121561041e5761041d6102c6565b5b5f61042b868287016103f3565b935050602084013567ffffffffffffffff81111561044c5761044b6102ca565b5b610458868287016102da565b92509250509250925092565b5f81905092915050565b828183375f83830152505050565b5f6104878385610464565b935061049483858461046e565b82840190509392505050565b5f6104ac82848661047c565b91508190509392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806104fc57607f821691505b60208210810361050f5761050e6104b8565b5b50919050565b5f819050815f5260205f209050919050565b5f8154610533816104e5565b61053d8186610464565b9450600182165f8114610557576001811461056c5761059e565b60ff198316865281151582028601935061059e565b61057585610515565b5f5b8381101561059657815481890152600182019150602081019050610577565b838801955050505b50505092915050565b5f6105b28284610527565b915081905092915050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026106507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610615565b61065a8683610615565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61069e61069961069484610672565b61067b565b610672565b9050919050565b5f819050919050565b6106b783610684565b6106cb6106c3826106a5565b848454610621565b825550505050565b5f5f905090565b6106e26106d3565b6106ed8184846106ae565b505050565b5b81811015610710576107055f826106da565b6001810190506106f3565b5050565b601f82111561075557610726816105f4565b61072f84610606565b8101602085101561073e578190505b61075261074a85610606565b8301826106f2565b50505b505050565b5f82821c905092915050565b5f6107755f198460080261075a565b1980831691505092915050565b5f61078d8383610766565b9150826002028217905092915050565b6107a6826105bd565b67ffffffffffffffff8111156107bf576107be6105c7565b5b6107c982546104e5565b6107d4828285610714565b5f60209050601f831160018114610805575f84156107f3578287015190505b6107fd8582610782565b865550610864565b601f198416610813866105f4565b5f5b8281101561083a57848901518255600182019150602085019450602081019050610815565b868310156108575784890151610853601f891682610766565b8355505b6001600288020188555050505b50505050505056fea264697066735822122056a6e6079f5f81a56a44fe8f9f53eb084108b67a1feb1b8997e24c597e0260c764736f6c637827302e382e32372d646576656c6f702e323032342e392e352b636f6d6d69742e34306133356130390058",
}

// VerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use VerifierMetaData.ABI instead.
var VerifierABI = VerifierMetaData.ABI

// VerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VerifierMetaData.Bin instead.
var VerifierBin = VerifierMetaData.Bin

// DeployVerifier deploys a new Ethereum contract, binding an instance of Verifier to it.
func DeployVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Verifier, error) {
	parsed, err := VerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// Verifier is an auto generated Go binding around an Ethereum contract.
type Verifier struct {
	VerifierCaller     // Read-only binding to the contract
	VerifierTransactor // Write-only binding to the contract
	VerifierFilterer   // Log filterer for contract events
}

// VerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type VerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VerifierSession struct {
	Contract     *Verifier         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VerifierCallerSession struct {
	Contract *VerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// VerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VerifierTransactorSession struct {
	Contract     *VerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// VerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type VerifierRaw struct {
	Contract *Verifier // Generic contract binding to access the raw methods on
}

// VerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VerifierCallerRaw struct {
	Contract *VerifierCaller // Generic read-only contract binding to access the raw methods on
}

// VerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VerifierTransactorRaw struct {
	Contract *VerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVerifier creates a new instance of Verifier, bound to a specific deployed contract.
func NewVerifier(address common.Address, backend bind.ContractBackend) (*Verifier, error) {
	contract, err := bindVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// NewVerifierCaller creates a new read-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierCaller(address common.Address, caller bind.ContractCaller) (*VerifierCaller, error) {
	contract, err := bindVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierCaller{contract: contract}, nil
}

// NewVerifierTransactor creates a new write-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifierTransactor, error) {
	contract, err := bindVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierTransactor{contract: contract}, nil
}

// NewVerifierFilterer creates a new log filterer instance of Verifier, bound to a specific deployed contract.
func NewVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifierFilterer, error) {
	contract, err := bindVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifierFilterer{contract: contract}, nil
}

// bindVerifier binds a generic wrapper to an already deployed contract.
func bindVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.VerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transact(opts, method, params...)
}

// IsAdded is a free data retrieval call binding the contract method 0x2e3d616e.
//
// Solidity: function isAdded(string ip) view returns(bool)
func (_Verifier *VerifierCaller) IsAdded(opts *bind.CallOpts, ip string) (bool, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "isAdded", ip)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAdded is a free data retrieval call binding the contract method 0x2e3d616e.
//
// Solidity: function isAdded(string ip) view returns(bool)
func (_Verifier *VerifierSession) IsAdded(ip string) (bool, error) {
	return _Verifier.Contract.IsAdded(&_Verifier.CallOpts, ip)
}

// IsAdded is a free data retrieval call binding the contract method 0x2e3d616e.
//
// Solidity: function isAdded(string ip) view returns(bool)
func (_Verifier *VerifierCallerSession) IsAdded(ip string) (bool, error) {
	return _Verifier.Contract.IsAdded(&_Verifier.CallOpts, ip)
}

// Verify is a free data retrieval call binding the contract method 0xe8fc9273.
//
// Solidity: function verify(address _addr, string _ip) view returns(bool)
func (_Verifier *VerifierCaller) Verify(opts *bind.CallOpts, _addr common.Address, _ip string) (bool, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "verify", _addr, _ip)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Verify is a free data retrieval call binding the contract method 0xe8fc9273.
//
// Solidity: function verify(address _addr, string _ip) view returns(bool)
func (_Verifier *VerifierSession) Verify(_addr common.Address, _ip string) (bool, error) {
	return _Verifier.Contract.Verify(&_Verifier.CallOpts, _addr, _ip)
}

// Verify is a free data retrieval call binding the contract method 0xe8fc9273.
//
// Solidity: function verify(address _addr, string _ip) view returns(bool)
func (_Verifier *VerifierCallerSession) Verify(_addr common.Address, _ip string) (bool, error) {
	return _Verifier.Contract.Verify(&_Verifier.CallOpts, _addr, _ip)
}

// Add is a paid mutator transaction binding the contract method 0xb0c8f9dc.
//
// Solidity: function add(string _ip) returns()
func (_Verifier *VerifierTransactor) Add(opts *bind.TransactOpts, _ip string) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "add", _ip)
}

// Add is a paid mutator transaction binding the contract method 0xb0c8f9dc.
//
// Solidity: function add(string _ip) returns()
func (_Verifier *VerifierSession) Add(_ip string) (*types.Transaction, error) {
	return _Verifier.Contract.Add(&_Verifier.TransactOpts, _ip)
}

// Add is a paid mutator transaction binding the contract method 0xb0c8f9dc.
//
// Solidity: function add(string _ip) returns()
func (_Verifier *VerifierTransactorSession) Add(_ip string) (*types.Transaction, error) {
	return _Verifier.Contract.Add(&_Verifier.TransactOpts, _ip)
}
