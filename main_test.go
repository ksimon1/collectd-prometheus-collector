package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	rawTestData = []byte("collectd_processes_ps_vm{processes=\"collectd\",instance=\"vmi-ephemeral\"} 835985408 1234567891234")

	mixedData = `collectd_processes_ps_rss{processes="collectd",instance="vmi-ephemeral"} 2985984 1536057993708
collectd_processes_ps_data{processes="libvirtd",instance="vmi-ephemeral"} 1498898432 1536057993708
# HELP collectd_processes_ps_stacksize write_prometheus plugin: 'processes' Type: 'ps_stacksize', Dstype: 'gauge', Dsname: 'value'
# TYPE collectd_processes_ps_stacksize gauge
collectd_processes_ps_stacksize{processes="collectd",instance="vmi-ephemeral"} 2224 1536057993708`

	afterUpdateMixedData = `collectd_processes_ps_rss{processes="collectd",instance="vmi-ephemeral"} 2985984
collectd_processes_ps_data{processes="libvirtd",instance="vmi-ephemeral"} 1498898432
# HELP collectd_processes_ps_stacksize write_prometheus plugin: 'processes' Type: 'ps_stacksize', Dstype: 'gauge', Dsname: 'value'
# TYPE collectd_processes_ps_stacksize gauge
collectd_processes_ps_stacksize{processes="collectd",instance="vmi-ephemeral"} 2224`

	immutableData = "# HELP collectd_processes_ps_stacksize write_prometheus plugin"

	dataWithRemovedLastColumn = "collectd_processes_ps_vm{processes=\"collectd\",instance=\"vmi-ephemeral\"} 835985408"
)

func getTestServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(rawTestData)
	}))

	return ts
}
func TestGetPrometheusData(t *testing.T) {
	ts := getTestServer()
	defer ts.Close()

	data, err := getPrometheusData(ts.URL)
	if err != nil {
		t.Error("it must not throw an error while getting prometheus data", err)
	}

	if string(data) != string(rawTestData) {
		t.Error("it must not update the data", err)
	}

	_, err = getPrometheusData("someFakeURL")
	if err == nil {
		t.Error("it must throw an error when passing inaccessible URL")
	}
}

func TestRemoveLastColumnFromData(t *testing.T) {
	updatedData := removeLastColumnFromData(rawTestData)
	if string(updatedData) == string(rawTestData) {
		t.Error("it must update the data")
	}

	if string(updatedData) != dataWithRemovedLastColumn {
		t.Error("it must have the last column removed")
	}

	updatedData = removeLastColumnFromData([]byte(immutableData))
	if string(updatedData) != immutableData {
		t.Error("it must not update the immutable data")
	}

	updatedMixedData := removeLastColumnFromData([]byte(mixedData))
	if string(updatedMixedData) != afterUpdateMixedData {
		t.Error("it must correctly update the mixed data")
	}
}
