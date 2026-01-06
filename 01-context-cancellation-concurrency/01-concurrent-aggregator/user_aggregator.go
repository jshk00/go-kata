package useraggr

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

type UserAggregator struct {
	timeout time.Duration
	logger  *slog.Logger
}

func NewUserAggregator(opts ...func(o *Options)) *UserAggregator {
	o := &Options{logger: slog.Default(), timeout: time.Second * 1}
	for _, oo := range opts {
		oo(o)
	}
	return &UserAggregator{logger: o.logger, timeout: o.timeout}
}

func (u *UserAggregator) Aggregate(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	g, gctx := errgroup.WithContext(ctx)
	var (
		profile string
		order   int
	)
	g.Go(func() error {
		m, err := ProfileService(gctx, id)
		if err != nil {
			return err
		}
		p, ok := m["Name"]
		if !ok {
			return errors.New("name is not found")
		}
		profile = p
		return nil
	})
	g.Go(func() error {
		m, err := UserService(gctx, id)
		if err != nil {
			return err
		}
		o, ok := m["Orders"]
		if !ok {
			return errors.New("orders is not found")
		}
		order = o
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	fmt.Printf("User: %s | Orders: %d\n", profile, order)
	return nil
}
