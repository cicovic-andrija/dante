package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/domain"
)

//
const (
	OperationalDataBucket = "operational-data"

	HealthEndpointPath = "/health"
)

// HealthReport
type HealthReport struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// EnsureOrganization
func (c *Client) EnsureOrganization(name string) (org *domain.Organization, err error) {
	var (
		orgAPI = c.influxClient.OrganizationsAPI()
	)

	lookupCtx, cancelLookup := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancelLookup()
	org, err = orgAPI.FindOrganizationByName(lookupCtx, name)
	if err == nil {
		return
	}

	if isNotFoundErr(err) {
		createCtx, cancelCreate := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancelCreate()
		org, err = orgAPI.CreateOrganizationWithName(createCtx, name)
		if err == nil {
			return
		}
	}

	return nil, err
}

// EnsureBucket
func (c *Client) EnsureBucket(name string) (bck *domain.Bucket, err error) {
	var (
		bckAPI = c.influxClient.BucketsAPI()
	)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	bck, err = bckAPI.FindBucketByName(ctx, name)
	if err == nil {
		return
	}

	if isNotFoundErr(err) {
		bck, err = c.CreateBucket(name)
		if err == nil {
			return
		}
	}

	return nil, err
}

// CreateBucket
// NOTE: Assumes c.Org is non-nil.
func (c *Client) CreateBucket(name string) (bck *domain.Bucket, err error) {
	var (
		bckAPI = c.influxClient.BucketsAPI()
	)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	if bck, err = bckAPI.CreateBucketWithName(ctx, c.Org, name); err != nil {
		bck = nil
		return
	}

	return
}

// ServerURL
func (c *Client) ServerURL() string {
	return c.influxClient.ServerURL()
}

// HealthEndpoint
func (c *Client) HealthEndpoint() string {
	return c.ServerURL() + HealthEndpointPath
}

// DataExplorerURL
func (c *Client) DataExplorerURL(bucket string) string {
	return fmt.Sprintf(
		"%s/orgs/%s/data-explorer?bucket=%s",
		c.ServerURL(),
		*c.Org.Id,
		bucket,
	)
}

func isNotFoundErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}
