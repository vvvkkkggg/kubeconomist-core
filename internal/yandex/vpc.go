package yandex

import (
	"context"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/vpc/v1"
)

func (c *Client) GetAddresses(ctx context.Context, folderID string) ([]*vpc.Address, error) {
	addresses := make([]*vpc.Address, 0)
	pageToken := EmptyPageToken
	for {
		resp, err := c.sdk.VPC().Address().List(ctx, &vpc.ListAddressesRequest{
			FolderId:  folderID,
			PageToken: pageToken,
			PageSize:  MaxPageSize,
		})
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, resp.Addresses...)
		pageToken = resp.NextPageToken

		if pageToken == EmptyPageToken {
			break
		}
	}

	return addresses, nil
}
