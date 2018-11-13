package pipeline

import (
	"encoding/json"
	"log"
	"time"

	"github.com/mariomac/nrpi/core/metrics"

	"github.com/mariomac/nrpi/core/api"
)

// todo: pass contexts to all the pipeline steps

// Join mixes metrics received from the collectors and forwards them next pipeline step
func Join(collectors []metrics.Collector, out chan<- metrics.Harvest) {
	go func() {
		// Receiver loop
		receiver := make(chan metrics.Harvest)
		for _, coll := range collectors {
			coll.Forward(receiver)
		}
		for {
			v := <-receiver
			log.Println("Joiner received", v)
			out <- v
		}
	}()
}

// clock is a signal
func Buffer(in <-chan metrics.Harvest, out chan<- []metrics.Harvest, clock <-chan time.Time) {
	go func() {
		batch := make([]metrics.Harvest, 0)
		for {
			select {
			case rh := <-in:
				log.Println("Buffer received", rh)
				if err := rh.Validate(); err == nil { // todo: test
					batch = append(batch, rh)
					log.Println("Buffering", batch)
				} else {
					log.Printf("Error receiving harvest %v: %s", rh, err.Error())
				}
			case <-clock:
				log.Println("Buffer sending", batch)
				out <- batch
				batch = make([]metrics.Harvest, 0)
			}
		}
	}()
}

func Marshaller(in <-chan []metrics.Harvest, out chan<- []byte) {
	go func() {
		for {
			payload := <- in
			log.Println("Marshaler received", payload)
			js, err := json.Marshal(payload)
			if err != nil {
				log.Printf("Error marshaling harvest: %s", err)
			}
			out <- js
		}
	}()
}

// Submit sends to the new relic api the received information
func Submit(in <-chan []byte, nr api.NewRelic) {
	go func() {
		for {
			nr.SendEvent(<-in)
		}
	}()
}
