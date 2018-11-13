package transport

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// HTTPServer holds an HTTP Server that would be used to receive metrics
type HTTPServer struct {
	mux *http.ServeMux
}

// NewHTTPServer returns an HTTPServer listening in the given port.
func NewHTTPServer(port int) HTTPServer {
	server := HTTPServer{
		mux: http.NewServeMux(),
	}
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), server.mux)
	}()
	return server
}

// Endpoint attaches a new Receiver to the given path
func (h *HTTPServer) Endpoint(path string) Receiver {
	reader, pwriter := io.Pipe()
	h.mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
		writer.WriteHeader(http.StatusOK)
		n, err := pwriter.Write(body)
		if err != nil {
			log.Println("Error", err)
		}
		log.Printf("written %d bytes (%s)", n, body)
		pwriter.Write([]byte{'\n'})
	})
	return func() io.Reader {
		return reader
	}
}
