package formatters

// Formatter represents a generic formatter with type information.
type Formatter[P any, R any] struct {
	Type    string
	Format  func(params P) R
	Exclude []string
}

// NewFormatter creates a new formatter with the specified type and format function.
func NewFormatter[P any, R any](formatterType string, format func(params P) R) *Formatter[P, R] {
	return &Formatter[P, R]{
		Type:   formatterType,
		Format: format,
	}
}

// WithExclude returns a new formatter that excludes the specified keys from output.
func (f *Formatter[P, R]) WithExclude(exclude []string) *Formatter[P, R] {
	return &Formatter[P, R]{
		Type:    f.Type,
		Format:  f.Format,
		Exclude: exclude,
	}
}

// BlockFormatter is the default block formatter.
var BlockFormatter = NewFormatter("block", FormatBlock)

// TransactionFormatter is the default transaction formatter.
var TransactionFormatter = NewFormatter("transaction", FormatTransaction)

// TransactionReceiptFormatter is the default transaction receipt formatter.
var TransactionReceiptFormatter = NewFormatter("transactionReceipt", FormatTransactionReceipt)

// TransactionRequestFormatter is the default transaction request formatter.
var TransactionRequestFormatter = NewFormatter("transactionRequest", FormatTransactionRequest)

// LogFormatter creates a log formatter with the given options.
func LogFormatter(opts *LogFormatOptions) func(RpcLog) Log {
	return func(log RpcLog) Log {
		return FormatLog(log, opts)
	}
}

// FeeHistoryFormatter is the default fee history formatter.
var FeeHistoryFormatter = NewFormatter("feeHistory", FormatFeeHistory)

// ProofFormatter is the default proof formatter.
var ProofFormatter = NewFormatter("proof", FormatProof)
