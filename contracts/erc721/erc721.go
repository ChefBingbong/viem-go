// Package erc721 provides bindings for the ERC721 NFT standard.
package erc721

import (
	"context"
	"math/big"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/contract"
	"github.com/ethereum/go-ethereum/common"
)

// ContractABI is the ABI of the ERC721 contract.
var ContractABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"getApproved","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"approve","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"operator","type":"address"},{"name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"safeTransferFrom","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"},{"name":"data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"tokenId","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"approved","type":"address"},{"indexed":true,"name":"tokenId","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"operator","type":"address"},{"indexed":false,"name":"approved","type":"bool"}],"name":"ApprovalForAll","type":"event"}]`

// ERC721 is a binding to an ERC721 NFT contract.
type ERC721 struct {
	contract *contract.Contract
}

// New creates a new ERC721 contract binding.
func New(address common.Address, c *client.Client) (*ERC721, error) {
	cont, err := contract.NewContract(address, []byte(ContractABI), c)
	if err != nil {
		return nil, err
	}
	return &ERC721{contract: cont}, nil
}

// MustNew creates a new ERC721 contract binding, panicking on error.
func MustNew(address common.Address, c *client.Client) *ERC721 {
	cont, err := New(address, c)
	if err != nil {
		panic(err)
	}
	return cont
}

// Address returns the contract address.
func (e *ERC721) Address() common.Address {
	return e.contract.Address()
}

// Contract returns the underlying contract instance.
func (e *ERC721) Contract() *contract.Contract {
	return e.contract
}

// Name returns the token name.
func (e *ERC721) Name(ctx context.Context) (string, error) {
	return e.contract.ReadString(ctx, "name")
}

// Symbol returns the token symbol.
func (e *ERC721) Symbol(ctx context.Context) (string, error) {
	return e.contract.ReadString(ctx, "symbol")
}

// TokenURI returns the URI for a token.
func (e *ERC721) TokenURI(ctx context.Context, tokenId *big.Int) (string, error) {
	return e.contract.ReadString(ctx, "tokenURI", tokenId)
}

// BalanceOf returns the number of NFTs owned by an address.
func (e *ERC721) BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	return e.contract.ReadBigInt(ctx, "balanceOf", owner)
}

// OwnerOf returns the owner of a token.
func (e *ERC721) OwnerOf(ctx context.Context, tokenId *big.Int) (common.Address, error) {
	return e.contract.ReadAddress(ctx, "ownerOf", tokenId)
}

// GetApproved returns the approved address for a token.
func (e *ERC721) GetApproved(ctx context.Context, tokenId *big.Int) (common.Address, error) {
	return e.contract.ReadAddress(ctx, "getApproved", tokenId)
}

// IsApprovedForAll returns if an operator is approved for all tokens of an owner.
func (e *ERC721) IsApprovedForAll(ctx context.Context, owner, operator common.Address) (bool, error) {
	return e.contract.ReadBool(ctx, "isApprovedForAll", owner, operator)
}

// Approve approves an address to transfer a token.
func (e *ERC721) Approve(ctx context.Context, opts contract.WriteOptions, to common.Address, tokenId *big.Int) (common.Hash, error) {
	return e.contract.Write(ctx, opts, "approve", to, tokenId)
}

// ApproveAndWait approves and waits for the transaction to be mined.
func (e *ERC721) ApproveAndWait(ctx context.Context, opts contract.WriteOptions, to common.Address, tokenId *big.Int) (*client.Receipt, error) {
	return e.contract.WriteAndWait(ctx, opts, "approve", to, tokenId)
}

// SetApprovalForAll sets approval for all tokens to an operator.
func (e *ERC721) SetApprovalForAll(ctx context.Context, opts contract.WriteOptions, operator common.Address, approved bool) (common.Hash, error) {
	return e.contract.Write(ctx, opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAllAndWait sets approval for all and waits for the transaction.
func (e *ERC721) SetApprovalForAllAndWait(ctx context.Context, opts contract.WriteOptions, operator common.Address, approved bool) (*client.Receipt, error) {
	return e.contract.WriteAndWait(ctx, opts, "setApprovalForAll", operator, approved)
}

// TransferFrom transfers a token from one address to another.
func (e *ERC721) TransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, tokenId *big.Int) (common.Hash, error) {
	return e.contract.Write(ctx, opts, "transferFrom", from, to, tokenId)
}

// TransferFromAndWait transfers a token and waits for the transaction.
func (e *ERC721) TransferFromAndWait(ctx context.Context, opts contract.WriteOptions, from, to common.Address, tokenId *big.Int) (*client.Receipt, error) {
	return e.contract.WriteAndWait(ctx, opts, "transferFrom", from, to, tokenId)
}

// SafeTransferFrom safely transfers a token from one address to another.
func (e *ERC721) SafeTransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, tokenId *big.Int) (common.Hash, error) {
	return e.contract.Write(ctx, opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFromAndWait safely transfers and waits for the transaction.
func (e *ERC721) SafeTransferFromAndWait(ctx context.Context, opts contract.WriteOptions, from, to common.Address, tokenId *big.Int) (*client.Receipt, error) {
	return e.contract.WriteAndWait(ctx, opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFromWithData safely transfers a token with additional data.
func (e *ERC721) SafeTransferFromWithData(ctx context.Context, opts contract.WriteOptions, from, to common.Address, tokenId *big.Int, data []byte) (common.Hash, error) {
	return e.contract.Write(ctx, opts, "safeTransferFrom", from, to, tokenId, data)
}

// TransferEvent represents a Transfer event.
type TransferEvent struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

// ApprovalEvent represents an Approval event.
type ApprovalEvent struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
}

// ApprovalForAllEvent represents an ApprovalForAll event.
type ApprovalForAllEvent struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
}

// ParseTransfer parses a Transfer event from a log.
func (e *ERC721) ParseTransfer(log client.Log) (*TransferEvent, error) {
	event, err := e.contract.DecodeEvent("Transfer", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &TransferEvent{
		From:    event["from"].(common.Address),
		To:      event["to"].(common.Address),
		TokenId: event["tokenId"].(*big.Int),
	}, nil
}

// ParseApproval parses an Approval event from a log.
func (e *ERC721) ParseApproval(log client.Log) (*ApprovalEvent, error) {
	event, err := e.contract.DecodeEvent("Approval", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &ApprovalEvent{
		Owner:    event["owner"].(common.Address),
		Approved: event["approved"].(common.Address),
		TokenId:  event["tokenId"].(*big.Int),
	}, nil
}

// ParseApprovalForAll parses an ApprovalForAll event from a log.
func (e *ERC721) ParseApprovalForAll(log client.Log) (*ApprovalForAllEvent, error) {
	event, err := e.contract.DecodeEvent("ApprovalForAll", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &ApprovalForAllEvent{
		Owner:    event["owner"].(common.Address),
		Operator: event["operator"].(common.Address),
		Approved: event["approved"].(bool),
	}, nil
}
