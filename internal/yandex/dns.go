package yandex

import (
	"context"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/dns/v1"
)

func (c *Client) GetDNSZones(ctx context.Context, folderID string) ([]*dns.DnsZone, error) {
	zones := make([]*dns.DnsZone, 0)
	pageToken := EmptyPageToken
	for {
		resp, err := c.sdk.DNS().DnsZone().List(ctx, &dns.ListDnsZonesRequest{
			FolderId:  folderID,
			PageToken: pageToken,
			PageSize:  MaxPageSize,
		})
		if err != nil {
			return nil, err
		}
		resp.GetDnsZones()

		zones = append(zones, resp.DnsZones...)
		pageToken = resp.NextPageToken

		if pageToken == EmptyPageToken {
			break
		}
	}

	return zones, nil
}

func (c *Client) GetDNSRecords(ctx context.Context, zoneID string) ([]*dns.RecordSet, error) {
	records := make([]*dns.RecordSet, 0)
	pageToken := EmptyPageToken
	for {
		resp, err := c.sdk.DNS().DnsZone().ListRecordSets(ctx, &dns.ListDnsZoneRecordSetsRequest{
			DnsZoneId: zoneID,
			PageToken: pageToken,
			PageSize:  MaxPageSize,
		})
		if err != nil {
			return nil, err
		}

		records = append(records, resp.RecordSets...)
		pageToken = resp.NextPageToken

		if pageToken == EmptyPageToken {
			break
		}
	}

	return records, nil
}

func (c *Client) IsDNSUsed(ctx context.Context, zoneID string) (bool, error) {
	records, err := c.GetDNSRecords(ctx, zoneID)
	if err != nil {
		return false, err
	}

	if len(records) == 0 {
		return false, nil
	}

	if len(records) != 2 {
		return true, nil
	}

	f, s := records[0].GetType(), records[1].GetType()
	if (f == "NS" && s == "SOA") || (f == "SOA" && s == "NS") {
		return false, nil
	}

	return true, nil
}
