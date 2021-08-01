package atlas

import "strings"

// Endpoint-related constants
const (
	URLBase              = "https://atlas.ripe.net:443/api/v2"
	CreditsEndpoint      = URLBase + "/credits"
	MeasurementsEndpoint = URLBase + "/measurements"
)

// HTTP headers
const (
	AuthorizationHeader = "Authorization"
	AuthorizationFmt    = "Key %s"
	ContentTypeHeader   = "Content-Type"
	ContentType         = "application/json"
)

// Measurement types
const (
	MeasHTTP       = "http"
	MeasPing       = "ping"
	MeasTraceroute = "traceroute"
	MeasDNS        = "dns"
	MeasSSL        = "sslcert"
	MeasNTP        = "ntp"
	MeasWiFi       = "wifi"
)

// Address families
const (
	IPv4 = 4
	IPv6 = 6
)

var (
	// Probe selection types
	ValidProbeRequestTypes    = []string{"area", "country", "asn"}
	ValidProbeRequestTypesStr = strings.Join(ValidProbeRequestTypes, ",")
)
