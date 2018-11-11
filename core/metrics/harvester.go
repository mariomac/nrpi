package metrics

type Harvester interface {
	Start(chan<- Harvest)
}


