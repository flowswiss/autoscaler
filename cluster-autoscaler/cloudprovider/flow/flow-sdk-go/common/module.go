package common

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
)

type Module struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Parent    *Module    `json:"parent"`
	Sorting   int        `json:"sorting"`
	Locations []Location `json:"locations"`
}

type ModuleList struct {
	flowgo.Pagination
	Items []Module
}

type ModuleService struct {
	client flowgo.Client
}

func NewModuleService(client flowgo.Client) ModuleService {
	return ModuleService{client: client}
}

func (l ModuleService) List(ctx context.Context, cursor flowgo.Cursor) (list ModuleList, err error) {
	list.Pagination, err = l.client.List(ctx, getModulesPath(), cursor, &list.Items)
	return
}

func (l ModuleService) Get(ctx context.Context, id int) (module Module, err error) {
	err = l.client.Get(ctx, getSpecificModulePath(id), &module)
	return
}

const modulesSegment = "/v4/entities/modules"

func getModulesPath() string {
	return modulesSegment
}

func getSpecificModulePath(moduleID int) string {
	return flowgo.Join(modulesSegment, moduleID)
}
