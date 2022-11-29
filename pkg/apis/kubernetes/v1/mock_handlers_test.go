package v1

import (
	"fmt"

	"github.com/onsi/gomega/ghttp"
)

func withExistingServer(srv *ghttp.Server, cb func(srv *ghttp.Server)) {
	if srv != nil {
		cb(srv)
	}
}

// Cluster handlers

func appendCreateClusterHandler(srv *ghttp.Server, req, res interface{}) {
	withExistingServer(srv, func(srv *ghttp.Server) {
		srv.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/api/kubernetes/v1/cluster.json"),
			ghttp.VerifyJSONRepresenting(req),
			ghttp.RespondWithJSONEncoded(200, res),
		))
	})
}

func appendGetClusterHandler(srv *ghttp.Server, clusterID string, resCode int, res interface{}) {
	withExistingServer(srv, func(srv *ghttp.Server) {
		srv.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("/api/kubernetes/v1/cluster.json/%s", clusterID)),
			ghttp.RespondWithJSONEncoded(resCode, res),
		))
	})
}

func appendDeleteClusterHandler(srv *ghttp.Server, clusterID string, resCode int) {
	withExistingServer(srv, func(srv *ghttp.Server) {
		srv.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("DELETE", fmt.Sprintf("/api/kubernetes/v1/cluster.json/%s", clusterID)),
			ghttp.RespondWithJSONEncoded(resCode, true),
		))
	})
}

type partialCluster struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

func appendListClustersHandler(srv *ghttp.Server, clusters ...partialCluster) {
	withExistingServer(srv, func(srv *ghttp.Server) {
		srv.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/api/kubernetes/v1/cluster.json"),
			ghttp.RespondWithJSONEncoded(200, map[string]interface{}{
				"data": map[string]interface{}{
					"page":        1,
					"total_items": len(clusters),
					"limit":       100,
					"data":        clusters,
				},
			}),
		))
	})
}

// NodePool handlers

func appendCreateNodePoolHandler(srv *ghttp.Server, req, res interface{}) {
	withExistingServer(srv, func(srv *ghttp.Server) {
		srv.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/api/kubernetes/v1/node_pool.json"),
			ghttp.VerifyJSONRepresenting(req),
			ghttp.RespondWithJSONEncoded(200, res),
		))
	})
}

func appendGetNodePoolHandler(srv *ghttp.Server, nodePoolID string, resCode int, res interface{}) {
	withExistingServer(srv, func(srv *ghttp.Server) {
		srv.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("/api/kubernetes/v1/node_pool.json/%s", nodePoolID)),
			ghttp.RespondWithJSONEncoded(resCode, res),
		))
	})
}

func appendDeleteNodePoolHandler(srv *ghttp.Server, nodePoolID string, resCode int) {
	withExistingServer(srv, func(srv *ghttp.Server) {
		srv.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("DELETE", fmt.Sprintf("/api/kubernetes/v1/node_pool.json/%s", nodePoolID)),
			ghttp.RespondWithJSONEncoded(resCode, true),
		))
	})
}

type partialNodePool struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

func appendListNodePoolsHandler(srv *ghttp.Server, nodePools ...partialNodePool) {
	withExistingServer(srv, func(srv *ghttp.Server) {
		srv.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/api/kubernetes/v1/node_pool.json"),
			ghttp.RespondWithJSONEncoded(200, map[string]interface{}{
				"data": map[string]interface{}{
					"page":        1,
					"total_items": len(nodePools),
					"limit":       100,
					"data":        nodePools,
				},
			}),
		))
	})
}
