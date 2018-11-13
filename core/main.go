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
		metrics.SystemHarvester(5 * time.Second),
	)

	toBuffer := make(chan metrics.Harvest)
	pipeline.Join([]metrics.Collector{staticCollector, httpCollector}, toBuffer)


	toMarshal := make(chan []metrics.Harvest)
	pipeline.Buffer(toBuffer, toMarshal, time.Tick(5*time.Second))

	toSubmit := make(chan []byte)
	pipeline.Marshaller(toMarshal, toSubmit)

	pipeline.Submit(toSubmit, api.New(cfg.AccountID, cfg.LicenseKey))

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
