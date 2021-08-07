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

func (c *measurementCache) insert(meas *measurement) {
	c.Lock()
	c.measurements[meas.Id] = meas
	c.Unlock()
}

func (c *measurementCache) get(id string) (meas *measurement, ok bool) {
	c.RLock()
	meas, ok = c.measurements[id]
	c.RUnlock()
	return
}

func (c *measurementCache) del(id string) {
	c.Lock()
	delete(c.measurements, id)
	c.Unlock()
}

func (c *measurementCache) getAll() []*measurement {
	c.RLock()
	defer c.RUnlock()
	measurements := make([]*measurement, 0, len(c.measurements))
	for _, meas := range c.measurements {
		measurements = append(measurements, meas)
	}
	return measurements
}
