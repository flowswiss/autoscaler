/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flow

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/kubernetes"
)

type nodeGroupClient interface {
	// ListClusterNodes lists all the node found in a Kubernetes cluster.
	ListClusterNodes(ctx context.Context, cursor flowgo.Cursor) (kubernetes.NodeList, error)

	// UpdateCluster updates the details of an existing kubernetes cluster.
	UpdateCluster(ctx context.Context, clusterID int, body kubernetes.ClusterUpdateFlavor) (cluster kubernetes.Cluster, err error)

	// DeleteClusterNode deletes a specific node in a kubernetes cluster.
	DeleteClusterNode(ctx context.Context, nodeID int) error
}

type Client struct {
	clusterID      int
	clusterService kubernetes.ClusterService
	nodeService    kubernetes.NodeService
}

var (
	version = "dev"
)

func newNodeGroupClient(clusterID int, apiToken string, apiURL string) *Client {
	opts := []flowgo.Option{}
	if apiURL != "" {
		opts = append(opts, flowgo.WithBase(apiURL))
	}

	opts = append(opts, flowgo.WithUserAgent("cluster-autoscaler-flow/"+version))
	opts = append(opts, flowgo.WithToken(apiToken))

	doClient := flowgo.NewClient(opts...)
	return &Client{
		clusterService: kubernetes.NewClusterService(doClient),
		nodeService:    kubernetes.NewNodeService(doClient, clusterID),
	}
}

func (c *Client) ListClusterNodes(ctx context.Context, cursor flowgo.Cursor) (kubernetes.NodeList, error) {
	return c.nodeService.List(ctx, cursor)
}

func (c *Client) UpdateCluster(ctx context.Context, clusterID int, req kubernetes.ClusterUpdateFlavor) (cluster kubernetes.Cluster, err error) {
	return c.clusterService.UpdateFlavor(ctx, clusterID, req)
}

func (c *Client) DeleteClusterNode(ctx context.Context, nodeID int) error {
	return c.nodeService.Delete(ctx, nodeID)
}
