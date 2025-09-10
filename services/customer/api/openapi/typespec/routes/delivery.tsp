import "@typespec/http";

using Http;

namespace CustomerService;

// ─── Shipment tracking (proxied from delivery service via gRPC) ──────────────

@route("/api/v1")
@tag("Delivery")
@useAuth(BearerAuth)
interface Delivery {
  @doc("Track a shipment by its uid")
  @route("/shipments/{uid}")
  @get
  track(@path uid: uuid): ShipmentResponse | NotFoundError;

  @doc("List shipments for an order")
  @route("/orders/{order_uid}/shipments")
  @get
  list(
    @path order_uid: uuid,
    @query status?: ShipmentStatus,
    @query size?: Size,
  ): ShipmentFindAllResponse | NotFoundError;
}
