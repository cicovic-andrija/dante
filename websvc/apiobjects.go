package websvc

import "github.com/influxdata/influxdb-client-go/v2/domain"

type control struct {
	Operation string `json:"operation"`
}

type status struct {
	Status string `json:"status"`
}

type creditResp struct {
	CurrentBalance int64  `json:"current_balance"`
	URL            string `json:"url"`
}

type measurementReq struct {
	Targets       []string   `json:"targets"`
	ProbeRequests []probeReq `json:"probe_requests"`
	Description   string     `json:"description"`
}

type probeReq struct {
	Requested int64  `json:"requested"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type measurement struct {
	Id          string  `json:"id"`
	BucketName  string  `json:"bucket_name"`
	BackendIDs  []int64 `json:"backend_ids"`
	Description string  `json:"description"`
	Status      string  `json:"status"`

	URL string `json:"url"`

	bucket *domain.Bucket `json:"-"`
}
