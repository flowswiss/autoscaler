package compute

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
)

type Route struct {
	ID          int    `json:"id"`
	Destination string `json:"destination"`
	NextHop     string `json:"nexthop"`
}

type RouteList struct {
	Items      []Route
	Pagination flowgo.Pagination
}

type RouteCreate struct {
	Destination string `json:"destination"`
	NextHop     string `json:"nexthop"`
}

type RouteService struct {
	client   flowgo.Client
	routerID int
}

func NewRouteService(client flowgo.Client, routerID int) RouteService {
	return RouteService{client: client, routerID: routerID}
}

func (r RouteService) List(ctx context.Context, cursor flowgo.Cursor) (list RouteList, err error) {
	list.Pagination, err = r.client.List(ctx, getRoutesPath(r.routerID), cursor, &list.Items)
	return
}

func (r RouteService) Create(ctx context.Context, body RouteCreate) (route Route, err error) {
	err = r.client.Create(ctx, getRoutesPath(r.routerID), body, &route)
	return
}

func (r RouteService) Delete(ctx context.Context, id int) (err error) {
	err = r.client.Delete(ctx, getSpecificRoutePath(r.routerID, id))
	return
}

const routesSegment = "routes"

func getRoutesPath(routerID int) string {
	return flowgo.Join(routersSegment, routerID, routesSegment)
}

func getSpecificRoutePath(routerID, routeID int) string {
	return flowgo.Join(routersSegment, routerID, routesSegment, routeID)
}
