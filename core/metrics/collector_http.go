package metrics

import (
	"bufio"
	"encoding/json"
	"github.com/mariomac/nrpi/core/transport"
	"io"
	"log"
)

type HttpCollector struct {
	receiver transport.Receiver
}

func NewHttpCollector(server transport.HttpServer) HttpCollector {
	return HttpCollector{
		receiver:server.Endpoint("/http"),
	}
}

func (h *HttpCollector) Receive(batcher chan<- Harvest) {
	// todo: validate format (e.g. content-type)
	// todo: integrate authentication/security
	go func() {
		reader := h.receiver()
		br := bufio.NewReader(reader)
		jsonbytes, err := br.ReadBytes(byte('\n'))
		for  err == nil  {
			payload := make(Harvest)
			json.Unmarshal(jsonbytes, &payload)
			log.Printf("HTTP Sending %#v", payload)
			batcher <- payload
			jsonbytes, err = br.ReadBytes(byte('\n'))
		}
		if err != io.EOF {
			log.Fatal(err)
		}
	}()
}

