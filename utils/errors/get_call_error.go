// Package errors provides error handling utilities.
package errors

import (
	"github.com/ethereum/go-ethereum/common"
)

// CallErrorParams contains parameters for wrapping call errors.
type CallErrorParams struct {
	To      *common.Address
	Data    []byte
	ChainID *int64
}

// CallExecutionError wraps an error with call context.
type CallExecutionError struct {
	Cause   error
	To      *common.Address
	Data    []byte
	ChainID *int64
}

func (e *CallExecutionError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return "call execution failed"
}

func (e *CallExecutionError) Unwrap() error {
	return e.Cause
}

// GetCallError wraps an error with call execution context.
// This provides better error messages with relevant call details.
func GetCallError(err error, params CallErrorParams) error {
	if err == nil {
		return nil
	}

	return &CallExecutionError{
		Cause:   err,
		To:      params.To,
		Data:    params.Data,
		ChainID: params.ChainID,
	}
}

// GetRevertErrorData extracts revert data from an error if available.
// Returns nil if the error doesn't contain revert data.
func GetRevertErrorData(err error) []byte {
	if err == nil {
		return nil
	}

	// Try to extract from known error types
	type dataError interface {
		ErrorData() interface{}
	}
	if e, ok := err.(dataError); ok {
		if data, ok := e.ErrorData().(string); ok {
			return common.FromHex(data)
		}
		if data, ok := e.ErrorData().([]byte); ok {
			return data
		}
	}

	return nil
}
