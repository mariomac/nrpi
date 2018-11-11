package metrics

import "testing"
import "github.com/stretchr/testify/assert"

func TestRaw_Validate(t *testing.T) {
	// Given a Harvest
	m := make(Harvest)

	// Which has a valid format
	m.EventType("myEventType")
	m["otherEvent"] = 123456

	// The data is validated
	assert.NoError(t, m.Validate())
}

func TestRaw_Validate_MissingEvent(t *testing.T) {
	// Given a Harvest
	m := make(Harvest)

	// Which has not all the required fields
	m["otherEvent"] = 123456

	// The data is validated with error
	assert.Error(t, m.Validate())
}

func TestRaw_Validate_NoEvent(t *testing.T) {
	// Given a Harvest
	m := make(Harvest)

	// Which has an invalid format
	m.EventType("")
	m["otherEvent"] = 123456

	// The data is validated with error
	assert.Error(t, m.Validate())
}
