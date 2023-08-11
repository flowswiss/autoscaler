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
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/kubernetes"
)

type flowClientMock struct {
	mock.Mock
}

func (f *flowClientMock) ListClusterNodes(ctx context.Context, cursor flowgo.Cursor) (kubernetes.NodeList, error) {
	args := f.Called(ctx, cursor)
	return args.Get(0).(kubernetes.NodeList), args.Error(1)
}

func (f *flowClientMock) UpdateCluster(ctx context.Context, clusterID int, req kubernetes.ClusterUpdateFlavor) (cluster kubernetes.Cluster, err error) {
	args := f.Called(ctx, clusterID, req)
	return args.Get(0).(kubernetes.Cluster), args.Error(1)
}

func (f *flowClientMock) DeleteClusterNode(ctx context.Context, nodeID int) error {
	args := f.Called(ctx, nodeID)
	return args.Error(0)
}

func testCloudProvider(t *testing.T, client *flowClientMock) *flowCloudProvider {
	cfg := `{"cluster_id": 123456, "api_token": "123-123-123", "api_url": "https://api.flow.ch/"}`

	nodeGroupSpecs := []string{"3:10:workers"}
	nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

	manager, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
	assert.NoError(t, err)
	rl := &cloudprovider.ResourceLimiter{}

	// fill the test provider with some example
	if client == nil {
		client = &flowClientMock{}

		client.On("ListClusterNodes",
			context.Background(),
			flowgo.Cursor{NoFilter: 1},
		).Return(
			kubernetes.NodeList{
				Items: []kubernetes.Node{
					{
						ID:   1,
						Name: "worker1",
						Status: kubernetes.NodeStatus{
							ID:   1,
							Key:  "healthy",
							Name: "Healthy",
						},
					},
					{
						ID:   2,
						Name: "worker2",
						Status: kubernetes.NodeStatus{
							ID:   1,
							Key:  "healthy",
							Name: "Healthy",
						},
					},
					{
						ID:   3,
						Name: "worker3",
						Status: kubernetes.NodeStatus{
							ID:   2,
							Key:  "creating",
							Name: "Creating",
						},
					},
				},
			},
			nil,
		).Once()
	}

	manager.client = client

	provider, err := newFlowCloudProvider(manager, rl)
	assert.NoError(t, err)
	return provider
}

func TestNewFlowCloudProvider(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_ = testCloudProvider(t, nil)
	})
}

func TestFlowCloudProvider_Name(t *testing.T) {
	provider := testCloudProvider(t, nil)

	t.Run("success", func(t *testing.T) {
		name := provider.Name()
		assert.Equal(t, cloudprovider.FlowProviderName, name, "provider name doesn't match")
	})
}

func TestFlowCloudProvider_NodeGroups(t *testing.T) {
	provider := testCloudProvider(t, nil)

	t.Run("success", func(t *testing.T) {
		nodegroups := provider.NodeGroups()
		assert.Equal(t, len(nodegroups), 1, "number of node groups does not match")
		nodes, _ := nodegroups[0].Nodes()
		assert.Equal(t, len(nodes), 3, "number of nodes in workers node group does not match")

	})

	t.Run("zero groups", func(t *testing.T) {
		provider.manager.nodeGroups = []*NodeGroup{}
		nodes := provider.NodeGroups()
		assert.Equal(t, len(nodes), 0, "number of nodes do not match")
	})
}

func TestFlowCloudProvider_NodeGroupForNode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &flowClientMock{}

		client.On("ListClusterNodes",
			context.Background(),
			flowgo.Cursor{NoFilter: 1},
		).Return(
			kubernetes.NodeList{
				Items: []kubernetes.Node{
					{
						ID:   11,
						Name: "worker11",
						Status: kubernetes.NodeStatus{
							ID:   1,
							Key:  "healthy",
							Name: "Healthy",
						},
					},
					{
						ID:   22,
						Name: "worker22",
						Status: kubernetes.NodeStatus{
							ID:   4,
							Key:  "draining",
							Name: "Draining",
						},
					},
				},
			},
			nil,
		).Once()

		provider := testCloudProvider(t, client)

		// let's get the nodeGroup for the node with ID 11
		node := &apiv1.Node{
			Spec: apiv1.NodeSpec{
				ProviderID: "flow://11",
			},
		}

		nodeGroup, err := provider.NodeGroupForNode(node)
		require.NoError(t, err)
		require.NotNil(t, nodeGroup)
		require.Equal(t, nodeGroup.Id(), "1", "node group ID does not match")
	})

	t.Run("node does not exist", func(t *testing.T) {
		client := &flowClientMock{}

		client.On("ListClusterNodes",
			context.Background(),
			flowgo.Cursor{NoFilter: 1},
		).Return(
			kubernetes.NodeList{
				Items: []kubernetes.Node{
					{
						ID:   11,
						Name: "worker11",
						Status: kubernetes.NodeStatus{
							ID:   1,
							Key:  "healthy",
							Name: "Healthy",
						},
					},
					{
						ID:   22,
						Name: "worker22",
						Status: kubernetes.NodeStatus{
							ID:   4,
							Key:  "draining",
							Name: "Draining",
						},
					},
				},
			},
			nil,
		).Once()

		provider := testCloudProvider(t, client)

		// let's get the nodeGroup for the node with ID 11
		node := &apiv1.Node{
			Spec: apiv1.NodeSpec{
				ProviderID: "flow://33",
			},
		}

		nodeGroup, err := provider.NodeGroupForNode(node)
		require.NoError(t, err)
		require.Nil(t, nodeGroup)
	})
}
