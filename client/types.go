package client

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// BlockNumber represents a block number or tag.
type BlockNumber interface {
	String() string
}

// BlockNumberTag represents a block tag like "latest", "pending", etc.
type BlockNumberTag string

const (
	// BlockLatest represents the latest mined block.
	BlockLatest BlockNumberTag = "latest"
	// BlockPending represents the pending state/transactions.
	BlockPending BlockNumberTag = "pending"
	// BlockEarliest represents the earliest/genesis block.
	BlockEarliest BlockNumberTag = "earliest"
	// BlockSafe represents the latest safe block.
	BlockSafe BlockNumberTag = "safe"
	// BlockFinalized represents the latest finalized block.
	BlockFinalized BlockNumberTag = "finalized"
)

func (b BlockNumberTag) String() string {
	return string(b)
}

// BlockNumberUint64 represents a specific block number.
type BlockNumberUint64 uint64

func (b BlockNumberUint64) String() string {
	return hexutil.EncodeUint64(uint64(b))
}

// CallRequest represents the parameters for an eth_call request.
type CallRequest struct {
	From     *common.Address `json:"from,omitempty"`
	To       common.Address  `json:"to"`
	Data     []byte          `json:"data,omitempty"`
	Value    *big.Int        `json:"value,omitempty"`
	Gas      uint64          `json:"gas,omitempty"`
	GasPrice *big.Int        `json:"gasPrice,omitempty"`
}

// MarshalJSON implements json.Marshaler for CallRequest.
func (c CallRequest) MarshalJSON() ([]byte, error) {
	type callRequestJSON struct {
		From     *common.Address `json:"from,omitempty"`
		To       common.Address  `json:"to"`
		Data     string          `json:"data,omitempty"`
		Value    string          `json:"value,omitempty"`
		Gas      string          `json:"gas,omitempty"`
		GasPrice string          `json:"gasPrice,omitempty"`
	}

	req := callRequestJSON{
		From: c.From,
		To:   c.To,
	}

	if len(c.Data) > 0 {
		req.Data = hexutil.Encode(c.Data)
	}
	if c.Value != nil {
		req.Value = hexutil.EncodeBig(c.Value)
	}
	if c.Gas > 0 {
		req.Gas = hexutil.EncodeUint64(c.Gas)
	}
	if c.GasPrice != nil {
		req.GasPrice = hexutil.EncodeBig(c.GasPrice)
	}

	return json.Marshal(req)
}

// Transaction represents a transaction to be sent.
type Transaction struct {
	From                 common.Address  `json:"from"`
	To                   *common.Address `json:"to,omitempty"`
	Data                 []byte          `json:"data,omitempty"`
	Value                *big.Int        `json:"value,omitempty"`
	Nonce                *uint64         `json:"nonce,omitempty"`
	Gas                  uint64          `json:"gas,omitempty"`
	GasPrice             *big.Int        `json:"gasPrice,omitempty"`
	MaxFeePerGas         *big.Int        `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int        `json:"maxPriorityFeePerGas,omitempty"`
	ChainID              *big.Int        `json:"chainId,omitempty"`
}

// MarshalJSON implements json.Marshaler for Transaction.
func (t Transaction) MarshalJSON() ([]byte, error) {
	type txJSON struct {
		From                 common.Address  `json:"from"`
		To                   *common.Address `json:"to,omitempty"`
		Data                 string          `json:"data,omitempty"`
		Value                string          `json:"value,omitempty"`
		Nonce                string          `json:"nonce,omitempty"`
		Gas                  string          `json:"gas,omitempty"`
		GasPrice             string          `json:"gasPrice,omitempty"`
		MaxFeePerGas         string          `json:"maxFeePerGas,omitempty"`
		MaxPriorityFeePerGas string          `json:"maxPriorityFeePerGas,omitempty"`
		ChainID              string          `json:"chainId,omitempty"`
	}

	tx := txJSON{
		From: t.From,
		To:   t.To,
	}

	if len(t.Data) > 0 {
		tx.Data = hexutil.Encode(t.Data)
	}
	if t.Value != nil {
		tx.Value = hexutil.EncodeBig(t.Value)
	}
	if t.Nonce != nil {
		tx.Nonce = hexutil.EncodeUint64(*t.Nonce)
	}
	if t.Gas > 0 {
		tx.Gas = hexutil.EncodeUint64(t.Gas)
	}
	if t.GasPrice != nil {
		tx.GasPrice = hexutil.EncodeBig(t.GasPrice)
	}
	if t.MaxFeePerGas != nil {
		tx.MaxFeePerGas = hexutil.EncodeBig(t.MaxFeePerGas)
	}
	if t.MaxPriorityFeePerGas != nil {
		tx.MaxPriorityFeePerGas = hexutil.EncodeBig(t.MaxPriorityFeePerGas)
	}
	if t.ChainID != nil {
		tx.ChainID = hexutil.EncodeBig(t.ChainID)
	}

	return json.Marshal(tx)
}

// Receipt represents a transaction receipt.
type Receipt struct {
	TransactionHash   common.Hash     `json:"transactionHash"`
	TransactionIndex  uint64          `json:"transactionIndex"`
	BlockHash         common.Hash     `json:"blockHash"`
	BlockNumber       uint64          `json:"blockNumber"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	GasUsed           uint64          `json:"gasUsed"`
	ContractAddress   *common.Address `json:"contractAddress"`
	Logs              []Log           `json:"logs"`
	Status            uint64          `json:"status"`
	LogsBloom         []byte          `json:"logsBloom"`
	EffectiveGasPrice *big.Int        `json:"effectiveGasPrice"`
}

