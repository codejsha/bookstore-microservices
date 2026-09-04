package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

var _ openapi.PointApi = (*pointController)(nil)

type pointController struct {
	customerUseCase usecase.CustomerUseCase
}

func NewPointController(customerUseCase usecase.CustomerUseCase) openapi.PointApi {
	return &pointController{customerUseCase: customerUseCase}
}

func (c *pointController) PointsGetBalance(ctx context.Context, uid string) (*openapi.PointBalanceResponse, error) {
	point, err := c.customerUseCase.GetPointBalance(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &openapi.PointBalanceResponse{
		UserUid: uid,
		Balance: point.Balance,
	}, nil
}

func (c *pointController) PointsEarn(
	ctx context.Context,
	uid string,
	req openapi.PointEarnRequest,
) (*openapi.PointBalanceResponse, error) {
	cmd := command.PointEarnCommand{
		UserUid: uid,
		Amount:  req.Amount,
		Reason:  req.Reason,
	}

	point, err := c.customerUseCase.EarnPoints(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, httpx.MapPointBalanceOverflow(ctx, err))
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "point", uid, "earned", logrus.Fields{"amount": req.Amount})

	return &openapi.PointBalanceResponse{
		UserUid: uid,
		Balance: point.Balance,
	}, nil
}

func (c *pointController) PointsSpend(
	ctx context.Context,
	uid string,
	req openapi.PointSpendRequest,
) (*openapi.PointBalanceResponse, error) {
	cmd := command.PointSpendCommand{
		UserUid: uid,
		Amount:  req.Amount,
		Reason:  req.Reason,
	}

	point, err := c.customerUseCase.SpendPoints(ctx, cmd)
	if err != nil {
		return nil, httpx.MapBusinessError(ctx, httpx.MapInsufficientPoints(ctx, err))
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "point", uid, "spent", logrus.Fields{"amount": req.Amount})

	return &openapi.PointBalanceResponse{
		UserUid: uid,
		Balance: point.Balance,
	}, nil
}

func (c *pointController) PointsHistory(
	ctx context.Context,
	uid string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.PointHistoryFindAllResponse, error) {
	pageOpt := pagination.NewPageOption(size, page, sort)

	total, entries, err := c.customerUseCase.GetPointHistory(ctx, uid, pageOpt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.PointHistoryFindResponse, len(entries))
	for i, e := range entries {
		items[i] = toPointHistoryResponse(e)
	}
	return &openapi.PointHistoryFindAllResponse{Items: items, Total: total}, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toPointHistoryResponse(e *aggregate.PointHistoryEntry) openapi.PointHistoryFindResponse {
	return openapi.PointHistoryFindResponse{
		Uid:        e.Uid,
		ChangeType: openapi.PointChangeType(e.ChangeType),
		Amount:     e.Amount,
		Reason:     e.Reason,
		CreatedAt:  e.CreatedAt,
	}
}
