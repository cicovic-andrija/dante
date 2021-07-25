package db

import (
	"time"

	"github.com/cicovic-andrija/dante/conf"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	DefaultTimeout = 10 * time.Second
)

type Client struct {
	Organization string
	influxClient influxdb2.Client
}

func NewClient(cfg *conf.InfluxDBConf) *Client {
	return &Client{
		Organization: cfg.Organization,
		influxClient: influxdb2.NewClient(cfg.Net.GetBaseURL(), cfg.Auth.Token),
	}
}
