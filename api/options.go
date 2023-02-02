package api

import (
	"errors"
	"time"
)

var (
	errPort    = errors.New("port should be positive")
	errTimeout = errors.New("timeout should be positive")
)

type options struct {
	port    *int
	timeout *time.Duration
}

type Option func(options *options) error

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errPort
		}
		options.port = &port
		return nil
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(options *options) error {
		if timeout < 0 {
			return errTimeout
		}
		options.timeout = &timeout
		return nil
	}
}
