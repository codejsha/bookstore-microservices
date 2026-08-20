package workflow

type ReserveStockRequest struct {
	OrderUid string                 `json:"orderUid"`
	Items    []StockReservationItem `json:"items"`
}

type StockReservationItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type ReleaseStockRequest = ReserveStockRequest
