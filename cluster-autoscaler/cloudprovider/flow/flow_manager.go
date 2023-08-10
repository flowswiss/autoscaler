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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/kubernetes"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"k8s.io/klog/v2"
)

// Manager handles Flow communication and data caching of
// node groups
type Manager struct {
	client        nodeGroupClient
	clusterID     int
	nodeGroups    []*NodeGroup
	discoveryOpts cloudprovider.NodeGroupDiscoveryOptions
}

// Config is the configuration of the Flow cloud provider
type Config struct {
	// ClusterID is the id associated with the cluster where Flow
	// Cluster Autoscaler is running.
	ClusterID int `json:"cluster_id" yaml:"cluster_id"`

	// Token is the User's Access Token associated with the cluster where
	// Flow Cluster Autoscaler is running.
	ApiToken string `json:"api_token" yaml:"api_token"`

	// URL points to Flow API. If empty, defaults to
	// https://api.flow.ch/
	ApiURL string `json:"api_url" yaml:"api_url"`
}

func newManager(configReader io.Reader, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions) (*Manager, error) {
	cfg := &Config{}
	if configReader != nil {
		body, err := io.ReadAll(configReader)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, cfg)
		if err != nil {
			return nil, err
		}
	} else {
		cfg.ApiURL = os.Getenv("FLOW_API_URL")
		cfg.ApiToken = os.Getenv("FLOW_API_TOKEN")

		clusterID, err := strconv.Atoi(os.Getenv("FLOW_CLUSTER_ID"))
		if err != nil {
			return nil, err
		}
		cfg.ClusterID = clusterID
	}

	if cfg.ApiToken == "" {
		return nil, errors.New("flow access token is not provided")
	}
	if cfg.ClusterID == 0 {
		return nil, errors.New("flow cluster ID is not provided")
	}

	m := &Manager{
		client:        newNodeGroupClient(cfg.ClusterID, cfg.ApiToken, cfg.ApiURL),
		clusterID:     cfg.ClusterID,
		nodeGroups:    make([]*NodeGroup, 0),
		discoveryOpts: discoveryOpts,
	}

	return m, nil
}

// Refresh refreshes the cache holding the nodegroups. This is called by the CA
// based on the `--scan-interval`. By default it's 10 seconds.
func (m *Manager) Refresh() error {
	var (
		minSize int
		maxSize int
	)

	klog.V(4).Infof("refreshing workers node group kubernetes cluster: %q", m.clusterID)

	for _, specString := range m.discoveryOpts.NodeGroupSpecs {
		spec, err := dynamic.SpecFromString(specString, true)
		if err != nil {
			return fmt.Errorf("failed to parse node group spec: %v", err)
		}

		if spec.Name == "workers" {
			minSize = spec.MinSize
			maxSize = spec.MaxSize

			klog.V(4).Infof("found configuration for workers node group: min: %d max: %d", minSize, maxSize)
		}
	}

	ctx := context.Background()
	nodeList, err := m.client.ListClusterNodes(ctx, flowgo.Cursor{NoFilter: 1})
	if err != nil {
		return fmt.Errorf("couldn't list Kubernetes cluster pools: %s", err)
	}

	var workerNodes []kubernetes.Node
	for _, node := range nodeList.Items {
		if isControlPlaneNode(node) {
			continue
		}
		workerNodes = append(workerNodes, node)
	}

	var group []*NodeGroup
	group = append(group, &NodeGroup{
		id:        1,
		clusterID: m.clusterID,
		client:    m.client,
		nodes:     workerNodes,
		minSize:   minSize,
		maxSize:   maxSize,
	})

	if len(group) == 0 {
		klog.V(4).Info("cluster-autoscaler is disabled. no node pools are configured")
	}

	m.nodeGroups = group
	return nil
}

func isControlPlaneNode(node kubernetes.Node) bool {
	for _, role := range node.Roles {
		if role.Key == "control-plane" {
			return true
		}
	}
	return false
}
