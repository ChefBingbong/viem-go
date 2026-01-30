package transaction

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ChefBingbong/viem-go/utils/encoding"
)

// SerializeTransaction serializes a transaction to RLP-encoded hex.
//
// Example:
//
//	serialized, err := SerializeTransaction(&Transaction{
//		Type:     TransactionTypeEIP1559,
//		ChainId:  1,
//		Nonce:    0,
//		MaxFeePerGas: big.NewInt(1000000000),
//		Gas:      big.NewInt(21000),
//		To:       "0x...",
//		Value:    big.NewInt(0),
//	}, nil)
func SerializeTransaction(tx *Transaction, signature *Signature) (string, error) {
	txType, err := GetTransactionType(tx)
	if err != nil {
		return "", err
	}

	switch txType {
	case TransactionTypeEIP7702:
		return serializeTransactionEIP7702(tx, signature)
	case TransactionTypeEIP4844:
		return serializeTransactionEIP4844(tx, signature)
	case TransactionTypeEIP1559:
		return serializeTransactionEIP1559(tx, signature)
	case TransactionTypeEIP2930:
		return serializeTransactionEIP2930(tx, signature)
	case TransactionTypeLegacy:
		return serializeTransactionLegacy(tx, signature)
	default:
		return "", fmt.Errorf("%w: unknown transaction type %s", ErrInvalidSerializableTransaction, txType)
	}
}

func serializeTransactionEIP7702(tx *Transaction, signature *Signature) (string, error) {
	if err := AssertTransactionEIP7702(tx); err != nil {
		return "", err
	}

	serializedAccessList, err := SerializeAccessList(tx.AccessList)
	if err != nil {
		return "", err
	}

	serializedAuthList, err := serializeAuthorizationList(tx.AuthorizationList)
	if err != nil {
		return "", err
	}

	fields := []any{
		numberToHexRlp(tx.ChainId),
		numberToHexRlp(tx.Nonce),
		bigIntToHexRlp(tx.MaxPriorityFeePerGas),
		bigIntToHexRlp(tx.MaxFeePerGas),
		bigIntToHexRlp(tx.Gas),
		addressOrEmpty(tx.To),
		bigIntToHexRlp(tx.Value),
		dataOrEmpty(tx.Data),
		serializedAccessList,
		serializedAuthList,
	}

	fields = appendSignature(fields, tx, signature)

	rlpEncoded, err := encoding.RlpEncodeToHex(fields)
	if err != nil {
		return "", err
	}

	return "0x04" + strings.TrimPrefix(rlpEncoded, "0x"), nil
}

func serializeTransactionEIP4844(tx *Transaction, signature *Signature) (string, error) {
	if err := AssertTransactionEIP4844(tx); err != nil {
		return "", err
	}

	serializedAccessList, err := SerializeAccessList(tx.AccessList)
	if err != nil {
		return "", err
	}

	fields := []any{
		numberToHexRlp(tx.ChainId),
		numberToHexRlp(tx.Nonce),
		bigIntToHexRlp(tx.MaxPriorityFeePerGas),
		bigIntToHexRlp(tx.MaxFeePerGas),
		bigIntToHexRlp(tx.Gas),
		addressOrEmpty(tx.To),
		bigIntToHexRlp(tx.Value),
		dataOrEmpty(tx.Data),
		serializedAccessList,
		bigIntToHexRlp(tx.MaxFeePerBlobGas),
		stringSliceToAny(tx.BlobVersionedHashes),
	}

	fields = appendSignature(fields, tx, signature)

	var rlpEncoded string
	var err2 error

	// If sidecars are present, use wrapper format
	if len(tx.Sidecars) > 0 {
		blobs := make([]any, len(tx.Sidecars))
		commitments := make([]any, len(tx.Sidecars))
		proofs := make([]any, len(tx.Sidecars))
		for i, sidecar := range tx.Sidecars {
			blobs[i] = sidecar.Blob
			commitments[i] = sidecar.Commitment
			proofs[i] = sidecar.Proof
		}

		wrapper := []any{fields, blobs, commitments, proofs}
		rlpEncoded, err2 = encoding.RlpEncodeToHex(wrapper)
	} else {
		rlpEncoded, err2 = encoding.RlpEncodeToHex(fields)
	}

	if err2 != nil {
		return "", err2
	}

	return "0x03" + strings.TrimPrefix(rlpEncoded, "0x"), nil
}

func serializeTransactionEIP1559(tx *Transaction, signature *Signature) (string, error) {
	if err := AssertTransactionEIP1559(tx); err != nil {
		return "", err
	}

	serializedAccessList, err := SerializeAccessList(tx.AccessList)
	if err != nil {
		return "", err
	}

	fields := []any{
		numberToHexRlp(tx.ChainId),
		numberToHexRlp(tx.Nonce),
		bigIntToHexRlp(tx.MaxPriorityFeePerGas),
		bigIntToHexRlp(tx.MaxFeePerGas),
		bigIntToHexRlp(tx.Gas),
		addressOrEmpty(tx.To),
		bigIntToHexRlp(tx.Value),
		dataOrEmpty(tx.Data),
		serializedAccessList,
	}

	fields = appendSignature(fields, tx, signature)

	rlpEncoded, err := encoding.RlpEncodeToHex(fields)
	if err != nil {
		return "", err
	}

	return "0x02" + strings.TrimPrefix(rlpEncoded, "0x"), nil
}

