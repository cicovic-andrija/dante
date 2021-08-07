package websvc

import (
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
	t.probes[probe.Id] = probe
	t.Unlock()
}
