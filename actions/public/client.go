// Package public provides standalone action functions for public (read-only) Ethereum JSON-RPC methods.
// These actions can be used directly or through a PublicClient.
//
// This mirrors viem's actions pattern where actions are standalone functions
// that take a client interface as their first parameter.
package public

import (
	"context"
	"time"

	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/types"
)

// Client is the interface that actions require from a client.
// This allows actions to be used with any client implementation
// that satisfies this interface.
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
	// Used for batch scheduler caching.
	UID() string
}

// BlockTag is an alias for types.BlockTag for convenience.
type BlockTag = types.BlockTag

// Block tag constants for convenience.
const (
	BlockTagLatest    BlockTag = types.BlockTagLatest
	BlockTagPending   BlockTag = types.BlockTagPending
	BlockTagEarliest  BlockTag = types.BlockTagEarliest
	BlockTagSafe      BlockTag = types.BlockTagSafe
	BlockTagFinalized BlockTag = types.BlockTagFinalized
)
