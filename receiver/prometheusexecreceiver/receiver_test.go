// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheusexecreceiver

import (
	"context"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/config/configtest"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/consumer/pdatautil"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.uber.org/zap"
)

// TestEndToEnd loads the test config.yaml and tests two things:
// 1. makes sure the prometheus_exec config without an exec key throws an error
// 2. An end-to-end test where metrics are scraped from a fake exporter that exposes Promtheus metrics,
// and makes sure that exporter subprocess is restarted correctly
func TestEndToEnd(t *testing.T) {
	// Load the config from the yaml file
	factories, err := componenttest.ExampleComponents()
	assert.NoError(t, err)

	factory := &Factory{}
	factories.Receivers[factory.Type()] = factory

	config, err := configtest.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)

	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Receiver without exec key, error expected and checked within function
	execErrorTest(t, config.Receivers["prometheus_exec"])

	// Normal test with port undefined by user
	endToEndScrapeTest(t, config.Receivers["prometheus_exec/end_to_end_test/2"], "end-to-end port not defined")
}

// execErrorTest makes sure the config passed throws an error, since it's missing the exec key
func execErrorTest(t *testing.T, errorReceiverConfig configmodels.Receiver) {
	_, err := new(component.ReceiverCreateParams{Logger: zap.NewNop()}, errorReceiverConfig.(*Config), nil)
	if err == nil {
		t.Errorf("end_to_end_test.go didn't get error, was expecting one since this config has no 'exec' key")
	}
}

// endToEndScrapeTest scrapes a test endpoint (test_prometheus_exporter.go) twice, and between each scrape yields the execution with Sleep() to wait for the subprocess (exporter) to restar
// - wait time is about 1s - to fail and restart, meaning it verifies three things: the scrape is successful (twice), the process was restarted correctly when failed and the underlying
// Prometheus receiver was correctly stopped and then restarted. For extra testing the metrics values are different every time the subprocess exporter is started
// And the uniqueness of the metric scraped is verified
func endToEndScrapeTest(t *testing.T, receiverConfig configmodels.Receiver, testName string) {
	sink := &exportertest.SinkMetricsExporter{}
	wrapper, err := new(component.ReceiverCreateParams{Logger: zap.NewNop()}, receiverConfig.(*Config), sink)
	if err != nil {
		t.Errorf("end_to_end_test.go got error = %w", err)
	}

	ctx := context.Background()
	err = wrapper.Start(ctx, componenttest.NewNopHost())
	if err != nil {
		t.Errorf("end_to_end_test.go got error = %w", err)
	}
	defer wrapper.Shutdown(ctx)

	var metrics []pdata.Metrics

	// Make sure a first scrape works by checking for metrics in the test metrics exporter "sink", only return true when there are metrics
	const waitFor = 15 * time.Second
	const tick = 500 * time.Millisecond
	require.Eventuallyf(t, func() bool {
		got := sink.AllMetrics()
		return len(got) != 0
	}, waitFor, tick, "No metrics were collected after %v for the first scrape (%v)", waitFor, testName)

	metrics = sink.AllMetrics()

	// Make sure the second scrape is successful, and validate that the metrics are different in the second scrape
	require.Eventuallyf(t, func() bool {
		got := sink.AllMetrics()
		if len(got) == 0 || len(got) == len(metrics) {
			return false
		}
		if validateMetrics(&got) {
			return true
		}

		metrics = got
		return false
	}, waitFor, tick, "No metrics were collected after %v for the second scrape (%v)", waitFor, testName)
}

// validateMetrics iterates over the found metrics and returns true if it finds at least 2 unique metrics, meaning the endpoint
// was successfully scraped twice AND the subprocess being handled was stopped and restarted
func validateMetrics(metricsSlice *[]pdata.Metrics) bool {
	var value float64
	for i, val := range *metricsSlice {
		temp := pdatautil.MetricsToMetricsData(val)[0].Metrics[0].Timeseries[0].Points[0].GetDoubleValue()
		if i != 0 && temp != value {
			return true
		}
		if temp != value {
			value = temp
		}
	}
	return false
}
