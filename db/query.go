package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

var (
	mdQuery = fmt.Sprintf(
		`from(bucket:"%s")|>range(start:0,stop:1)|>filter(fn:(r)=>r["_measurement"]=="%s")`,
		SystemBucket,
		MetadataMeasurement,
	)

	errCorrupted = errors.New("measurement metadata corrupted")
)

// QueryMeasurementMetadata reads measurement metadata from the SystemBucket.
// It assumes c.Org is not nil.
func (c *Client) QueryMeasurementMetadata() ([]MeasurementMetadata, error) {
	var (
		queryAPI = c.influxClient.QueryAPI(c.Org.Name)
		result   *api.QueryTableResult
		ok       bool
		err      error
	)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	if result, err = queryAPI.Query(ctx, mdQuery); err != nil {
		return nil, err
	}

	md := []MeasurementMetadata{}
	for result.Next() {
		mdPart := MeasurementMetadata{}
		if mdPart.ID, ok = result.Record().ValueByKey(tagID).(string); !ok {
			return nil, errCorrupted
		}
		if mdPart.Description, ok = result.Record().ValueByKey(tagDescription).(string); !ok {
			return nil, errCorrupted
		}
		if mdPart.BackendIDsStr, ok = result.Record().ValueByKey(tagBackendIDs).(string); !ok {
			return nil, errCorrupted
		}
		md = append(md, mdPart)
	}
	if result.Err() != nil {
		return nil, err
	}

	return md, nil
}
