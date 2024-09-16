package v1_test

import (
	"net/http"

	. "github.com/onsi/gomega/ghttp"
)

const (
	mockProgressIdentifier = "PROGRESSg1234567aaaaaaaabbbbbbbb"
)

var (
	mockGetProvisioned = false
)

func prepareGetProgress() {
	if isIntegrationTest {
		return
	}
	var response http.HandlerFunc

	progress := 0
	status := ""
	id := ""

	if mockGetProvisioned {
		progress = 100
		status = "1"
		id = mockVMIdentifier
	}

	response = RespondWithJSONEncoded(200, map[string]interface{}{
		"errors":        []string{},
		"identifier":    mockProgressIdentifier,
		"progress":      progress,
		"queued":        false,
		"status":        status,
		"vm_identifier": id,
	})

	mock.AppendHandlers(CombineHandlers(
		VerifyRequest("GET", "/api/vsphere/v1/provisioning/progress.json/"+mockProgressIdentifier),
		response,
	))
}

func prepareEventuallyProvisioned() {
	prepareGetProgress()
	mockGetProvisioned = true
}
