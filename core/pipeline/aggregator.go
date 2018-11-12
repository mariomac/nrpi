package pipeline

import (
	"github.com/mariomac/nrpi/core/metrics"
	"log"
	"time"

	"github.com/mariomac/nrpi/core/api"
)

func Aggregate(nr api.NewRelic, collectors ...metrics.Collector) {

	// Receiver loop
	receiver := make(chan metrics.Harvest)
	for _, coll := range collectors {
		coll.Receive(receiver)
	}
	batch := make([]metrics.Harvest, 0)
	submitTicker := time.Tick(5 * time.Second) // todo: make configurable
	for {
		select {
		case rh := <-receiver:
			if err := rh.Validate() ; err == nil { // todo: test
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
