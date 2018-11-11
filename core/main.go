package main

import (
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
	metrics.Aggregate(client,
		&metrics.StaticCollector{},
		&httpCollector,
	)
}
