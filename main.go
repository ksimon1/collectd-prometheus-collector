package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		collectdURL      string
		fileLocation     string
		pollInterval     int
		removeLastColumn bool
	)

	flag.StringVar(&collectdURL, "collectdURL", "", "collectd URL *required")
	flag.StringVar(&fileLocation, "fileLocation", "", "location of file to write in *required")
	flag.IntVar(&pollInterval, "pollInterval", 15, "poll interval")
	flag.BoolVar(&removeLastColumn, "removeLastColumn", true, "remove last column with timestamp")

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
			log.Println("could not get data from collectd: ", err.Error())
		}

		if len(data) > 0 {
			if removeLastColumn {
				data = removeLastColumnFromData(data)
			}

			err = writeToFile(fileLocation, data)
			if err != nil {
				log.Println("could not write data to file: ", err.Error())
			}
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

//getPrometheusData retrieves prometheus data from server
func getPrometheusData(collectdURL string) ([]byte, error) {
	r, err := http.Get(collectdURL)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	return ioutil.ReadAll(r.Body)
}

//removeLastColumnFromData deletes timestamp from each line
func removeLastColumnFromData(data []byte) []byte {
	lines := strings.Split(string(data), "\n")

	updatedLines := make([]string, len(lines))
	for i, line := range lines {
		//if line is empty or starts with # skip it
		if line == "" || string(line[0]) == "#" {
			updatedLines[i] = line
			continue
		}
		updatedLines[i] = strings.TrimSuffix(line, line[strings.LastIndex(line, " "):])
	}

	return []byte(strings.Join(updatedLines, "\n"))
}

//writeToFile checks if file path exists and writes data to file
func writeToFile(fileLocation string, prometheusData []byte) error {
	pathToFile := strings.TrimSuffix(fileLocation, fileLocation[strings.LastIndex(fileLocation, "/"):])
	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		log.Fatalln("file location is not accessible")
	}

	return ioutil.WriteFile(fileLocation, prometheusData, 0644)
}
