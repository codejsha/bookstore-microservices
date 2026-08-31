package restcontroller

import (
	"context"
	"errors"
	"testing"

	"github.com/codejsha/bookstore-microservices/inventory/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/option"
)

func failingUseCase(t *testing.T) *stubUseCase {
	t.Helper()
	return &stubUseCase{
		findAllStocks: func(context.Context, option.StockQueryOption) (int64, []*aggregate.StockAggregate, error) {
			t.Error("use case must not be called for an invalid uid")
			return 0, nil, nil
		},
		findStock: func(context.Context, string) (*aggregate.StockAggregate, error) {
			t.Error("use case must not be called for an invalid uid")
			return nil, nil
		},
		findWarehouse: func(context.Context, string) (*aggregate.WarehouseAggregate, error) {
			t.Error("use case must not be called for an invalid uid")
			return nil, nil
		},
		updateWarehouse: func(context.Context, string, command.WarehouseUpdateCommand) (*aggregate.WarehouseAggregate, error) {
			t.Error("use case must not be called for an invalid uid")
			return nil, nil
		},
		findTransfer: func(context.Context, string) (*aggregate.StockTransferAggregate, error) {
			t.Error("use case must not be called for an invalid uid")
			return nil, nil
		},
		completeTransfer: func(context.Context, string) (*aggregate.StockTransferAggregate, error) {
			t.Error("use case must not be called for an invalid uid")
			return nil, nil
		},
		findAudit: func(context.Context, string) (*aggregate.StockAuditAggregate, error) {
			t.Error("use case must not be called for an invalid uid")
			return nil, nil
		},
		findClosing: func(context.Context, string) (*aggregate.MonthlyClosingAggregate, error) {
			t.Error("use case must not be called for an invalid uid")
			return nil, nil
		},
		findStockBalance: func(context.Context, option.BalanceQueryOption) (int64, []*aggregate.StockBalanceEntry, error) {
			t.Error("use case must not be called for an invalid uid")
			return 0, nil, nil
		},
	}
}

// Every uid-taking read is driven through the same subtest, named after the method under test.
func TestControllers_WhenPathParamNotUuid_ReturnBadRequest(t *testing.T) {
	use := failingUseCase(t)
	bad := "not-a-uuid"

	calls := map[string]func() error{
		"StocksRead": func() error {
			_, err := NewStockController(use).StocksRead(context.Background(), bad)
			return err
		},
		"StocksHistory": func() error {
			_, err := NewStockController(use).StocksHistory(context.Background(), bad, nil, nil, nil, nil)
			return err
		},
		"WarehousesRead": func() error {
			_, err := NewWarehouseController(use).WarehousesRead(context.Background(), bad)
			return err
		},
		"WarehousesUpdate": func() error {
			_, err := NewWarehouseController(use).WarehousesUpdate(context.Background(), bad, openapi.WarehouseUpdateRequest{})
			return err
		},
		"TransfersRead": func() error {
			_, err := NewTransferController(use).TransfersRead(context.Background(), bad)
			return err
		},
		"TransfersComplete": func() error {
			_, err := NewTransferController(use).TransfersComplete(context.Background(), bad)
			return err
		},
		"AuditsRead": func() error {
			_, err := NewAuditController(use).AuditsRead(context.Background(), bad)
			return err
		},
		"ClosingsRead": func() error {
			_, err := NewClosingController(use).ClosingsRead(context.Background(), bad)
			return err
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			if err := call(); !errors.Is(err, command.ErrInvalidCommand) {
				t.Errorf("err = %v, want ErrInvalidCommand", err)
			}
		})
	}
}

func TestControllers_WhenQueryParamNotUuid_ReturnBadRequest(t *testing.T) {
	use := failingUseCase(t)
	bad := "not-a-uuid"

	calls := map[string]func() error{
		"StocksGetAll": func() error {
			_, err := NewStockController(use).StocksGetAll(context.Background(), &bad, nil, nil, nil, nil)
			return err
		},
		"BalanceGetAll": func() error {
			_, err := NewBalanceController(use).BalanceGetAll(context.Background(), nil, &bad, nil, nil, nil, nil)
			return err
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			if err := call(); !errors.Is(err, command.ErrInvalidCommand) {
				t.Errorf("err = %v, want ErrInvalidCommand", err)
			}
		})
	}
}

func TestOptionalUidParam_WhenParamEmpty_SkipsTheFilter(t *testing.T) {
	empty := ""
	if err := optionalUidParam(context.Background(), "warehouse_uid", nil); err != nil {
		t.Errorf("nil must be accepted, got %v", err)
	}
	if err := optionalUidParam(context.Background(), "warehouse_uid", &empty); err != nil {
		t.Errorf("empty must be accepted, got %v", err)
	}
}
