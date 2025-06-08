package yandex

import (
	"context"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/storage/v1"
)

func (c *Client) GetBuckets(ctx context.Context, folderID string) ([]*storage.Bucket, error) {
	resp, err := c.sdk.StorageAPI().Bucket().List(ctx, &storage.ListBucketsRequest{
		FolderId: folderID,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetBuckets(), nil
}
