package longnet

import (
	"errors"
)

var (
	ErrSessionNotExist     = errors.New("session not exist")
	ErrSessionClosed       = errors.New("session closed")
	ErrSessionWriteTimeout = errors.New("session write timeout")
	ErrPacketTooLarge      = errors.New("packet too large")
)
