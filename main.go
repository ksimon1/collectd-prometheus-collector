package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	var (
		collectdURL  string
		fileLocation string
		graceTimeout int
	)

	flag.StringVar(&collectdURL, "collectdURL", "", "collectd URL *required")
	flag.StringVar(&fileLocation, "fileLocation", "", "location of file to write in *required")
	flag.IntVar(&graceTimeout, "wait", 15, "collectd URL")
	flag.Parse()

	if collectdURL == "" {
		log.Fatalln("collectdURL is required param")
	}

	if fileLocation == "" {
		log.Fatalln("fileLocation is required param")
	}

	for {
		data, err := getPrometheusData(collectdURL)
		if err != nil {
			fmt.Println("could not get data from collectd: ", err.Error())
		}

		if len(data) > 0 {
			err = writeToFile(fileLocation, data)
			if err != nil {
				fmt.Println("could not write data to file: ", err.Error())
			}
		}

		time.Sleep(time.Duration(graceTimeout) * time.Second)
	}
}

func getPrometheusData(collectdURL string) ([]byte, error) {
	r, err := http.Get(collectdURL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return ioutil.ReadAll(r.Body)
}

func writeToFile(fileLocation string, prometheusData []byte) error {
	return ioutil.WriteFile(fileLocation, prometheusData, 0644)
}
