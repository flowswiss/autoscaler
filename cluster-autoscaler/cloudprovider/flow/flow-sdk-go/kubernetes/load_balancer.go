package kubernetes

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/compute"
)

type LoadBalancer = compute.LoadBalancer
type LoadBalancerList = compute.LoadBalancerList

type LoadBalancerService struct {
	client    flowgo.Client
	clusterID int
}

func NewLoadBalancerService(client flowgo.Client, clusterID int) LoadBalancerService {
	return LoadBalancerService{
		client:    client,
		clusterID: clusterID,
	}
}

func (v LoadBalancerService) List(ctx context.Context, cursor flowgo.Cursor) (list LoadBalancerList, err error) {
	list.Pagination, err = v.client.List(ctx, getLoadBalancerPath(v.clusterID), cursor, &list.Items)
	return
}

const loadBalancerSegment = "load-balancers"

func getLoadBalancerPath(clusterID int) string {
	return flowgo.Join(clusterSegment, clusterID, loadBalancerSegment)
}