func serializeTransactionEIP2930(tx *Transaction, signature *Signature) (string, error) {
	if err := AssertTransactionEIP2930(tx); err != nil {
		return "", err
	}

	serializedAccessList, err := SerializeAccessList(tx.AccessList)
	if err != nil {
		return "", err
	}

	fields := []any{
		numberToHexRlp(tx.ChainId),
		numberToHexRlp(tx.Nonce),
		bigIntToHexRlp(tx.GasPrice),
		bigIntToHexRlp(tx.Gas),
		addressOrEmpty(tx.To),
		bigIntToHexRlp(tx.Value),
		dataOrEmpty(tx.Data),
		serializedAccessList,
	}

	fields = appendSignature(fields, tx, signature)

	rlpEncoded, err := encoding.RlpEncodeToHex(fields)
	if err != nil {
		return "", err
	}

	return "0x01" + strings.TrimPrefix(rlpEncoded, "0x"), nil
}

func serializeTransactionLegacy(tx *Transaction, signature *Signature) (string, error) {
	if err := AssertTransactionLegacy(tx); err != nil {
		return "", err
	}

	chainId := tx.ChainId

	fields := []any{
		numberToHexRlp(tx.Nonce),
		bigIntToHexRlp(tx.GasPrice),
		bigIntToHexRlp(tx.Gas),
		addressOrEmpty(tx.To),
		bigIntToHexRlp(tx.Value),
		dataOrEmpty(tx.Data),
	}

	// Use signature from parameter or from transaction
	sig := signature
	if sig == nil && tx.HasSignature() {
		sig = tx.GetSignature()
	}

	if sig != nil {
		v := calculateLegacyV(sig, chainId)
		r := trimHex(sig.R)
		s := trimHex(sig.S)

		fields = append(fields,
			bigIntToHexRlp(v),
			trimmedHexOrEmpty(r),
			trimmedHexOrEmpty(s),
		)
	} else if chainId > 0 {
		// EIP-155: include chainId for unsigned transactions
		fields = append(fields,
			numberToHexRlp(chainId),
			"0x",
			"0x",
		)
	}

	return encoding.RlpEncodeToHex(fields)
}

// calculateLegacyV calculates the v value for legacy transactions.
func calculateLegacyV(sig *Signature, chainId int) *big.Int {
	if sig.V != nil {
		v := sig.V.Int64()

		// EIP-155 (inferred chainId)
		if v >= 35 {
			inferredChainId := (v - 35) / 2
			if inferredChainId > 0 {
				return sig.V
			}
			if v == 35 {
				return big.NewInt(27)
			}
			return big.NewInt(28)
		}

		// EIP-155 (explicit chainId)
		if chainId > 0 {
			return big.NewInt(int64(chainId*2) + 35 + v - 27)
		}

		// Pre-EIP-155
		if v == 27 {
			return big.NewInt(27)
		}
		return big.NewInt(28)
	}

	// Derive from yParity
	if chainId > 0 {
		return big.NewInt(int64(chainId*2) + 35 + int64(sig.YParity))
	}
	return big.NewInt(27 + int64(sig.YParity))
}

// appendSignature appends signature fields to the transaction fields.
func appendSignature(fields []any, tx *Transaction, signature *Signature) []any {
	sig := signature
	if sig == nil && tx.HasSignature() {
		sig = tx.GetSignature()
	}

	if sig == nil || sig.R == "" || sig.S == "" {
		return fields
	}

	yParity := sig.YParity
	if sig.V != nil {
		v := sig.V.Int64()
		if v == 0 || v == 1 {
			yParity = int(v)
		} else if v == 27 {
			yParity = 0
		} else if v == 28 {
			yParity = 1
		}
	}

	yParityHex := "0x"
	if yParity == 1 {
		yParityHex = "0x01"
	}

	r := trimHex(sig.R)
	s := trimHex(sig.S)

	return append(fields,
		yParityHex,
		trimmedHexOrEmpty(r),
		trimmedHexOrEmpty(s),
	)
}

// serializeAuthorizationList serializes an EIP-7702 authorization list.
func serializeAuthorizationList(authList []SignedAuthorization) ([]any, error) {
	if len(authList) == 0 {
		return []any{}, nil
	}

	result := make([]any, len(authList))
	for i, auth := range authList {
		yParityHex := "0x"
		if auth.YParity == 1 {
			yParityHex = "0x01"
		}

		result[i] = []any{
			numberToHexRlp(auth.ChainId),
			auth.Address,
			numberToHexRlp(auth.Nonce),
			yParityHex,
			auth.R,
			auth.S,
		}
	}

	return result, nil
}

// Helper functions

func numberToHexRlp(n int) string {
	if n == 0 {
		return "0x"
	}
	return "0x" + fmt.Sprintf("%x", n)
}

func bigIntToHexRlp(n *big.Int) string {
	if n == nil || n.Sign() == 0 {
		return "0x"
	}
	return "0x" + n.Text(16)
}

func addressOrEmpty(addr string) string {
	if addr == "" {
		return "0x"
	}
	return addr
}

func dataOrEmpty(data string) string {
	if data == "" {
		return "0x"
	}
	return data
}

func stringSliceToAny(s []string) []any {
	if s == nil {
		return []any{}
	}
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

func trimHex(h string) string {
	h = strings.TrimPrefix(h, "0x")
	h = strings.TrimPrefix(h, "0X")

	// Remove leading zeros
	h = strings.TrimLeft(h, "0")
	if h == "" {
		return "0x00"
	}

	// Ensure even length
	if len(h)%2 != 0 {
		h = "0" + h
	}

	return "0x" + h
}

func trimmedHexOrEmpty(h string) string {
	if h == "0x00" || h == "0x" || h == "" {
		return "0x"
	}
	return h
}

// BytesToHex converts bytes to hex string.
func BytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

// HexToBytes converts hex string to bytes.
func HexToBytes(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}
