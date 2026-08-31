package service

import (
	"context"
	"errors"
	"testing"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
)

// ─── TrackShipment / ListOrderShipments (delivery proxy) ─────────────────────

func ownerOf(ownerUid string) *stubOrderClient {
	return &stubOrderClient{
		findOrderFn: func(_ context.Context, uid string) (*aggregate.OrderAggregate, error) {
			return &aggregate.OrderAggregate{Uid: uid, UserUid: ownerUid}, nil
		},
	}
}

const callerUid = "u-1"

var theCaller = usecase.Caller{UserUid: callerUid}

func TestTrackShipment_WhenDeliveryReturnsShipment_ReturnsAggregate(t *testing.T) {
	t.Run("whenDeliveryReturnsShipment_mapsEveryField", func(t *testing.T) {
		svc := &customerService{orderClient: ownerOf(callerUid), deliveryClient: &stubDeliveryClient{
			trackShipmentFn: func(_ context.Context, uid string, gotCaller string) (*aggregate.ShipmentAggregate, error) {
				if uid != "sh-1" {
					t.Errorf("TrackShipment uid = %q, want sh-1", uid)
				}
				if gotCaller != callerUid {
					t.Errorf("TrackShipment callerUid = %q, want %q", gotCaller, callerUid)
				}
				return &aggregate.ShipmentAggregate{
					Uid:      "sh-1",
					OrderUid: "o-1",
					Status:   aggregate.SHIPMENTSTATUS_DELIVERED,
				}, nil
			},
		}}
		got, err := svc.TrackShipment(context.Background(), "sh-1", theCaller)
		if err != nil {
			t.Fatalf("TrackShipment err: %v", err)
		}
		if got == nil || got.Uid != "sh-1" || got.OrderUid != "o-1" {
			t.Errorf("got = %+v", got)
		}
		if got.Status != aggregate.SHIPMENTSTATUS_DELIVERED {
			t.Errorf("status = %v, want DELIVERED", got.Status)
		}
	})

	t.Run("whenShipmentNil_returnsNilAggregate", func(t *testing.T) {
		svc := &customerService{orderClient: ownerOf(callerUid), deliveryClient: &stubDeliveryClient{
			trackShipmentFn: func(context.Context, string, string) (*aggregate.ShipmentAggregate, error) {
				return nil, nil
			},
		}}
		got, err := svc.TrackShipment(context.Background(), "missing", theCaller)
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		if got != nil {
			t.Errorf("got = %+v, want nil", got)
		}
	})

	t.Run("whenDeliveryFails_returnsWrappedError", func(t *testing.T) {
		wantErr := errors.New("delivery down")
		svc := &customerService{orderClient: ownerOf(callerUid), deliveryClient: &stubDeliveryClient{
			trackShipmentFn: func(context.Context, string, string) (*aggregate.ShipmentAggregate, error) {
				return nil, wantErr
			},
		}}
		if _, err := svc.TrackShipment(context.Background(), "sh-1", theCaller); !errors.Is(err, wantErr) {
			t.Errorf("err = %v, want wrapped %v", err, wantErr)
		}
	})
}

func TestListOrderShipments_WhenFiltersGiven_ForwardsFiltersAndMapsShipments(t *testing.T) {
	t.Run("whenFiltersGiven_forwardsFiltersAndMapsShipments", func(t *testing.T) {
		svc := &customerService{orderClient: ownerOf(callerUid), deliveryClient: &stubDeliveryClient{
			listShipmentsFn: func(_ context.Context, orderUid string, status *string, pageSize *int32, _ string) (int64, []*aggregate.ShipmentAggregate, error) {
				if orderUid != "o-1" {
					t.Errorf("ListShipments orderUid = %q, want o-1", orderUid)
				}
				if status == nil || *status != "IN_TRANSIT" {
					t.Errorf("ListShipments status = %v, want IN_TRANSIT", status)
				}
				if pageSize == nil || *pageSize != 25 {
					t.Errorf("ListShipments pageSize = %v, want 25", pageSize)
				}
				return 2, []*aggregate.ShipmentAggregate{
					{Uid: "sh-1", OrderUid: "o-1"},
					{Uid: "sh-2", OrderUid: "o-1"},
				}, nil
			},
		}}
		status := "IN_TRANSIT"
		size := int32(25)
		total, aggs, err := svc.ListOrderShipments(context.Background(), "o-1", &status, &size, theCaller)
		if err != nil {
			t.Fatalf("ListOrderShipments err: %v", err)
		}
		if total != 2 || len(aggs) != 2 {
			t.Fatalf("total=%d len=%d, want 2/2", total, len(aggs))
		}
		if aggs[0].Uid != "sh-1" || aggs[1].Uid != "sh-2" {
			t.Errorf("aggs = %+v", aggs)
		}
	})

	t.Run("whenStatusAndSizeNil_omitsFilters", func(t *testing.T) {
		svc := &customerService{orderClient: ownerOf(callerUid), deliveryClient: &stubDeliveryClient{
			listShipmentsFn: func(_ context.Context, _ string, status *string, pageSize *int32, _ string) (int64, []*aggregate.ShipmentAggregate, error) {
				if status != nil {
					t.Errorf("status = %v, want nil", status)
				}
				if pageSize != nil {
					t.Errorf("pageSize = %v, want nil", pageSize)
				}
				return 0, nil, nil
			},
		}}
		total, aggs, err := svc.ListOrderShipments(context.Background(), "o-1", nil, nil, theCaller)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if total != 0 || len(aggs) != 0 {
			t.Errorf("total=%d aggs=%+v, want 0/empty", total, aggs)
		}
	})

	t.Run("whenDeliveryFails_returnsWrappedError", func(t *testing.T) {
		wantErr := errors.New("delivery down")
		svc := &customerService{orderClient: ownerOf(callerUid), deliveryClient: &stubDeliveryClient{
			listShipmentsFn: func(context.Context, string, *string, *int32, string) (int64, []*aggregate.ShipmentAggregate, error) {
				return 0, nil, wantErr
			},
		}}
		if _, _, err := svc.ListOrderShipments(context.Background(), "o-1", nil, nil, theCaller); !errors.Is(err, wantErr) {
			t.Errorf("err = %v, want wrapped %v", err, wantErr)
		}
	})
}

