package macbaremetal

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
)

type NetworkInterface struct {
	ID                int           `json:"id"`
	PrivateIP         string        `json:"private_ip"`
	MacAddress        string        `json:"mac_address"`
	Network           Network       `json:"network"`
	SecurityGroup     SecurityGroup `json:"security_group"`
	AttachedElasticIP ElasticIP     `json:"attached_elastic_ip"`
}

type NetworkInterfaceList struct {
	Items      []NetworkInterface
	Pagination flowgo.Pagination
}

type NetworkInterfaceSecurityGroupUpdate struct {
	SecurityGroupID int `json:"security_group_id"`
}

type NetworkInterfaceService struct {
	client   flowgo.Client
	deviceID int
}

func NewNetworkInterfaceService(client flowgo.Client, deviceID int) NetworkInterfaceService {
	return NetworkInterfaceService{client: client, deviceID: deviceID}
}

func (n NetworkInterfaceService) List(ctx context.Context, cursor flowgo.Cursor) (list NetworkInterfaceList, err error) {
	list.Pagination, err = n.client.List(ctx, getNetworkInterfacesPath(n.deviceID), cursor, &list.Items)
	return
}

func (n NetworkInterfaceService) UpdateSecurityGroup(ctx context.Context, id int, body NetworkInterfaceSecurityGroupUpdate) (networkInterface NetworkInterface, err error) {
	err = n.client.Update(ctx, getSpecificNetworkInterfacePath(n.deviceID, id), body, &networkInterface)
	return
}

const networkInterfacesSegment = "network-interfaces"

func getNetworkInterfacesPath(deviceID int) string {
	return flowgo.Join(getSpecificDevicePath(deviceID), networkInterfacesSegment)
}

func getSpecificNetworkInterfacePath(deviceID, networkInterfaceID int) string {
	return flowgo.Join(getNetworkInterfacesPath(deviceID), networkInterfaceID)
}
