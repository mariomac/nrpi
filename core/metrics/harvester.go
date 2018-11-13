package metrics

// Harvester defines the interface for any item that collects (harvests) metrics, and submits them
// through the channel passed in the Start function
type Harvester interface {
	// Start makes the harvester collecting new metrics. Each individual harvested metric will be
	// submitted through the channel that is passed by argument.
	Start(chan<- Harvest)
}
