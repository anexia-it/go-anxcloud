//go:build !integration

package v1_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/client"
)

type statusFlags struct {
	SetAssigned bool // Engine takes some time to assign a name, emulate with this
	SetActive   bool // Engine takes some time to have a new VLAN active, emulate with this
	SetDeleting bool // After destroying, the Engine needs some time to actually delete it - emulate with this
	SetDeleted  bool // This is set to true to return a 404 response when retrieving our test VLAN
	SetInactive bool
}

type mockAPI struct {
	api.API
	srv *Server

	statusFlags statusFlags
}

func GetTestAPIClient() api.API {
	GinkgoHelper()

	s := NewServer()
	// s.Writer = GinkgoWriter
	DeferCleanup(s.Close)

	a, err := api.NewAPI(
		api.WithLogger(GinkgoLogr),
		api.WithClientOptions(
			client.BaseURL(s.URL()),
			client.IgnoreMissingToken(),
		),
	)
	Expect(err).ToNot(HaveOccurred())

	return &mockAPI{API: a, srv: s}
}
