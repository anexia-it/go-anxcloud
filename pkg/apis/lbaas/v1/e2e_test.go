//go:build integration
// +build integration

package v1_test

import (
	"fmt"
	"math/rand"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	lbaasv1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"
	testutil "go.anx.io/go-anxcloud/pkg/utils/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test logic is most of the time
//   * create Object with some given parameters
//     - create Object
//     - wait until Object is ready
//     - check if List-ing Objects includes our newly created one (only checking Identifier)
//   * update some parameters of original Object
//     - update the Object
//     - wait until Object is ready
//     - run the optional validate methods
//   * run tests for Objects depending on the currently tested Object
//     - LoadBalancer -> Backend
//     - Backend      -> Frontend
//     - Frontend     -> Bind
//     - Bind         -> Server (depends on Backend, but this way we have a simple chain)
//   * destroy Object
//     - anonymous function returned by createObject helper
//     - deferred call, making sure the Object get's deleted either in test spec
//       (`It("is destroyed successfully",...)`) or in a `DeferCleanup`
//
// Verbose Engine request/response logging is enabled and shown with either -v given to ginkgo
// or when a test fails.
//
// Maybe we can extract some of those e2e helpers for use by other API bindings?
//   -- Mara @LittleFox94 Grosch, 2022-05-03

func ruleChecks(testrun LBaaSE2ETestRun, frontend *lbaasv1.Frontend, acl *lbaasv1.ACL, testURL string) {
	Context("with a fresh Rule", Ordered, func() {
		var rule lbaasv1.Rule

		defer createObject(func() types.Object {
			rule = lbaasv1.Rule{
				Name:          fmt.Sprintf("go-anxcloud-%s", testrun.Name),
				ParentType:    "frontend",
				Index:         pointer.Int(0),
				Frontend:      *frontend,
				Condition:     "if",
				ConditionTest: acl.Name,
				Type:          "connection",
				Action:        "reject",
			}
			return &rule
		}, true)()

		Context("rule blocks port", func() {
			connectionResetByPeerCheck(testURL)
		})

		Context("rule allows port", func() {
			updateObject(func() types.Object {
				rule.Action = "accept"
				return &rule
			}, true)
			successfulConnectionCheck(testURL)
		})
	})
}

func aclChecks(testrun LBaaSE2ETestRun, frontend *lbaasv1.Frontend, testURL string) {
	Context("with a fresh ACL", Ordered, func() {
		var acl lbaasv1.ACL

		defer createObject(func() types.Object {
			acl = lbaasv1.ACL{
				Name:       fmt.Sprintf("go-anxcloud-%s", testrun.Name),
				ParentType: "frontend",
				Index:      pointer.Int(0),
				Criterion:  "dst_port",
				Value:      fmt.Sprintf("%d", testrun.Port),
				Frontend:   *frontend,
			}
			return &acl
		}, true)()

		ruleChecks(testrun, frontend, &acl, testURL)
	})
}

func serverChecks(testrun LBaaSE2ETestRun, backend *lbaasv1.Backend, frontend *lbaasv1.Frontend) {
	Context("with a fresh Server", Ordered, func() {
		var server lbaasv1.Server

		defer createObject(func() types.Object {
			server = lbaasv1.Server{
				Name:    fmt.Sprintf("go-anxcloud-%s", testrun.Name),
				IP:      "127.0.0.1",
				Port:    8080,
				Check:   "enabled",
				Backend: *backend,
			}
			return &server
		}, true)()

		url := fmt.Sprintf("http://go-anxcloud-lbaas-e2e.se.anx.io:%d", testrun.Port)

		Context("correct server port", func() {
			successfulConnectionCheck(url)
		})

		aclChecks(testrun, frontend, url)

		Context("invalid server port", Ordered, func() {
			updateObject(func() types.Object {
				server.Port = 8081
				return &server
			}, true)

			unavailableServerConnectionCheck(url)
		})
	})
}

func bindChecks(testrun LBaaSE2ETestRun, frontend *lbaasv1.Frontend, backend *lbaasv1.Backend) {
	Context("with a fresh Bind", Ordered, func() {
		var bind lbaasv1.Bind

		defer createObject(func() types.Object {
			bind = lbaasv1.Bind{
				Name:     fmt.Sprintf("go-anxcloud-%s", testrun.Name),
				Port:     testrun.Port,
				Frontend: *frontend,
			}
			return &bind
		}, true)()

		serverChecks(testrun, backend, frontend)

		updateObject(func() types.Object {
			bind.Port = testrun.Port + 1
			return &bind
		}, true, func(o types.Object) {
			Expect(o.(*lbaasv1.Bind).Port).To(Equal(testrun.Port + 1))
		})
	})
}

func frontendChecks(testrun LBaaSE2ETestRun, lb *lbaasv1.LoadBalancer, backend *lbaasv1.Backend) {
	Context("with a fresh Frontend", Ordered, func() {
		var frontend lbaasv1.Frontend

		defer createObject(func() types.Object {
			frontend = lbaasv1.Frontend{
				Name:           fmt.Sprintf("go-anxcloud-%s", testrun.Name),
				Mode:           lbaasv1.TCP,
				LoadBalancer:   lb,
				DefaultBackend: backend,
			}
			return &frontend
		}, true)()

		updateObject(func() types.Object {
			frontend.Mode = lbaasv1.HTTP
			return &frontend
		}, true, func(o types.Object) {
			Expect(o.(*lbaasv1.Frontend).Mode).To(Equal(lbaasv1.HTTP))
		})

		bindChecks(testrun, &frontend, backend)
	})
}

func backendChecks(testrun LBaaSE2ETestRun, lb *lbaasv1.LoadBalancer) {
	Context("with a fresh Backend", Ordered, func() {
		var backend lbaasv1.Backend

		// create Backend instance in test execution phase as lb is not filled with an Identifier before
		defer createObject(func() types.Object {
			backend = lbaasv1.Backend{
				Name:         fmt.Sprintf("go-anxcloud-%s", testrun.Name),
				Mode:         lbaasv1.TCP,
				LoadBalancer: *lb,
			}
			return &backend
		}, true)()

		updateObject(func() types.Object {
			backend.Mode = lbaasv1.HTTP
			return &backend
		}, true, func(o types.Object) {
			Expect(o.(*lbaasv1.Backend).Mode).To(Equal(lbaasv1.HTTP))
		})

		frontendChecks(testrun, lb, &backend)
	})
}

func loadbalancerChecks(testrun LBaaSE2ETestRun) {
	Context("with a fresh LoadBalancer", Ordered, func() {
		var lb lbaasv1.LoadBalancer

		defer createObject(func() types.Object {
			lb = lbaasv1.LoadBalancer{
				Name:      fmt.Sprintf("go-anxcloud-%s", testrun.Name),
				IpAddress: "go-anxcloud-lbaas-e2e.se.anx.io",
			}
			return &lb
		}, false)()

		updateObject(func() types.Object {
			lb.Name += " (updated)"
			return &lb
		}, false, func(o types.Object) {
			Expect(o.(*lbaasv1.LoadBalancer).Name).To(HaveSuffix(" (updated)"))
		})

		backendChecks(testrun, &lb)
	})
}

var _ = Describe("lbaas/v1 bindings", Ordered, func() {
	if ac, err := api.NewAPI(
		api.WithClientOptions(
			client.AuthFromEnv(false),
		),
		api.WithLogger(testutil.NewGinkgor()),
	); err != nil {
		panic(fmt.Sprintf("error creating API client: %v", err))
	} else {
		apiClient = ac
	}

	rand.Seed(GinkgoRandomSeed())

	testrun := LBaaSE2ETestRun{
		Name: testutil.RandomHostname(),

		// there might come a time where we have to check if a port is already in use
		// and throw another set of dice.
		Port: 32000 + rand.Intn(1000),
	}

	loadbalancerChecks(testrun)
})
