package pipeline

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mariomac/nrpi/core/metrics"
)

// Simple harvester that just forwards a sample
type forwarder struct {
	h  metrics.Harvest
	wg *sync.WaitGroup
}

func (f *forwarder) Start(ch chan<- metrics.Harvest) {
	ch <- f.h
	f.wg.Done()
}

type fakeApi struct {
	invocations []string
	submitted   *sync.WaitGroup
}

func (f *fakeApi) SendEvent(body []byte) error {
	f.invocations = append(f.invocations, string(body))
	f.submitted.Done()
	return nil
}

func TestAgent(t *testing.T) {
	apiWait := sync.WaitGroup{}
	apiWait.Add(1)
	hvstWait := sync.WaitGroup{}
	hvstWait.Add(2)

	s1 := metrics.Harvest{}
	s1.EventType("H1Sample")
	s1["val"] = "hello"

	s2 := metrics.Harvest{}
	s2.EventType("H2Sample")
	s2["val"] = 0xdeadbeef

	// Given a set of harvesters
	coll := []metrics.Harvester{&forwarder{s1, &hvstWait}, &forwarder{s2, &hvstWait}}

	// And an agent pipeline
	bufferClock := make(chan time.Time)
	nrApi := fakeApi{
		invocations: make([]string, 0),
		submitted:   &apiWait,
	}
	Agent(bufferClock, []metrics.Collector{metrics.StaticCollector(coll...)}, &nrApi)

	// Whose harvests forwarded any data
	hvstWait.Wait()

	// Synchronously the harvested data is submitted in a single request
	bufferClock <- time.Now()

	apiWait.Wait()
	require.Equal(t, 1, len(nrApi.invocations))

	s1b, _ := json.Marshal(s1)
	require.Contains(t, nrApi.invocations[0], string(s1b))
	s2b, _ := json.Marshal(s1)
	require.Contains(t, nrApi.invocations[0], string(s2b))
}
