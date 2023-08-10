package common

import (
	"context"

	flowgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/flow/flow-sdk-go"
)

type ProductType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type ProductUsageCycle struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

type ProductItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
}

type ProductAvailability struct {
	Location  Location `json:"location"`
	Available int      `json:"available"`
}

type DeploymentFee struct {
	Location        Location `json:"location"`
	Price           float64  `json:"price"`
	FreeDeployments int      `json:"free_deployments"`
}

type Product struct {
	ID             int                   `json:"id"`
	Name           string                `json:"product_name"`
	Type           ProductType           `json:"type"`
	Visibility     string                `json:"visibility"`
	UsageCycle     ProductUsageCycle     `json:"usage_cycle"`
	Items          []ProductItem         `json:"items"`
	Price          float64               `json:"price"`
	Availability   []ProductAvailability `json:"availability"`
	Category       string                `json:"category"`
	DeploymentFees []DeploymentFee       `json:"deployment_fees"`
}

type BriefProduct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ProductList struct {
	Items      []Product
	Pagination flowgo.Pagination
}

type ProductTypeList struct {
	Items      []ProductType
	Pagination flowgo.Pagination
}

type ProductService struct {
	client flowgo.Client
}

func NewProductService(client flowgo.Client) ProductService {
	return ProductService{client: client}
}

func (p ProductService) List(ctx context.Context, cursor flowgo.Cursor) (list ProductList, err error) {
	list.Pagination, err = p.client.List(ctx, getProductsPath(), cursor, &list.Items)
	return
}

func (p ProductService) ListByType(ctx context.Context, productType string, cursor flowgo.Cursor) (list ProductList, err error) {
	list.Pagination, err = p.client.List(ctx, getProductsByTypePath(productType), cursor, &list.Items)
	return
}

func (p ProductService) Get(ctx context.Context, id int) (product Product, err error) {
	err = p.client.Get(ctx, getSpecificProductPath(id), &product)
	return
}

func (p ProductService) ListTypes(ctx context.Context, cursor flowgo.Cursor) (list ProductTypeList, err error) {
	list.Pagination, err = p.client.List(ctx, getProductTypesPath(), cursor, &list.Items)
	return
}

const (
	productsSegment     = "/v4/products"
	productTypesSegment = "/v4/entities/product-types"
)

func getProductsPath() string {
	return productsSegment
}

func getProductsByTypePath(productType string) string {
	return flowgo.Join(productsSegment, productType)
}

func getSpecificProductPath(id int) string {
	return flowgo.Join(productsSegment, id)
}

func getProductTypesPath() string {
	return productTypesSegment
}
