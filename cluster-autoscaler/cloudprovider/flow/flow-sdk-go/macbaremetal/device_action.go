package macbaremetal

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
)

type DeviceAction struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Command string `json:"command"`
	Sorting int    `json:"sorting"`
}

type DeviceStatus struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Key     string         `json:"key"`
	Actions []DeviceAction `json:"actions"`
}

type DeviceRunAction struct {
	Action string `json:"action"`
}

type DeviceWorkflow struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Command string `json:"command"`
	Sorting int    `json:"sorting"`
}

type DeviceWorkflowList struct {
	Items      []DeviceWorkflow
	Pagination flowgo.Pagination
}

type DeviceRunWorkflow struct {
	Workflow string `json:"workflow"`
}

type DeviceActionService struct {
	client   flowgo.Client
	deviceID int
}

func NewDeviceActionService(client flowgo.Client, deviceID int) DeviceActionService {
	return DeviceActionService{client: client, deviceID: deviceID}
}

func (d DeviceActionService) Run(ctx context.Context, body DeviceRunAction) (device Device, err error) {
	err = d.client.Create(ctx, getDeviceActionPath(d.deviceID), body, &device)
	return
}

type DeviceWorkflowService struct {
	client   flowgo.Client
	deviceID int
}

func NewDeviceWorkflowService(client flowgo.Client, deviceID int) DeviceWorkflowService {
	return DeviceWorkflowService{client: client, deviceID: deviceID}
}

func (d DeviceWorkflowService) List(ctx context.Context, cursor flowgo.Cursor) (list DeviceWorkflowList, err error) {
	list.Pagination, err = d.client.List(ctx, getDeviceWorkflowPath(d.deviceID), cursor, &list.Items)
	return
}

func (d DeviceWorkflowService) Run(ctx context.Context, body DeviceRunWorkflow) (device Device, err error) {
	err = d.client.Create(ctx, getDeviceWorkflowPath(d.deviceID), body, &device)
	return
}

const (
	deviceActionSegment   = "actions"
	deviceWorkflowSegment = "workflows"
)

func getDeviceActionPath(id int) string {
	return flowgo.Join(devicesSegment, id, deviceActionSegment)
}

func getDeviceWorkflowPath(id int) string {
	return flowgo.Join(devicesSegment, id, deviceWorkflowSegment)
}
