//go:build integration
// +build integration

package v1

import (
	"io/ioutil"
	"net/http"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// The checks in this file run on a configured LBaaS to check if we can actually access things with it as
// expected with the given configuration. These shouldn't be part of go-anxcloud E2E tests but in LBaaS, but
// here we are, testing it with not a lot of code.
//
// This is in no way intended to allow or wish for **full** E2E tests in go-anxcloud. We normally only test API
// bindings. See SYSENG-1239 and the linked AIS site for more.

func successfulConnectionCheck(url string) {
	It("delivers HAProxy status page", func() {
		Eventually(func(g Gomega) {
			resp, err := http.Get(url)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(string(body)).To(ContainSubstring("<title>Statistics Report for HAProxy</title>"))
		}, 5*time.Second, 1*time.Second).Should(Succeed())
	})
}

func unavailableServerConnectionCheck(url string) {
	It("delivers a 503 error", func() {
		Eventually(func(g Gomega) {
			resp, err := http.Get(url)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(resp.StatusCode).To(Equal(http.StatusServiceUnavailable))
		}, 5*time.Second, 1*time.Second).Should(Succeed())
	})
}

func connectionResetByPeerCheck(url string) {
	It("resets connection", func() {
		_, err := http.Get(url)
		Expect(err).To(MatchError(syscall.ECONNRESET))
	})
}