// UnmarshalJSON implements json.Unmarshaler for Receipt.
func (r *Receipt) UnmarshalJSON(data []byte) error {
	type receiptJSON struct {
		TransactionHash   common.Hash     `json:"transactionHash"`
		TransactionIndex  string          `json:"transactionIndex"`
		BlockHash         common.Hash     `json:"blockHash"`
		BlockNumber       string          `json:"blockNumber"`
		From              common.Address  `json:"from"`
		To                *common.Address `json:"to"`
		CumulativeGasUsed string          `json:"cumulativeGasUsed"`
		GasUsed           string          `json:"gasUsed"`
		ContractAddress   *common.Address `json:"contractAddress"`
		Logs              []Log           `json:"logs"`
		Status            string          `json:"status"`
		LogsBloom         string          `json:"logsBloom"`
		EffectiveGasPrice string          `json:"effectiveGasPrice"`
	}

	var raw receiptJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.TransactionHash = raw.TransactionHash
	r.BlockHash = raw.BlockHash
	r.From = raw.From
	r.To = raw.To
	r.ContractAddress = raw.ContractAddress
	r.Logs = raw.Logs

	if raw.TransactionIndex != "" {
		idx, err := hexutil.DecodeUint64(raw.TransactionIndex)
		if err != nil {
			return fmt.Errorf("invalid transaction index: %w", err)
		}
		r.TransactionIndex = idx
	}

	if raw.BlockNumber != "" {
		bn, err := hexutil.DecodeUint64(raw.BlockNumber)
		if err != nil {
			return fmt.Errorf("invalid block number: %w", err)
		}
		r.BlockNumber = bn
	}

	if raw.CumulativeGasUsed != "" {
		cgu, err := hexutil.DecodeUint64(raw.CumulativeGasUsed)
		if err != nil {
			return fmt.Errorf("invalid cumulative gas used: %w", err)
		}
		r.CumulativeGasUsed = cgu
	}

	if raw.GasUsed != "" {
		gu, err := hexutil.DecodeUint64(raw.GasUsed)
		if err != nil {
			return fmt.Errorf("invalid gas used: %w", err)
		}
		r.GasUsed = gu
	}

	if raw.Status != "" {
		status, err := hexutil.DecodeUint64(raw.Status)
		if err != nil {
			return fmt.Errorf("invalid status: %w", err)
		}
		r.Status = status
	}

	if raw.LogsBloom != "" {
		bloom, err := hexutil.Decode(raw.LogsBloom)
		if err != nil {
			return fmt.Errorf("invalid logs bloom: %w", err)
		}
		r.LogsBloom = bloom
	}

	if raw.EffectiveGasPrice != "" {
		egp, err := hexutil.DecodeBig(raw.EffectiveGasPrice)
		if err != nil {
			return fmt.Errorf("invalid effective gas price: %w", err)
		}
		r.EffectiveGasPrice = egp
	}

	return nil
}

// IsSuccess returns true if the transaction was successful.
func (r *Receipt) IsSuccess() bool {
	return r.Status == 1
}

