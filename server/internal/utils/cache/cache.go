package cache

import "context"

type Client interface {
	Set(ctx context.Context, key, value string) (err error)
	Get(ctx context.Context, key string) (value string, err error)
}

