package rules

import (
	"context"
	"errors"

	"go.anx.io/go-anxcloud/pkg/client"
)

type AutomationResult struct {
	State    AutomationState `json:"state"`
	Messages []string        `json:"messages"`
	Data     any             `json:"data"`
}

var ErrNoAutomationSuccess = errors.New("firing of automation rule did not return success")

// Validate validates the automation result for success.
func (ar *AutomationResult) Validate() error {
	if ar.State != StateSuccess {
		errs := []error{ErrNoAutomationSuccess}

		for _, m := range ar.Messages {
			errs = append(errs, errors.New(m))
		}

		return errors.Join(errs...)
	}

	return nil
}

type AutomationState string

const (
	StateSuccess AutomationState = "success"
)

type API interface {
	FireSingle(ctx context.Context, ruleIdentifier string, objectIdentifier string) (AutomationResult, error)
}

type api struct {
	client client.Client
}

func NewAPI(c client.Client) API {
	return &api{c}
}
