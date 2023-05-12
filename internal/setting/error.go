package setting

import (
	"errors"
)

var (
	// internal
	MissingConditionErr      error
	TransactionInProgressErr error
	TransactionNotStartedErr error

	// chain service
	ErrClientConnectionFailure error
	ErrNotSupportedChain       error
)

func init() {
	MissingConditionErr = errors.New("missing conditions")
	TransactionInProgressErr = errors.New("transaction already in progress")
	TransactionNotStartedErr = errors.New("transaction not started")

	ErrClientConnectionFailure = errors.New("failed to connect to node rpc endpoint")
	ErrNotSupportedChain = errors.New("not supported chain")
}
