package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Documentation: https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/send-custom-events-event-api

const (
	url       = "https://insights-collector.newrelic.com/v1/accounts/%s/events"
	keyHeader = "X-Insert-Key"
)

// NewRelic describes the interface to the New Relic API
type NewRelic interface {
	SendEvent(event interface{}) error
}

type newRelicImpl struct {
	eventsURL  string
	licenseKey string
}

// New returns a new NewRelic API implementation, given an account ID and a License Key
func New(accountID, licenseKey string) NewRelic {
	return &newRelicImpl{
		eventsURL:  fmt.Sprintf(url, accountID),
		licenseKey: licenseKey,
	}
}

func (n newRelicImpl) SendEvent(event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	log.Print("sending ", string(body))

	var gzippedBody bytes.Buffer
	// Todo: configure gzip level
	w := gzip.NewWriter(&gzippedBody)
	w.Write(body)
	w.Close()

	// TODO: persistent connection
	req, err := http.NewRequest(http.MethodPost, n.eventsURL, bytes.NewReader(gzippedBody.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set(keyHeader, n.licenseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	rbody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(rbody))

	return nil
}
