from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ShipmentStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    SHIPMENT_STATUS_UNSPECIFIED: _ClassVar[ShipmentStatus]
    SHIPMENT_STATUS_PLANNED: _ClassVar[ShipmentStatus]
    SHIPMENT_STATUS_DISPATCHED: _ClassVar[ShipmentStatus]
    SHIPMENT_STATUS_PICKED_UP: _ClassVar[ShipmentStatus]
    SHIPMENT_STATUS_IN_TRANSIT: _ClassVar[ShipmentStatus]
    SHIPMENT_STATUS_OUT_FOR_DELIVERY: _ClassVar[ShipmentStatus]
    SHIPMENT_STATUS_DELIVERED: _ClassVar[ShipmentStatus]
    SHIPMENT_STATUS_FAILED: _ClassVar[ShipmentStatus]
    SHIPMENT_STATUS_CANCELLED: _ClassVar[ShipmentStatus]
SHIPMENT_STATUS_UNSPECIFIED: ShipmentStatus
SHIPMENT_STATUS_PLANNED: ShipmentStatus
SHIPMENT_STATUS_DISPATCHED: ShipmentStatus
SHIPMENT_STATUS_PICKED_UP: ShipmentStatus
SHIPMENT_STATUS_IN_TRANSIT: ShipmentStatus
SHIPMENT_STATUS_OUT_FOR_DELIVERY: ShipmentStatus
SHIPMENT_STATUS_DELIVERED: ShipmentStatus
SHIPMENT_STATUS_FAILED: ShipmentStatus
SHIPMENT_STATUS_CANCELLED: ShipmentStatus

class Shipment(_message.Message):
    __slots__ = ("uid", "order_uid", "carrier_uid", "status", "tracking_number", "destination_city", "destination_state", "destination_country_code", "destination_postal_code", "planned_delivery_at", "actual_delivery_at", "created_at", "updated_at")
    UID_FIELD_NUMBER: _ClassVar[int]
    ORDER_UID_FIELD_NUMBER: _ClassVar[int]
    CARRIER_UID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    TRACKING_NUMBER_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_CITY_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_STATE_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_COUNTRY_CODE_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_POSTAL_CODE_FIELD_NUMBER: _ClassVar[int]
    PLANNED_DELIVERY_AT_FIELD_NUMBER: _ClassVar[int]
    ACTUAL_DELIVERY_AT_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    UPDATED_AT_FIELD_NUMBER: _ClassVar[int]
    uid: str
    order_uid: str
    carrier_uid: str
    status: ShipmentStatus
    tracking_number: str
    destination_city: str
    destination_state: str
    destination_country_code: str
    destination_postal_code: str
    planned_delivery_at: str
    actual_delivery_at: str
    created_at: str
    updated_at: str
    def __init__(self, uid: _Optional[str] = ..., order_uid: _Optional[str] = ..., carrier_uid: _Optional[str] = ..., status: _Optional[_Union[ShipmentStatus, str]] = ..., tracking_number: _Optional[str] = ..., destination_city: _Optional[str] = ..., destination_state: _Optional[str] = ..., destination_country_code: _Optional[str] = ..., destination_postal_code: _Optional[str] = ..., planned_delivery_at: _Optional[str] = ..., actual_delivery_at: _Optional[str] = ..., created_at: _Optional[str] = ..., updated_at: _Optional[str] = ...) -> None: ...

class TrackShipmentRequest(_message.Message):
    __slots__ = ("uid",)
    UID_FIELD_NUMBER: _ClassVar[int]
    uid: str
    def __init__(self, uid: _Optional[str] = ...) -> None: ...

class TrackShipmentResponse(_message.Message):
    __slots__ = ("shipment",)
    SHIPMENT_FIELD_NUMBER: _ClassVar[int]
    shipment: Shipment
    def __init__(self, shipment: _Optional[_Union[Shipment, _Mapping]] = ...) -> None: ...

class ListShipmentsRequest(_message.Message):
    __slots__ = ("order_uid", "status", "page_size", "page_token", "order_by")
    ORDER_UID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    PAGE_SIZE_FIELD_NUMBER: _ClassVar[int]
    PAGE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    ORDER_BY_FIELD_NUMBER: _ClassVar[int]
    order_uid: str
    status: ShipmentStatus
    page_size: int
    page_token: str
    order_by: str
    def __init__(self, order_uid: _Optional[str] = ..., status: _Optional[_Union[ShipmentStatus, str]] = ..., page_size: _Optional[int] = ..., page_token: _Optional[str] = ..., order_by: _Optional[str] = ...) -> None: ...

class ListShipmentsResponse(_message.Message):
    __slots__ = ("shipments", "next_page_token", "total_size")
    SHIPMENTS_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    TOTAL_SIZE_FIELD_NUMBER: _ClassVar[int]
    shipments: _containers.RepeatedCompositeFieldContainer[Shipment]
    next_page_token: str
    total_size: int
    def __init__(self, shipments: _Optional[_Iterable[_Union[Shipment, _Mapping]]] = ..., next_page_token: _Optional[str] = ..., total_size: _Optional[int] = ...) -> None: ...
