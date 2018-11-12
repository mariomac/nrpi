package metrics

import "time"

// StaticCollector collects metrics from Harvesters that are statically embedded into the Agent:
// SystemHarvester
type StaticCollector struct{}

func NewStaticCollector(
	systemHarvestPeriod time.Duration,
	) {

}
// Receive forwards by the channel the payloads from the Static Harvesters
func (*StaticCollector) Receive(ch chan<- Harvest) {
	go func() {
		sh := systemHarvester{}
		sh.Collect(ch)
	}()
}
