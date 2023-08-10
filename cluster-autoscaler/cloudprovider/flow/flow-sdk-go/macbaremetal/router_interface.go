package macbaremetal

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

const routerInterfacesSegment = "router-interfaces"

func getRouterInterfacesPath(id int) string {
	return flowgo.Join(routersSegment, id, routerInterfacesSegment)
}
