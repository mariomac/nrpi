package metrics

import "fmt"

const fieldEventType = "eventType"

var requiredFields = []string{fieldEventType}

// Harvest represents a key-value map of a collected set of metrics. To ve valid, it must have an
// "eventType" field.
type Harvest map[string]interface{}

// EventType sets the event type of the harvest
func (p Harvest) EventType(eventType string) {
	p[fieldEventType] = eventType
}

// Validate checks that all the required fields are set with valid values
func (p Harvest) Validate() error {
	for _, f := range requiredFields {
		if found, ok := p[f]; !ok || found == "" {
			return fmt.Errorf("missing field %q", f)
		}
	}
	return nil
}

// todo: fromJSON

// Collector implementors harvest metrics from a given source and forward them by the channel
// specified in the Forward function
type Collector interface {
	// Forward makes the collector start collecting Harvests. When the Collector receives a Harvest,
	// it forwards it from the channel passed as argument
	Forward(chan<- Harvest)
}
