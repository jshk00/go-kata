package useraggr

import (
	"log/slog"
	"time"
)

type Options struct {
	logger  *slog.Logger
	timeout time.Duration
}

func WithLogger(l *slog.Logger) func(o *Options) {
	return func(o *Options) {
		o.logger = l
	}
}

func WithTimeout(t time.Duration) func(o *Options) {
	return func(o *Options) {
		o.timeout = t
	}
}
