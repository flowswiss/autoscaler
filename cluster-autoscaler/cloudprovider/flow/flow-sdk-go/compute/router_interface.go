package compute

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
)

type RouterInterface struct {
	ID        int     `json:"id"`
	PrivateIP string  `json:"private_ip"`
	Network   Network `json:"network"`
}

type RouterInterfaceList struct {
	Items      []RouterInterface
	Pagination flowgo.Pagination
}

type RouterInterfaceCreate struct {
	NetworkID int    `json:"network_id"`
	PrivateIP string `json:"private_ip,omitempty"`
}

type RouterInterfaceService struct {
	client   flowgo.Client
	routerID int
}

func NewRouterInterfaceService(client flowgo.Client, routerID int) RouterInterfaceService {
	return RouterInterfaceService{client: client, routerID: routerID}
}

func (r RouterInterfaceService) List(ctx context.Context, cursor flowgo.Cursor) (list RouterInterfaceList, err error) {
	list.Pagination, err = r.client.List(ctx, getRouterInterfacesPath(r.routerID), cursor, &list.Items)
	return
}

func (r RouterInterfaceService) Create(ctx context.Context, body RouterInterfaceCreate) (routerInterface RouterInterface, err error) {
	err = r.client.Create(ctx, getRouterInterfacesPath(r.routerID), body, &routerInterface)
	return
}

func (r RouterInterfaceService) Delete(ctx context.Context, id int) (err error) {
	err = r.client.Delete(ctx, getSpecificRouterInterfacePath(r.routerID, id))
	return
}

const routerInterfacesSegment = "interfaces"

func getRouterInterfacesPath(routerID int) string {
	return flowgo.Join(routersSegment, routerID, routerInterfacesSegment)
}

func getSpecificRouterInterfacePath(routerID, routerInterfaceID int) string {
	return flowgo.Join(routersSegment, routerID, routerInterfacesSegment, routerInterfaceID)
}
