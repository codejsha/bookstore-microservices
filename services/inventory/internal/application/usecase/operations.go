package usecase

import (
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/condition"
)

type (
	StockCommand       interface{ isStockCommand() }
	StockCreateCommand struct{ Ent entity.StockEntity }
	StockUpdateCommand struct{ Ent entity.StockEntity }
	StockDeleteCommand struct{ Id int64 }
)

func (StockCreateCommand) isStockCommand() {}
func (StockUpdateCommand) isStockCommand() {}
func (StockDeleteCommand) isStockCommand() {}

type (
	StockQuery        interface{ isStockQuery() }
	StockFindAllQuery struct{ Cond condition.StockCondition }
	StockFindQuery    struct{ Id int64 }
)

func (StockFindAllQuery) isStockQuery() {}
func (StockFindQuery) isStockQuery()    {}

type (
	WarehouseCommand       interface{ isWarehouseCommand() }
	WarehouseCreateCommand struct{ Ent entity.WarehouseEntity }
	WarehouseUpdateCommand struct{ Ent entity.WarehouseEntity }
	WarehouseDeleteCommand struct{ Id int64 }
)

func (WarehouseCreateCommand) isWarehouseCommand() {}
func (WarehouseUpdateCommand) isWarehouseCommand() {}
func (WarehouseDeleteCommand) isWarehouseCommand() {}

type (
	WarehouseQuery        interface{ isWarehouseQuery() }
	WarehouseFindAllQuery struct{ Cond condition.WarehouseCondition }
	WarehouseFindQuery    struct{ Id int64 }
)

func (WarehouseFindAllQuery) isWarehouseQuery() {}
func (WarehouseFindQuery) isWarehouseQuery()    {}
