package pipeline

import (
	"log"
	"time"

	"github.com/mariomac/nrpi/core/metrics"

	"github.com/mariomac/nrpi/core/api"
)

// Aggregate mixes and buffers the metrics received from the collectors and forwards them
// periodically to the New Relic api
func Aggregate(nr api.NewRelic, collectors ...metrics.Collector) {

	// Receiver loop
	receiver := make(chan metrics.Harvest)
	for _, coll := range collectors {
		coll.Forward(receiver)
	}
	batch := make([]metrics.Harvest, 0)
	submitTicker := time.Tick(5 * time.Second) // todo: make configurable
	for {
		select {
		case rh := <-receiver:
			if err := rh.Validate(); err == nil { // todo: test
				batch = append(batch, rh)
			} else {
				log.Printf("Error receiving harvest %v: %s", rh, err.Error())
			}
		case <-submitTicker:
			nr.SendEvent(batch)
			batch = batch[:0]
		}
	}
}
