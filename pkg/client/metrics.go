package client

import (
	"net/http"
	"strconv"
	"time"
)

// Metric is a named ... metric. We have them named in prometheus format, but using this type and the given
// constants, non-prometheus users don't have to care about the actual name.
type Metric string

const (
	// MetricRequestDuration is the time in seconds it took to send the request until a response was received.
	MetricRequestDuration Metric = "http_request_duration_seconds"

	// MetricRequestCount is the number of requests sent for the given labels. It is a counter, delivered
	// as 1 for increment and -1 for decrement.
	MetricRequestCount Metric = "http_request_total"

	// MetricRequestInflight is the number of requests currently waiting for a response. It is a counter,
	// delivered as 1 for increment and -1 for decrement.
	MetricRequestInflight Metric = "http_requests_in_flight"
)

// MetricLabel is the key for the labels-map we give when passing metrics to the receiver. It again is just a
// prometheus-like label name but non-prometheus users don't have to care about the actual name this way.
type MetricLabel string

const (
	// MetricLabelResource contains the name of the resource the given metric is about - this might be a
	// request URI, the name of a type or something similar. It will be the same for all requests for
	// the same resource-kind.
	MetricLabelResource MetricLabel = "resource"

	// MetricLabelMethod contains the HTTP verb used in the request.
	MetricLabelMethod MetricLabel = "method"

	// MetricLabelStatus contains the status code we received for the request.
	MetricLabelStatus MetricLabel = "status"
)

// MetricReceiver receives a bunch of metrics with the same labels. Counter metrics will be delivered
// as 1 for increment and -1 for decrement.
type MetricReceiver func(metrics map[Metric]float64, labels map[MetricLabel]string)

type metricsTransport struct {
	baseTransport http.RoundTripper
	receiver      MetricReceiver
}

func (m metricsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	labels := make(map[MetricLabel]string, 3)

	labels[MetricLabelResource] = req.URL.Path
	labels[MetricLabelMethod] = req.Method

	start := time.Now()
	m.receiver(
		map[Metric]float64{
			MetricRequestInflight: 1,
		},
		labels,
	)

	response, err := m.baseTransport.RoundTrip(req)

	if response != nil {
		labels[MetricLabelStatus] = strconv.Itoa(response.StatusCode)
	} else {
		labels[MetricLabelStatus] = "-1"
	}

	m.receiver(
		map[Metric]float64{
			MetricRequestDuration: time.Since(start).Seconds(),
			MetricRequestInflight: -1,
			MetricRequestCount:    1,
		},
		labels,
	)

	return response, err
}

func wrapClientForMetrics(c *http.Client, r MetricReceiver) *http.Client {
	transport := http.DefaultTransport

	if c.Transport != nil {
		transport = c.Transport
	}

	return &http.Client{
		Transport: metricsTransport{
			baseTransport: transport,
			receiver:      r,
		},
		CheckRedirect: c.CheckRedirect,
		Jar:           c.Jar,
		Timeout:       c.Timeout,
	}
}
