// Package erc20 provides bindings for the ERC20 token standard.
package erc20

import (
	"context"
	"math/big"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/contract"
	"github.com/ethereum/go-ethereum/common"
)

// ContractABI is the ABI of the ERC20 contract.
var ContractABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}]`

// ERC20 is a binding to an ERC20 token contract.
type ERC20 struct {
	contract *contract.Contract
}

// New creates a new ERC20 contract binding.
func New(address common.Address, c *client.Client) (*ERC20, error) {
	cont, err := contract.NewContract(address, []byte(ContractABI), c)
	if err != nil {
		return nil, err
	}
	return &ERC20{contract: cont}, nil
}

// MustNew creates a new ERC20 contract binding, panicking on error.
func MustNew(address common.Address, c *client.Client) *ERC20 {
	cont, err := New(address, c)
	if err != nil {
		panic(err)
	}
	return cont
}

// Address returns the contract address.
func (e *ERC20) Address() common.Address {
	return e.contract.Address()
}

// Contract returns the underlying contract instance.
func (e *ERC20) Contract() *contract.Contract {
	return e.contract
}

// Name returns the token name.
func (e *ERC20) Name(ctx context.Context) (string, error) {
	return e.contract.ReadString(ctx, "name")
}

// Symbol returns the token symbol.
func (e *ERC20) Symbol(ctx context.Context) (string, error) {
	return e.contract.ReadString(ctx, "symbol")
}

// Decimals returns the token decimals.
func (e *ERC20) Decimals(ctx context.Context) (uint8, error) {
	return e.contract.ReadUint8(ctx, "decimals")
}

// TotalSupply returns the total token supply.
func (e *ERC20) TotalSupply(ctx context.Context) (*big.Int, error) {
	return e.contract.ReadBigInt(ctx, "totalSupply")
}

// BalanceOf returns the token balance of an address.
func (e *ERC20) BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	return e.contract.ReadBigInt(ctx, "balanceOf", owner)
}

// Allowance returns the allowance of a spender for an owner.
func (e *ERC20) Allowance(ctx context.Context, owner, spender common.Address) (*big.Int, error) {
	return e.contract.ReadBigInt(ctx, "allowance", owner, spender)
}

// Transfer transfers tokens to a recipient.
func (e *ERC20) Transfer(ctx context.Context, opts contract.WriteOptions, to common.Address, amount *big.Int) (common.Hash, error) {
	return e.contract.Write(ctx, opts, "transfer", to, amount)
}

// TransferAndWait transfers tokens and waits for the transaction to be mined.
func (e *ERC20) TransferAndWait(ctx context.Context, opts contract.WriteOptions, to common.Address, amount *big.Int) (*client.Receipt, error) {
	return e.contract.WriteAndWait(ctx, opts, "transfer", to, amount)
}

// Approve approves a spender to transfer tokens.
func (e *ERC20) Approve(ctx context.Context, opts contract.WriteOptions, spender common.Address, amount *big.Int) (common.Hash, error) {
	return e.contract.Write(ctx, opts, "approve", spender, amount)
}

// ApproveAndWait approves a spender and waits for the transaction to be mined.
func (e *ERC20) ApproveAndWait(ctx context.Context, opts contract.WriteOptions, spender common.Address, amount *big.Int) (*client.Receipt, error) {
	return e.contract.WriteAndWait(ctx, opts, "approve", spender, amount)
}

// TransferFrom transfers tokens from one address to another.
func (e *ERC20) TransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, amount *big.Int) (common.Hash, error) {
	return e.contract.Write(ctx, opts, "transferFrom", from, to, amount)
}

// TransferFromAndWait transfers tokens from one address to another and waits.
func (e *ERC20) TransferFromAndWait(ctx context.Context, opts contract.WriteOptions, from, to common.Address, amount *big.Int) (*client.Receipt, error) {
	return e.contract.WriteAndWait(ctx, opts, "transferFrom", from, to, amount)
}

// TransferEvent represents a Transfer event.
type TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// ApprovalEvent represents an Approval event.
type ApprovalEvent struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}

// ParseTransfer parses a Transfer event from a log.
func (e *ERC20) ParseTransfer(log client.Log) (*TransferEvent, error) {
	event, err := e.contract.DecodeEvent("Transfer", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &TransferEvent{
		From:  event["from"].(common.Address),
		To:    event["to"].(common.Address),
		Value: event["value"].(*big.Int),
	}, nil
}

// ParseApproval parses an Approval event from a log.
func (e *ERC20) ParseApproval(log client.Log) (*ApprovalEvent, error) {
	event, err := e.contract.DecodeEvent("Approval", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &ApprovalEvent{
		Owner:   event["owner"].(common.Address),
		Spender: event["spender"].(common.Address),
		Value:   event["value"].(*big.Int),
	}, nil
}
