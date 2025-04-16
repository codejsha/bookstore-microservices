package protosvc

import (
	"context"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/pb/stockpb"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/mapper/stockmapper"
)

var _ stockpb.StockServiceServer = (*StockGrpcServer)(nil)

func NewStockGrpcServer(
	stockUseCase usecase.StockUseCase,
) *StockGrpcServer {
	return &StockGrpcServer{
		stockUseCase: stockUseCase,
	}
}

type StockGrpcServer struct {
	stockUseCase usecase.StockUseCase
	stockpb.UnimplementedStockServiceServer
}

func (s StockGrpcServer) FindAllStocks(ctx context.Context, req *stockpb.StockFindAllProtoReq) (*stockpb.StockFindAllProtoResp, error) {
	cond, err := stockmapper.ToFilterCondition(req)
	if err != nil {
		return nil, err
	}
	total, stockAggs, err := s.stockUseCase.FindAllStocks(ctx, cond)
	if err != nil {
		return nil, err
	}

	resps := make([]*stockpb.StockFindProtoResp, len(stockAggs))
	for _, stockAgg := range stockAggs {
		resps = append(resps, stockAgg.ToStockFindProtoResp())
	}
	response := &stockpb.StockFindAllProtoResp{Total: total, Items: resps}
	return response, nil
}

func (s StockGrpcServer) FindStock(ctx context.Context, req *stockpb.StockFindProtoReq) (*stockpb.StockFindProtoResp, error) {
	stockAgg, err := s.stockUseCase.FindStock(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	response := stockAgg.ToStockFindProtoResp()
	return response, nil
}
