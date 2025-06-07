package yandex

import (
	"context"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
)

func (c *Client) GetImages(ctx context.Context, folderID string) (*compute.ListImagesResponse, error) {
	return c.sdk.Compute().Image().List(ctx, &compute.ListImagesRequest{
		FolderId: folderID,
	})
}
