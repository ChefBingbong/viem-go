package transaction

import (
	"errors"
	"math/big"
)

// TransactionType represents the type of Ethereum transaction.
type TransactionType string

const (
	TransactionTypeLegacy  TransactionType = "legacy"
	TransactionTypeEIP2930 TransactionType = "eip2930"
	TransactionTypeEIP1559 TransactionType = "eip1559"
	TransactionTypeEIP4844 TransactionType = "eip4844"
	TransactionTypeEIP7702 TransactionType = "eip7702"
)

// Common errors
var (
	ErrInvalidAddress                   = errors.New("invalid address")
	ErrInvalidChainId                   = errors.New("invalid chain ID")
	ErrFeeCapTooHigh                    = errors.New("fee cap too high")
	ErrTipAboveFeeCap                   = errors.New("tip above fee cap")
	ErrInvalidSerializedTransactionType = errors.New("invalid serialized transaction type")
	ErrInvalidSerializableTransaction   = errors.New("invalid serializable transaction")
	ErrInvalidSerializedTransaction     = errors.New("invalid serialized transaction")
	ErrInvalidStorageKeySize            = errors.New("invalid storage key size")
	ErrInvalidLegacyV                   = errors.New("invalid legacy v value")
	ErrEmptyBlob                        = errors.New("empty blob")
	ErrInvalidVersionedHashSize         = errors.New("invalid versioned hash size")
	ErrInvalidVersionedHashVersion      = errors.New("invalid versioned hash version")
	ErrMaxFeePerGasNotAllowed           = errors.New("maxFeePerGas/maxPriorityFeePerGas is not allowed for this transaction type")
)

// MaxUint256 is 2^256 - 1
var MaxUint256 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

// VersionedHashVersionKzg is the version byte for KZG commitments (EIP-4844)
const VersionedHashVersionKzg = 0x01

// AccessListItem represents an item in an access list.
type AccessListItem struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

// AccessList is a list of access list items.
type AccessList []AccessListItem

// Authorization represents an EIP-7702 authorization.
type Authorization struct {
	Address string `json:"address"`
	ChainId int    `json:"chainId"`
	Nonce   int    `json:"nonce"`
}

// SignedAuthorization includes signature fields.
type SignedAuthorization struct {
	Authorization
	R       string `json:"r"`
	S       string `json:"s"`
	YParity int    `json:"yParity"`
}

// Signature represents a transaction signature.
type Signature struct {
	R       string   `json:"r"`
	S       string   `json:"s"`
	V       *big.Int `json:"v,omitempty"`
	YParity int      `json:"yParity"`
}

// TransactionBase contains common transaction fields.
type TransactionBase struct {
	Type    TransactionType `json:"type,omitempty"`
	ChainId int             `json:"chainId,omitempty"`
	Nonce   int             `json:"nonce,omitempty"`
	To      string          `json:"to,omitempty"`
	Value   *big.Int        `json:"value,omitempty"`
	Data    string          `json:"data,omitempty"`
	Gas     *big.Int        `json:"gas,omitempty"`

	// Signature fields
	R       string   `json:"r,omitempty"`
	S       string   `json:"s,omitempty"`
	V       *big.Int `json:"v,omitempty"`
	YParity int      `json:"yParity,omitempty"`
}

// TransactionLegacy represents a legacy (pre-EIP-2718) transaction.
type TransactionLegacy struct {
	TransactionBase
	GasPrice *big.Int `json:"gasPrice,omitempty"`
}

// TransactionEIP2930 represents an EIP-2930 transaction.
type TransactionEIP2930 struct {
	TransactionBase
	GasPrice   *big.Int   `json:"gasPrice,omitempty"`
	AccessList AccessList `json:"accessList,omitempty"`
}

// TransactionEIP1559 represents an EIP-1559 transaction.
type TransactionEIP1559 struct {
	TransactionBase
	MaxFeePerGas         *big.Int   `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int   `json:"maxPriorityFeePerGas,omitempty"`
	AccessList           AccessList `json:"accessList,omitempty"`
}

// TransactionEIP4844 represents an EIP-4844 (blob) transaction.
type TransactionEIP4844 struct {
	TransactionEIP1559
	MaxFeePerBlobGas    *big.Int `json:"maxFeePerBlobGas,omitempty"`
	BlobVersionedHashes []string `json:"blobVersionedHashes,omitempty"`
	// Sidecars fields (for wrapper format)
	Blobs       []string `json:"blobs,omitempty"`
	Commitments []string `json:"commitments,omitempty"`
	Proofs      []string `json:"proofs,omitempty"`
}

// TransactionEIP7702 represents an EIP-7702 transaction.
type TransactionEIP7702 struct {
	TransactionEIP1559
	AuthorizationList []SignedAuthorization `json:"authorizationList,omitempty"`
}

// Transaction is a generic transaction that can be any type.
type Transaction struct {
	Type TransactionType `json:"type,omitempty"`

	// Common fields
	ChainId int      `json:"chainId,omitempty"`
	Nonce   int      `json:"nonce,omitempty"`
	To      string   `json:"to,omitempty"`
	Value   *big.Int `json:"value,omitempty"`
	Data    string   `json:"data,omitempty"`
	Gas     *big.Int `json:"gas,omitempty"`

	// Legacy/EIP-2930 fields
	GasPrice *big.Int `json:"gasPrice,omitempty"`

	// EIP-1559 fields
	MaxFeePerGas         *big.Int `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas,omitempty"`

	// EIP-2930 fields
	AccessList AccessList `json:"accessList,omitempty"`

	// EIP-4844 fields
	MaxFeePerBlobGas    *big.Int      `json:"maxFeePerBlobGas,omitempty"`
	BlobVersionedHashes []string      `json:"blobVersionedHashes,omitempty"`
	Blobs               []string      `json:"blobs,omitempty"`
	Sidecars            []BlobSidecar `json:"sidecars,omitempty"`

	// EIP-7702 fields
	AuthorizationList []SignedAuthorization `json:"authorizationList,omitempty"`

	// Signature fields
	R       string   `json:"r,omitempty"`
	S       string   `json:"s,omitempty"`
	V       *big.Int `json:"v,omitempty"`
	YParity int      `json:"yParity,omitempty"`
}

// BlobSidecar represents a blob sidecar (EIP-4844).
type BlobSidecar struct {
	Blob       string `json:"blob"`
	Commitment string `json:"commitment"`
	Proof      string `json:"proof"`
}

// HasSignature returns true if the transaction has signature fields.
func (t *Transaction) HasSignature() bool {
	return t.R != "" && t.S != ""
}

// GetSignature returns the transaction signature if present.
func (t *Transaction) GetSignature() *Signature {
	if !t.HasSignature() {
		return nil
	}
	return &Signature{
		R:       t.R,
		S:       t.S,
		V:       t.V,
		YParity: t.YParity,
	}
}
