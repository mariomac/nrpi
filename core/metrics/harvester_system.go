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

func SystemHarvester(interval time.Duration) Harvester {
	return &systemHarvester{
		interval: interval,
	}
}

func (sh *systemHarvester) Start(ch chan<- Harvest) { // todo: test
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
		ch <- h
		time.Sleep(sh.interval)
	}
}
