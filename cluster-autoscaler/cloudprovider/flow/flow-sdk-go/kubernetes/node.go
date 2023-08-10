package kubernetes

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/common"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/compute"
)

type Node struct {
	ID      int                             `json:"id"`
	Name    string                          `json:"name"`
	Roles   []NodeRole                      `json:"roles"`
	Product common.Product                  `json:"product"`
	Network compute.ServerNetworkAttachment `json:"network"`
	Status  NodeStatus                      `json:"status"`
}

type NodeList struct {
	Items      []Node
	Pagination flowgo.Pagination
}

type NodeRole struct {
	ID   int    `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type NodeStatus struct {
	ID      int          `json:"id"`
	Key     string       `json:"key"`
	Name    string       `json:"name"`
	Actions []NodeAction `json:"actions"`
}

type NodeAction struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Command string `json:"command"`
	Sorting int    `json:"sorting"`
}

type NodePerformAction struct {
	Action string `json:"action"`
}

type NodeService struct {
	client    flowgo.Client
	clusterID int
}

func NewNodeService(client flowgo.Client, clusterID int) NodeService {
	return NodeService{
		client:    client,
		clusterID: clusterID,
	}
}

func (n NodeService) List(ctx context.Context, cursor flowgo.Cursor) (list NodeList, err error) {
	list.Pagination, err = n.client.List(ctx, getNodePath(n.clusterID), cursor, &list.Items)
	return
}

func (n NodeService) Delete(ctx context.Context, id int) (err error) {
	err = n.client.Delete(ctx, getSpecificNodePath(n.clusterID, id))
	return
}

func (n NodeService) PerformAction(ctx context.Context, id int, body NodePerformAction) (node Node, err error) {
	err = n.client.Create(ctx, getNodeActionPath(n.clusterID, id), body, &node)
	return
}

const (
	nodeSegment       = "nodes"
	nodeActionSegment = "action"
)

func getNodePath(clusterID int) string {
	return flowgo.Join(clusterSegment, clusterID, nodeSegment)
}

func getSpecificNodePath(clusterID, nodeID int) string {
	return flowgo.Join(clusterSegment, clusterID, nodeSegment, nodeID)
}

func getNodeActionPath(clusterID, nodeID int) string {
	return flowgo.Join(clusterSegment, clusterID, nodeSegment, nodeID, nodeActionSegment)
}
