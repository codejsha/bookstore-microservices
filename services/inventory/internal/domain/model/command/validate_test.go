package command

import (
	"errors"
	"strings"
	"testing"
)

const (
	uidEdition1   = "aaaaaaa1-0000-4000-8000-000000000001"
	uidWarehouse1 = "bbbbbbb1-0000-4000-8000-000000000001"
	uidWarehouse2 = "bbbbbbb1-0000-4000-8000-000000000002"
)

func ptr[T any](v T) *T { return &v }

// Receive, release, reserve and adjust share one rule set, so every command runs here.
func TestStockCommands_WhenUidBlankOrNotUuid_ReturnErrInvalidCommand(t *testing.T) {
	if err := (StockReceiveCommand{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 5}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if err := (StockReceiveCommand{EditionUid: " ", WarehouseUid: uidWarehouse1, Quantity: 1}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("blank edition_uid: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockReceiveCommand{EditionUid: "ed-1", WarehouseUid: uidWarehouse1, Quantity: 1}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("non-uuid edition_uid: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockReleaseCommand{EditionUid: uidEdition1, WarehouseUid: ""}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("blank warehouse_uid: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockAdjustCommand{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 1, Reason: " "}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("blank reason: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockReserveCommand{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 1, Reason: ptr("")}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("blank optional reason: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockReserveCommand{
		EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 1, Reason: ptr(strings.Repeat("r", 256)),
	}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("over-long reason: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockAdjustCommand{
		EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 1, Reason: strings.Repeat("r", 256),
	}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("over-long required reason: want ErrInvalidCommand, got %v", err)
	}
}

// Zero and negative quantities are rejected on each of the four stock commands.
func TestStockCommands_WhenQuantityNotPositive_ReturnErrInvalidCommand(t *testing.T) {
	if err := (StockReceiveCommand{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 0}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("zero receive quantity: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockReleaseCommand{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: -1}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("negative release quantity: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockReserveCommand{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 0}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("zero reserve quantity: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockAdjustCommand{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: 0, Reason: "noop"}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("zero adjust quantity: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockAdjustCommand{EditionUid: uidEdition1, WarehouseUid: uidWarehouse1, Quantity: -3, Reason: "damaged"}).Validate(); err != nil {
		t.Fatalf("negative adjust quantity must be valid, got %v", err)
	}
}

func TestStockOrderCommand_WhenOrderUidOrEditionIdInvalid_ReturnsErrInvalidCommand(t *testing.T) {
	if err := (StockOrderReserveCommand{OrderUid: "o-1", EditionId: 42, Quantity: 1}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if err := (StockOrderReserveCommand{OrderUid: " ", EditionId: 42}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("blank order_uid: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockOrderReserveCommand{OrderUid: strings.Repeat("o", 201), EditionId: 42}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("over-long order_uid: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockOrderReleaseCommand{OrderUid: "o-1", EditionId: 0}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("zero edition_id: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockOrderReleaseCommand{OrderUid: strings.Repeat("o", 201), EditionId: 1}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("over-long order_uid: want ErrInvalidCommand, got %v", err)
	}
}

func TestStockTransferCommand_WhenWarehousesOrQuantityInvalid_ReturnsErrInvalidCommand(t *testing.T) {
	if err := (StockTransferCommand{
		EditionUid: uidEdition1, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: uidWarehouse2, Quantity: 1,
	}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if err := (StockTransferCommand{
		EditionUid: uidEdition1, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: " ", Quantity: 1,
	}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("blank target: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockTransferCommand{
		EditionUid: uidEdition1, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: uidWarehouse1, Quantity: 1,
	}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("same warehouse: want ErrInvalidCommand, got %v", err)
	}
	if err := (StockTransferCommand{
		EditionUid: uidEdition1, SourceWarehouseUid: uidWarehouse1, TargetWarehouseUid: uidWarehouse2, Quantity: 0,
	}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("zero quantity: want ErrInvalidCommand, got %v", err)
	}
}

func TestWarehouseCommand_WhenNameOrAddressInvalid_ReturnsErrInvalidCommand(t *testing.T) {
	if err := (WarehouseCreateCommand{Name: "Main", Capacity: 100}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if err := (WarehouseCreateCommand{Name: " ", Capacity: 100}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("blank name: want ErrInvalidCommand, got %v", err)
	}
	if err := (WarehouseCreateCommand{Name: strings.Repeat("n", 256), Capacity: 100}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("over-long name: want ErrInvalidCommand, got %v", err)
	}
	if err := (WarehouseCreateCommand{
		Name: "Main", Address: ptr(strings.Repeat("a", 256)), Capacity: 100,
	}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("over-long address: want ErrInvalidCommand, got %v", err)
	}
	if err := (WarehouseCreateCommand{Name: "Main", Capacity: -1}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("negative capacity: want ErrInvalidCommand, got %v", err)
	}
	if err := (WarehouseUpdateCommand{}).Validate(); err != nil {
		t.Fatalf("empty update must be valid, got %v", err)
	}
	if err := (WarehouseUpdateCommand{Name: ptr(strings.Repeat("n", 256))}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("over-long name: want ErrInvalidCommand, got %v", err)
	}
	if err := (WarehouseUpdateCommand{Capacity: ptr(int32(-5))}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("negative capacity: want ErrInvalidCommand, got %v", err)
	}
}

func TestStockAuditCreateCommand_Validate(t *testing.T) {
	valid := StockAuditCreateCommand{
		WarehouseUid: uidWarehouse1,
		Items:        []StockAuditItemCommand{{EditionUid: uidEdition1, ActualQuantity: 3}},
	}
	tests := []struct {
		name    string
		mutate  func(c *StockAuditCreateCommand)
		wantErr bool
	}{
		{"whenEveryFieldValid_returnsNil", func(c *StockAuditCreateCommand) {}, false},
		{"whenWarehouseUidBlank_returnsErrInvalidCommand", func(c *StockAuditCreateCommand) { c.WarehouseUid = " " }, true},
		{"whenWarehouseUidNotUuid_returnsErrInvalidCommand", func(c *StockAuditCreateCommand) { c.WarehouseUid = "wh-1" }, true},
		{"whenItemsEmpty_returnsErrInvalidCommand", func(c *StockAuditCreateCommand) { c.Items = nil }, true},
		{"whenItemEditionUidBlank_returnsErrInvalidCommand", func(c *StockAuditCreateCommand) {
			c.Items = []StockAuditItemCommand{{EditionUid: " ", ActualQuantity: 1}}
		}, true},
		{"whenItemEditionUidNotUuid_returnsErrInvalidCommand", func(c *StockAuditCreateCommand) {
			c.Items = []StockAuditItemCommand{{EditionUid: "ed-1", ActualQuantity: 1}}
		}, true},
		{"whenActualQuantityNegative_returnsErrInvalidCommand", func(c *StockAuditCreateCommand) {
			c.Items = []StockAuditItemCommand{{EditionUid: uidEdition1, ActualQuantity: -1}}
		}, true},
		{"whenItemEditionsDuplicated_returnsErrInvalidCommand", func(c *StockAuditCreateCommand) {
			c.Items = []StockAuditItemCommand{
				{EditionUid: uidEdition1, ActualQuantity: 1},
				{EditionUid: uidEdition1, ActualQuantity: 2},
			}
		}, true},
		{"whenNotesBlank_returnsErrInvalidCommand", func(c *StockAuditCreateCommand) { c.Notes = ptr(" ") }, true},
		{"whenNotesOverLimit_returnsErrInvalidCommand", func(c *StockAuditCreateCommand) { c.Notes = ptr(strings.Repeat("n", 501)) }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := valid
			tt.mutate(&cmd)
			err := cmd.Validate()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidCommand) {
					t.Fatalf("want ErrInvalidCommand, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("want nil, got %v", err)
			}
		})
	}
}

func TestMonthlyClosingCommand_WhenPeriodOutOfRange_ReturnsErrInvalidCommand(t *testing.T) {
	if err := (MonthlyClosingCommand{WarehouseUid: uidWarehouse1, Year: 2026, Month: 8}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	for _, bad := range []MonthlyClosingCommand{
		{WarehouseUid: " ", Year: 2026, Month: 8},
		{WarehouseUid: "wh-1", Year: 2026, Month: 8},
		{WarehouseUid: uidWarehouse1, Year: 1999, Month: 8},
		{WarehouseUid: uidWarehouse1, Year: 2026, Month: 0},
		{WarehouseUid: uidWarehouse1, Year: 2026, Month: 13},
	} {
		if err := bad.Validate(); !errors.Is(err, ErrInvalidCommand) {
			t.Fatalf("%+v: want ErrInvalidCommand, got %v", bad, err)
		}
	}
}

// A nil optional is absent, not invalid; the same helpers still reject a present-but-bad value.
func TestOptionalValidationHelpers_WhenValueNil_ReturnNil(t *testing.T) {
	if err := optionalMaxLen("reason", nil, 1); err != nil {
		t.Errorf("nil value must be accepted, got %v", err)
	}
	if err := optionalUUID("warehouse_uid", nil); err != nil {
		t.Errorf("nil value must be accepted, got %v", err)
	}
	if err := optionalUUID("warehouse_uid", ptr(uidWarehouse1)); err != nil {
		t.Errorf("valid uuid must be accepted, got %v", err)
	}
	if err := optionalUUID("warehouse_uid", ptr("wh-1")); !errors.Is(err, ErrInvalidCommand) {
		t.Errorf("err = %v, want ErrInvalidCommand", err)
	}
}
