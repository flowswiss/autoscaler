package objectstorage

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/common"
)

type Instance struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Location common.Location `json:"location"`
}

type InstanceList struct {
	Items      []Instance
	Pagination flowgo.Pagination
}

type InstanceCreate struct {
	LocationID int `json:"location_id"`
}

type InstanceService struct {
	client flowgo.Client
}

func NewInstanceService(client flowgo.Client) InstanceService {
	return InstanceService{
		client: client,
	}
}

func (i InstanceService) List(ctx context.Context, cursor flowgo.Cursor) (list InstanceList, err error) {
	list.Pagination, err = i.client.List(ctx, getInstancePath(), cursor, &list.Items)
	return
}

func (i InstanceService) Create(ctx context.Context, body InstanceCreate) (instance Instance, err error) {
	err = i.client.Create(ctx, getInstancePath(), body, &instance)
	return
}

func (i InstanceService) Delete(ctx context.Context, id int) (err error) {
	err = i.client.Delete(ctx, getSpecificInstancePath(id))
	return
}

const instanceSegment = "/v4/object-storage/instances"

func getInstancePath() string {
	return instanceSegment
}

func getSpecificInstancePath(loadBalancerID int) string {
	return flowgo.Join(instanceSegment, loadBalancerID)
}
