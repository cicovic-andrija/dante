package api

// ResourceListBase (Object List Resource)
type ResourceListBase struct {
	Count int64  `json:"count"`
	Next  string `json:"next"`
	Prev  string `json:"previous"`
}

// ResourceBase (Object Detail Resource)
type ResourceBase struct {
	Id   int64  `json:"id"`
	Type string `json:"type"`
}

// Probe
type ProbeResource struct {
	ResourceBase
}

// ProbeResourceList
type ProbeResourceList struct {
	ResourceListBase
	Probes []ProbeResource `json:"results"`
}

// Error
type Error struct {
	Detail string `json:"detail"`
	Title  string `json:"title"`
	Status int64  `json:"status"`
}

// Credit
type Credit struct {
	Error          Error `json:"error"`
	CurrentBalance int64 `json:"current_balance"`
}
