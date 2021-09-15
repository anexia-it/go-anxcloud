package common

import "github.com/anexia-it/go-anxcloud/pkg/lbaas/pagination"

type Mode string

const (
	TCP  = Mode("tcp")
	HTTP = Mode("http")

	OptNameSearch = "search"
	OptNameFilter = "filter"
)

type State string

const (
	Updating        = State("0")
	Updated         = State("1")
	DeploymentError = State("2")
	Deployed        = State("3")
	NewlyCreated    = State("4")
)

func SearchOpt(value string) pagination.Option {
	return pagination.Option{
		Name:  OptNameSearch,
		Value: value,
	}
}
