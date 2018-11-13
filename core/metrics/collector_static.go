package metrics

// nativeCollector collects metrics from go language static code or dynamic libraries
type nativeCollector struct {
	harvesters []Harvester
}

// StaticCollector creates a staticCollector from a set of harvesters that are statically available
// at compile time
func StaticCollector(harvesters ...Harvester) Collector {
	return &nativeCollector{
		harvesters: harvesters,
	}
}

// Forward forwards by the channel the payloads from the Static Harvesters
func (c *nativeCollector) Forward(ch chan<- Harvest) {
	go func() {
		for _, h := range c.harvesters {
			// todo: avoid blockings
			h.Start(ch)
		}
	}()
}
