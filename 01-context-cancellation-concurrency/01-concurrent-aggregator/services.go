package useraggr

import "context"

var ProfileService = func(ctx context.Context, id int) (map[string]string, error) { //nolint
	return map[string]string{"Name": "Alice"}, nil
}

var UserService = func(ctx context.Context, id int) (map[string]int, error) { //nolint
	return map[string]int{"Orders": 5}, nil
}
