package grpc

import (
	"astigo/internal/domain/handler"
	"astigo/pkg/proto"
	"context"
	"fmt"
	"github.com/google/uuid"
)

var (
	_ proto.FooServiceServer = (*FooService)(nil)
)

type FooService struct {
	proto.UnimplementedFooServiceServer
	svc handler.IFooHandler
}

func (s *FooService) List(ctx context.Context, req *proto.ListFoosRequest) (*proto.ListFoosResponse, error) {
	foos, err := s.svc.GetAll(ctx, handler.PaginationInput{
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to get all foos: %w", err)
	}

	foosProto := make([]*proto.Foo, len(foos))
	for i, foo := range foos {
		foosProto[i] = &proto.Foo{
			Id:    foo.Id.String(),
			Label: foo.Label,
		}
	}

	return &proto.ListFoosResponse{Foos: foosProto}, nil
}

func (s *FooService) Get(ctx context.Context, req *proto.GetFooRequest) (*proto.FooResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("fail to parse id: %w", err)
	}

	foo, err := s.svc.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fail to get foo by id: %w", err)
	}

	return &proto.FooResponse{
		Foo: &proto.Foo{
			Id:    foo.Id.String(),
			Label: foo.Label,
		},
	}, nil
}

func NewFooService(svc handler.IFooHandler) proto.FooServiceServer {
	return &FooService{
		svc: svc,
	}
}
