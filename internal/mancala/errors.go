package mancala

import (
	"errors"
)

var (
	ErrInvalidPlay  = errors.New("invalid play")
	ErrGameNotFound = errors.New("game not found")
)
