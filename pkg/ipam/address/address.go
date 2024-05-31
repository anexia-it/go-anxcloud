// Package address implements API functions residing under /ipam/address.
// This path contains methods for managing IPs.
package address

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.anx.io/go-anxcloud/pkg/utils/param"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	pathAddressPrefix         = "/api/ipam/v1/address.json"
	pathReserveAddressPrefix  = "/api/ipam/v1/address/reserve/ip/count.json"
	pathFilteredAddressPrefix = "/api/ipam/v1/address/filtered.json"
)

// Filters that can be applied to the GetFiltered request
var (
	PrefixFilter       = param.ParameterBuilder("prefix")
	VlanFilter         = param.ParameterBuilder("vlan")
	VersionFilter      = param.ParameterBuilder("version")
	RoleTextFilter     = param.ParameterBuilder("role_text")
	StatusFilter       = param.ParameterBuilder("status")
	LocationFilter     = param.ParameterBuilder("location")
	OrganizationFilter = param.ParameterBuilder("organization_identifier")
)

// Address contains all the information about a specific address.
type Address struct {
	ID                  string `json:"identifier"`
	Name                string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
	DescriptionInternal string `json:"description_internal"`
	Role                string `json:"role_text"`
	Version             int    `json:"version"`
	Status              string `json:"status"`
	VLANID              string `json:"vlan"`
	PrefixID            string `json:"prefix"`
}

// Summary is the address information returned by a listing.
type Summary struct {
	ID                  string `json:"identifier"`
	Name                string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
	Role                string `json:"role_text"`
}

// Update contains fields to change on a prefix.
type Update struct {
	Name                string `json:"name,omitempty"`
	DescriptionCustomer string `json:"description_customer,omitempty"`
	Role                string `json:"role,omitempty"`
}

// Create defines meta data of an address to create.
type Create struct {
	PrefixID            string `json:"prefix"`
	Address             string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
	Role                string `json:"role"`
	Organization        string `json:"organization"`
}

// IPReserveVersionLimit limits the IP version for address reservations
type IPReserveVersionLimit int

const (
	// IPReserveVersionLimitIPv4 specifies that only IPv4 addresses should be reserved
	IPReserveVersionLimitIPv4 IPReserveVersionLimit = 4

	// IPReserveVersionLimitIPv4 specifies that only IPv6 addresses should be reserved
	IPReserveVersionLimitIPv6 IPReserveVersionLimit = 6
)

// ReserveRandom defines metadata of addresses to reserve randomly.
type ReserveRandom struct {
	LocationID string `json:"location_identifier"`
	VlanID     string `json:"vlan_identifier"`
	Count      int    `json:"count"`

	// PrefixID limits the potential addresses to a specific prefix
	// within the specified VLAN. Defaults to all prefixes within the VLAN.
	PrefixID string `json:"prefix_identifier,omitempty"`

	// IPVersion limits the potential addresses to a specific IP version.
	// Defaults to v4 or v6, depending on availability.
	IPVersion IPReserveVersionLimit `json:"ip_version,omitempty"`

	// ReservationPeriod specifies how many seconds the addresses should be reserved.
	// If IP addresses haven't been assigned to resources within the period, they are released again.
	// Defaults to 1800 seconds (= 30 minutes).
	ReservationPeriod uint `json:"reservation_period,omitempty"`
}

// ReserveRandomSummary is the reserved IPs information returned by list request.
type ReserveRandomSummary struct {
	Limit      int          `json:"limit"`
	Page       int          `json:"page"`
	TotalItems int          `json:"total_items"`
	TotalPages int          `json:"total_pages"`
	Data       []ReservedIP `json:"data"`
}

// ReservedIP returns details about reserved ip.
type ReservedIP struct {
	ID      string `json:"identifier"`
	Address string `json:"text"`
	Prefix  string `json:"prefix"`
}

type listResponse struct {
	Data struct {
		Data []Summary `json:"data"`
	} `json:"data"`
}

