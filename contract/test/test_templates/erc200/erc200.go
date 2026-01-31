package contract_test

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/contract"
)

// Suppress unused import warnings
var (
	_ = big.NewInt
	_ = common.Address{}
)

// ContractABI is the ABI of the Erc200 contract.
var ContractABI = `[
    {
        "constant": true,
        "inputs": [],
        "name": "name",
        "outputs": [
            {
                "name": "",
                "type": "string"
            }
        ],
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "symbol",
        "outputs": [
            {
                "name": "",
                "type": "string"
            }
        ],
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "decimals",
        "outputs": [
            {
                "name": "",
                "type": "uint8"
            }
        ],
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "totalSupply",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "owner",
                "type": "address"
            }
        ],
        "name": "balanceOf",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "owner",
                "type": "address"
            },
            {
                "name": "spender",
                "type": "address"
            }
        ],
        "name": "allowance",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "to",
                "type": "address"
            },
            {
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "transfer",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "spender",
                "type": "address"
            },
            {
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "approve",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "from",
                "type": "address"
            },
            {
                "name": "to",
                "type": "address"
            },
            {
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "transferFrom",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "type": "function"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "from",
                "type": "address"
            },
            {
                "indexed": true,
                "name": "to",
                "type": "address"
            },
            {
                "indexed": false,
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "Transfer",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "owner",
                "type": "address"
            },
            {
                "indexed": true,
                "name": "spender",
                "type": "address"
            },
            {
                "indexed": false,
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "Approval",
        "type": "event"
    }
]`

// ============================================================================
// Typed Method Descriptors
// ============================================================================

// Erc200Methods defines typed method descriptors for the Erc200 contract.
// Use these with contract.ReadTyped() and contract.WriteTyped() for type-safe calls.
type Erc200Methods struct {
	Decimals     contract.ReadUint8
	TotalSupply  contract.ReadBigInt
	BalanceOf    contract.ReadBigInt
	Transfer     contract.WriteMethod
	Name         contract.ReadString
	Symbol       contract.ReadString
	Approve      contract.WriteMethod
	Allowance    contract.ReadBigInt
	TransferFrom contract.WriteMethod
}

// Methods is the typed method descriptors instance for Erc200.
// Use with contract.ReadTyped(c.Contract(), ctx, Methods.MethodName, args...)
var Methods = Erc200Methods{
	Decimals:     contract.ReadUint8{Name: "decimals"},
	TotalSupply:  contract.ReadBigInt{Name: "totalSupply"},
	BalanceOf:    contract.ReadBigInt{Name: "balanceOf"},
	Transfer:     contract.WriteMethod{Name: "transfer"},
	Name:         contract.ReadString{Name: "name"},
	Symbol:       contract.ReadString{Name: "symbol"},
	Approve:      contract.WriteMethod{Name: "approve"},
	Allowance:    contract.ReadBigInt{Name: "allowance"},
	TransferFrom: contract.WriteMethod{Name: "transferFrom"},
}

// ============================================================================
// Contract Binding
// ============================================================================

// Erc200 is a binding to the Erc200 contract.
type Erc200 struct {
	contract *contract.Contract
	M        Erc200Methods // Typed method descriptors
}

// New creates a new Erc200 contract binding.
func New(address common.Address, c *client.Client) (*Erc200, error) {
	cont, err := contract.NewContract(address, []byte(ContractABI), c)
	if err != nil {
		return nil, err
	}
	return &Erc200{contract: cont, M: Methods}, nil
}

// MustNew creates a new Erc200 contract binding, panicking on error.
func MustNew(address common.Address, c *client.Client) *Erc200 {
	cont, err := New(address, c)
	if err != nil {
		panic(err)
	}
	return cont
}

// Address returns the contract address.
func (c *Erc200) Address() common.Address {
	return c.contract.Address()
}

