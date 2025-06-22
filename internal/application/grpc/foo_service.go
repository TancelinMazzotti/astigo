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

func (s *FooService) Create(ctx context.Context, req *proto.CreateFooRequest) (*proto.FooResponse, error) {
	foo, err := s.svc.Create(ctx, handler.FooCreateInput{
		Label:  req.Label,
		Secret: req.Secret,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to create foo: %w", err)
	}

	return &proto.FooResponse{
		Foo: &proto.Foo{
			Id:    foo.Id.String(),
			Label: foo.Label,
		},
	}, nil
}

func (s *FooService) Update(ctx context.Context, req *proto.UpdateFooRequest) (*proto.FooResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("fail to parse id: %w", err)
	}

	if err := s.svc.Update(ctx, handler.FooUpdateInput{
		Id:     id,
		Label:  req.Label,
		Secret: req.Secret,
	}); err != nil {
		return nil, fmt.Errorf("fail to update foo: %w", err)
	}

	return &proto.FooResponse{
		Foo: &proto.Foo{
			Id:    id.String(),
			Label: req.Label,
		},
	}, nil
}

func (s *FooService) Delete(ctx context.Context, req *proto.DeleteFooRequest) (*proto.DeleteFooResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("fail to parse id: %w", err)
	}
	if err := s.svc.DeleteByID(ctx, id); err != nil {
		return nil, fmt.Errorf("fail to delete foo: %w", err)
	}

	return &proto.DeleteFooResponse{
		Success: true,
	}, nil
}

func NewFooService(svc handler.IFooHandler) proto.FooServiceServer {
	return &FooService{
		svc: svc,
	}
}
