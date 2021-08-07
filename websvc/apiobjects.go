package websvc

import "github.com/influxdata/influxdb-client-go/v2/domain"

type control struct {
	Operation string `json:"operation"`
}

type status struct {
	Status string `json:"status"`
	ID     string `json:"id,omitempty"`
}

type creditResp struct {
	CurrentBalance int64  `json:"current_balance"`
	URL            string `json:"url"`
}

type measurementReq struct {
	// mandatory fields
	Targets       []string   `json:"targets"`
	ProbeRequests []probeReq `json:"probe_requests"`
	Description   string     `json:"description"`

	// optional fields
	//StartTimeUTC *string `json:"start_time_utc"`
	//EndTimeUTC   *string `json:"end_time_utc"`
}

type probeReq struct {
	Requested int64  `json:"requested"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type measurement struct {
	ID                  string                `json:"id"`
	Status              string                `json:"status"`
	BucketName          string                `json:"bucket_name,omitempty"`
	Description         string                `json:"description,omitempty"`
	BackendMeasurements []*backendMeasurement `json:"backend_measurements,omitempty"`
	Reason              string                `json:"reason,omitempty"`
	URL                 string                `json:"url,omitempty"`

	bucket *domain.Bucket `json:"-"`
}

type backendMeasurement struct {
	ID       int64  `json:"backend_id"`
	Target   string `json:"target"`
	TargetIP string `json:"target_ip"`

	status string `json:"-"`
}
