package metrics

import (
	"bufio"
	"encoding/json"
	"io"
	"log"

	"github.com/mariomac/nrpi/core/transport"
)

// httpCollector implements a Collector that receives the Metrics as JSONs sent through an http
// connection.
type httpCollector struct {
	receiver transport.Receiver
}

// NewHTTPCollector attaches an httpCollector to the "/http" path of the passed server, and returns
// it.
func NewHTTPCollector(server transport.HTTPServer) Collector {
	return &httpCollector{
		receiver: server.Endpoint("/http"),
	}
}

// Forward forwards to the channel the received metrics
func (h *httpCollector) Forward(ch chan<- Harvest) {
	// todo: validate format (e.g. content-type)
	// todo: integrate authentication/security
	go func() {
		reader := h.receiver()
		br := bufio.NewReader(reader)
		jsonbytes, err := br.ReadBytes(byte('\n'))
		for err == nil {
			payload := make(Harvest)
			json.Unmarshal(jsonbytes, &payload)
			ch <- payload
			jsonbytes, err = br.ReadBytes(byte('\n'))
		}
		if err != io.EOF {
			log.Fatal(err)
		}
	}()
}
