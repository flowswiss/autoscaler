package compute

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/common"
)

type SecurityGroup struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Location    common.Location `json:"location"`
	Default     bool            `json:"default"`
	Immutable   bool            `json:"immutable"`
}

type SecurityGroupList struct {
	Items      []SecurityGroup
	Pagination flowgo.Pagination
}

type SecurityGroupCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LocationID  int    `json:"location_id"`
}

type SecurityGroupUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SecurityGroupService struct {
	client flowgo.Client
}

func NewSecurityGroupService(client flowgo.Client) SecurityGroupService {
	return SecurityGroupService{client: client}
}

func (s SecurityGroupService) Rules(securityGroupID int) SecurityGroupRuleService {
	return NewSecurityGroupRuleService(s.client, securityGroupID)
}

func (s SecurityGroupService) List(ctx context.Context, cursor flowgo.Cursor) (list SecurityGroupList, err error) {
	list.Pagination, err = s.client.List(ctx, getSecurityGroupsPath(), cursor, &list.Items)
	return
}

func (s SecurityGroupService) Create(ctx context.Context, body SecurityGroupCreate) (securityGroup SecurityGroup, err error) {
	err = s.client.Create(ctx, getSecurityGroupsPath(), body, &securityGroup)
	return
}

func (s SecurityGroupService) Get(ctx context.Context, id int) (securityGroup SecurityGroup, err error) {
	err = s.client.Get(ctx, getSpecificSecurityGroupPath(id), &securityGroup)
	return
}

func (s SecurityGroupService) Update(ctx context.Context, id int, body SecurityGroupUpdate) (securityGroup SecurityGroup, err error) {
	err = s.client.Update(ctx, getSpecificSecurityGroupPath(id), body, &securityGroup)
	return
}

func (s SecurityGroupService) Delete(ctx context.Context, id int) (err error) {
	err = s.client.Delete(ctx, getSpecificSecurityGroupPath(id))
	return
}

const securityGroupsSegment = "/v4/compute/security-groups"

func getSecurityGroupsPath() string {
	return securityGroupsSegment
}

func getSpecificSecurityGroupPath(securityGroupID int) string {
	return flowgo.Join(securityGroupsSegment, securityGroupID)
}