// NewCreate creates a new address definition with required vlaues.
func NewCreate(prefixID string, address string) Create {
	return Create{
		PrefixID: prefixID,
		Address:  address,
		Role:     "Default",
	}
}

func (a api) List(ctx context.Context, page, limit int, search string) ([]Summary, error) {
	url := fmt.Sprintf(
		"%s%s?page=%d&limit=%d&search=%s",
		a.client.BaseURL(),
		pathAddressPrefix, page, limit, search,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create address list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute address list request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute address list request, got response %s", httpResponse.Status)
	}

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)

	if err != nil {
		return nil, fmt.Errorf("could not decode address list response: %w", err)
	}

	return responsePayload.Data.Data, err
}

func (a api) Get(ctx context.Context, id string) (Address, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathAddressPrefix,
		id,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Address{}, fmt.Errorf("could not create address get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Address{}, fmt.Errorf("could not execute address get request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Address{}, fmt.Errorf("could not execute address get request, got response %s", httpResponse.Status)
	}

	var responsePayload Address
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)

	if err != nil {
		return Address{}, fmt.Errorf("could not decode address get response: %w", err)
	}

	return responsePayload, err
}

func (a api) Delete(ctx context.Context, id string) error {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathAddressPrefix,
		id,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create address delete request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute address delete request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute address delete request, got response %s", httpResponse.Status)
	}

	return nil
}

func (a api) Create(ctx context.Context, create Create) (Summary, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathAddressPrefix,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(create); err != nil {
		panic(fmt.Sprintf("could not create request data for vlan creation: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return Summary{}, fmt.Errorf("could not create vlan post request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Summary{}, fmt.Errorf("could not execute vlan post request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Summary{}, fmt.Errorf("could not execute vlan post request, got response %s", httpResponse.Status)
	}

	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	if err != nil {
		return Summary{}, fmt.Errorf("could not decode vlan post response: %w", err)
	}

	return summary, nil
}

func (a api) Update(ctx context.Context, id string, update Update) (Summary, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathAddressPrefix, id,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(update); err != nil {
		panic(fmt.Sprintf("could not create request data for vlan update: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, &requestData)
	if err != nil {
		return Summary{}, fmt.Errorf("could not create vlan update request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Summary{}, fmt.Errorf("could not execute vlan update request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Summary{}, fmt.Errorf("could not execute vlan update request, got response %s", httpResponse.Status)
	}

	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	if err != nil {
		return summary, fmt.Errorf("could not decode vlan update response: %w", err)
	}

	return summary, err
}

func (a api) ReserveRandom(ctx context.Context, reserve ReserveRandom) (ReserveRandomSummary, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathReserveAddressPrefix,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(reserve); err != nil {
		panic(fmt.Sprintf("could not create request data for IP address reservation: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return ReserveRandomSummary{}, fmt.Errorf("could not create IP address reserve random post request: %w", err)
	}

	// Workaround to avoid race-conditions on IP reservations for the same VLAN
	randomDelay := time.Duration(rand.Intn(1000))
	time.Sleep(randomDelay * time.Millisecond)

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return ReserveRandomSummary{}, fmt.Errorf("could not execute IP address reserve random post request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return ReserveRandomSummary{}, fmt.Errorf("could not execute IP address reserve random post request, got response %s", httpResponse.Status)
	}

	var summary ReserveRandomSummary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	if err != nil {
		return ReserveRandomSummary{}, fmt.Errorf("could not decode IP address reserve random post response: %w", err)
	}

	return summary, nil
}

func (a api) GetFiltered(ctx context.Context, page, limit int, filters ...param.Parameter) ([]Summary, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = pathFilteredAddressPrefix
	query := endpoint.Query()
	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(limit))
	for _, filter := range filters {
		filter(query)
	}

	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when executing request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not get filtered addresses %s", response.Status)
	}

	var payload struct {
		Data struct {
			Data []Summary `json:"data"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&payload)
	return payload.Data.Data, err
}
