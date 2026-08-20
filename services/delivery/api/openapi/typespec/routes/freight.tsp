import "@typespec/http";

using Http;

namespace DeliveryService;

@route("/api/v1/freights")
@tag("Freight")
@useAuth(BearerAuth)
interface Freights {
  @doc("List freights")
  @get
  getAll(
    @query shipment_uid?: uuid,
    @query carrier_uid?: uuid,
    @query status?: FreightStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): FreightFindAllResponse;

  @doc("Create freight")
  @post
  create(@body body: FreightCreateRequest): CreatedResponse;

  @doc("Find freight")
  @route("{uid}")
  @get
  read(@path uid: uuid): FreightFindResponse | NotFoundError;

  @doc("Find freight by shipment")
  @route("shipment/{shipment_uid}")
  @get
  readByShipment(@path shipment_uid: uuid): FreightFindResponse | NotFoundError;

  @doc("Update freight status")
  @route("{uid}/status")
  @patch(#{ implicitOptionality: false })
  updateStatus(
    @path uid: uuid,
    @body body: FreightUpdateStatusRequest,
  ): FreightUpdateResponse | NotFoundError;
}
