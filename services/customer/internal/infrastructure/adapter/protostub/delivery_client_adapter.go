package protostub

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/deliverypb"
	port "github.com/codejsha/bookstore-microservices/customer/internal/application/port/protostub"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
)

var _ port.DeliveryClient = (*deliveryClient)(nil)

const deliveryCallerMetadataKey = "x-user-id"

type deliveryClient struct {
	client deliverypb.DeliveryServiceClient
}

func NewDeliveryClient(c *DeliveryGrpcClient) port.DeliveryClient {
	return &deliveryClient{client: c.Client}
}

func withDeliveryCaller(ctx context.Context, callerUid string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, deliveryCallerMetadataKey, callerUid)
}

func (g *deliveryClient) TrackShipment(ctx context.Context, uid string, callerUid string) (*aggregate.ShipmentAggregate, error) {
	ctx = withDeliveryCaller(ctx, callerUid)
	resp, err := g.client.TrackShipment(ctx, &deliverypb.TrackShipmentRequest{Uid: uid})
	if err != nil {
		return nil, err
	}
	sh := resp.GetShipment()
	if sh == nil {
		return nil, nil
	}
	return toShipmentAggregate(sh), nil
}

func (g *deliveryClient) ListShipments(ctx context.Context, orderUid string, status *string, pageSize *int32, callerUid string) (int64, []*aggregate.ShipmentAggregate, error) {
	req := &deliverypb.ListShipmentsRequest{
		OrderUid: orderUid,
		Status:   stringToShipmentStatusProto(status),
	}
	if pageSize != nil {
		req.PageSize = *pageSize
	}
	ctx = withDeliveryCaller(ctx, callerUid)
	resp, err := g.client.ListShipments(ctx, req)
	if err != nil {
		return 0, nil, err
	}
	out := make([]*aggregate.ShipmentAggregate, len(resp.GetShipments()))
	for i, sh := range resp.GetShipments() {
		out[i] = toShipmentAggregate(sh)
	}
	return int64(resp.GetTotalSize()), out, nil
}

func toShipmentAggregate(sh *deliverypb.Shipment) *aggregate.ShipmentAggregate {
	return &aggregate.ShipmentAggregate{
		Uid:                    sh.GetUid(),
		OrderUid:               sh.GetOrderUid(),
		CarrierUid:             sh.GetCarrierUid(),
		Status:                 toShipmentStatus(sh.GetStatus()),
		TrackingNumber:         sh.GetTrackingNumber(),
		DestinationCity:        sh.GetDestinationCity(),
		DestinationState:       sh.GetDestinationState(),
		DestinationCountryCode: sh.GetDestinationCountryCode(),
		DestinationPostalCode:  sh.GetDestinationPostalCode(),
		PlannedDeliveryAt:      sh.GetPlannedDeliveryAt(),
		ActualDeliveryAt:       sh.GetActualDeliveryAt(),
		CreatedAt:              sh.GetCreatedAt(),
		UpdatedAt:              sh.GetUpdatedAt(),
	}
}

func toShipmentStatus(s deliverypb.ShipmentStatus) aggregate.ShipmentStatus {
	switch s {
	case deliverypb.ShipmentStatus_SHIPMENT_STATUS_PLANNED:
		return aggregate.SHIPMENTSTATUS_PLANNED
	case deliverypb.ShipmentStatus_SHIPMENT_STATUS_DISPATCHED:
		return aggregate.SHIPMENTSTATUS_DISPATCHED
	case deliverypb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP:
		return aggregate.SHIPMENTSTATUS_PICKED_UP
	case deliverypb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT:
		return aggregate.SHIPMENTSTATUS_IN_TRANSIT
	case deliverypb.ShipmentStatus_SHIPMENT_STATUS_OUT_FOR_DELIVERY:
		return aggregate.SHIPMENTSTATUS_OUT_FOR_DELIVERY
	case deliverypb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED:
		return aggregate.SHIPMENTSTATUS_DELIVERED
	case deliverypb.ShipmentStatus_SHIPMENT_STATUS_FAILED:
		return aggregate.SHIPMENTSTATUS_FAILED
	case deliverypb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED:
		return aggregate.SHIPMENTSTATUS_CANCELLED
	default:
		return aggregate.SHIPMENTSTATUS_UNKNOWN
	}
}

func stringToShipmentStatusProto(s *string) deliverypb.ShipmentStatus {
	if s == nil {
		return deliverypb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED
	}
	switch *s {
	case "PLANNED":
		return deliverypb.ShipmentStatus_SHIPMENT_STATUS_PLANNED
	case "DISPATCHED":
		return deliverypb.ShipmentStatus_SHIPMENT_STATUS_DISPATCHED
	case "PICKED_UP":
		return deliverypb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP
	case "IN_TRANSIT":
		return deliverypb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT
	case "OUT_FOR_DELIVERY":
		return deliverypb.ShipmentStatus_SHIPMENT_STATUS_OUT_FOR_DELIVERY
	case "DELIVERED":
		return deliverypb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED
	case "FAILED":
		return deliverypb.ShipmentStatus_SHIPMENT_STATUS_FAILED
	case "CANCELLED":
		return deliverypb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED
	default:
		return deliverypb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED
	}
}
