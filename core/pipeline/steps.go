package pipeline

import (
	"encoding/json"
	"log"
	"time"

	"github.com/mariomac/nrpi/core/metrics"

	"github.com/mariomac/nrpi/core/api"
)

func Agent(bufferClock <-chan time.Time, collectors []metrics.Collector, nr api.NewRelic) {
	toBuffer := make(chan metrics.Harvest)
	join(collectors, toBuffer)

	toMarshal := make(chan []metrics.Harvest)
	buffer(toBuffer, toMarshal, bufferClock)

	toSubmit := make(chan []byte)
	marshal(toMarshal, toSubmit)

	submit(toSubmit, nr)
}

// todo: pass contexts to all the pipeline steps

// joins metrics received from the collectors and forwards them next pipeline step
func join(collectors []metrics.Collector, out chan<- metrics.Harvest) {
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
func buffer(in <-chan metrics.Harvest, out chan<- []metrics.Harvest, clock <-chan time.Time) {
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

func marshal(in <-chan []metrics.Harvest, out chan<- []byte) {
	go func() {
		for {
			payload := <-in
			log.Println("Marshaler received", payload)
			js, err := json.Marshal(payload)
			if err != nil {
				log.Printf("Error marshaling harvest: %s", err)
			}
			out <- js
		}
	}()
}

// submit sends to the new relic api the received information
func submit(in <-chan []byte, nr api.NewRelic) {
	go func() {
		for {
			nr.SendEvent(<-in)
		}
	}()
}
