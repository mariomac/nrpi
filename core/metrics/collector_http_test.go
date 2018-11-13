package metrics

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/mariomac/nrpi/core/transport"
	"github.com/stretchr/testify/require"
)

const timeout = 3 * time.Second

func TestHttpCollector_Forward(t *testing.T) {
	const port = 56789 // todo: find automatically an unused port
	// Given an HTTP collector
	coll := NewHTTPCollector(transport.NewHTTPServer(port))
	fwd := make(chan Harvest)
	coll.Forward(fwd)

	// When it receives a correct JSON payload
	resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d/http", port), "application/json",
		strings.NewReader(`{"eventType":"TestSample","value":1234}`))
	require.NoError(t, err)

	// It returns OK
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// And forwards the Harvest
	select {
	case hvst := <-fwd:
		require.Equal(t, "TestSample", hvst["eventType"])
		require.NoError(t, hvst.Validate())
		require.InDelta(t, 1234, hvst["value"], 0.01)
	case <-time.After(timeout):
		require.Fail(t, "operation timed out")
	}

}

// TODO: testValidate
// TODO: testContentType
