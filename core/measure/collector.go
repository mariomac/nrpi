package measure

import "fmt"

const fieldEventType = "eventType"

var requiredFields = []string{fieldEventType}

type Harvest map[string]interface{}

func (p Harvest) EventType(eventType string) {
	p[fieldEventType] = eventType
}

func (p Harvest) Validate() error {
	for _, f := range requiredFields {
		if found, ok := p[f]; !ok || found == "" {
			return fmt.Errorf("missing field %q", f)
		}
	}
	return nil
}

type Collector interface {
	Receive(chan<- Harvest)
}

