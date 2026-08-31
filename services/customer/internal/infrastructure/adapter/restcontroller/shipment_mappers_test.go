package restcontroller

import (
	"testing"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
)

func TestToShipmentResponse_WhenAggregateFullyPopulated_MapsEveryField(t *testing.T) {
	t.Run("whenAggregateFullyPopulated_mapsEveryField", func(t *testing.T) {
		a := &aggregate.ShipmentAggregate{
			Uid:                    "sh-1",
			OrderUid:               "o-1",
			CarrierUid:             "car-1",
			Status:                 aggregate.SHIPMENTSTATUS_OUT_FOR_DELIVERY,
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
		got := toShipmentResponse(a)

		if got.Uid != "sh-1" {
			t.Errorf("Uid = %q, want sh-1", got.Uid)
		}
		if got.Status == nil || *got.Status != openapi.SHIPMENTSTATUS_OUT_FOR_DELIVERY {
			t.Errorf("Status = %v, want OUT_FOR_DELIVERY", got.Status)
		}
		assertStrPtr(t, "OrderUid", got.OrderUid, "o-1")
		assertStrPtr(t, "CarrierUid", got.CarrierUid, "car-1")
		assertStrPtr(t, "TrackingNumber", got.TrackingNumber, "1Z-TRACK-9")
		assertStrPtr(t, "DestinationCity", got.DestinationCity, "Seattle")
		assertStrPtr(t, "DestinationState", got.DestinationState, "WA")
		assertStrPtr(t, "DestinationCountryCode", got.DestinationCountryCode, "US")
		assertStrPtr(t, "DestinationPostalCode", got.DestinationPostalCode, "98101")
		assertStrPtr(t, "PlannedDeliveryAt", got.PlannedDeliveryAt, "2026-07-01T10:00:00Z")
		assertStrPtr(t, "ActualDeliveryAt", got.ActualDeliveryAt, "2026-07-02T14:30:00Z")
		assertStrPtr(t, "CreatedAt", got.CreatedAt, "2026-06-28T08:00:00Z")
		assertStrPtr(t, "UpdatedAt", got.UpdatedAt, "2026-06-29T09:00:00Z")
	})

	t.Run("whenOptionalStringsEmpty_returnsNilPointers", func(t *testing.T) {
		got := toShipmentResponse(&aggregate.ShipmentAggregate{Uid: "sh-2"})
		if got.Uid != "sh-2" {
			t.Errorf("Uid = %q, want sh-2", got.Uid)
		}
		nilable := map[string]*string{
			"OrderUid":               got.OrderUid,
			"CarrierUid":             got.CarrierUid,
			"TrackingNumber":         got.TrackingNumber,
			"DestinationCity":        got.DestinationCity,
			"DestinationState":       got.DestinationState,
			"DestinationCountryCode": got.DestinationCountryCode,
			"DestinationPostalCode":  got.DestinationPostalCode,
			"PlannedDeliveryAt":      got.PlannedDeliveryAt,
			"ActualDeliveryAt":       got.ActualDeliveryAt,
			"CreatedAt":              got.CreatedAt,
			"UpdatedAt":              got.UpdatedAt,
		}
		for field, ptr := range nilable {
			if ptr != nil {
				t.Errorf("%s = %q, want nil for empty string", field, *ptr)
			}
		}
		if got.Status == nil || *got.Status != openapi.SHIPMENTSTATUS_UNKNOWN {
			t.Errorf("Status = %v, want UNKNOWN", got.Status)
		}
	})
}

func TestToShipmentStatusRest_WhenDomainStatusGiven_ReturnsRestStatus(t *testing.T) {
	cases := []struct {
		in   aggregate.ShipmentStatus
		want openapi.ShipmentStatus
	}{
		{aggregate.SHIPMENTSTATUS_PLANNED, openapi.SHIPMENTSTATUS_PLANNED},
		{aggregate.SHIPMENTSTATUS_DISPATCHED, openapi.SHIPMENTSTATUS_DISPATCHED},
		{aggregate.SHIPMENTSTATUS_PICKED_UP, openapi.SHIPMENTSTATUS_PICKED_UP},
		{aggregate.SHIPMENTSTATUS_IN_TRANSIT, openapi.SHIPMENTSTATUS_IN_TRANSIT},
		{aggregate.SHIPMENTSTATUS_OUT_FOR_DELIVERY, openapi.SHIPMENTSTATUS_OUT_FOR_DELIVERY},
		{aggregate.SHIPMENTSTATUS_DELIVERED, openapi.SHIPMENTSTATUS_DELIVERED},
		{aggregate.SHIPMENTSTATUS_FAILED, openapi.SHIPMENTSTATUS_FAILED},
		{aggregate.SHIPMENTSTATUS_CANCELLED, openapi.SHIPMENTSTATUS_CANCELLED},
		{aggregate.SHIPMENTSTATUS_UNKNOWN, openapi.SHIPMENTSTATUS_UNKNOWN},
	}
	for _, c := range cases {
		if got := toShipmentStatusRest(c.in); got != c.want {
			t.Errorf("toShipmentStatusRest(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func assertStrPtr(t *testing.T, field string, got *string, want string) {
	t.Helper()
	if got == nil || *got != want {
		t.Errorf("%s = %v, want %q", field, got, want)
	}
}
