package native

import (
	"log"
	"time"

	"github.com/mariomac/nrpi/core/measure"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type Collector struct{}

func (*Collector) Receive(ch chan<- measure.Harvest) {
	go func() {
		sh := systemHarvester{}
		sh.Collect(ch)
	}()
}

type Harvester interface {
	Collect(chan<- measure.Harvest)
}

type systemHarvester struct{}

const interval = 5 * time.Second // todo: make configurable

func (*systemHarvester) Collect(ch chan<- measure.Harvest) { // todo: test
	for { // TODO: stop when channel is closed
		h := measure.Harvest{}
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
		time.Sleep(interval)
	}
}
