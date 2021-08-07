package atlas

// ResourceBase (Object Detail Resource in RIPE Atlas Documentation).
// Common fields in all Atlas API resources.
type ResourceBase struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// Error represents a JSON object returned in an HTTP response
// to an invalid Atlas API request.
type Error struct {
	Detail string `json:"detail"`
	Title  string `json:"title"`
	Status int64  `json:"status"`
}

// Probe represents a probe resource on the Atlas platform.
type Probe struct {
	ResourceBase
	CountryCode string `json:"country_code"`
	ASNv4       int64  `json:"asn_v4"`
}

// MeasurementDefinition
// TODO
type MeasurementDefinition struct {
	Type          string `json:"type"`
	AddressFamily int32  `json:"af"`
	Target        string `json:"target"`
	Description   string `json:"description"`
	IsPublic      bool   `json:"is_public"`
	IsOneOff      bool   `json:"is_oneoff"`
	StartTime     *int64 `json:"start_time"`
	StopTime      *int64 `json:"stop_time"`
	Interval      *int64 `json:"interval"`
}

// ProbeRequest specifies how the probes are to be selected
// for a new measurement, in a measurement request.
type ProbeRequest struct {
	Requested int64  `json:"requested"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

// MeasurementRequest contains measurement definitions and
// probe requests needed to create a measurement specification
// on the Atlas platform.
type MeasurementRequest struct {
	Definitions []*MeasurementDefinition `json:"definitions"`
	Probes      []*ProbeRequest          `json:"probes"`
}

// MeasurementReqResponse contains IDs of created measurements,
// returned as a response to a measurement request
// from the Atlas API.
type MeasurementReqResponse struct {
	Measurements []int64 `json:"measurements"`

	Error *Error `json:"error"`
}

// Measurement represents a measurement resource on the Atlas platform.
type Measurement struct {
	ResourceBase
	AddressFamily int32             `json:"af"`
	Description   string            `json:"description"`
	Target        string            `json:"target"`
	TargetIP      string            `json:"target_ip"`
	Status        MeasurementStatus `json:"status"`

	Error *Error `json:"error"`
}

// MeasurementStatus contains status values of a measurement resource.
type MeasurementStatus struct {
	ID    int32  `json:"id"`
	Value string `json:"name"`
}

// MeasurementResults contains an array of single measurement results
// performed by a single probe.
type MeasurementResults []ProbeMeasurementResults

// ProbeMeasurementResults contains an array of measurements performed
// by a single probe, for a single measurement.
type ProbeMeasurementResults struct {
	FirmwareVersion int32    `json:"fw"`
	Timestamp       int64    `json:"timestamp"`
	ProbeID         int64    `json:"prb_id"`
	Results         []Result `json:"result"`
}

// Result represents a result of a performed HTTP measurement,
// done by a probe.
type Result struct {
	// All firmware versions
	BodySize   int64   `json:"bsize"`
	HeaderSize int64   `json:"hsize"`
	Result     int32   `json:"res"`
	RT         float64 `json:"rt"`

	// Firmware version 4400
	Src  string `json:"srcaddr"`
	Addr string `json:"addr"`
	Mode string `json:"mode"`

	// Firmware version 4460,4540,4570,4610,4750,5000
	SrcAddr string `json:"src_addr"`
	DstAddr string `json:"dst_addr"`
	Method  string `json:"method"`
}

// Credit represents a credit report object,
// fetched from the Atlas API.
type Credit struct {
	CurrentBalance int64 `json:"current_balance"`

	Error *Error `json:"error"`
}
