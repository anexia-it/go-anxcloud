package v2

import (
	"context"
	"net/http"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// EndpointURL returns the URL where to retrieve objects of type Tenant and the identifier of the given Tenant.
// It implements the api.Object interface on *Tenant, making it usable with the generic API client.
func (t *Tenant) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/object_storage/v2/tenant")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		query := u.Query()

		// Add attributes parameter to get all fields
		query.Add("attributes", "name,state,remote_id,description,user,password,quota,usage,backend,reseller,customer")

		filters := make(url.Values)

		if t.State != nil && t.State.ID != "" {
			filters.Add("state", t.State.ID)
		}

		if t.Backend.Identifier != "" {
			filters.Add("backend", t.Backend.Identifier)
		}

		if t.UserName != "" {
			filters.Add("user", t.UserName)
		}

		if t.CustomerIdentifier != "" {
			filters.Add("customer", t.CustomerIdentifier)
		}

		if t.ResellerIdentifier != "" {
			filters.Add("reseller", t.ResellerIdentifier)
		}

		if len(filters) > 0 {
			query.Add("filters", filters.Encode())
		}

		u.RawQuery = query.Encode()
	}

	return u, err
}

// FilterAPIRequest modifies the URL based on operation type and implements RequestFilterHook interface
func (t *Tenant) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// For UPDATE operations, use PATCH method (Object Storage API uses PATCH, not PUT)
	if op == types.OperationUpdate {
		// Keep the identifier in the path but change method from PUT to PATCH
		req.Method = "PATCH"
	}

	return req, nil
}

// FilterAPIRequestBody generates the request body for Tenants, replacing linked Objects with just their identifier.
func (t *Tenant) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	body := requestBody(ctx, func() interface{} {
		// For UPDATE operations, only send the fields that are actually being updated
		if op == types.OperationUpdate {
			updateBody := make(map[string]interface{})

			// Only include fields that have been explicitly set for update
			if t.Name != "" {
				updateBody["name"] = t.Name
			}
			if t.Description != "" {
				updateBody["description"] = t.Description
			}
			if t.UserName != "" {
				updateBody["user"] = t.UserName
			}
			if t.Password != "" {
				updateBody["password"] = t.Password
			}
			if t.Quota != nil {
				updateBody["quota"] = *t.Quota
			}

			return updateBody
		}

		// For CREATE operations, use the full request body
		return &struct {
			// Exclude the Backend field to avoid conflict
			CustomerIdentifier string                 `json:"customer_identifier,omitempty"`
			ResellerIdentifier string                 `json:"reseller_identifier,omitempty"`
			Identifier         string                 `json:"identifier,omitempty"`
			Tags               gs.PartialResourceList `json:"tags,omitempty"`
			Reseller           string                 `json:"reseller,omitempty"`
			Customer           string                 `json:"customer,omitempty"`
			Share              bool                   `json:"share,omitempty"`
			AutomationRules    []AutomationRule       `json:"automation_rules,omitempty"`
			Name               string                 `json:"name"`
			State              *GenericAttributeState `json:"state,omitempty"`
			RemoteID           *string                `json:"remote_id,omitempty"`
			Description        string                 `json:"description"`
			UserName           string                 `json:"user"`
			Password           string                 `json:"password,omitempty"`
			Quota              *float64               `json:"quota,omitempty"`
			Usage              *float64               `json:"usage,omitempty"`

			// Backend as string identifier
			Backend string `json:"backend"`
		}{
			CustomerIdentifier: t.CustomerIdentifier,
			ResellerIdentifier: t.ResellerIdentifier,
			Identifier:         t.Identifier,
			Tags:               t.Tags,
			Reseller:           t.Reseller,
			Customer:           t.Customer,
			Share:              t.Share,
			AutomationRules:    t.AutomationRules,
			Name:               t.Name,
			State:              t.State,
			RemoteID:           t.RemoteID,
			Description:        t.Description,
			UserName:           t.UserName,
			Password:           t.Password,
			Quota:              t.Quota,
			Usage:              t.Usage,
			Backend:            t.Backend.Identifier,
		}
	})

	return body, nil
}
