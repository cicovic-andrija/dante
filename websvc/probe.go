package websvc

import (
	"net/http"
	"sync"

	"github.com/cicovic-andrija/dante/atlas"
)

type probeTable struct {
	sync.RWMutex

	probes map[int64]*atlas.Probe
}

func newProbeTable() *probeTable {
	return &probeTable{
		probes: make(map[int64]*atlas.Probe),
	}
}

func (t *probeTable) lookup(id int64) (probe *atlas.Probe, ok bool) {
	t.RLock()
	probe, ok = t.probes[id]
	t.RUnlock()
	return
}

func (t *probeTable) insert(probe *atlas.Probe) {
	t.Lock()
	t.probes[probe.ID] = probe
	t.Unlock()
}

func (s *server) getProbe(id int64) (*atlas.Probe, error) {
	var (
		probe *atlas.Probe
		req   *http.Request
		ok    bool
		err   error
	)

	if probe, ok = s.probeInfo.lookup(id); ok {
		return probe, nil
	}

	req, err = atlas.PrepareRequest(
		atlas.ProbeURL(id),
		&atlas.ReqParams{
			Method: http.MethodGet,
			Key:    cfg.Atlas.Auth.Key,
		},
	)
	if err != nil {
		return nil, err
	}

	probe = &atlas.Probe{}
	if err = s.makeRequest(req, probe); err != nil {
		return nil, err
	}

	// update probe cache
	s.probeInfo.insert(probe)
	s.log.info("[mgmt] probe info cached: id=%d country=%s asn=%d",
		probe.ID, probe.CountryCode, probe.ASNv4)

	return probe, nil
}
