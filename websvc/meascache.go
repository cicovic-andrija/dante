package websvc

import "sync"

type measurementCache struct {
	sync.RWMutex

	measurements map[string]*measurement
}

func newMeasurementCache() *measurementCache {
	return &measurementCache{
		measurements: make(map[string]*measurement),
	}
}

func (m *measurementCache) insert(meas *measurement) {
	m.Lock()
	m.measurements[meas.Id] = meas
	m.Unlock()
}

func (m *measurementCache) get(id string) (meas *measurement, ok bool) {
	m.RLock()
	meas, ok = m.measurements[id]
	m.RUnlock()
	return
}

func (m *measurementCache) del(id string) {
	m.Lock()
	delete(m.measurements, id)
	m.Unlock()
}

func (m *measurementCache) getAll() []*measurement {
	m.RLock()
	defer m.RUnlock()
	measurements := make([]*measurement, 0, len(m.measurements))
	for _, meas := range m.measurements {
		measurements = append(measurements, meas)
	}
	return measurements
}
