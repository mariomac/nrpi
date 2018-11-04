package main

import (
	"fmt"
	"github.com/mariomac/nrpi/core/api"

	"github.com/mariomac/nrpi/core/config"
)

type SystemSample struct {
	AgentName string `json:"agentName"`
	AgentVersion string `json:"agentVersion"`
	HostName string `json:"hostname"`
	EventType string `json:"eventType"`
}

func main() {
	cfg, err := config.Load("test.yml")
	if err != nil {
		panic(err)
	}
	client := api.New(cfg.AccountId, cfg.LicenseKey)
	err = client.SendEvent(SystemSample{
		AgentName:"IoT-agent",
		AgentVersion:"0.0.1",
		HostName:"my-test-host",
		EventType:"SystemSample",
	})
	if err != nil {
		fmt.Println(err)
	}
}
