package db

import (
	"time"

	"github.com/cicovic-andrija/dante/conf"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

const (
	DefaultTimeout = 10 * time.Second
)

type Client struct {
	Org            *domain.Organization
	MeasBucket     *domain.Bucket
	OperDataBucket *domain.Bucket
	influxClient   influxdb2.Client
}

func NewClient(cfg *conf.InfluxDBConf) *Client {
	return &Client{
		influxClient: influxdb2.NewClient(cfg.Net.GetURLBase(), cfg.Auth.Token),
	}
}

func (c *Client) Close() {
	c.influxClient.Close()
}
