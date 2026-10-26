class ConflictError(Exception):
    pass


class InventoryError(Exception):
    pass


class InvalidQuantityError(InventoryError):
    pass


class InsufficientStockError(InventoryError):
    pass


class SameWarehouseError(InventoryError):
    pass


class StockNotFoundError(InventoryError):
    pass


class WarehouseNotFoundError(InventoryError):
    pass


class DuplicateClosingError(InventoryError):
    pass


class InvalidTransitionError(InventoryError):
    pass
