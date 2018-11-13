package metrics

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type systemHarvester struct {
	interval time.Duration
}

// SystemHarvester returns a harvester that collects system-level metrics (cpu, memory consumption...)
// with the interval passed as parameter
func SystemHarvester(interval time.Duration) Harvester {
	return &systemHarvester{
		interval: interval,
	}
}

func (sh *systemHarvester) Start(ch chan<- Harvest) { // todo: test
	log.Println("starting harvester")
	for { // TODO: stop when channel is closed
		h := Harvest{}
		h.EventType("SystemHarvest")
		vm, err := mem.VirtualMemory()
		if err == nil {
			h["memoryTotal"] = vm.Total
			h["memoryUsed"] = vm.Used
		} else {
			log.Println("harvesting memory: ", err.Error())
		}
		c, err := cpu.Percent(0, false)
		if err == nil {
			h["cpuPercent"] = c[0]
		} else {
			log.Println("Harvesting cpu: ", err.Error())
		}
		log.Println("submitting ")
		ch <- h
		log.Println("submitted")

		time.Sleep(sh.interval)
	}
}
