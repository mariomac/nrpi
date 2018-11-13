package main

import (
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
	client := api.New(cfg.AccountId, cfg.LicenseKey)

	server := transport.NewHttpServer(8080)
	httpCollector := metrics.NewHttpCollector(server)
	pipeline.Aggregate(client,
		metrics.StaticCollector(
			metrics.SystemHarvester(5*time.Second),
		),
		&httpCollector,
	)
}
