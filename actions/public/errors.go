package public

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// BlockNotFoundError is returned when a block is not found.
type BlockNotFoundError struct {
	BlockHash   *common.Hash
	BlockNumber *uint64
}

func (e *BlockNotFoundError) Error() string {
	if e.BlockHash != nil {
		return fmt.Sprintf("block not found: hash=%s", e.BlockHash.Hex())
	}
	if e.BlockNumber != nil {
		return fmt.Sprintf("block not found: number=%d", *e.BlockNumber)
	}
	return "block not found"
}

// TransactionNotFoundError is returned when a transaction is not found.
type TransactionNotFoundError struct {
	Hash        *common.Hash
	BlockHash   *common.Hash
	BlockNumber *uint64
	Index       *int
}

func (e *TransactionNotFoundError) Error() string {
	if e.Hash != nil {
		return fmt.Sprintf("transaction not found: hash=%s", e.Hash.Hex())
	}
	if e.BlockHash != nil && e.Index != nil {
		return fmt.Sprintf("transaction not found: blockHash=%s, index=%d", e.BlockHash.Hex(), *e.Index)
	}
	if e.BlockNumber != nil && e.Index != nil {
		return fmt.Sprintf("transaction not found: blockNumber=%d, index=%d", *e.BlockNumber, *e.Index)
	}
	return "transaction not found"
}
