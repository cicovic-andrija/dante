package atlas

import (
	"fmt"
	"strings"
)

// Atlas API endpoint-related constants.
const (
	URLBase                       = "https://atlas.ripe.net:443/api/v2"
	CreditsEndpoint               = URLBase + "/credits"
	MeasurementsEndpoint          = URLBase + "/measurements"
	MeasurementEndpointFmt        = MeasurementsEndpoint + "/%d"
	MeasurementResultsEndpointFmt = MeasurementsEndpoint + "/%d" + MeasurementResultsPathSuffix
	MeasurementResultsPathSuffix  = "/results"
)

// HTTP header constants.
const (
	AuthorizationHeader = "Authorization"
	AuthorizationFmt    = "Key %s"
	ContentTypeHeader   = "Content-Type"
	ContentType         = "application/json"
)

// Measurement type constants.
const (
	MeasHTTP       = "http"
	MeasPing       = "ping"
	MeasTraceroute = "traceroute"
	MeasDNS        = "dns"
	MeasSSL        = "sslcert"
	MeasNTP        = "ntp"
	MeasWiFi       = "wifi"
)

// Address family constants.
const (
	IPv4 = 4
	IPv6 = 6
)

// Probe selection types.
var (
	ValidProbeRequestTypes    = []string{"area", "country", "asn"}
	ValidProbeRequestTypesStr = strings.Join(ValidProbeRequestTypes, ",")
)

// MeasurementResultsURL returns an endpoint for fetching measurement results.
func MeasurementResultsURL(measurementId int64) string {
	return fmt.Sprintf(MeasurementResultsEndpointFmt, measurementId)
}

// MeasurementURL returns an endpoint for working with a measurement resource.
func MeasurementURL(measurementId int64) string {
	return fmt.Sprintf(MeasurementEndpointFmt, measurementId)
}
