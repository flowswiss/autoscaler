package objectstorage

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/common"
)

type Credential struct {
	ID        int             `json:"id"`
	Location  common.Location `json:"location"`
	Endpoint  string          `json:"endpoint"`
	AccessKey string          `json:"access_key"`
	SecretKey string          `json:"secret_key"`
}

type CredentialList struct {
	Items      []Credential
	Pagination flowgo.Pagination
}

type CredentialService struct {
	client flowgo.Client
}

func NewCredentialService(client flowgo.Client) CredentialService {
	return CredentialService{
		client: client,
	}
}

func (i CredentialService) List(ctx context.Context, cursor flowgo.Cursor) (list CredentialList, err error) {
	list.Pagination, err = i.client.List(ctx, getCredentialSegment(), cursor, &list.Items)
	return
}

const credentialSegment = "/v4/object-storage/credentials"

func getCredentialSegment() string {
	return credentialSegment
}
