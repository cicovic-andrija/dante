package api

// ResourceListBase (Object List Resource)
type ResourceListBase struct {
	Count int    `json:"count"`
	Next  string `json:"next"`
	Prev  string `json:"previous"`
}

// ResourceBase (Object Detail Resource)
type ResourceBase struct {
	Id   int    `json:"id"`
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
