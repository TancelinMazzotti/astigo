package grpc

import (
	"astigo/internal/domain/contract/data"
	"astigo/internal/domain/contract/service"
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
	svc service.IFooService
}

func (s *FooService) List(ctx context.Context, req *proto.ListFoosRequest) (*proto.ListFoosResponse, error) {
	foos, err := s.svc.GetAll(ctx, data.FooReadListInput{
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to get all foos: %w", err)
	}

	foosProto := make([]*proto.Foo, len(foos))
	for i, foo := range foos {
		foosProto[i] = &proto.Foo{
			Id:     foo.Id.String(),
			Label:  foo.Label,
			Value:  int32(foo.Value),
			Weight: foo.Weight,
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
			Id:     foo.Id.String(),
			Label:  foo.Label,
			Value:  int32(foo.Value),
			Weight: foo.Weight,
		},
	}, nil
}

func (s *FooService) Create(ctx context.Context, req *proto.CreateFooRequest) (*proto.FooResponse, error) {
	foo, err := s.svc.Create(ctx, data.FooCreateInput{
		Label:  req.Label,
		Secret: req.Secret,
		Value:  int(req.Value),
		Weight: req.Weight,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to create foo: %w", err)
	}

	return &proto.FooResponse{
		Foo: &proto.Foo{
			Id:     foo.Id.String(),
			Label:  foo.Label,
			Value:  int32(foo.Value),
			Weight: foo.Weight,
		},
	}, nil
}

func (s *FooService) Update(ctx context.Context, req *proto.UpdateFooRequest) (*proto.FooResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("fail to parse id: %w", err)
	}

	if err := s.svc.Update(ctx, &data.FooUpdateInput{
		Id:     id,
		Label:  req.Label,
		Secret: req.Secret,
		Value:  int(req.Value),
		Weight: req.Weight,
	}); err != nil {
		return nil, fmt.Errorf("fail to update foo: %w", err)
	}

	return &proto.FooResponse{
		Foo: &proto.Foo{
			Id:     id.String(),
			Label:  req.Label,
			Value:  req.Value,
			Weight: req.Weight,
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

func NewFooService(svc service.IFooService) proto.FooServiceServer {
	return &FooService{
		svc: svc,
	}
}
