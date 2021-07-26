package db

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	CreditBalanceMeasurement = "credit-balance"
)

// WriteCreditBalance
// NOTE: Assumes c.Org is non-nil.
func (c *Client) WriteCreditBalance(creditBalance int64) error {
	var (
		writeAPI = c.influxClient.WriteAPIBlocking(c.Org.Name, OperationalDataBucket)
	)

	dataPoint := influxdb2.NewPoint(
		CreditBalanceMeasurement,
		nil, /* tags */
		map[string]interface{}{"value": creditBalance},
		time.Now(),
	)

	context, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return writeAPI.WritePoint(context, dataPoint)
}
