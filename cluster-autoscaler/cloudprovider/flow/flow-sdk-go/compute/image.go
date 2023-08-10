package compute

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go/common"
)

type Image struct {
	ID                 int              `json:"id"`
	OperatingSystem    string           `json:"os"`
	Version            string           `json:"version"`
	Key                string           `json:"key"`
	Category           string           `json:"category"`
	Type               string           `json:"type"`
	Username           string           `json:"username"`
	MinRootDiskSize    int              `json:"min_root_disk_size"`
	Sorting            int              `json:"sorting"`
	RequiredLicenses   []common.Product `json:"required_licenses"`
	AvailableLocations []int            `json:"available_locations"`
}

type ImageList struct {
	Items      []Image
	Pagination flowgo.Pagination
}

type ImageService struct {
	client flowgo.Client
}

func NewImageService(client flowgo.Client) ImageService {
	return ImageService{client: client}
}

func (i ImageService) List(ctx context.Context, cursor flowgo.Cursor) (list ImageList, err error) {
	list.Pagination, err = i.client.List(ctx, getImagesPath(), cursor, &list.Items)
	return
}

func (i ImageService) Get(ctx context.Context, id int) (image Image, err error) {
	err = i.client.Get(ctx, getSpecificImagePath(id), &image)
	return
}

const imagesSegment = "/v4/entities/compute/images"

func getImagesPath() string {
	return imagesSegment
}

func getSpecificImagePath(id int) string {
	return flowgo.Join(imagesSegment, id)
}
