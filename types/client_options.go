package types

import (
	"context"
	"time"
)

// MulticallBatchOptions contains options for multicall batching.
type MulticallBatchOptions struct {
	// BatchSize is the maximum size (in bytes) for each calldata chunk.
	// Default: 1024
	BatchSize int

	// Deployless enables deployless multicall (doesn't require deployed multicall3).
	Deployless bool

	// Wait is the duration to wait before sending a batch.
	// Default: 0 (send immediately)
	Wait time.Duration
}

// BatchOptions contains batch settings for the client.
type BatchOptions struct {
	// Multicall enables eth_call multicall aggregation.
	// Set to true for default options, or provide MulticallBatchOptions for custom config.
	Multicall *MulticallBatchOptions
}

// CCIPReadOptions contains CCIP-Read configuration.
type CCIPReadOptions struct {
	// Request is a custom CCIP gateway request function.
	// If nil, the default HTTP request is used.
	Request func(ctx context.Context, data []byte, sender string, urls []string) ([]byte, error)
}