// ─── Ownership (BOLA) ────────────────────────────────────────────────────────

func TestTrackShipment_WhenShipmentBelongsToAnotherCustomer_ReturnsNil(t *testing.T) {
	delivery := &stubDeliveryClient{
		trackShipmentFn: func(context.Context, string, string) (*aggregate.ShipmentAggregate, error) {
			return &aggregate.ShipmentAggregate{Uid: "sh-1", OrderUid: "o-1"}, nil
		},
	}

	t.Run("whenOwnedByAnotherCustomer_returnsNilLikeAMissingShipment", func(t *testing.T) {
		svc := &customerService{orderClient: ownerOf("someone-else"), deliveryClient: delivery}

		got, err := svc.TrackShipment(context.Background(), "sh-1", theCaller)

		if err != nil {
			t.Fatalf("err = %v, want nil (a 403 would confirm the shipment exists)", err)
		}
		if got != nil {
			t.Fatalf("got = %+v, want nil — leaked another customer's shipment", got)
		}
	})

	t.Run("whenCallerIsAdmin_returnsAnyShipment", func(t *testing.T) {
		svc := &customerService{orderClient: ownerOf("someone-else"), deliveryClient: delivery}

		got, err := svc.TrackShipment(context.Background(), "sh-1", usecase.Caller{UserUid: "ops", IsAdmin: true})

		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if got == nil || got.Uid != "sh-1" {
			t.Fatalf("got = %+v, want the shipment", got)
		}
	})

	t.Run("whenCallerUnauthenticated_returnsNil", func(t *testing.T) {
		svc := &customerService{orderClient: ownerOf(callerUid), deliveryClient: delivery}

		got, err := svc.TrackShipment(context.Background(), "sh-1", usecase.Caller{})

		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		if got != nil {
			t.Fatalf("got = %+v, want nil — a zero Caller must own nothing", got)
		}
	})
}

func TestListOrderShipments_WhenOrderBelongsToAnotherCustomer_ReturnsEmptyWithoutCallingDelivery(t *testing.T) {
	t.Run("whenOrderOwnedByAnotherCustomer_returnsEmptyWithoutCallingDelivery", func(t *testing.T) {
		listed := false
		svc := &customerService{
			orderClient: ownerOf("someone-else"),
			deliveryClient: &stubDeliveryClient{
				listShipmentsFn: func(context.Context, string, *string, *int32, string) (int64, []*aggregate.ShipmentAggregate, error) {
					listed = true
					return 1, []*aggregate.ShipmentAggregate{{Uid: "sh-1", OrderUid: "o-1"}}, nil
				},
			},
		}

		total, aggs, err := svc.ListOrderShipments(context.Background(), "o-1", nil, nil, theCaller)

		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		if total != 0 || len(aggs) != 0 {
			t.Fatalf("total=%d aggs=%+v, want 0/empty — leaked another customer's shipments", total, aggs)
		}
		if listed {
			t.Error("delivery was queried for an order the caller does not own")
		}
	})

	t.Run("whenOrderUnknown_returnsEmpty", func(t *testing.T) {
		svc := &customerService{
			orderClient: &stubOrderClient{
				findOrderFn: func(context.Context, string) (*aggregate.OrderAggregate, error) {
					return nil, nil
				},
			},
			deliveryClient: &stubDeliveryClient{},
		}

		total, aggs, err := svc.ListOrderShipments(context.Background(), "nope", nil, nil, theCaller)

		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		if total != 0 || len(aggs) != 0 {
			t.Errorf("total=%d aggs=%+v, want 0/empty", total, aggs)
		}
	})
}
