package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"

	"github.com/onsi/gomega/ghttp"

	lbaasv1 "github.com/anexia-it/go-anxcloud/pkg/apis/lbaas/v1"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
)

func newMockServer() *ghttp.Server {
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

	backends := []lbaasv1.Backend{
		{
			Name:       "Example-Backend",
			Identifier: "bogus identifier 1",
			Mode:       lbaasv1.TCP,
			LoadBalancer: lbaasv1.LoadBalancer{
				Identifier: "bogus identifier 2",
			},
		},
		{
			Name:       "backend-01",
			Identifier: "bogus identifier 3",
			Mode:       lbaasv1.TCP,
		},
		{
			Name:       "test-backend-01",
			Identifier: "test identifier 1",
			Mode:       lbaasv1.TCP,
		},
		{
			Name:       "test-backend-02",
			Identifier: "test identifier 2",
			Mode:       lbaasv1.TCP,
			LoadBalancer: lbaasv1.LoadBalancer{
				Identifier: "bogus identifier 2",
			},
		},
		{
			Name:       "test-backend-03",
			Identifier: "test identifier 3",
			Mode:       lbaasv1.TCP,
		},
		{
			Name:       "test-backend-04",
			Identifier: "test identifier 4",
			Mode:       lbaasv1.TCP,
			LoadBalancer: lbaasv1.LoadBalancer{
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

		ret := make([]backend.Backend, 0, len(backends))

		for _, b := range backends {
			if lb_filter == nil || b.LoadBalancer.Identifier == *lb_filter {
				ret = append(ret, b)
			}
		}

		if limit > 0 {
			idxStart := (page - 1) * limit
			idxEnd := idxStart + limit

			if idxStart >= len(ret) {
				ret = make([]backend.Backend, 0)
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
			backend.Backend{
				Name:       "backend-01",
				Identifier: "generated identifier " + strconv.Itoa(identifierGenerateCounter),
				Mode:       lbaasv1.TCP,
			},
		)

		identifierGenerateCounter++
	})

	server.RouteToHandler("GET", singleBackendPath, func(res http.ResponseWriter, req *http.Request) {

		identifier := path.Base(req.URL.Path)

		for _, b := range backends {
			if b.Identifier == identifier {
				res.Header().Add("Content-Type", "application/json; charset=utf-8")
				_ = json.NewEncoder(res).Encode(b)
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

		update := struct {
			lbaasv1.Backend
			LoadBalancer string `json:"load_balancer"`
		}{}
		if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
			fmt.Printf("Invalid request body: %v\n", err)
			errorResponse(500, fmt.Sprintf("invalid request body: %v", err), res)
			return
		}

		identifier := path.Base(req.URL.Path)

		newBackends := make([]backend.Backend, 0, len(backends))
		found := false

		for _, b := range backends {
			if b.Identifier == identifier {
				update.Backend.LoadBalancer.Identifier = update.LoadBalancer
				newBackends = append(newBackends, update.Backend)
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
			_ = json.NewEncoder(res).Encode(update.Backend)
		}
	})

	server.RouteToHandler("DELETE", singleBackendPath, func(res http.ResponseWriter, req *http.Request) {
		identifier := path.Base(req.URL.Path)

		newBackends := make([]backend.Backend, 0, len(backends))

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
