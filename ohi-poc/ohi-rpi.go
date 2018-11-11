package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func vcgencmd(args ...string) string {
	cmd := exec.Command("vcgencmd", args...)

	out, err := cmd.Output()
	if err != nil {
		log.Println("ERROR cmd.Output() ", err)
	}

	return string(out)
}

func main() {

	for {
		sample := make(map[string]interface{})
		sample["eventType"] = "RaspberryPiSample"
		temp := vcgencmd("measure_temp")
		degrees := strings.Split(temp, "=")[1]
		sample["cpu_temperature"], _ = strconv.ParseFloat(strings.Split(degrees, "'")[0], 32)

		freq, err := ioutil.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
		if err != nil {
			log.Println("ERROR reading scaling_cur_freq", err)
		}

		sample["cpu_frequency"], err = strconv.Atoi(strings.Trim(string(freq), " \n"))
		if err != nil {
			log.Println("ERROR reading strconv.Atoi", err)
		}

		jsonSample, _ := json.Marshal(sample)

		fmt.Println(string(jsonSample))
		http.Post("http://localhost:8080/http", "application/json", bytes.NewReader(jsonSample))

		time.Sleep(15 * time.Second)
	}
}
