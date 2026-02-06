// Package wallet provides standalone action functions for wallet (write) Ethereum JSON-RPC methods.
// These actions can be used directly or through a WalletClient.
//
// This mirrors viem's actions pattern where actions are standalone functions
// that take a client interface as their first parameter.
package wallet

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/types"
	"github.com/ChefBingbong/viem-go/utils/signature"
	utiltx "github.com/ChefBingbong/viem-go/utils/transaction"
)

// Client is the interface that wallet actions require from a client.
// This embeds the same methods as public.Client so wallet actions can
// delegate to public actions (e.g., getChainId, getTransactionCount)
// without requiring a separate public client.
type Client interface {
	// Request sends a raw JSON-RPC request.
	Request(ctx context.Context, method string, params ...any) (*transport.RPCResponse, error)

	// Chain returns the chain configuration, if set.
	Chain() *chain.Chain

	// CacheTime returns the cache duration for cached data.
	CacheTime() time.Duration

	// ExperimentalBlockTag returns the default block tag for RPC requests.
	ExperimentalBlockTag() types.BlockTag

	// Batch returns the batch options, if configured.
	// Returns nil if batching is not enabled.
	Batch() *types.BatchOptions

	// CCIPRead returns the CCIP-Read options.
	// Returns nil if CCIP-Read should use defaults, or false to disable.
	CCIPRead() *types.CCIPReadOptions

	// UID returns the unique identifier for this client instance.
	UID() string

	// Account returns the account associated with this client, if any.
	// Returns nil if no account is set.
	Account() Account
}

// Account represents an account that can be used with the client.
// This mirrors the client package's Account interface.
type Account interface {
	// Address returns the account address.
	Address() common.Address
}

// LocalAccount extends Account to indicate a locally-managed account (private key, HD, etc.).
// When present, actions like GetAddresses can return the local address directly
// without making an RPC call, mirroring viem's `client.account?.type === 'local'` check.
type LocalAccount interface {
	Account
	// IsLocal is a marker method indicating this is a local account.
	IsLocal()
}

// SignableAccount represents an account that can sign messages locally.
// This mirrors viem's account.signMessage capability.
type SignableAccount interface {
	Account
	// SignMessage signs a message and returns the signature as a hex string.
	SignMessage(message signature.SignableMessage) (string, error)
}

// TypedDataSignableAccount represents an account that can sign EIP-712 typed data locally.
// This mirrors viem's account.signTypedData capability.
type TypedDataSignableAccount interface {
	Account
	// SignTypedData signs EIP-712 typed data and returns the signature as a hex string.
	SignTypedData(data signature.TypedDataDefinition) (string, error)
}

// TransactionSignableAccount represents an account that can sign transactions locally.
// This mirrors viem's account.signTransaction capability.
type TransactionSignableAccount interface {
	Account
	// SignTransaction signs a transaction and returns the serialized signed transaction hex.
	SignTransaction(tx *utiltx.Transaction) (string, error)
}

// AuthorizationSignableAccount represents an account that can sign EIP-7702 authorizations locally.
// This mirrors viem's account.signAuthorization capability.
type AuthorizationSignableAccount interface {
	Account
	// SignAuthorization signs an EIP-7702 authorization and returns the signed authorization.
	SignAuthorization(auth types.AuthorizationRequest) (*types.SignedAuthorization, error)
}
