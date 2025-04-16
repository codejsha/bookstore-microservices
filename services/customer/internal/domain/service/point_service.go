package service

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/handler"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/condition"
)

var _ usecase.PointUseCase = (*pointService)(nil)

func NewPointService(
	pointHandler handler.PointHandler,
) usecase.PointUseCase {
	return &pointService{
		pointHandler: pointHandler,
	}
}

type pointService struct {
	pointHandler handler.PointHandler
}

func (s pointService) FindPoint(ctx context.Context, cond condition.PointCondition) (entity.PointEntity, error) {
	// find point
	_, points, err := s.pointHandler.HandleFindAll(ctx, usecase.PointFindAllQuery{Cond: cond})
	if err != nil {
		return entity.PointEntity{}, err
	}
	if len(points) == 0 {
		return entity.PointEntity{}, nil
	}
	return points[0], nil
}

func (s pointService) UpdatePoint(ctx context.Context, ent entity.PointEntity) (entity.PointEntity, error) {
	// find point
	cond := condition.PointCondition{UserId: &ent.UserId}
	_, points, err := s.pointHandler.HandleFindAll(ctx, usecase.PointFindAllQuery{Cond: cond})
	if err != nil {
		return entity.PointEntity{}, err
	}
	if len(points) == 0 {
		return entity.PointEntity{}, nil
	}

	// update point
	ent.Point += points[0].Point
	point, err := s.pointHandler.Handle(ctx, usecase.PointUpdateCommand{Ent: ent})
	if err != nil {
		return entity.PointEntity{}, err
	}

	return point, nil
}
