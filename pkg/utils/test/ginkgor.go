package test

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/onsi/ginkgo/v2"
)

// NewGinkgor creates a new logr.Logger logging to GinkgoWriter. Log messages are displayed when a test fails but
// can optionally be displayed always by passing `-v` to ginkgo.
//
// Verbosity is hard-coded to 3, enabling pkg/client request/response logging.
func NewGinkgor() logr.Logger {
	stdr.SetVerbosity(3)

	g := ginkgor{}
	return stdr.New(g)
}

type ginkgor struct{}

func (g ginkgor) Output(calldepth int, logline string) error {
	ginkgo.GinkgoWriter.Println(logline)
	return nil
}