// Log represents a log entry from a transaction receipt.
type Log struct {
	Address          common.Address `json:"address"`
	Topics           []common.Hash  `json:"topics"`
	Data             []byte         `json:"data"`
	BlockNumber      uint64         `json:"blockNumber"`
	TransactionHash  common.Hash    `json:"transactionHash"`
	TransactionIndex uint64         `json:"transactionIndex"`
	BlockHash        common.Hash    `json:"blockHash"`
	LogIndex         uint64         `json:"logIndex"`
	Removed          bool           `json:"removed"`
}

// UnmarshalJSON implements json.Unmarshaler for Log.
func (l *Log) UnmarshalJSON(data []byte) error {
	type logJSON struct {
		Address          common.Address `json:"address"`
		Topics           []common.Hash  `json:"topics"`
		Data             string         `json:"data"`
		BlockNumber      string         `json:"blockNumber"`
		TransactionHash  common.Hash    `json:"transactionHash"`
		TransactionIndex string         `json:"transactionIndex"`
		BlockHash        common.Hash    `json:"blockHash"`
		LogIndex         string         `json:"logIndex"`
		Removed          bool           `json:"removed"`
	}

	var raw logJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	l.Address = raw.Address
	l.Topics = raw.Topics
	l.TransactionHash = raw.TransactionHash
	l.BlockHash = raw.BlockHash
	l.Removed = raw.Removed

	if raw.Data != "" {
		d, err := hexutil.Decode(raw.Data)
		if err != nil {
			return fmt.Errorf("invalid log data: %w", err)
		}
		l.Data = d
	}

	if raw.BlockNumber != "" {
		bn, err := hexutil.DecodeUint64(raw.BlockNumber)
		if err != nil {
			return fmt.Errorf("invalid block number: %w", err)
		}
		l.BlockNumber = bn
	}

	if raw.TransactionIndex != "" {
		ti, err := hexutil.DecodeUint64(raw.TransactionIndex)
		if err != nil {
			return fmt.Errorf("invalid transaction index: %w", err)
		}
		l.TransactionIndex = ti
	}

	if raw.LogIndex != "" {
		li, err := hexutil.DecodeUint64(raw.LogIndex)
		if err != nil {
			return fmt.Errorf("invalid log index: %w", err)
		}
		l.LogIndex = li
	}

	return nil
}

// Block represents an Ethereum block.
type Block struct {
	Number           uint64         `json:"number"`
	Hash             common.Hash    `json:"hash"`
	ParentHash       common.Hash    `json:"parentHash"`
	Nonce            uint64         `json:"nonce"`
	Miner            common.Address `json:"miner"`
	Difficulty       *big.Int       `json:"difficulty"`
	TotalDifficulty  *big.Int       `json:"totalDifficulty"`
	GasLimit         uint64         `json:"gasLimit"`
	GasUsed          uint64         `json:"gasUsed"`
	Timestamp        uint64         `json:"timestamp"`
	BaseFeePerGas    *big.Int       `json:"baseFeePerGas"`
	Transactions     []common.Hash  `json:"transactions"`
	TransactionsRoot common.Hash    `json:"transactionsRoot"`
	StateRoot        common.Hash    `json:"stateRoot"`
	ReceiptsRoot     common.Hash    `json:"receiptsRoot"`
}

// RPCRequest represents a JSON-RPC request.
type RPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

// RPCResponse represents a JSON-RPC response.
type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      uint64          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Error implements the error interface.
func (e *RPCError) Error() string {
	if e.Data != nil {
		return fmt.Sprintf("RPC error %d: %s (data: %v)", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

// FilterQuery represents parameters for eth_getLogs.
type FilterQuery struct {
	FromBlock BlockNumber     `json:"fromBlock,omitempty"`
	ToBlock   BlockNumber     `json:"toBlock,omitempty"`
	Addresses []common.Address `json:"address,omitempty"`
	Topics    [][]common.Hash  `json:"topics,omitempty"`
}

// MarshalJSON implements json.Marshaler for FilterQuery.
func (f FilterQuery) MarshalJSON() ([]byte, error) {
	type filterJSON struct {
		FromBlock string           `json:"fromBlock,omitempty"`
		ToBlock   string           `json:"toBlock,omitempty"`
		Address   []common.Address `json:"address,omitempty"`
		Topics    [][]common.Hash  `json:"topics,omitempty"`
	}

	fj := filterJSON{
		Address: f.Addresses,
		Topics:  f.Topics,
	}

	if f.FromBlock != nil {
		fj.FromBlock = f.FromBlock.String()
	}
	if f.ToBlock != nil {
		fj.ToBlock = f.ToBlock.String()
	}

	return json.Marshal(fj)
}
