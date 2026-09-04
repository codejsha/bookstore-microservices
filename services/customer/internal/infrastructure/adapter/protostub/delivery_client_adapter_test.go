package protostub

import (
	"testing"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/deliverypb"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
)

func TestToShipmentAggregate_WhenProtoFullyPopulated_MapsEveryField(t *testing.T) {
	t.Run("whenProtoFullyPopulated_mapsEveryField", func(t *testing.T) {
		in := &deliverypb.Shipment{
			Uid:                    "sh-1",
			OrderUid:               "o-1",
			CarrierUid:             "car-1",
			Status:                 deliverypb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT,
			TrackingNumber:         "1Z-TRACK-9",
			DestinationCity:        "Seattle",
			DestinationState:       "WA",
			DestinationCountryCode: "US",
			DestinationPostalCode:  "98101",
			PlannedDeliveryAt:      "2026-07-01T10:00:00Z",
			ActualDeliveryAt:       "2026-07-02T14:30:00Z",
			CreatedAt:              "2026-06-28T08:00:00Z",
			UpdatedAt:              "2026-06-29T09:00:00Z",
		}
		got := toShipmentAggregate(in)
		want := &aggregate.ShipmentAggregate{
			Uid:                    "sh-1",
			OrderUid:               "o-1",
			CarrierUid:             "car-1",
			Status:                 aggregate.SHIPMENTSTATUS_IN_TRANSIT,
			TrackingNumber:         "1Z-TRACK-9",
			DestinationCity:        "Seattle",
			DestinationState:       "WA",
			DestinationCountryCode: "US",
			DestinationPostalCode:  "98101",
			PlannedDeliveryAt:      "2026-07-01T10:00:00Z",
			ActualDeliveryAt:       "2026-07-02T14:30:00Z",
			CreatedAt:              "2026-06-28T08:00:00Z",
			UpdatedAt:              "2026-06-29T09:00:00Z",
		}
		if *got != *want {
			t.Errorf("toShipmentAggregate mismatch\n got = %+v\nwant = %+v", *got, *want)
		}
	})

	t.Run("whenOptionalStringsEmpty_passesThemThroughEmpty", func(t *testing.T) {
		got := toShipmentAggregate(&deliverypb.Shipment{Uid: "sh-2"})
		if got.Uid != "sh-2" {
			t.Errorf("Uid = %q, want sh-2", got.Uid)
		}
		if got.OrderUid != "" || got.CarrierUid != "" || got.TrackingNumber != "" {
			t.Errorf("expected empty optionals, got order=%q carrier=%q track=%q",
				got.OrderUid, got.CarrierUid, got.TrackingNumber)
		}
		if got.DestinationCity != "" || got.DestinationState != "" ||
			got.DestinationCountryCode != "" || got.DestinationPostalCode != "" {
			t.Errorf("expected empty destination fields, got %+v", got)
		}
		if got.PlannedDeliveryAt != "" || got.ActualDeliveryAt != "" ||
			got.CreatedAt != "" || got.UpdatedAt != "" {
			t.Errorf("expected empty timestamp fields, got %+v", got)
		}
		if got.Status != aggregate.SHIPMENTSTATUS_UNKNOWN {
			t.Errorf("Status = %v, want UNKNOWN for unspecified proto status", got.Status)
		}
	})
}

func TestToShipmentStatus_WhenProtoStatusGiven_ReturnsDomainStatus(t *testing.T) {
	cases := []struct {
		in   deliverypb.ShipmentStatus
		want aggregate.ShipmentStatus
	}{
		{deliverypb.ShipmentStatus_SHIPMENT_STATUS_PLANNED, aggregate.SHIPMENTSTATUS_PLANNED},
		{deliverypb.ShipmentStatus_SHIPMENT_STATUS_DISPATCHED, aggregate.SHIPMENTSTATUS_DISPATCHED},
		{deliverypb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP, aggregate.SHIPMENTSTATUS_PICKED_UP},
		{deliverypb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT, aggregate.SHIPMENTSTATUS_IN_TRANSIT},
		{deliverypb.ShipmentStatus_SHIPMENT_STATUS_OUT_FOR_DELIVERY, aggregate.SHIPMENTSTATUS_OUT_FOR_DELIVERY},
		{deliverypb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED, aggregate.SHIPMENTSTATUS_DELIVERED},
		{deliverypb.ShipmentStatus_SHIPMENT_STATUS_FAILED, aggregate.SHIPMENTSTATUS_FAILED},
		{deliverypb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED, aggregate.SHIPMENTSTATUS_CANCELLED},
		{deliverypb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED, aggregate.SHIPMENTSTATUS_UNKNOWN},
	}
	for _, c := range cases {
		if got := toShipmentStatus(c.in); got != c.want {
			t.Errorf("toShipmentStatus(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestStringToShipmentStatusProto_WhenDomainStatusGiven_ReturnsProtoStatus(t *testing.T) {
	cases := []struct {
		in   *string
		want deliverypb.ShipmentStatus
	}{
		{nil, deliverypb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED},
		{ptrStr("PLANNED"), deliverypb.ShipmentStatus_SHIPMENT_STATUS_PLANNED},
		{ptrStr("DISPATCHED"), deliverypb.ShipmentStatus_SHIPMENT_STATUS_DISPATCHED},
		{ptrStr("PICKED_UP"), deliverypb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP},
		{ptrStr("IN_TRANSIT"), deliverypb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT},
		{ptrStr("OUT_FOR_DELIVERY"), deliverypb.ShipmentStatus_SHIPMENT_STATUS_OUT_FOR_DELIVERY},
		{ptrStr("DELIVERED"), deliverypb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED},
		{ptrStr("FAILED"), deliverypb.ShipmentStatus_SHIPMENT_STATUS_FAILED},
		{ptrStr("CANCELLED"), deliverypb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED},
		{ptrStr(""), deliverypb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED},
		{ptrStr("garbage"), deliverypb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED},
	}
	for _, c := range cases {
		if got := stringToShipmentStatusProto(c.in); got != c.want {
			t.Errorf("stringToShipmentStatusProto(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}
