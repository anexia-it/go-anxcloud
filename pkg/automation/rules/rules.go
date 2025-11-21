package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	utils "path"
)

const path = "/automation/rules/"
const fireSingle = "fire-single"

func (a *api) FireSingle(ctx context.Context, ruleIdentifier string, objectIdentifier string) (AutomationResult, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return AutomationResult{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, ruleIdentifier, fireSingle, objectIdentifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), nil)
	if err != nil {
		return AutomationResult{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return AutomationResult{}, fmt.Errorf("error when firing automation rule '%s': %w", ruleIdentifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return AutomationResult{}, fmt.Errorf("could not fire automation rule '%s': %s", ruleIdentifier,
			response.Status)
	}

	var payload AutomationResult

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return AutomationResult{}, fmt.Errorf("could not parse automation result for '%s' : %w",
			ruleIdentifier, err)
	}

	return payload, nil
}
