package main

import (
	"sync"
	"time"

	"github.com/mariomac/nrpi/core/pipeline"

	"github.com/mariomac/nrpi/core/api"
	"github.com/mariomac/nrpi/core/config"
	"github.com/mariomac/nrpi/core/metrics"
	"github.com/mariomac/nrpi/core/transport"
)

func main() {

	cfg, err := config.Load("test.yml")
	if err != nil {
		panic(err)
	}
	server := transport.NewHTTPServer(8080)
	httpCollector := metrics.NewHTTPCollector(server)
	staticCollector := metrics.StaticCollector(
		metrics.SystemHarvester(5 * time.Second), // todo: make configurable
	)

	pipeline.Agent(
		time.Tick(5*time.Second), // todo: make configurable
		[]metrics.Collector{staticCollector, httpCollector},
		api.New(cfg.AccountID, cfg.LicenseKey),
	)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
