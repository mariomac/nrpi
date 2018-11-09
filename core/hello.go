package main

import (
	"github.com/mariomac/nrpi/core/api"
	"github.com/mariomac/nrpi/core/measure"
	"github.com/mariomac/nrpi/core/measure/native"

	"github.com/mariomac/nrpi/core/config"
)

func main() {
	cfg, err := config.Load("test.yml")
	if err != nil {
		panic(err)
	}
	client := api.New(cfg.AccountId, cfg.LicenseKey)
	measure.Aggregate(client, &native.Collector{})
}
