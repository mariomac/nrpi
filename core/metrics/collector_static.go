package metrics

type StaticCollector struct{}

func (*StaticCollector) Receive(ch chan<- Harvest) {
	go func() {
		sh := systemHarvester{}
		sh.Collect(ch)
	}()
}
