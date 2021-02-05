// Package zone implements API functions residing under /zone.
// This path contains methods for querying and setting the DNS zones and records.
package zone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	pathPrefix string = "api/clouddns/v1/zone.json"
	zoneDataPayloadKey string = "zoneData"
)

type listResponse struct {
	Results []Response `json:"results"`
}

type Record struct {
	Identifier string `json:"identifier"`
	Immutable bool `json:"immutable"`
	Name string `json:"name"`
	RData string `json:"rdata"`
	Region string `json:"region"`
	TTL string `json:"ttl"`
	Type string `json:"Type"`
}

type Revision struct {
	CreatedAt time.Time `json:"created_at"`
	Identifier string `json:"identifier"`
	ModifiedAt time.Time `json:"modified_at"`
	Records []Record `json:"records"`
}

type DNSServer struct {
	// Required - DNS Server name (FQDN).
	Server string `json:"server"`

	// DNS Server alias
	Alias string
}

type Definition struct {

	// Required - Definition name
	Name string `json:"name"`

	// Required - Is master flag
	// Flag designating if CloudDNS operates as master or slave.
	IsMaster bool `json:"master"`

	// Required - DNSSEC mode
	// DNSSEC mode (master-only) ["managed" or "unvalidated"].
	DNSSecMode string `json:"dnssec_mode"`

	// Required - Admin email address
	// Admin email address used in SOA record.
	AdminEmail string `json:"admin_email"`

	// Required - Refresh value
	// Refresh value used in SOA record.
	Refresh int `json:"refresh"`

	// Required - Retry value
	//Retry value used in SOA record.
	Retry int `json:"retry"`

	// Required - Expire value
	// Expire value used in SOA record.
	Expire string `json:"expire"`

	// Required - Time to live
	// Default TTL for NS records.
	TTL int `json:"ttl"`

	// Master Name Server
	MasterNS string `json:"master_ns"`

	// IP addresses allowed to initiate domain transfer (DNS NOTIFY).
	NotifyAllowedIPs []string `json:"notify_allowed_ips"`

	// Configured DNS servers (empty means default servers).
	DNSServers []DNSServer `json:"dns_servers"`
}

type Response struct {
	*Definition
	Customer string `json:"customer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
	IsEditable bool `json:"is_editable"`
	ValidationLevel int `json:"validation_level"`
	Revisions []Revision `json:"revisions"`
}

type ResourceRecord struct {
	Name string `json:"name"`
	Type string `json:"Type"`
	Region string `json:"region"`
	RData string `json:"rdata"`
	TTL string `json:"ttl"`
}

type Create struct {
	Name string `json:"zoneName"`
	Master bool `json:"master"`

}

type ChangeSet struct {
	Create ResourceRecord `json:"create"`
	Delete ResourceRecord `json:"delete"`
}

type Import struct {
	ZoneData string `json:"zoneData"`
}

func NewResourceRecord(name, recordType, region, rdata, ttl string) ResourceRecord {
	return ResourceRecord{
		Name:   name,
		Type:   recordType,
		Region: region,
		RData:  rdata,
		TTL:    ttl,
	}
}

// List Zones API method´
func (a api) List(ctx context.Context) ([]Response, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathPrefix,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create zone list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute zone list request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute zone list request, got response %s", httpResponse.Status)
	}

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not decode zone list respone: %w", err)
	}

	return responsePayload.Results, nil
}

// Get zone details API method
func (a api) Get(ctx context.Context, name string) (Response, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		name,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Response{}, fmt.Errorf("could not create zone get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("could not execute zone get request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Response{}, fmt.Errorf("could not execute zone get request, got response %s", httpResponse.Status)
	}

	var responsePayload Response
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Response{}, fmt.Errorf("could not decode zone get respone: %w", err)
	}

	return responsePayload, nil
}

// create
func (a api) Create(ctx context.Context, create Definition) error {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(create); err != nil {
		panic(fmt.Sprintf("could not create request data for create zone: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return fmt.Errorf("could not create zone create request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute zone create request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute zone create request, got response %s", httpResponse.Status)
	}

	// TODO correct to real response from server
	return httpResponse.Body.Close()
}

// update zone
// delete zone

func (a api) Delete(ctx context.Context, name string) error {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		name,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create zone delete request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute zone delete request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute zone delete request, got response %s", httpResponse.Status)
	}

	return httpResponse.Body.Close()
}

// apply (changeset)
func (a api) Apply(ctx context.Context, name string, changeset ChangeSet) error {
	url := fmt.Sprintf(
		"%s%s/%s/changeset",
		a.client.BaseURL(),
		pathPrefix,
		name,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(changeset); err != nil {
		panic(fmt.Sprintf("could not create request data for applying changeset: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return fmt.Errorf("could not create zone changeset request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute zone changeset request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute zone changeset request, got response %s", httpResponse.Status)
	}

	// TODO correct to real response from server
	return httpResponse.Body.Close()
}

// import
func (a api) Import(ctx context.Context, name string, zoneData Import) error {
	url := fmt.Sprintf(
		"%s%s/%s/import",
		a.client.BaseURL(),
		pathPrefix,
		name,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(zoneData); err != nil {
		panic(fmt.Sprintf("could not create request data for import zone request: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return fmt.Errorf("could not create zone import request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute zone import request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute zone import request, got response %s", httpResponse.Status)
	}

	// TODO correct to real response from server
	return httpResponse.Body.Close()
}
