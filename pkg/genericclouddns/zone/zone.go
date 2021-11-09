package zone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/api"
	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

type Record struct {
	Identifier string `json:"identifier,omitempty" anxcloud:"identifier"`
	ZoneName   string
	Immutable  bool   `json:"immutable,omitempty"`
	Name       string `json:"name"`
	RData      string `json:"rdata"`
	Region     string `json:"region"`
	TTL        int    `json:"ttl"`
	Type       string `json:"type"`
}

type Revision struct {
	CreatedAt  time.Time `json:"created_at"`
	Identifier string    `json:"identifier"`
	ModifiedAt time.Time `json:"modified_at"`
	Records    []Record  `json:"records"`
	Serial     int       `json:"serial"`
	State      string    `json:"state"`
}

type DNSServer struct {
	// Required - DNS Server name (FQDN).
	Server string `json:"server"`

	// DNS Server alias
	Alias string
}

type Zone struct {
	// Zone name
	Name string `json:"name,omitempty" anxcloud:"identifier"`

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
	Expire int `json:"expire"`

	// Required - Time to live
	// Default TTL for NS records.
	TTL int `json:"ttl"`

	// Master Name Server
	MasterNS string `json:"master_ns,omitempty"`

	// IP addresses allowed to initiate domain transfer (DNS NOTIFY).
	NotifyAllowedIPs []string `json:"notify_allowed_ips,omitempty"`

	// Configured DNS servers (empty means default servers).
	DNSServers []DNSServer `json:"dns_servers,omitempty"`

	Customer        string     `json:"customer"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	PublishedAt     time.Time  `json:"published_at"`
	IsEditable      bool       `json:"is_editable"`
	ValidationLevel int        `json:"validation_level"`
	DeploymentLevel int        `json:"deployment_level"`
	Revisions       []Revision `json:"revisions"`
	CurrentRevision string     `json:"current_revision,omitempty"`
}

func (z *Zone) EndpointURL(ctx context.Context) (*url.URL, error) {
	u, err := url.ParseRequestURI("/api/clouddns/v1/zone.json/")
	return u, err
}

func (z *Zone) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	// Declare a custom decoder which allows unknown fields - the Zone struct is not modelling all the fields
	d := json.NewDecoder(data)
	return d.Decode(z)
}

func (z *Zone) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// The Update endpoint is NOT at ".../zone.json/{zoneName}", but simply ".../zone.json"
	if op == types.OperationUpdate {
		// Strip the appended zoneName from the URL
		req.URL.Path = path.Dir(req.URL.Path)
	}

	return req, nil
}

func (z *Zone) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// The Create and Update endpoints expect the Zone's name to be in the request body under the key "zoneName"
	if op == types.OperationCreate || op == types.OperationUpdate {
		zWithZoneName := struct {
			Zone
			ZoneName string `json:"zoneName"`
		}{*z, z.Name}

		return zWithZoneName, nil
	}

	return z, nil
}

func (z *Zone) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// CloudDNS API's List response contains some non-functional pagination remnants, which are stripped here
	// Actual array of Zones is in the key 'results'
	if op == types.OperationList {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var m map[string]json.RawMessage
		err = json.Unmarshal(data, &m)
		if err != nil {
			return nil, err
		}

		data = m["results"]
		res.Body = ioutil.NopCloser(bytes.NewReader(data))
		res.ContentLength = int64(len(data))
	}
	return res, nil
}

func (z *Zone) HasPagination(ctx context.Context) (bool, error) {
	return false, nil
}

func (r *Record) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.ParseRequestURI(fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records", r.ZoneName))
	if err != nil {
		return nil, err
	}

	// There is no endpoint to get details of a single record
	if op == types.OperationGet {
		return nil, api.ErrOperationNotSupported
	}

	if op == types.OperationList {
		query := u.Query()

		if r.Name != "" {
			query.Add("name", r.Name)
		}

		if r.RData != "" {
			query.Add("data", r.RData)
		}

		if r.Type != "" {
			query.Add("type", r.Type)
		}
		u.RawQuery = query.Encode()
	}
	return u, err
}

func (r *Record) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	// Response to POST and PUT are the _Zone_ details, which contain some of the updated Record details, but not all
	// To work around these inconsistencies, we just leave the receiver as it is
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return err
	}
	if op == types.OperationCreate || op == types.OperationUpdate {
		return nil
	}

	d := json.NewDecoder(data)
	err = d.Decode(r)
	if err != nil {
		return err
	}

	// Get zoneName from URL and put that into r.ZoneName
	if op == types.OperationList {
		url, err := types.URLFromContext(ctx)
		if err != nil {
			return err
		}
		r.ZoneName = path.Base(path.Dir(url.Path))
	}
	return nil
}

func (r *Record) HasPagination(ctx context.Context) (bool, error) {
	return false, nil
}
