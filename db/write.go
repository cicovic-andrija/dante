package db

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	CreditBalanceMeasurement = "credit-balance"
	HTTPMeasurement          = "http"
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

//
func (c *Client) WriteSingleMeasurementResult(bucketName string, backendId int64, rt float64, ts time.Time) {
	var (
		writeAPI = c.influxClient.WriteAPIBlocking(c.Org.Name, bucketName)
	)

	dataPoint := influxdb2.NewPoint(
		HTTPMeasurement,
		nil,
		map[string]interface{}{
			"backend-measurement-id": backendId,
			"rt":                     rt,
		},
		ts,
	)

	context, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	writeAPI.WritePoint(context, dataPoint)
}
