package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/domain"
)

// Constants related to management database entities.
const (
	SystemBucket = "system"

	HealthEndpointPath = "/health"
	HealthStatusPass   = "pass"
)

// MeasurementMetadata specifies measurement details
// that will be persisted in the database.
type MeasurementMetadata struct {
	ID            string
	Description   string
	BackendIDsStr string
}

// HealthReport represents a response object returned
// by the database API health endpoint.
type HealthReport struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// EnsureOrganization looks up an organization by name and returns it.
// If the organization does not exists, it is created.
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

// LookupBucket looks up a bucket by name and returns it.
func (c *Client) LookupBucket(name string) (bck *domain.Bucket, err error) {
	var (
		bckAPI = c.influxClient.BucketsAPI()
	)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	bck, err = bckAPI.FindBucketByName(ctx, name)
	return
}

// EnsureBucket looks up a bucket by name and returns it.
// If the bucket does not exists, it is created.
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

// CreateBucket creates a bucket.
// It assumes that c.Org is not nil.
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

// DeleteBucket deletes a bucket.
func (c *Client) DeleteBucket(bck *domain.Bucket) error {
	var (
		bckAPI = c.influxClient.BucketsAPI()
	)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return bckAPI.DeleteBucket(ctx, bck)
}

// ServerURL returns the database API base URL.
func (c *Client) ServerURL() string {
	return c.influxClient.ServerURL()
}

// HealthEndpoint the database API health endpoint.
func (c *Client) HealthEndpoint() string {
	return c.ServerURL() + HealthEndpointPath
}

// DataExplorerURL returns a formatted URL to InfluxDB Data Explorer,
// for a specified bucket.
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
