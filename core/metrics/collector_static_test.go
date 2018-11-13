package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type h1 struct{}

func (*h1) Start(ch chan<- Harvest) {
	h := Harvest{}
	h.EventType("H1Sample")
	h["val"] = "hello"
	ch <- h
}

type h2 struct{}

func (*h2) Start(ch chan<- Harvest) {
	h := Harvest{}
	h.EventType("H2Sample")
	h["val"] = 0xb1e
	ch <- h
}

func TestStaticCollector(t *testing.T) {
	const timeout = 1 * time.Second

	// Given a set of static harvesters
	// And a static collector
	coll := StaticCollector(&h1{}, &h2{})

	ch := make(chan Harvest)

	// When the collector starts receiving the harvests
	coll.Receive(ch)

	// both harvests are received
	harvests := make([]Harvest, 0, 2)
	for i := 0; i < 2; i++ {
		select {
		case h := <-ch:
			harvests = append(harvests, h)
		case <-time.After(timeout):
			require.Fail(t, "timeout while waiting for a harvest")
		}
	}
	require.ElementsMatch(t, harvests,
		[]Harvest{
			{
				"eventType": "H1Sample",
				"val":       "hello",
			},
			{
				"eventType": "H2Sample",
				"val":       0xb1e,
			},
		})
}
