package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Documentation: https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/send-custom-events-event-api

const (
	url       = "https://insights-collector.newrelic.com/v1/accounts/%s/events"
	keyHeader = "X-Insert-Key"
)

type NewRelic interface {
	SendEvent(event interface{}) error
}

type newRelicImpl struct {
	eventsUrl  string
	licenseKey string
}

func New(accountId, licenseKey string) NewRelic {
	return &newRelicImpl{
		eventsUrl:  fmt.Sprintf(url, accountId),
		licenseKey: licenseKey,
	}
}

func (n newRelicImpl) SendEvent(event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	var gzippedBody bytes.Buffer
	// Todo: configure gzip level
	w := gzip.NewWriter(&gzippedBody)
	w.Write(body)
	w.Close()

	// TODO: persistent connection
	req, err := http.NewRequest(http.MethodPost, n.eventsUrl, bytes.NewReader(gzippedBody.Bytes()))
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
