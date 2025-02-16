package blueprint

import (
	"errors"
)

const (
	InvalidRedisUrlMsg = "invalid redis url"
	UnknownMsg         = "unknown"
	NotSupportedMsg    = "not supported"
)

var (
	ErrInvalidRedisUrl = errors.New(InvalidRedisUrlMsg)
	ErrUnknown         = errors.New(UnknownMsg)
	ErrNotSupported    = errors.New(NotSupportedMsg)
)
