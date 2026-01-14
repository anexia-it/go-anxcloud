// Package automation implements API functions residing under /automation.
// This path contains methods for managing automations.
package automation

import (
	"errors"
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
