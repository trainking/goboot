package gameapi

import "github.com/pkg/errors"

var (
	ErrNoImplementNetwork = errors.New("no implement network: ")
	ErrConnClosing        = errors.New("session connect is closing")
)
