package useraggr

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestUserAggregatorServiceTimeout(t *testing.T) {
	ProfileService = func(ctx context.Context, id int) (map[string]string, error) { //nolint
		select {
		case <-time.After(2 * time.Second):
			return map[string]string{"Name": "Alice"}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if err := NewUserAggregator().Aggregate(context.Background(), 1); err != nil {
		if !strings.Contains(err.Error(), "context deadline") {
			t.Errorf("got: %v, want: context deadline exited", err)
		}
		t.Log(err)
	} else {
		t.Error("error should not be nil")
	}
}

func TestUserAggregatorServiceParentCtx(t *testing.T) {
	ProfileService = func(ctx context.Context, id int) (map[string]string, error) { //nolint
		select {
		case <-time.After(5 * time.Second):
			return map[string]string{"Name": "Alice"}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	ctx, cancel := context.WithCancelCause(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		cancel(errors.New("early cancellation"))
	}()
	if err := NewUserAggregator(WithTimeout(5*time.Second)).Aggregate(ctx, 1); err != nil {
		if !strings.Contains(err.Error(), "context canceled") {
			t.Errorf("got: %v, want: context canceled", err)
		}
		t.Log(err)
	} else {
		t.Error("error should not be nil")
	}
}

func TestUserAggregatorServiceError(t *testing.T) {
	ProfileService = func(ctx context.Context, id int) (map[string]string, error) { //nolint
		return nil, errors.New("id not found")
	}
	UserService = func(ctx context.Context, id int) (map[string]int, error) { //nolint
		select {
		case <-time.After(10 * time.Second):
			return map[string]int{"Orders": 5}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if err := NewUserAggregator(WithTimeout(15*time.Second)).Aggregate(context.Background(), 1); err == nil {
		t.Error("error should not be nil")
	} else {
		t.Log(err)
	}
}
