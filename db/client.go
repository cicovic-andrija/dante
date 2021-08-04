package db

import (
	"time"

	"github.com/cicovic-andrija/dante/conf"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

// Database client-related constants.
const (
	DefaultTimeout = 10 * time.Second
)

// Client represents an interface to the InfluxDB database.
// It wraps the library-provided database client.
type Client struct {
	Org          *domain.Organization
	MeasBucket   *domain.Bucket
	SystemBucket *domain.Bucket

	influxClient influxdb2.Client
}

// NewClient returns a new Client.
func NewClient(cfg *conf.InfluxDBConf) *Client {
	return &Client{
		influxClient: influxdb2.NewClient(cfg.Net.GetURLBase(), cfg.Auth.Token),
	}
}

// Close closes the underlying library-provided database client.
func (c *Client) Close() {
	c.influxClient.Close()
}
