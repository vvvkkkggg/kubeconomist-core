package yandex

import (
	"context"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/k8s/v1"
)

func (c *Client) GetNodeGroups(ctx context.Context, folderID string) ([]*k8s.NodeGroup, error) {
	nodeGroups := make([]*k8s.NodeGroup, 0)
	pageToken := EmptyPageToken
	for {
		resp, err := c.sdk.Kubernetes().NodeGroup().List(ctx, &k8s.ListNodeGroupsRequest{
			FolderId:  folderID,
			PageToken: pageToken,
			PageSize:  MaxPageSize,
		})
		if err != nil {
			return nil, err
		}

		nodeGroups = append(nodeGroups, resp.NodeGroups...)
		pageToken = resp.NextPageToken

		if pageToken == EmptyPageToken {
			break
		}
	}

	return nodeGroups, nil
}
