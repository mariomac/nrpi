package main

import (
	"bufio"
	"fmt"
	"github.com/mariomac/nrpi/core/iot"
	"io"
	"log"
)

func main() {
	/*
	cfg, err := config.Load("test.yml")
	if err != nil {
		panic(err)
	}
	client := api.New(cfg.AccountId, cfg.LicenseKey)
	measure.Aggregate(client, &native.Collector{})
*/


	server := iot.NewHttpServer(8080)
	connector := server.Connector("/hello")
	reader := connector()
	br := bufio.NewReader(reader)
	str, err := br.ReadString(byte('\n'))
	for  err == nil  {
		fmt.Println(str)
		str, err = br.ReadString(byte('\n'))
	}
	if err != io.EOF {
		log.Fatal(err)
	}
}
