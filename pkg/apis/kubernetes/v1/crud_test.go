package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	corev1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"
	"go.anx.io/go-anxcloud/pkg/apis/internal/gs"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"
	testutils "go.anx.io/go-anxcloud/pkg/utils/test"
)

const (
	mockKubernetesVersion = "1.23.6"
	nonExistingIdentifier = "non-existing-identifier"
)

var _ = Describe("CRUD", Ordered, func() {
	var (
		a   api.API
		srv *ghttp.Server

		clusterIdentifier  string
		nodePoolIdentifier string

		clusterName  = "go-anxcloud-" + testutils.RandomHostname()
		nodePoolName = "go-anxcloud-" + testutils.RandomHostname()
		location     = corev1.Location{Identifier: "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"} // ANX04
	)

	mockStateOK := map[string]interface{}{"id": "0"}
	mockStateError := map[string]interface{}{"id": "1"}
	mockStatePending := map[string]interface{}{"id": "2"}

	BeforeEach(func() {
		srv = ghttp.NewServer()

		var err error

		if isIntegrationTest {
			a, err = api.NewAPI(api.WithClientOptions(client.AuthFromEnv(false)))
			Expect(err).ToNot(HaveOccurred())
		} else {
			a, err = api.NewAPI(api.WithClientOptions(
				client.BaseURL(srv.URL()),
				client.IgnoreMissingToken(),
			))
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Context("Cluster object", Ordered, func() {
		Context("Create operation", Ordered, func() {
			It("Create new cluster", func() {
				appendCreateClusterHandler(srv,
					map[string]interface{}{"name": clusterName, "location": location.Identifier},
					mockedClusterResponse(Cluster{Identifier: "mocked-cluster-identifier", Name: clusterName}),
				)

				cluster := Cluster{Name: clusterName, Location: location}

				err := a.Create(context.TODO(), &cluster)
				Expect(err).ToNot(HaveOccurred())
				Expect(cluster.Identifier).NotTo(BeEmpty())
				clusterIdentifier = cluster.Identifier
				Expect(cluster.StateProgressing()).To(BeTrue())
				Expect(pointer.BoolVal(cluster.NeedsServiceVMs)).To(BeTrue()) // Default if not set on Create
			})
		})

		Context("AwaitCompletion", Ordered, func() {
			It("can wait until state is ready", func() {
				appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStatePending})
				appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStatePending})
				appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStateOK})

				cluster := Cluster{Identifier: clusterIdentifier}
				err := cluster.AwaitCompletion(context.TODO(), a)
				Expect(err).ToNot(HaveOccurred())
				Expect(cluster.StateSuccess()).To(BeTrue())

				if !isIntegrationTest {
					Expect(srv.ReceivedRequests()).To(HaveLen(3))
				}
			})

			It("returns an error when updating the state fails", func() {
				// Updating a Cluster is currently not supported at all
				cluster := Cluster{}
				err := cluster.AwaitCompletion(context.TODO(), a)
				Expect(err).To(MatchError(types.ErrUnidentifiedObject))
			})

			Context("mock test", Ordered, func() {
				BeforeEach(func() {
					if isIntegrationTest {
						Skip("Can't be tested with integration tests")
					}
				})

				It("returns ErrClusterProvisioning when state is Error", func() {
					appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStatePending})
					appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStateError})

					cluster := Cluster{Identifier: clusterIdentifier}
					err := cluster.AwaitCompletion(context.TODO(), a)
					Expect(err).To(MatchError(gs.ErrStateError))
					Expect(cluster.StateFailure()).To(BeTrue())
					Expect(srv.ReceivedRequests()).To(HaveLen(2))
				})

				It("supports Context cancelation", func() {
					appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStatePending})

					ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
					defer cancel()

					cluster := Cluster{Identifier: clusterIdentifier}
					err := cluster.AwaitCompletion(ctx, a)
					Expect(err).To(MatchError(context.DeadlineExceeded))
				})
			})
		})

		Context("Get operation", Ordered, func() {
			It("can Get existing Clusters", func() {
				appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, mockedClusterResponse(Cluster{
					Identifier: clusterIdentifier,
					Name:       clusterName,
				}))

				cluster := Cluster{Identifier: clusterIdentifier}

				err := a.Get(context.TODO(), &cluster)
				Expect(err).ToNot(HaveOccurred())
				Expect(cluster.Name).To(Equal(clusterName))
				Expect(pointer.BoolVal(cluster.NeedsServiceVMs)).To(BeTrue())
			})

			It("returns api.ErrNotFound if cluster does not exist", func() {
				appendGetClusterHandler(srv, nonExistingIdentifier, http.StatusNotFound, nil)

				err := a.Get(context.TODO(), &Cluster{Identifier: nonExistingIdentifier})
				Expect(err).To(MatchError(api.ErrNotFound))
			})
		})

		Context("List operation", Ordered, func() {
			BeforeEach(func() {
				if isIntegrationTest {
					Skip("ENGSUP-6236")
				}
			})

			It("can list existing Clusters", func() {
				appendListClustersHandler(srv,
					partialCluster{"id-0", "name-0"},
					partialCluster{"id-1", "name-1"},
					partialCluster{"id-2", "name-2"},
				)

				for i := 0; i < 3; i++ {
					appendGetClusterHandler(srv, fmt.Sprintf("id-%d", i), http.StatusOK, mockedClusterResponse(Cluster{
						Identifier: fmt.Sprintf("id-%d", i),
						Name:       fmt.Sprintf("name-%d", i),
					}))
				}

				var pi types.PageInfo
				err := a.List(context.TODO(), &Cluster{}, api.Paged(1, 1, &pi), api.FullObjects(true))
				Expect(err).ToNot(HaveOccurred())

				var clusters []Cluster
				Expect(pi.Next(&clusters)).To(BeTrue())
				Expect(pi.Error()).ToNot(HaveOccurred())

				Expect(clusters).To(HaveLen(3))

				Expect(clusters[0].Name).To(Equal("name-0"))
				Expect(clusters[1].Name).To(Equal("name-1"))
				Expect(clusters[2].Name).To(Equal("name-2"))

				Expect(clusters[0].Version).To(Equal(mockKubernetesVersion))
			})
		})

		Context("Update operation", Ordered, func() {
			It("responds with api.ErrOperationNotSupported", func() {
				err := a.Update(context.TODO(), &Cluster{Identifier: clusterIdentifier})
				Expect(err).To(MatchError(api.ErrOperationNotSupported))
			})
		})

	})

	Context("NodePool object", Ordered, func() {
		Context("Create operation", Ordered, func() {
			It("can Create new node pool", func() {
				appendCreateNodePoolHandler(srv,
					map[string]interface{}{
						"name":      nodePoolName,
						"cluster":   clusterIdentifier,
						"replicas":  1,
						"memory":    2 * 1073741824,
						"disk_size": 10 * 1073741824,
					},
					mockedNodePoolResponse(NodePool{
						Identifier: "mocked-cluster-identifier",
						Name:       nodePoolName,
						Cluster:    Cluster{Identifier: clusterIdentifier, Name: clusterName},
					}),
				)

				nodePool := NodePool{
					Name:     nodePoolName,
					Cluster:  Cluster{Identifier: clusterIdentifier},
					Replicas: pointer.Int(1),
					Memory:   2 * 1073741824,  // 2 GiB
					DiskSize: 10 * 1073741824, // 10 GiB
				}

				err := a.Create(context.TODO(), &nodePool)
				Expect(err).ToNot(HaveOccurred())
				Expect(nodePool.Identifier).NotTo(BeEmpty())
				nodePoolIdentifier = nodePool.Identifier
				Expect(nodePool.StateProgressing()).To(BeTrue())

				appendGetNodePoolHandler(srv, nodePoolIdentifier, http.StatusOK, map[string]interface{}{"state": mockStateOK})

				err = nodePool.AwaitCompletion(context.TODO(), a)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Get operation", Ordered, func() {
			It("can Get existing NodePools", func() {
				appendGetNodePoolHandler(srv, nodePoolIdentifier, http.StatusOK, mockedNodePoolResponse(NodePool{
					Identifier: nodePoolIdentifier,
					Name:       nodePoolName,
					Cluster:    Cluster{Identifier: clusterIdentifier, Name: clusterName},
				}))

				nodePool := NodePool{Identifier: nodePoolIdentifier}

				err := a.Get(context.TODO(), &nodePool)
				Expect(err).ToNot(HaveOccurred())
				Expect(nodePool.Name).To(Equal(nodePoolName))
				Expect(nodePool.Memory).To(Equal(2 * 1073741824))
				Expect(nodePool.DiskSize).To(Equal(10 * 1073741824))
			})

			It("returns api.ErrNotFound if node pool does not exist", func() {
				appendGetNodePoolHandler(srv, nonExistingIdentifier, http.StatusNotFound, nil)

				err := a.Get(context.TODO(), &NodePool{Identifier: nonExistingIdentifier})
				Expect(err).To(MatchError(api.ErrNotFound))
			})
		})

		Context("List operation", Ordered, func() {
			BeforeEach(func() {
				if isIntegrationTest {
					Skip("ENGSUP-6236")
				}
			})

			It("can list existing NodePools", func() {
				appendListNodePoolsHandler(srv,
					partialNodePool{"id-0", "name-0"},
					partialNodePool{"id-1", "name-1"},
					partialNodePool{"id-2", "name-2"},
				)

				for i := 0; i < 3; i++ {
					appendGetNodePoolHandler(srv, fmt.Sprintf("id-%d", i), http.StatusOK, mockedNodePoolResponse(NodePool{Identifier: fmt.Sprintf("id-%d", i), Name: fmt.Sprintf("name-%d", i)}))
				}

				var pi types.PageInfo
				err := a.List(context.TODO(), &NodePool{}, api.Paged(1, 1, &pi), api.FullObjects(true))
				Expect(err).ToNot(HaveOccurred())

				var nodePools []NodePool
				Expect(pi.Next(&nodePools)).To(BeTrue())
				Expect(pi.Error()).ToNot(HaveOccurred())

				Expect(nodePools).To(HaveLen(3))

				Expect(nodePools[0].Name).To(Equal("name-0"))
				Expect(nodePools[1].Name).To(Equal("name-1"))
				Expect(nodePools[2].Name).To(Equal("name-2"))
			})
		})

		Context("Update operation", Ordered, func() {
			// Updating a NodePool is currently not supported at all
			It("responds with api.ErrOperationNotSupported", func() {
				err := a.Update(context.TODO(), &NodePool{Identifier: nodePoolIdentifier})
				Expect(err).To(MatchError(api.ErrOperationNotSupported))
			})
		})
	})

	Context("Cluster and NodePool deletion", Ordered, func() {
		It("can destroy existing NodePools", func() {
			appendDeleteNodePoolHandler(srv, nodePoolIdentifier, http.StatusOK)
			appendGetNodePoolHandler(srv, nodePoolIdentifier, http.StatusOK, map[string]interface{}{"state": mockStateOK})

			nodePool := NodePool{Identifier: nodePoolIdentifier}

			err := a.Destroy(context.TODO(), &nodePool)
			Expect(err).ToNot(HaveOccurred())

			err = nodePool.AwaitCompletion(context.TODO(), a)
			Expect(api.IgnoreNotFound(err)).ToNot(HaveOccurred())
		})

		It("returns api.ErrNotFound if node pool does not exist", func() {
			appendDeleteNodePoolHandler(srv, nonExistingIdentifier, http.StatusNotFound)

			err := a.Destroy(context.TODO(), &NodePool{Identifier: nonExistingIdentifier})
			Expect(err).To(MatchError(api.ErrNotFound))
		})

		It("can destroy existing Clusters", func() {
			appendDeleteClusterHandler(srv, clusterIdentifier, http.StatusOK)
			appendGetClusterHandler(srv, clusterIdentifier, http.StatusOK, map[string]interface{}{"state": mockStateOK})

			cluster := Cluster{Identifier: clusterIdentifier}

			err := a.Destroy(context.TODO(), &cluster)
			Expect(err).ToNot(HaveOccurred())

			err = cluster.AwaitCompletion(context.TODO(), a)
			Expect(api.IgnoreNotFound(err)).ToNot(HaveOccurred())
		})

		It("returns api.ErrNotFound if cluster does not exist", func() {
			appendDeleteClusterHandler(srv, nonExistingIdentifier, http.StatusNotFound)

			err := a.Destroy(context.TODO(), &Cluster{Identifier: nonExistingIdentifier})
			Expect(err).To(MatchError(api.ErrNotFound))
		})
	})
})

func mockedClusterResponse(cluster Cluster) map[string]interface{} {
	return map[string]interface{}{
		"identifier": cluster.Identifier,
		"name":       cluster.Name,
		"state": map[string]interface{}{
			"text":          "Pending",
			"id":            "2",
			"type":          2,
			"toStringValue": "Pending",
		},
		"location":          cluster.Location,
		"version":           mockKubernetesVersion,
		"kubeconfig":        nil,
		"needs_service_vms": true,
	}
}

func mockedNodePoolResponse(nodePool NodePool) map[string]interface{} {
	return map[string]interface{}{
		"identifier": nodePool.Identifier,
		"name":       nodePool.Name,
		"state": map[string]interface{}{
			"text":          "Pending",
			"id":            "2",
			"type":          2,
			"toStringValue": "Pending",
		},
		"cluster": map[string]interface{}{
			"identifier": nodePool.Cluster.Identifier,
			"name":       nodePool.Cluster.Name,
		},
		"replicas":         1,
		"cpus":             2,
		"memory":           2147483648,
		"disk_size":        10737418240,
		"operating_system": "Flatcar Linux",
	}
}
