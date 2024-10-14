//go:build !integration
// +build !integration

package progress

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/client"
)

var (
	mock *ghttp.Server
)

var _ = Describe("vsphere/provisioning/progress API client", func() {
	var (
		cli client.Client

		identifier string
		result     Progress
		requestErr error
	)

	JustBeforeEach(func() {
		result, requestErr = NewAPI(cli).Get(context.TODO(), identifier)
	})

	Context("retrieving progress", func() {
		BeforeEach(func() {
			mock = ghttp.NewServer()
			var err error
			cli, err = client.New(client.BaseURL(mock.URL()), client.IgnoreMissingToken())
			Expect(err).ToNot(HaveOccurred())
		})

		When("progress has failed", func() {
			BeforeEach(func() {
				identifier = "statusFailedTest"
				prepareGet(identifier, []string{}, StatusFailed)
			})

			It("it returns status failed", func() {
				Expect(requestErr).NotTo(HaveOccurred())
				Expect(result.Status).To(Equal(StatusFailed))
			})
		})

		When("progress has succeeded", func() {
			BeforeEach(func() {
				identifier = "statusSuccessTest"
				prepareGet(identifier, []string{}, StatusSuccess)
			})

			It("it returns status success", func() {
				Expect(requestErr).NotTo(HaveOccurred())
				Expect(result.Status).To(Equal(StatusSuccess))
			})
		})

		When("progress has errors", func() {
			BeforeEach(func() {
				identifier = "statusInProgressTest"
				prepareGet(identifier, []string{"some error"}, StatusInProgress)
			})

			It("it returns an error", func() {
				Expect(requestErr).To(HaveOccurred())
				Expect(result.Status).To(Equal(StatusInProgress))
				Expect(result.Errors).To(HaveLen(1))
			})
		})

		When("progress is cancelled", func() {
			BeforeEach(func() {
				identifier = "statusCancelledTest"
				prepareGet(identifier, []string{}, StatusCancelled)
			})

			It("it returns status cancelled", func() {
				Expect(requestErr).NotTo(HaveOccurred())
				Expect(result.Status).To(Equal(StatusCancelled))
			})
		})
	})
})

func prepareGet(identifier string, errors []string, status Status) {
	mock.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", "/api/vsphere/v1/provisioning/progress.json/"+identifier),
		ghttp.RespondWithJSONEncoded(http.StatusOK, Progress{
			TaskIdentifier: identifier,
			Queued:         false,
			Progress:       0,
			VMIdentifier:   "",
			Errors:         errors,
			Status:         status,
		}),
	))

}
