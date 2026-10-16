package restcontroller

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

func callerFrom(ctx context.Context) usecase.Caller {
	c := httpx.CallerFromContext(ctx)
	return usecase.Caller{UserUid: c.UserUid, IsAdmin: c.IsAdmin}
}

var _ openapi.DeliveryApi = (*deliveryController)(nil)

type deliveryController struct {
	customerUseCase usecase.CustomerUseCase
}

func NewDeliveryController(customerUseCase usecase.CustomerUseCase) openapi.DeliveryApi {
	return &deliveryController{customerUseCase: customerUseCase}
}

func (c *deliveryController) DeliveryTrack(
	ctx context.Context,
	_ openapi.DeliveryApiCredentials,
	uid string,
) (openapi.DeliveryTrackResponse, error) {
	shipment, err := c.customerUseCase.TrackShipment(ctx, uid, callerFrom(ctx))
	if err != nil {
		return nil, httpx.MapGrpcStatus(ctx, httpx.MapNotFound(ctx, err))
	}
	if shipment == nil {
		return nil, httpx.MapNotFound(ctx, httpx.ErrNotFound)
	}
	return openapi.DeliveryTrackOk{Body: toShipmentResponse(shipment)}, nil
}

func (c *deliveryController) DeliveryList(
	ctx context.Context,
	_ openapi.DeliveryApiCredentials,
	orderUid string,
	status *openapi.ShipmentStatus,
	size *int32,
) (openapi.DeliveryListResponse, error) {
	var statusStr *string
	if status != nil {
		s := string(*status)
		statusStr = &s
	}

	total, shipments, err := c.customerUseCase.ListOrderShipments(ctx, orderUid, statusStr, size, callerFrom(ctx))
	if err != nil {
		return nil, httpx.MapGrpcStatus(ctx, err)
	}

	items := make([]openapi.ShipmentResponse, len(shipments))
	for i, shipment := range shipments {
		items[i] = toShipmentResponse(shipment)
	}
	return openapi.DeliveryListOk{Body: openapi.ShipmentFindAllResponse{Items: items, Total: total}}, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toShipmentResponse(a *aggregate.ShipmentAggregate) openapi.ShipmentResponse {
	status := toShipmentStatusRest(a.Status)
	return openapi.ShipmentResponse{
		Uid:                    a.Uid,
		OrderUid:               optionalRestString(a.OrderUid),
		CarrierUid:             optionalRestString(a.CarrierUid),
		Status:                 &status,
		TrackingNumber:         optionalRestString(a.TrackingNumber),
		DestinationCity:        optionalRestString(a.DestinationCity),
		DestinationState:       optionalRestString(a.DestinationState),
		DestinationCountryCode: optionalRestString(a.DestinationCountryCode),
		DestinationPostalCode:  optionalRestString(a.DestinationPostalCode),
		PlannedDeliveryAt:      optionalRestString(a.PlannedDeliveryAt),
		ActualDeliveryAt:       optionalRestString(a.ActualDeliveryAt),
		CreatedAt:              optionalRestString(a.CreatedAt),
		UpdatedAt:              optionalRestString(a.UpdatedAt),
	}
}

func toShipmentStatusRest(s aggregate.ShipmentStatus) openapi.ShipmentStatus {
	switch s {
	case aggregate.SHIPMENTSTATUS_PLANNED:
		return openapi.ShipmentStatusPlanned
	case aggregate.SHIPMENTSTATUS_DISPATCHED:
		return openapi.ShipmentStatusDispatched
	case aggregate.SHIPMENTSTATUS_PICKED_UP:
		return openapi.ShipmentStatusPickedUp
	case aggregate.SHIPMENTSTATUS_IN_TRANSIT:
		return openapi.ShipmentStatusInTransit
	case aggregate.SHIPMENTSTATUS_OUT_FOR_DELIVERY:
		return openapi.ShipmentStatusOutForDelivery
	case aggregate.SHIPMENTSTATUS_DELIVERED:
		return openapi.ShipmentStatusDelivered
	case aggregate.SHIPMENTSTATUS_FAILED:
		return openapi.ShipmentStatusFailed
	case aggregate.SHIPMENTSTATUS_CANCELLED:
		return openapi.ShipmentStatusCancelled
	default:
		return openapi.ShipmentStatusUnknown
	}
}

func optionalRestString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
