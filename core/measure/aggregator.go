package measure

import (
	"log"
	"time"

	"github.com/mariomac/nrpi/core/api"
)

func Aggregate(nr api.NewRelic, collectors ...Collector) {

	// Receiver loop
	receiver := make(chan Harvest)
	for _, coll := range collectors {
		coll.Receive(receiver)
	}
	batch := make([]Harvest, 0)
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
