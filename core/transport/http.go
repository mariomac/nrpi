package transport

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpServer struct {
	mux *http.ServeMux
}

func NewHttpServer(port int) HttpServer {
	server := HttpServer{
		mux: http.NewServeMux(),
	}
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), server.mux)
	}()
	return server
}

func (h *HttpServer) Endpoint(path string) Receiver {
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

