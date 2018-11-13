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
	client := api.New(cfg.AccountID, cfg.LicenseKey)

	server := transport.NewHTTPServer(8080)
	httpCollector := metrics.NewHTTPCollector(server)
	pipeline.Aggregate(client,
		metrics.StaticCollector(
			metrics.SystemHarvester(5*time.Second),
		),
		&httpCollector,
	)
}
