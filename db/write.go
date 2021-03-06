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
	tagBackendIDs  = "backend-ids"
	tagBackendID   = "backend-id"
	tagProbeID     = "probe-id"
	tagASN         = "asn"
	tagCountry     = "country"
	tagTarget      = "target"
	tagTargetIP    = "target-ip"

	fieldValue      = "value"
	fieldRT         = "rt"
	fieldBodySize   = "body-size"
	fieldHeaderSize = "header-size"
	fieldStatusCode = "status-code"
)

var nullTimestamp = time.Unix(0, 0)

// HTTPData specifies values of a data point
// that is written to the HTTPMeasurement bucket.
type HTTPData struct {
	BackendID     int64
	ProbeID       int64
	ASN           int64
	Country       string
	Target        string
	TargetIP      string
	RoundTripTime float64
	BodySize      int64
	HeaderSize    int64
	StatusCode    int32
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
			tagASN:       strconv.FormatInt(httpData.ASN, 10),
			tagCountry:   httpData.Country,
			tagTarget:    httpData.Target,
			tagTargetIP:  httpData.TargetIP,
		},
		map[string]interface{}{
			fieldRT:         httpData.RoundTripTime,
			fieldBodySize:   httpData.BodySize,
			fieldHeaderSize: httpData.HeaderSize,
			fieldStatusCode: httpData.StatusCode,
		},
		httpData.Timestamp,
	)

	return c.write(bucketName, dataPoint)
}

// WriteMeasurementMetadata writes a MetadataMeasurement data point to the SystemBucket.
// It assumes c.Org is not nil.
func (c *Client) WriteMeasurementMetadata(md MeasurementMetadata) error {
	dataPoint := influxdb2.NewPoint(
		MetadataMeasurement,
		map[string]string{
			tagID:          md.ID,
			tagDescription: md.Description,
			tagBackendIDs:  md.BackendIDsStr,
		},
		map[string]interface{}{
			"dummy-value": 42,
		},
		nullTimestamp,
	)

	return c.write(SystemBucket, dataPoint)
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
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return writeAPI.WritePoint(ctx, dataPoint)
}
