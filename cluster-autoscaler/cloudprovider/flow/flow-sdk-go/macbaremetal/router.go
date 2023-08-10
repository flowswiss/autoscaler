package macbaremetal

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/common"
)

type Router struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Location    common.Location `json:"location"`
	Public      bool            `json:"public"`
	SourceNAT   bool            `json:"snat"`
	PublicIP    string          `json:"public_ip"`
}

type RouterList struct {
	Items      []Router
	Pagination flowgo.Pagination
}

type RouterUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RouterService struct {
	client flowgo.Client
}

func NewRouterService(client flowgo.Client) RouterService {
	return RouterService{client: client}
}

func (r RouterService) RouterInterfaces(routerID int) RouterInterfaceService {
	return NewRouterInterfaceService(r.client, routerID)
}

func (r RouterService) List(ctx context.Context, cursor flowgo.Cursor) (list RouterList, err error) {
	list.Pagination, err = r.client.List(ctx, getRoutersPath(), cursor, &list.Items)
	return
}

func (r RouterService) Get(ctx context.Context, id int) (router Router, err error) {
	err = r.client.Get(ctx, getSpecificRouterPath(id), &router)
	return
}

func (r RouterService) Update(ctx context.Context, id int, body RouterUpdate) (router Router, err error) {
	err = r.client.Update(ctx, getSpecificRouterPath(id), body, &router)
	return
}

const routersSegment = "/v4/macbaremetal/routers"

func getRoutersPath() string {
	return routersSegment
}

func getSpecificRouterPath(id int) string {
	return flowgo.Join(routersSegment, id)
}
