package atlas

// ResourceBase (Object Detail Resource in RIPE Atlas Documentation)
type ResourceBase struct {
	Id   int64  `json:"id"`
	Type string `json:"type"`
}

// Error
type Error struct {
	Detail string `json:"detail"`
	Title  string `json:"title"`
	Status int64  `json:"status"`
}

// MeasurementDefinition
type MeasurementDefinition struct {
	Type          string `json:"type"`
	AddressFamily int32  `json:"af"`
	Target        string `json:"target"`
	Description   string `json:"description"`
	//StartTime     int64  `json:"start_time"`
	//StopTime      int64  `json:"stop_time"`
	//Interval      int64  `json:"interval"`
	IsPublic bool `json:"is_public"`
	IsOneOff bool `json:"is_oneoff"`
}

// ProbeRequest
type ProbeRequest struct {
	Requested int64  `json:"requested"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

// MeasurementRequest
type MeasurementRequest struct {
	Definitions []*MeasurementDefinition `json:"definitions"`
	Probes      []*ProbeRequest          `json:"probes"`
}

// MeasurementReqResponse
type MeasurementReqResponse struct {
	Measurements []int64 `json:"measurements"`
}

// Measurement
type Measurement struct {
	ResourceBase
	AddressFamily int32  `json:"af"`
	Description   string `json:"description"`
	IsPublic      bool   `json:"is_public"`
	Error         Error  `json:"error"`
}

// MeasurementResults
type MeasurementResults []SingleMeasurementResult

// SingleMeasurementResult
type SingleMeasurementResult struct {
	FirmwareVersion int32    `json:"fw"`
	Timestamp       int64    `json:"timestamp"`
	Results         []Result `json:"result"`
}

// Result
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

// Credit
type Credit struct {
	Error          Error `json:"error"`
	CurrentBalance int64 `json:"current_balance"`
}
