package yandex

import (
	"context"
	"fmt"

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

func (c *Client) GetAllFolders(ctx context.Context, cloudID, folderID string) ([]*resourcemanager.Folder, error) {
	if folderID != "" {
		folder, err := c.sdk.ResourceManager().Folder().Get(ctx, &resourcemanager.GetFolderRequest{
			FolderId: folderID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get folder %s: %w", folderID, err)
		}

		return []*resourcemanager.Folder{folder}, nil
	}

	clouds, err := c.GetClouds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get clouds: %w", err)
	}

	res := make([]*resourcemanager.Folder, 0)

	for _, cloud := range clouds {
		folders, err := c.GetFolders(ctx, cloud.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get folders %s: %w", cloud.Id, err)
		}

		res = append(res, folders...)
	}

	return res, nil
}
