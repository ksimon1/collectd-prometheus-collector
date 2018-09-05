# collectd-prometheus-collector

[![Go Report Card](https://goreportcard.com/badge/github.com/ksimon1/collectd-prometheus-collector)](https://goreportcard.com/report/github.com/ksimon1/collectd-prometheus-collector)

**collectd-prometheus-collector** is a script which retrives data from collectd prometheus server and saves it to file. 
By default it removes last column with [timestamp](https://prometheus.io/docs/instrumenting/exposition_formats/#line-format) from data due to node-exporter
