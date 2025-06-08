package yandex

import (
	"context"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
)

func (c *Client) GetInstances(ctx context.Context, folderID string) ([]*compute.Instance, error) {
	instances := make([]*compute.Instance, 0)
	pageToken := EmptyPageToken
	for {
		resp, err := c.sdk.Compute().Instance().List(ctx, &compute.ListInstancesRequest{
			FolderId:  folderID,
			PageToken: pageToken,
			PageSize:  MaxPageSize,
		})
		if err != nil {
			return nil, err
		}

		instances = append(instances, resp.GetInstances()...)
		pageToken = resp.NextPageToken

		if pageToken == EmptyPageToken {
			break
		}
	}

	return instances, nil
}
