package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"

	"github.com/onsi/gomega/ghttp"

	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
)

type backend v1.Backend
type state gs.State

func NewMockServer() *ghttp.Server {
	type errorResponseData struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	type errorResponseType struct {
		E errorResponseData `json:"error"`
	}

	errorResponse := func(code int, msg string, res http.ResponseWriter) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")
		res.WriteHeader(code)

		if msg == "" {
			msg = http.StatusText(code)
		}

		_ = json.NewEncoder(res).Encode(errorResponseType{
			E: errorResponseData{
				Message: msg,
				Code:    code,
			},
		})
	}

	server := ghttp.NewServer()
	server.SetAllowUnhandledRequests(true)

	backends := []v1.Backend{
		{
			Name:       "Example-Backend",
			Identifier: "bogus identifier 1",
			Mode:       v1.TCP,
			LoadBalancer: v1.LoadBalancer{
				Identifier: "bogus identifier 2",
			},
		},
		{
			Name:       "backend-01",
			Identifier: "bogus identifier 3",
			Mode:       v1.TCP,
		},
		{
			Name:       "test-backend-01",
			Identifier: "test identifier 1",
			Mode:       v1.TCP,
		},
		{
			Name:       "test-backend-02",
			Identifier: "test identifier 2",
			Mode:       v1.TCP,
			LoadBalancer: v1.LoadBalancer{
				Identifier: "bogus identifier 2",
			},
		},
		{
			Name:       "test-backend-03",
			Identifier: "test identifier 3",
			Mode:       v1.TCP,
		},
		{
			Name:       "test-backend-04",
			Identifier: "test identifier 4",
			Mode:       v1.TCP,
			LoadBalancer: v1.LoadBalancer{
				Identifier: "bogus identifier 2",
			},
		},
	}

	const backendBasePath = "/api/LBaaS/v1/backend.json"
	singleBackendPath, _ := regexp.Compile(backendBasePath + "/.*")

	server.RouteToHandler("GET", backendBasePath, func(res http.ResponseWriter, req *http.Request) {
		var lb_filter *string

		if filters, err := url.ParseQuery(req.URL.Query().Get("filters")); err != nil {
			errorResponse(500, "invalid filter parameter", res)
			return
		} else {
			if lb := filters.Get("load_balancer"); lb != "" {
				lb_filter = &lb
			}
		}

		page := 1
		limit := 0

		if p := req.URL.Query().Get("page"); p != "" {
			if pp, err := strconv.ParseInt(p, 10, 32); err != nil || pp <= 0 {
				errorResponse(500, "invalid page parameter", res)
				return
			} else {
				page = int(pp)
			}
		}

		if l := req.URL.Query().Get("limit"); l != "" {
			if pl, err := strconv.ParseInt(l, 10, 32); err != nil || pl < 0 {
				errorResponse(500, "invalid limit parameter", res)
				return
			} else {
				limit = int(pl)
			}
		}

		ret := make([]map[string]string, 0, len(backends))

		for _, b := range backends {
			if lb_filter == nil || b.LoadBalancer.Identifier == *lb_filter {
				ret = append(ret, map[string]string{
					"name":       b.Name,
					"identifier": b.Identifier,
				})
			}
		}

		if limit > 0 {
			idxStart := (page - 1) * limit
			idxEnd := idxStart + limit

			if idxStart >= len(ret) {
				ret = []map[string]string{}
			} else {
				if idxEnd > len(ret) {
					idxEnd = len(ret)
				}

				ret = ret[idxStart:idxEnd]
			}
		}

		res.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(res).Encode(ret)
	})

	identifierGenerateCounter := 1

	server.RouteToHandler("POST", backendBasePath, func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(res).Encode(
			backend{
				Name:       "backend-01",
				Identifier: "generated identifier " + strconv.Itoa(identifierGenerateCounter),
				Mode:       v1.TCP,
			},
		)

		identifierGenerateCounter++
	})

	server.RouteToHandler("GET", singleBackendPath, func(res http.ResponseWriter, req *http.Request) {

		identifier := path.Base(req.URL.Path)

		for _, b := range backends {
			if b.Identifier == identifier {
				res.Header().Add("Content-Type", "application/json; charset=utf-8")
				_ = json.NewEncoder(res).Encode(backend(b))
				return
			}
		}

		errorResponse(404, "", res)
	})

	server.RouteToHandler("PUT", singleBackendPath, func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("content-type") != "application/json; charset=utf-8" {
			errorResponse(500, "Content-Type header on request not set", res)
			return
		}

		var update backend
		if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
			fmt.Printf("Invalid request body: %v\n", err)
			errorResponse(500, fmt.Sprintf("invalid request body: %v", err), res)
			return
		}

		identifier := path.Base(req.URL.Path)

		newBackends := make([]v1.Backend, 0, len(backends))
		found := false

		for _, b := range backends {
			if b.Identifier == identifier {
				newBackends = append(newBackends, v1.Backend(update))
				found = true
			} else {
				newBackends = append(newBackends, b)
			}
		}

		backends = newBackends

		if !found {
			errorResponse(404, "", res)
		} else {
			res.Header().Add("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(res).Encode(update)
		}
	})

	server.RouteToHandler("DELETE", singleBackendPath, func(res http.ResponseWriter, req *http.Request) {
		identifier := path.Base(req.URL.Path)

		newBackends := make([]v1.Backend, 0, len(backends))

		found := false
		for _, b := range backends {
			if b.Identifier != identifier {
				newBackends = append(newBackends, b)
			} else {
				found = true
			}
		}

		backends = newBackends

		if !found {
			errorResponse(404, "", res)
		} else {
			res.Header().Add("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(res).Encode(map[string]interface{}{})
		}
	})

	return server
}

func (b *backend) UnmarshalJSON(bytes []byte) error {
	var clientData struct {
		v1.Backend
		LoadBalancer string `json:"load_balancer,omitempty"`
		State        string `json:"state"`
	}

	err := json.Unmarshal(bytes, &clientData)
	if err != nil {
		return err
	}
	clientData.Backend.State = getStateByID(clientData.State)
	clientData.Backend.LoadBalancer.Identifier = clientData.LoadBalancer

	*b = backend(clientData.Backend)
	return nil
}

func (b backend) MarshalJSON() ([]byte, error) {
	clientData := struct {
		v1.Backend
		State        state                  `json:"state"`
		LoadBalancer map[string]interface{} `json:"load_balancer"`
	}{
		Backend: v1.Backend(b),
		State:   state(getStateByID(b.State.ID)),
		LoadBalancer: map[string]interface{}{
			"name":       b.LoadBalancer.Name,
			"identifier": b.LoadBalancer.Name,
		},
	}

	return json.Marshal(clientData)
}

// MarshalJSON overwrites the original MarshalJSON from lbaasv1.State so that we can properly use it in the mocked server
func (s state) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":   s.ID,
		"type": s.Type,
		"text": s.Text,
	})
}

func getStateByID(stateID string) gs.State {
	switch stateID {
	case v1.NewlyCreated.ID, "":
		return v1.NewlyCreated
	case v1.Updating.ID:
		return v1.Updating
	case v1.Updated.ID:
		return v1.Updated
	case v1.Deployed.ID:
		return v1.Deployed
	case v1.DeploymentError.ID:
		return v1.DeploymentError
	default:
		panic(fmt.Sprintf("unknown id '%s'", stateID))
	}
}
