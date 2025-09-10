import "@typespec/http";

using Http;

namespace DeliveryService;

@route("/api/v1/shipments")
@tag("Shipment")
@useAuth(BearerAuth)
interface Shipments {
  @doc("List shipments")
  @get
  getAll(
    @query order_uid?: uuid,
    @query carrier_uid?: uuid,
    @query status?: ShipmentStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): ShipmentFindAllResponse;

  @doc("Create shipment")
  @post
  create(@body body: ShipmentCreateRequest): CreatedResponse;

  @doc("Find shipment")
  @route("{uid}")
  @get
  read(@path uid: uuid): ShipmentFindResponse | NotFoundError;

  @doc("Assign carrier to shipment")
  @route("{uid}/carrier")
  @patch(#{ implicitOptionality: false })
  assignCarrier(
    @path uid: uuid,
    @body body: ShipmentAssignCarrierRequest,
  ): ShipmentUpdateResponse | NotFoundError | BadRequestError;

  @doc("Dispatch shipment")
  @route("{uid}/dispatch")
  @patch(#{ implicitOptionality: false })
  dispatch(@path uid: uuid): ShipmentUpdateResponse | NotFoundError | BadRequestError;

  @doc("Mark shipment as picked up")
  @route("{uid}/pickup")
  @patch(#{ implicitOptionality: false })
  pickUp(@path uid: uuid): ShipmentUpdateResponse | NotFoundError | BadRequestError;

  @doc("Mark shipment as delivered")
  @route("{uid}/deliver")
  @patch(#{ implicitOptionality: false })
  deliver(@path uid: uuid): ShipmentUpdateResponse | NotFoundError | BadRequestError;

  @doc("Cancel shipment")
  @route("{uid}/cancel")
  @patch(#{ implicitOptionality: false })
  cancel(@path uid: uuid): ShipmentUpdateResponse | NotFoundError | BadRequestError;

  @doc("Add tracking event")
  @route("{shipment_uid}/tracking")
  @post
  addTracking(
    @path shipment_uid: uuid,
    @body body: TrackingCreateRequest,
  ): CreatedResponse | NotFoundError;

  @doc("Get tracking history")
  @route("{shipment_uid}/tracking")
  @get
  getTrackingHistory(@path shipment_uid: uuid): TrackingFindAllResponse;
}
