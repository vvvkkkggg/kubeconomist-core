package yandex

import (
	"context"

	ycsdk "github.com/yandex-cloud/go-sdk"
)

const (
	MaxPageSize    = 100
	EmptyPageToken = ""
)

type Client struct {
	sdk *ycsdk.SDK
}

func New(ctx context.Context, token string) (*Client, error) {
	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: ycsdk.OAuthToken(token),
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		sdk: sdk,
	}, nil
}
