package db

import (
	"context"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/domain"
)

const (
	MeasurementsBucket    = "measurements"
	OperationalDataBucket = "operational-data"
)

type HealthReport struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// EnsureOrganization
// NOTE: Assumes c.Org is nil.
func (c *Client) EnsureOrganization(name string) error {
	var (
		orgAPI = c.influxClient.OrganizationsAPI()
	)

	lookupCtx, cancelLookup := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancelLookup()
	org, err := orgAPI.FindOrganizationByName(lookupCtx, name)
	if err == nil {
		c.Org = org
		return nil
	}

	if isNotFound(err) {
		createCtx, cancelCreate := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancelCreate()
		org, err = orgAPI.CreateOrganizationWithName(createCtx, name)
		if err == nil {
			c.Org = org
			return nil
		}
	}

	return err
}

// EnsureBuckets
func (c *Client) EnsureBuckets() error {
	var (
		bck *domain.Bucket
		err error
	)

	if bck, err = c.EnsureBucket(MeasurementsBucket); err == nil {
		c.MeasBucket = bck
	} else {
		return err
	}

	if bck, err = c.EnsureBucket(OperationalDataBucket); err == nil {
		c.OperDataBucket = bck
	} else {
		return err
	}

	return nil
}

// EnsureBucket
// NOTE: Assumes c.Org is non-nil.
func (c *Client) EnsureBucket(name string) (bck *domain.Bucket, err error) {
	var (
		bckAPI = c.influxClient.BucketsAPI()
	)

	lookupCtx, cancelLookup := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancelLookup()
	bck, err = bckAPI.FindBucketByName(lookupCtx, name)
	if err == nil {
		return
	}

	if isNotFound(err) {
		createCtx, cancelCreate := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancelCreate()
		bck, err = bckAPI.CreateBucketWithName(createCtx, c.Org, name)
		if err == nil {
			return
		}
	}

	return nil, err
}

func (c *Client) ServerURL() string {
	return c.influxClient.ServerURL()
}

func (c *Client) HealthEndpoint() string {
	return c.ServerURL() + "/health"
}

func isNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}
