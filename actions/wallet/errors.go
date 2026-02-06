package wallet

import "fmt"

// AccountNotFoundError is returned when no account is provided to an action that requires one.
// This mirrors viem's AccountNotFoundError.
type AccountNotFoundError struct {
	DocsPath string
}

func (e *AccountNotFoundError) Error() string {
	msg := "could not find an Account to execute with this Action"
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AccountTypeNotSupportedError is returned when an action requires a local account
// but a JSON-RPC account was provided.
// This mirrors viem's AccountTypeNotSupportedError.
type AccountTypeNotSupportedError struct {
	DocsPath     string
	MetaMessages []string
}

func (e *AccountTypeNotSupportedError) Error() string {
	msg := "account type not supported"
	for _, m := range e.MetaMessages {
		msg += "\n" + m
	}
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}
