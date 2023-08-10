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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/kubernetes"
)

func TestNewManager(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := `{"cluster_id": 123456, "api_token": "123-123-123", "api_url": "https://api.flow.ch/"}`

		nodeGroupSpecs := []string{"1:10:workers"}
		nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

		manager, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
		assert.NoError(t, err)
		assert.Equal(t, 123456, manager.clusterID, "cluster ID does not match")
		assert.Equal(t, nodeGroupDiscoveryOptions, manager.discoveryOpts, "node group discovery options do not match")
	})

	t.Run("empty api_token", func(t *testing.T) {
		cfg := `{"cluster_id": 123456, "api_token": "", "api_url": "https://api.flow.ch/"}`

		nodeGroupSpecs := []string{"1:10:workers"}
		nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

		_, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
		assert.EqualError(t, err, errors.New("flow access token is not provided").Error())
	})

	t.Run("empty cluster ID", func(t *testing.T) {
		cfg := `{"api_token": "123-123-123", "api_url": "https://api.flow.ch/"}`

		nodeGroupSpecs := []string{"1:10:workers"}
		nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

		_, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
		assert.EqualError(t, err, errors.New("flow cluster ID is not provided").Error())
	})
}

func TestFlowManager_Refresh(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := `{"cluster_id": 123456, "api_token": "123-123-123", "api_url": "https://api.flow.ch/"}`

		nodeGroupSpecs := []string{"3:10:workers"}
		nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

		manager, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
		assert.NoError(t, err)

		client := &flowClientMock{}
		ctx := context.Background()
		cursor := flowgo.Cursor{NoFilter: 1}

		client.On("ListClusterNodes", ctx, cursor).Return(
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
						ID:   1,
						Name: "worker2",
						Status: kubernetes.NodeStatus{
							ID:   1,
							Key:  "healthy",
							Name: "Healthy",
						},
					},
				},
			},
			nil,
		).Once()

		manager.client = client
		err = manager.Refresh()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(manager.nodeGroups), "number of node groups do not match")
	})
}

func TestFlowManager_RefreshWithNodeSpec(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := `{"cluster_id": 123456, "api_token": "123-123-123", "api_url": "https://api.flow.ch/"}`

		nodeGroupSpecs := []string{"3:10:workers"}
		nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

		manager, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
		assert.NoError(t, err)

		client := &flowClientMock{}
		ctx := context.Background()
		cursor := flowgo.Cursor{NoFilter: 1}

		client.On("ListClusterNodes", ctx, cursor).Return(
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
						ID:   1,
						Name: "worker2",
						Status: kubernetes.NodeStatus{
							ID:   1,
							Key:  "healthy",
							Name: "Healthy",
						},
					},
				},
			},
			nil,
		).Once()

		manager.client = client
		err = manager.Refresh()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(manager.nodeGroups), "number of node groups do not match")
		assert.Equal(t, 3, manager.nodeGroups[0].minSize, "minimum node for node group does not match")
		assert.Equal(t, 10, manager.nodeGroups[0].maxSize, "maximum node for node group does not match")
	})
}
