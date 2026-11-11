package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/utils/object/compare"
)

var (
	// ErrModifyRevisionNotFound is returned for Create and Update requests when the zones CurrentRevision is not
	// found in the set of Revisions. This is probably an Engine problem and not your code.
	ErrModifyRevisionNotFound = errors.New("revision not found")

	// ErrModifyRecordNotFound is returned for Create and Update requests when the modified Record is not found in
	// the zones current Revision. This is probably an Engine problem and not your code, but might be a problem in
	// these API bindings.
	ErrModifyRecordNotFound = errors.New("record not found")

	// ErrEmptyRecordNameNotSupported is returned for Create requests when the record has no name set.
	// Instead "@" must be used to target the domain root.
	ErrEmptyRecordNameNotSupported = errors.New("empty record name not supported - use \"@\" instead")

	// ErrInvalidZoneName is returned when a Record's ZoneName contains a "/", which cannot be sent as a
	// single path segment.
	ErrInvalidZoneName = errors.New("zone name must not contain \"/\"")
)

func (r *Record) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if strings.Contains(r.ZoneName, "/") {
		return nil, fmt.Errorf("%w: %q", ErrInvalidZoneName, r.ZoneName)
	}

	// Set as the decoded path so URL encoding escapes characters like "?" and "%" instead of
	// reading them as URL syntax, and DecodeAPIResponse reads the same ZoneName back.
	u := &url.URL{Path: "/api/clouddns/v1/zone.json/" + r.ZoneName + "/records"}

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
	return u, nil
}

func (r *Record) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return err
	}

	// Response to POST and PUT are the _Zone_ details, which contain some of the updated Record details, but not all.
	// We have to find our Record in the current revision of the zone and grab some values from it, notably the Identifier
	// and RData (as the Engine might change its format).
	if op == types.OperationCreate || op == types.OperationUpdate {
		var zone Zone
		err := json.NewDecoder(data).Decode(&zone)
		if err != nil {
			return nil
		}

		responseRecord, err := r.findInZone(&zone)
		if err != nil {
			return err
		}

		r.Identifier = responseRecord.Identifier

		// Engine changes RData sometimes
		r.RData = responseRecord.RData

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

// FilterAPIRequest checks whether a name was set on create operations
// and returns ErrEmptyRecordNameNotSupported if not.
func (r *Record) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate && r.Name == "" {
		return nil, ErrEmptyRecordNameNotSupported
	}

	return req, nil
}

func (r *Record) findInZone(zone *Zone) (*Record, error) {
	var rev *Revision
	for _, r := range zone.Revisions {
		if r.Identifier == zone.CurrentRevision {
			rev = &r
		}
	}

	if rev == nil {
		return nil, ErrModifyRevisionNotFound
	}

	changedR := *r

	if r.Type == "TXT" {
		// Engine returns TXT RData enclosed in quotes
		changedR.RData = fmt.Sprintf("%q", r.RData)
	}

	idx, err := compare.Search(changedR, rev.Records, "Name", "Type", "RData", "TTL")
	if err != nil {
		return nil, fmt.Errorf("error searching record in response: %w", err)
	} else if idx == -1 {
		return nil, ErrModifyRecordNotFound
	}

	return &rev.Records[idx], nil
}
