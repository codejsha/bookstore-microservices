import "@typespec/protobuf";

using TypeSpec.Protobuf;

@package({
  name: "delivery.v1",
})
namespace ProtoDelivery;

// ─── Enum ───────────────────────────────────────────────────────────────────

enum ShipmentStatus {
  SHIPMENT_STATUS_UNSPECIFIED: 0,
  SHIPMENT_STATUS_PLANNED: 1,
  SHIPMENT_STATUS_DISPATCHED: 2,
  SHIPMENT_STATUS_PICKED_UP: 3,
  SHIPMENT_STATUS_IN_TRANSIT: 4,
  SHIPMENT_STATUS_OUT_FOR_DELIVERY: 5,
  SHIPMENT_STATUS_DELIVERED: 6,
  SHIPMENT_STATUS_FAILED: 7,
  SHIPMENT_STATUS_CANCELLED: 8,
}

// ─── Entity ─────────────────────────────────────────────────────────────────

model Shipment {
  @field(1) uid: string;
  @field(2) order_uid: string;
  @field(3) carrier_uid?: string;
  @field(4) status: ShipmentStatus;
  @field(5) tracking_number?: string;
  @field(6) destination_city: string;
  @field(7) destination_state: string;
  @field(8) destination_country_code: string;
  @field(9) destination_postal_code: string;
  @field(10) planned_delivery_at?: string;
  @field(11) actual_delivery_at?: string;
  @field(12) created_at: string;
  @field(13) updated_at?: string;
}

// ─── Messages ───────────────────────────────────────────────────────────────

model TrackShipmentRequest {
  @field(1) uid: string;
}

model TrackShipmentResponse {
  @field(1) shipment: Shipment;
}

model ListShipmentsRequest {
  @field(1) order_uid?: string;
  @field(2) status?: ShipmentStatus;
  @field(3) page_size?: int32;
  @field(4) page_token?: string;
  @field(5) order_by?: string;
}

model ListShipmentsResponse {
  @field(1) shipments: Shipment[];
  @field(2) next_page_token?: string;
  @field(3) total_size: int32;
}

// ─── Service ────────────────────────────────────────────────────────────────

@Protobuf.service
interface DeliveryService {
  TrackShipment(...TrackShipmentRequest): TrackShipmentResponse;
  ListShipments(...ListShipmentsRequest): ListShipmentsResponse;
}
