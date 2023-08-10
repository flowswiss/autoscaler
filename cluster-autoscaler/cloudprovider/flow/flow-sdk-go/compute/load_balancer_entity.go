package compute

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
)

const (
	LoadBalancerStatusActive = iota + 1
	LoadBalancerStatusDisabled
	LoadBalancerStatusWorking
	LoadBalancerStatusDegraded
	LoadBalancerStatusError
)

type LoadBalancerAlgorithm struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type LoadBalancerAlgorithmList struct {
	Items      []LoadBalancerAlgorithm
	Pagination flowgo.Pagination
}

type LoadBalancerProtocol struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type LoadBalancerProtocolList struct {
	Items      []LoadBalancerProtocol
	Pagination flowgo.Pagination
}

type LoadBalancerHealthCheckType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type LoadBalancerHealthCheckTypeList struct {
	Items      []LoadBalancerHealthCheckType
	Pagination flowgo.Pagination
}

type LoadBalancerStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type LoadBalancerEntityService struct {
	client flowgo.Client
}

func NewLoadBalancerEntityService(client flowgo.Client) LoadBalancerEntityService {
	return LoadBalancerEntityService{client: client}
}

func (l LoadBalancerEntityService) ListAlgorithms(ctx context.Context, cursor flowgo.Cursor) (list LoadBalancerAlgorithmList, err error) {
	list.Pagination, err = l.client.List(ctx, "/v4/entities/compute/load-balancer-algorithms", cursor, &list.Items)
	return
}

func (l LoadBalancerEntityService) ListProtocols(ctx context.Context, cursor flowgo.Cursor) (list LoadBalancerProtocolList, err error) {
	list.Pagination, err = l.client.List(ctx, "/v4/entities/compute/load-balancer-protocols", cursor, &list.Items)
	return
}

func (l LoadBalancerEntityService) ListHealthCheckTypes(ctx context.Context, cursor flowgo.Cursor) (list LoadBalancerHealthCheckTypeList, err error) {
	list.Pagination, err = l.client.List(ctx, "/v4/entities/compute/load-balancer-health-check-types", cursor, &list.Items)
	return
}
