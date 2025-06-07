package clusterapi

import (
	"context"
	"fmt"
	"os"

	ycsdk "github.com/yandex-cloud/go-sdk"
)

var (
	err  error
	ysdk *ycsdk.SDK
)

func init() {
	token := os.Getenv("YC_TOKEN")
	if ysdk, err = ycsdk.Build(context.Background(), ycsdk.Config{
		Credentials: ycsdk.OAuthToken(token),
	}); err != nil {
		panic(fmt.Sprintf("Can't init ya-sdk client %s", err.Error()))
	}
}
