package client

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

type metricTestTransport struct {
	numRequests int
}

func (m *metricTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.numRequests++

	if m.numRequests > 1 {
		return nil, errors.New("test error")
	}

	return http.DefaultTransport.RoundTrip(req)
}

var _ = Describe("client metrics", func() {
	type metricReceiverCall struct {
		metrics map[Metric]float64
		labels  map[MetricLabel]string
	}

	var server *ghttp.Server
	var client Client
	var hc *http.Client

	var receivedMetrics []metricReceiverCall
	var countRequests int
	var countInflight int
	var requestDurationTotal float64

	BeforeEach(func() {
		hc = http.DefaultClient
	})

	JustBeforeEach(func() {
		countRequests = 0
		countInflight = 0
		requestDurationTotal = 0

		receivedMetrics = make([]metricReceiverCall, 0, 2)

		server = ghttp.NewServer()

		client, _ = New(
			WithClient(hc),
			ParseEngineErrors(false),
			IgnoreMissingToken(),
			BaseURL(server.URL()),
			WithMetricReceiver(func(m map[Metric]float64, l map[MetricLabel]string) {
				if v, o := m[MetricRequestCount]; o {
					countRequests += int(v)
				}

				if v, o := m[MetricRequestInflight]; o {
					countInflight += int(v)
				}

				if v, o := m[MetricRequestDuration]; o {
					requestDurationTotal += v
				}

				mrc := metricReceiverCall{
					metrics: make(map[Metric]float64, len(m)),
					labels:  make(map[MetricLabel]string, len(l)),
				}

				for key, val := range m {
					mrc.metrics[key] = val
				}

				for key, val := range l {
					mrc.labels[key] = val
				}

				receivedMetrics = append(receivedMetrics, mrc)
			}),
		)
	})

	It("delivers counter metrics as expected", func() {
		type testRequest struct {
			method string
			status int
		}

		testRequests := []testRequest{
			{"GET", 200},
			{"GET", 404},
			{"POST", 200},
			{"DELETE", 200},
			{"GET", 500},
		}

		for i, tr := range testRequests {
			r, err := http.NewRequest(tr.method, client.BaseURL()+"/", nil)
			Expect(err).NotTo(HaveOccurred())

			server.AppendHandlers(ghttp.RespondWith(tr.status, nil))
			_, err = client.Do(r)
			Expect(err).NotTo(HaveOccurred())

			firstMetric := receivedMetrics[i*2]
			secondMetric := receivedMetrics[i*2+1]

			Expect(firstMetric).To(BeEquivalentTo(
				metricReceiverCall{
					metrics: map[Metric]float64{
						MetricRequestInflight: 1,
					},
					labels: map[MetricLabel]string{
						MetricLabelResource: "/",
						MetricLabelMethod:   tr.method,
					},
				},
			))

			Expect(secondMetric.labels).To(BeEquivalentTo(
				map[MetricLabel]string{
					MetricLabelResource: "/",
					MetricLabelMethod:   tr.method,
					MetricLabelStatus:   strconv.Itoa(tr.status),
				},
			))
		}

		Expect(countRequests).To(Equal(len(testRequests)))
		Expect(countInflight).To(Equal(0))
	})

	It("delivers request duration metric as expected", func() {
		server.AppendHandlers(ghttp.CombineHandlers(
			func(res http.ResponseWriter, req *http.Request) {
				time.Sleep(500 * time.Millisecond)
			},
			ghttp.RespondWith(200, nil),
		))

		r, err := http.NewRequest("GET", client.BaseURL()+"/", nil)
		Expect(err).NotTo(HaveOccurred())

		_, err = client.Do(r)
		Expect(err).NotTo(HaveOccurred())

		Expect(countRequests).To(Equal(1))
		Expect(countInflight).To(Equal(0))
		Expect(requestDurationTotal).To(BeNumerically("~", 0.5, 0.1))
	})

	Context("when configured with a custom http.Client", func() {
		var transport metricTestTransport
		BeforeEach(func() {
			transport = metricTestTransport{}

			hc = &http.Client{
				Transport: &transport,
			}
		})

		It("uses the transport of the custom http.Client", func() {
			server.AppendHandlers(ghttp.RespondWith(200, nil))

			r, err := http.NewRequest("GET", client.BaseURL()+"/", nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = client.Do(r)
			Expect(err).NotTo(HaveOccurred())

			Expect(countRequests).To(Equal(1))
			Expect(countInflight).To(Equal(0))
			Expect(requestDurationTotal).To(BeNumerically("~", 0, 0.1))

			Expect(transport.numRequests).To(Equal(1))
		})

		It("uses status code -1 for transport returning an error", func() {
			server.AppendHandlers(
				ghttp.RespondWith(200, nil),
				// only first request comes through to the server, second is catched in metricTestTransport.RoundTrip
			)

			r, err := http.NewRequest("GET", client.BaseURL()+"/", nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = client.Do(r)
			Expect(err).NotTo(HaveOccurred())

			Expect(countRequests).To(Equal(1))
			Expect(countInflight).To(Equal(0))
			Expect(requestDurationTotal).To(BeNumerically("~", 0, 0.1))

			Expect(transport.numRequests).To(Equal(1))

			r, err = http.NewRequest("GET", client.BaseURL()+"/", nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = client.Do(r)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))

			Expect(countRequests).To(Equal(2))
			Expect(countInflight).To(Equal(0))
			Expect(requestDurationTotal).To(BeNumerically("~", 0, 0.1))
			Expect(transport.numRequests).To(Equal(2))

			Expect(receivedMetrics[1].labels[MetricLabelStatus]).To(Equal("200"))
			Expect(receivedMetrics[3].labels[MetricLabelStatus]).To(Equal("-1"))
		})
	})
})
