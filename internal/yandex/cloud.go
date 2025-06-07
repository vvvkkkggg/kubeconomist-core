package yandex

import (
	"context"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/resourcemanager/v1"
)

func (c *Client) GetClouds(ctx context.Context) ([]*resourcemanager.Cloud, error) {
	clouds := make([]*resourcemanager.Cloud, 0)
	page := EmptyPageToken
	for {
		resp, err := c.sdk.ResourceManager().Cloud().List(ctx, &resourcemanager.ListCloudsRequest{
			PageToken: page,
			PageSize:  MaxPageSize,
		})
		if err != nil {
			return nil, err
		}

		clouds = append(clouds, resp.Clouds...)
		page = resp.NextPageToken

		if page == EmptyPageToken {
			break
		}
	}

	return clouds, nil
}

func (c *Client) GetFolders(ctx context.Context, cloudID string) ([]*resourcemanager.Folder, error) {
	folders := make([]*resourcemanager.Folder, 0)
	pageToken := EmptyPageToken
	for {
		resp, err := c.sdk.ResourceManager().Folder().List(ctx, &resourcemanager.ListFoldersRequest{
			CloudId:   cloudID,
			PageToken: pageToken,
			PageSize:  MaxPageSize,
		})
		if err != nil {
			return nil, err
		}

		folders = append(folders, resp.Folders...)
		pageToken = resp.NextPageToken

		if pageToken == EmptyPageToken {
			break
		}
	}

	return folders, nil
}
