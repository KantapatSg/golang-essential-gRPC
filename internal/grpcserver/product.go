package grpcserver

import (
	"context"
	"errors"

	productv1 "github.com/KantapatSg/golang-essential-gRPC/gen/product"
	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/KantapatSg/golang-essential-gRPC/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductServer struct {
	productv1.UnimplementedProductServiceServer
	service *service.ProductService
}

func NewProductServer(productService *service.ProductService) *ProductServer {
	return &ProductServer{service: productService}
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.Product, error) {
	p, err := s.service.Create(ctx, req.GetName(), req.GetDescription(), req.GetPrice(), req.GetStock())
	if err != nil {
		return nil, toStatus(err)
	}
	return toProto(p), nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.Product, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "id must be a valid UUID")
	}
	p, err := s.service.Get(ctx, id)
	if err != nil {
		return nil, toStatus(err)
	}
	return toProto(p), nil
}

func (s *ProductServer) ListProducts(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsResponse, error) {
	products, total, err := s.service.List(ctx, int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, toStatus(err)
	}
	result := &productv1.ListProductsResponse{Page: req.GetPage(), PageSize: req.GetPageSize(), Total: total, Products: make([]*productv1.Product, 0, len(products))}
	if result.Page < 1 {
		result.Page = 1
	}
	if result.PageSize < 1 {
		result.PageSize = 20
	}
	if result.PageSize > 100 {
		result.PageSize = 100
	}
	for _, p := range products {
		result.Products = append(result.Products, toProto(p))
	}
	return result, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.Product, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "id must be a valid UUID")
	}
	p, err := s.service.Update(ctx, id, req.GetName(), req.GetDescription(), req.GetPrice(), req.GetStock())
	if err != nil {
		return nil, toStatus(err)
	}
	return toProto(p), nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *productv1.DeleteProductRequest) (*emptypb.Empty, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "id must be a valid UUID")
	}
	if err := s.service.Delete(ctx, id); err != nil {
		return nil, toStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func parseID(value string) (uuid.UUID, error) { return uuid.Parse(value) }

func toStatus(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

func toProto(p domain.Product) *productv1.Product {
	return &productv1.Product{Id: p.ID.String(), Name: p.Name, Description: p.Description, Price: p.Price, Stock: p.Stock, CreatedAt: timestamppb.New(p.CreatedAt), UpdatedAt: timestamppb.New(p.UpdatedAt)}
}