// Contract returns the underlying contract instance.
func (c *Erc200) Contract() *contract.Contract {
	return c.contract
}

// Decimals calls the decimals function.
// Solidity: decimals()
func (c *Erc200) Decimals(ctx context.Context) (uint8, error) {
	result, err := c.contract.Read(ctx, "decimals")
	if err != nil {
		return 0, err
	}

	return result[0].(uint8), nil

}

// TotalSupply calls the totalSupply function.
// Solidity: totalSupply()
func (c *Erc200) TotalSupply(ctx context.Context) (*big.Int, error) {
	result, err := c.contract.Read(ctx, "totalSupply")
	if err != nil {
		return nil, err
	}

	return result[0].(*big.Int), nil

}

// BalanceOf calls the balanceOf function.
// Solidity: balanceOf(address)
func (c *Erc200) BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	result, err := c.contract.Read(ctx, "balanceOf", owner)
	if err != nil {
		return nil, err
	}

	return result[0].(*big.Int), nil

}

// Transfer sends a transaction to the transfer function.
// Solidity: transfer(address,uint256)
func (c *Erc200) Transfer(ctx context.Context, opts contract.WriteOptions, to common.Address, value *big.Int) (common.Hash, error) {
	return c.contract.Write(ctx, opts, "transfer", to, value)
}

// TransferAndWait sends a transaction and waits for the receipt.
func (c *Erc200) TransferAndWait(ctx context.Context, opts contract.WriteOptions, to common.Address, value *big.Int) (*client.Receipt, error) {
	return c.contract.WriteAndWait(ctx, opts, "transfer", to, value)
}

// Name calls the name function.
// Solidity: name()
func (c *Erc200) Name(ctx context.Context) (string, error) {
	result, err := c.contract.Read(ctx, "name")
	if err != nil {
		return "", err
	}

	return result[0].(string), nil

}

// Symbol calls the symbol function.
// Solidity: symbol()
func (c *Erc200) Symbol(ctx context.Context) (string, error) {
	result, err := c.contract.Read(ctx, "symbol")
	if err != nil {
		return "", err
	}

	return result[0].(string), nil

}

// Approve sends a transaction to the approve function.
// Solidity: approve(address,uint256)
func (c *Erc200) Approve(ctx context.Context, opts contract.WriteOptions, spender common.Address, value *big.Int) (common.Hash, error) {
	return c.contract.Write(ctx, opts, "approve", spender, value)
}

// ApproveAndWait sends a transaction and waits for the receipt.
func (c *Erc200) ApproveAndWait(ctx context.Context, opts contract.WriteOptions, spender common.Address, value *big.Int) (*client.Receipt, error) {
	return c.contract.WriteAndWait(ctx, opts, "approve", spender, value)
}

// Allowance calls the allowance function.
// Solidity: allowance(address,address)
func (c *Erc200) Allowance(ctx context.Context, owner common.Address, spender common.Address) (*big.Int, error) {
	result, err := c.contract.Read(ctx, "allowance", owner, spender)
	if err != nil {
		return nil, err
	}

	return result[0].(*big.Int), nil

}

// TransferFrom sends a transaction to the transferFrom function.
// Solidity: transferFrom(address,address,uint256)
func (c *Erc200) TransferFrom(ctx context.Context, opts contract.WriteOptions, from common.Address, to common.Address, value *big.Int) (common.Hash, error) {
	return c.contract.Write(ctx, opts, "transferFrom", from, to, value)
}

// TransferFromAndWait sends a transaction and waits for the receipt.
func (c *Erc200) TransferFromAndWait(ctx context.Context, opts contract.WriteOptions, from common.Address, to common.Address, value *big.Int) (*client.Receipt, error) {
	return c.contract.WriteAndWait(ctx, opts, "transferFrom", from, to, value)
}

// Event types

// TransferEvent represents a Transfer event.
type TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// ApprovalEvent represents a Approval event.
type ApprovalEvent struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}
