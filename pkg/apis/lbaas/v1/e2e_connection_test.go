//go:build integration
// +build integration

package v1_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
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
		pollCheck := func() error {
			resp, err := http.Get(url) // #nosec G107 -- URL is controlled test endpoint
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusOK {
				return errors.New("unexpected status code")
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			if !strings.Contains(string(body), "<title>Statistics Report for HAProxy</title>") {
				return errors.New("HAProxy status page title not found")
			}
			return nil
		}

		Eventually(pollCheck, 5*time.Second, 1*time.Second).Should(Succeed())
	})
}

func unavailableServerConnectionCheck(url string) {
	It("delivers a 503 error", func() {
		pollCheck := func() error {
			resp, err := http.Get(url) // #nosec G107 -- URL is controlled test endpoint
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusServiceUnavailable {
				return errors.New("expected 503 Service Unavailable")
			}
			return nil
		}

		Eventually(pollCheck, 60*time.Second, 5*time.Second).Should(Succeed())
	})
}

func connectionResetByPeerCheck(url string) {
	It("resets connection", func() {
		pollCheck := func() error {
			_, err := http.Get(url) // #nosec G107 -- URL is controlled test endpoint
			if err == nil {
				return errors.New("expected connection reset error")
			}
			if !errors.Is(err, syscall.ECONNRESET) {
				return err
			}
			return nil
		}

		Eventually(pollCheck, 30*time.Second, 1*time.Second).Should(Succeed())
	})
}
