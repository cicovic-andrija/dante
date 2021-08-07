package db

import (
	"context"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// Measurement names.
const (
	MetadataMeasurement      = "md"
	HTTPMeasurement          = "http"
	CreditBalanceMeasurement = "credit-balance"
)

const (
	tagID          = "id"
	tagDescription = "description"
	tagStatus      = "status"
	tagBackendID   = "backend-id"
	tagProbeID     = "probe-id"

	fieldBackendID = "backend-id"
	fieldRT        = "rt"
	fieldValue     = "value"
)

var nullTimestamp = time.Unix(0, 0)

// HTTPData specifies values of a data point
// that is written to the HTTPMeasurement bucket.
type HTTPData struct {
	BackendID     int64
	ProbeID       int64
	RoundTripTime float64
	ASN           int64
	Country       string
	Timestamp     time.Time
}

// WriteMeasurementResult writes a single data point
// of the HTTPMeasurement measurement to a specified bucket.
// It assumes c.Org is not nil.
func (c *Client) WriteHTTPMeasurementResult(bucketName string, httpData *HTTPData) error {
	// specify data point
	dataPoint := influxdb2.NewPoint(
		HTTPMeasurement,
		map[string]string{
			tagBackendID: strconv.FormatInt(httpData.BackendID, 10),
			tagProbeID:   strconv.FormatInt(httpData.ProbeID, 10),
		},
		map[string]interface{}{
			fieldRT: httpData.RoundTripTime,
		},
		httpData.Timestamp,
	)

	return c.write(bucketName, dataPoint)
}

// WriteMeasurementMetadata writes one or more data points
// of the MetadataMeasurement measurement.
// It assumes c.Org is not nil.
func (c *Client) WriteMeasurementMetadata(measID string, description string, status string, backendIDs []int64) error {
	// write one data point per backend ID
	for _, id := range backendIDs {
		// specify data point
		dataPoint := influxdb2.NewPoint(
			MetadataMeasurement,
			map[string]string{
				tagID:          measID,
				tagDescription: description,
				tagStatus:      status,
			},
			map[string]interface{}{
				fieldBackendID: id,
			},
			time.Now(),
		)

		err := c.write(SystemBucket, dataPoint)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteCreditBalance writes a single data point
// of the CreditBalanceMeasurement measurement.
// It assumes c.Org is not nil.
func (c *Client) WriteCreditBalance(creditBalance int64) error {
	// specify data point
	dataPoint := influxdb2.NewPoint(
		CreditBalanceMeasurement,
		nil, /* tags */
		map[string]interface{}{fieldValue: creditBalance},
		time.Now(),
	)

	return c.write(SystemBucket, dataPoint)
}

func (c *Client) write(bucketName string, dataPoint *write.Point) error {
	writeAPI := c.influxClient.WriteAPIBlocking(c.Org.Name, bucketName)
	context, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return writeAPI.WritePoint(context, dataPoint)
}
