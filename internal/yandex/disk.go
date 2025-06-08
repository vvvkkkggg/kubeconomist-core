package yandex

import (
	"context"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
)

func (c *Client) GetDisk(ctx context.Context, folderID string) ([]*compute.Disk, error) {
	disks := make([]*compute.Disk, 0)
	pageToken := EmptyPageToken
	for {
		resp, err := c.sdk.Compute().Disk().List(ctx, &compute.ListDisksRequest{
			FolderId:  folderID,
			PageToken: pageToken,
			PageSize:  MaxPageSize,
		})
		if err != nil {
			return nil, err
		}

		disks = append(disks, resp.GetDisks()...)
		pageToken = resp.NextPageToken

		if pageToken == EmptyPageToken {
			break
		}
	}

	return disks, nil
}
