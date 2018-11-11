package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/host"
	"log"
	"time"
)

func main() {

	for {
		time.Sleep(5 * time.Second)

		temps, err := host.SensorsTemperatures()
		if err != nil {
			log.Println("Error", err)
			continue
		}

		payload := make(map[string]interface{})
		payload["eventType"] = "Temperatures"
		for _, t := range temps {
			payload[t.SensorKey] = t.Temperature
		}
		bytes, err := json.Marshal(payload)
		if err != nil {
			log.Println("Error", err)
		}
		fmt.Println(string(bytes))
	}
}
