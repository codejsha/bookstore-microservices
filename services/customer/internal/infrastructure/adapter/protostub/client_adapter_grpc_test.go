package protostub

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/fulfillmentpb"
	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/orderpb"
	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/paymentpb"
	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/config"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
)

const (
	callerUid = "018f0000-0000-7000-8000-0000000000c1"
	targetUid = "018f0000-0000-7000-8000-0000000000e1"
)

func dialInProcess(t *testing.T, register func(*grpc.Server)) *grpc.ClientConn {
	t.Helper()
	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer()
	register(srv)
	go func() { _ = srv.Serve(lis) }()

	opts := append(baseDialOptions(), grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.DialContext(ctx)
	}))
	conn, err := grpc.NewClient("passthrough:///in-process", opts...)
	if err != nil {
		t.Fatalf("create in-process client: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
		srv.Stop()
	})
	return conn
}

func requestContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func incoming(ctx context.Context) metadata.MD {
	md, _ := metadata.FromIncomingContext(ctx)
	return md
}

func assertGrpcCode(t *testing.T, err error, want codes.Code, wantMessage string) {
	t.Helper()
	st, ok := status.FromError(err)
	if err == nil || !ok {
		t.Fatalf("err = %v, want a grpc status error", err)
	}
	if st.Code() != want || st.Message() != wantMessage {
		t.Fatalf("status = %v %q, want %v %q", st.Code(), st.Message(), want, wantMessage)
	}
}

var upstreamCodes = []codes.Code{
	codes.NotFound, codes.PermissionDenied, codes.Unauthenticated, codes.InvalidArgument, codes.Unavailable, codes.Internal,
}

type fakeFulfillmentServer struct {
	fulfillmentpb.UnimplementedFulfillmentServiceServer
	err      error
	delay    time.Duration
	track    *fulfillmentpb.TrackShipmentResponse
	list     *fulfillmentpb.ListShipmentsResponse
	gotTrack *fulfillmentpb.TrackShipmentRequest
	gotList  *fulfillmentpb.ListShipmentsRequest
	gotMD    metadata.MD
}

func (s *fakeFulfillmentServer) TrackShipment(ctx context.Context, in *fulfillmentpb.TrackShipmentRequest) (*fulfillmentpb.TrackShipmentResponse, error) {
	s.gotTrack, s.gotMD = in, incoming(ctx)
	if s.delay > 0 {
		select {
		case <-time.After(s.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if s.err != nil {
		return nil, s.err
	}
	return s.track, nil
}

func (s *fakeFulfillmentServer) ListShipments(ctx context.Context, in *fulfillmentpb.ListShipmentsRequest) (*fulfillmentpb.ListShipmentsResponse, error) {
	s.gotList, s.gotMD = in, incoming(ctx)
	if s.err != nil {
		return nil, s.err
	}
	return s.list, nil
}

func newDeliveryAdapter(t *testing.T, srv *fakeFulfillmentServer) *deliveryClient {
	t.Helper()
	conn := dialInProcess(t, func(s *grpc.Server) { fulfillmentpb.RegisterFulfillmentServiceServer(s, srv) })
	adapter, ok := NewDeliveryClient(&DeliveryGrpcClient{Client: fulfillmentpb.NewFulfillmentServiceClient(conn)}).(*deliveryClient)
	if !ok {
		t.Fatal("NewDeliveryClient did not build the grpc adapter")
	}
	return adapter
}

func TestDeliveryClientTrackShipment_shipmentFound_aggregateMappedAndCallerSentAsUserId(t *testing.T) {
	srv := &fakeFulfillmentServer{track: &fulfillmentpb.TrackShipmentResponse{Shipment: &fulfillmentpb.Shipment{
		Uid: targetUid, OrderUid: "o-1", Status: fulfillmentpb.ShipmentStatus_SHIPMENT_STATUS_OUT_FOR_DELIVERY, TrackingNumber: "1Z-9",
	}}}
	adapter := newDeliveryAdapter(t, srv)

	got, err := adapter.TrackShipment(requestContext(t), targetUid, callerUid)
	if err != nil {
		t.Fatalf("track shipment: %v", err)
	}

	want := &aggregate.ShipmentAggregate{Uid: targetUid, OrderUid: "o-1", Status: aggregate.SHIPMENTSTATUS_OUT_FOR_DELIVERY, TrackingNumber: "1Z-9"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %+v, want %+v", got, want)
	}
	if srv.gotTrack.GetUid() != targetUid {
		t.Errorf("requested uid = %q, want %q", srv.gotTrack.GetUid(), targetUid)
	}
	if ids := srv.gotMD.Get(deliveryCallerMetadataKey); !reflect.DeepEqual(ids, []string{callerUid}) {
		t.Errorf("x-user-id metadata = %v, want [%s]", ids, callerUid)
	}
}

func TestDeliveryClientTrackShipment_responseWithoutShipment_nilWithoutError(t *testing.T) {
	adapter := newDeliveryAdapter(t, &fakeFulfillmentServer{track: &fulfillmentpb.TrackShipmentResponse{}})

	got, err := adapter.TrackShipment(requestContext(t), targetUid, callerUid)
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil", got, err)
	}
}

func TestDeliveryClient_upstreamStatus_codeAndMessagePreserved(t *testing.T) {
	for _, code := range upstreamCodes {
		t.Run(code.String(), func(t *testing.T) {
			adapter := newDeliveryAdapter(t, &fakeFulfillmentServer{err: status.Error(code, "fulfillment says no")})

			tracked, err := adapter.TrackShipment(requestContext(t), targetUid, callerUid)
			assertGrpcCode(t, err, code, "fulfillment says no")
			if tracked != nil {
				t.Errorf("tracked = %+v, want nil", tracked)
			}

			total, listed, err := adapter.ListShipments(requestContext(t), "o-1", nil, nil, callerUid)
			assertGrpcCode(t, err, code, "fulfillment says no")
			if total != 0 || listed != nil {
				t.Errorf("list = %d/%v, want zero values", total, listed)
			}
		})
	}
}

func TestDeliveryClientTrackShipment_upstreamSlowerThanDeadline_deadlineExceeded(t *testing.T) {
	adapter := newDeliveryAdapter(t, &fakeFulfillmentServer{delay: 3 * time.Second, track: &fulfillmentpb.TrackShipmentResponse{}})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := adapter.TrackShipment(ctx, targetUid, callerUid)

	if status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("err = %v, want DeadlineExceeded", err)
	}
}

func TestDeliveryClientListShipments_statusAndPageSizeGiven_requestBuiltAndPageMapped(t *testing.T) {
	srv := &fakeFulfillmentServer{list: &fulfillmentpb.ListShipmentsResponse{
		TotalSize: 9,
		Shipments: []*fulfillmentpb.Shipment{
			{Uid: "sh-1", Status: fulfillmentpb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED},
			{Uid: "sh-2", Status: fulfillmentpb.ShipmentStatus_SHIPMENT_STATUS_FAILED},
		},
	}}
	adapter := newDeliveryAdapter(t, srv)
	delivered, size := "DELIVERED", int32(25)

	total, got, err := adapter.ListShipments(requestContext(t), "o-1", &delivered, &size, callerUid)
	if err != nil {
		t.Fatalf("list shipments: %v", err)
	}

	if total != 9 || len(got) != 2 || got[0].Uid != "sh-1" || got[0].Status != aggregate.SHIPMENTSTATUS_DELIVERED ||
		got[1].Uid != "sh-2" || got[1].Status != aggregate.SHIPMENTSTATUS_FAILED {
		t.Errorf("total = %d, shipments = %+v", total, got)
	}
	req := srv.gotList
	if req.GetOrderUid() != "o-1" || req.GetStatus() != fulfillmentpb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED || req.GetPageSize() != 25 {
		t.Errorf("request = %v", req)
	}
	if ids := srv.gotMD.Get(deliveryCallerMetadataKey); !reflect.DeepEqual(ids, []string{callerUid}) {
		t.Errorf("x-user-id metadata = %v, want [%s]", ids, callerUid)
	}
}

func TestDeliveryClientListShipments_noFilters_unspecifiedStatusAndServerDefaultPageSize(t *testing.T) {
	srv := &fakeFulfillmentServer{list: &fulfillmentpb.ListShipmentsResponse{}}
	adapter := newDeliveryAdapter(t, srv)

	total, got, err := adapter.ListShipments(requestContext(t), "o-1", nil, nil, callerUid)
	if err != nil {
		t.Fatalf("list shipments: %v", err)
	}

	if total != 0 || got == nil || len(got) != 0 {
		t.Errorf("total = %d, shipments = %v, want an empty non-nil page", total, got)
	}
	if srv.gotList.GetStatus() != fulfillmentpb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED || srv.gotList.GetPageSize() != 0 {
		t.Errorf("request = %v, want no status filter and no page size", srv.gotList)
	}
}

type fakeOrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	err     error
	list    *orderpb.ListOrdersResponse
	find    *orderpb.FindOrderResponse
	gotList *orderpb.ListOrdersRequest
	gotFind *orderpb.FindOrderRequest
	gotMD   metadata.MD
}

func (s *fakeOrderServer) ListOrders(ctx context.Context, in *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	s.gotList, s.gotMD = in, incoming(ctx)
	if s.err != nil {
		return nil, s.err
	}
	return s.list, nil
}

func (s *fakeOrderServer) FindOrder(ctx context.Context, in *orderpb.FindOrderRequest) (*orderpb.FindOrderResponse, error) {
	s.gotFind, s.gotMD = in, incoming(ctx)
	if s.err != nil {
		return nil, s.err
	}
	return s.find, nil
}

func newOrderAdapter(t *testing.T, srv *fakeOrderServer) *orderClient {
	t.Helper()
	conn := dialInProcess(t, func(s *grpc.Server) { orderpb.RegisterOrderServiceServer(s, srv) })
	adapter, ok := NewOrderClient(&OrderGrpcClient{Client: orderpb.NewOrderServiceClient(conn)}).(*orderClient)
	if !ok {
		t.Fatal("NewOrderClient did not build the grpc adapter")
	}
	return adapter
}

func assertActor(t *testing.T, md metadata.MD, wantUid []string, wantAdmin []string) {
	t.Helper()
	if got := md.Get(actorUidMetadataKey); !reflect.DeepEqual(got, wantUid) {
		t.Errorf("x-actor-uid metadata = %v, want %v", got, wantUid)
	}
	if got := md.Get(actorAdminMetadataKey); !reflect.DeepEqual(got, wantAdmin) {
		t.Errorf("x-actor-admin metadata = %v, want %v", got, wantAdmin)
	}
}

func TestOrderClientListOrders_customerCaller_requestBuiltPageMappedAndActorSentWithoutAdminFlag(t *testing.T) {
	srv := &fakeOrderServer{list: &orderpb.ListOrdersResponse{
		TotalSize: 31,
		Orders: []*orderpb.Order{
			{Uid: "o-1", UserUid: callerUid, TotalAmount: 12.5, Status: orderpb.OrderStatus_ORDER_STATUS_PROCESSING,
				LineItems: []*orderpb.OrderLineItem{{ProductUid: "b-1", Quantity: 2}}},
			{Uid: "o-2", UserUid: callerUid, Status: orderpb.OrderStatus_ORDER_STATUS_REFUNDED},
		},
	}}
	adapter := newOrderAdapter(t, srv)
	ctx := httpx.WithCaller(requestContext(t), httpx.Caller{UserUid: callerUid})
	paid := "PAID"

	total, got, err := adapter.ListOrders(ctx, callerUid, &paid, pageOf(10, 2, "created_at,desc"))
	if err != nil {
		t.Fatalf("list orders: %v", err)
	}

	if total != 31 || len(got) != 2 {
		t.Fatalf("total = %d, orders = %d, want 31 and 2", total, len(got))
	}
	want := &aggregate.OrderAggregate{
		Uid: "o-1", UserUid: callerUid, TotalAmount: 12.5, Status: aggregate.ORDERSTATUS_PAID,
		OrderItems: []aggregate.OrderItem{{BookUid: "b-1", Quantity: 2}},
	}
	if !reflect.DeepEqual(got[0], want) {
		t.Errorf("first order = %+v, want %+v", got[0], want)
	}
	if got[1].Status != aggregate.ORDERSTATUS_CANCELLED || len(got[1].OrderItems) != 0 {
		t.Errorf("second order = %+v, want cancelled without items", got[1])
	}
	req := srv.gotList
	if req.GetUserUid() != callerUid || req.GetStatus() != orderpb.OrderStatus_ORDER_STATUS_CONFIRMED ||
		req.GetPageSize() != 10 || req.GetPageToken() != "2" || req.GetOrderBy() != "created_at,desc" {
		t.Errorf("request = %v", req)
	}
	assertActor(t, srv.gotMD, []string{callerUid}, nil)
}

func TestOrderClientFindOrder_adminCaller_orderMappedAndActorSentWithAdminFlag(t *testing.T) {
	srv := &fakeOrderServer{find: &orderpb.FindOrderResponse{Order: &orderpb.Order{
		Uid: targetUid, UserUid: "u-9", Status: orderpb.OrderStatus_ORDER_STATUS_DELIVERED, TotalAmount: 40,
	}}}
	adapter := newOrderAdapter(t, srv)
	ctx := httpx.WithCaller(requestContext(t), httpx.Caller{UserUid: callerUid, IsAdmin: true})

	got, err := adapter.FindOrder(ctx, targetUid)
	if err != nil {
		t.Fatalf("find order: %v", err)
	}

	if got == nil || got.Uid != targetUid || got.UserUid != "u-9" || got.Status != aggregate.ORDERSTATUS_COMPLETED || got.TotalAmount != 40 {
		t.Errorf("got = %+v", got)
	}
	if srv.gotFind.GetUid() != targetUid {
		t.Errorf("requested uid = %q, want %q", srv.gotFind.GetUid(), targetUid)
	}
	assertActor(t, srv.gotMD, []string{callerUid}, []string{actorAdminMetadataValue})
}

func TestOrderClientFindOrder_noCallerInContext_noActorMetadataSent(t *testing.T) {
	srv := &fakeOrderServer{find: &orderpb.FindOrderResponse{}}
	adapter := newOrderAdapter(t, srv)

	got, err := adapter.FindOrder(requestContext(t), targetUid)
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil for a response without an order", got, err)
	}
	assertActor(t, srv.gotMD, nil, nil)
}

func TestOrderClient_upstreamStatus_codeAndMessagePreserved(t *testing.T) {
	for _, code := range upstreamCodes {
		t.Run(code.String(), func(t *testing.T) {
			adapter := newOrderAdapter(t, &fakeOrderServer{err: status.Error(code, "order says no")})

			found, err := adapter.FindOrder(requestContext(t), targetUid)
			assertGrpcCode(t, err, code, "order says no")
			if found != nil {
				t.Errorf("found = %+v, want nil", found)
			}

			total, listed, err := adapter.ListOrders(requestContext(t), callerUid, nil, pageOf(10, 0, ""))
			assertGrpcCode(t, err, code, "order says no")
			if total != 0 || listed != nil {
				t.Errorf("list = %d/%v, want zero values", total, listed)
			}
		})
	}
}

type fakePaymentServer struct {
	paymentpb.UnimplementedPaymentServiceServer
	err     error
	list    *paymentpb.ListPaymentsResponse
	find    *paymentpb.FindPaymentResponse
	gotList *paymentpb.ListPaymentsRequest
	gotFind *paymentpb.FindPaymentRequest
	gotMD   metadata.MD
}

func (s *fakePaymentServer) ListPayments(ctx context.Context, in *paymentpb.ListPaymentsRequest) (*paymentpb.ListPaymentsResponse, error) {
	s.gotList, s.gotMD = in, incoming(ctx)
	if s.err != nil {
		return nil, s.err
	}
	return s.list, nil
}

func (s *fakePaymentServer) FindPayment(ctx context.Context, in *paymentpb.FindPaymentRequest) (*paymentpb.FindPaymentResponse, error) {
	s.gotFind, s.gotMD = in, incoming(ctx)
	if s.err != nil {
		return nil, s.err
	}
	return s.find, nil
}

func newPaymentAdapter(t *testing.T, srv *fakePaymentServer) *paymentClient {
	t.Helper()
	conn := dialInProcess(t, func(s *grpc.Server) { paymentpb.RegisterPaymentServiceServer(s, srv) })
	adapter, ok := NewPaymentClient(&PaymentGrpcClient{Client: paymentpb.NewPaymentServiceClient(conn)}).(*paymentClient)
	if !ok {
		t.Fatal("NewPaymentClient did not build the grpc adapter")
	}
	return adapter
}

func TestPaymentClientListPayments_customerCaller_requestBuiltPageMappedAndActorSent(t *testing.T) {
	srv := &fakePaymentServer{list: &paymentpb.ListPaymentsResponse{
		TotalSize: 3,
		Payments: []*paymentpb.Payment{
			{Uid: "p-1", PaymentUid: "px-1", CustomerUid: callerUid, Currency: "USD", Amount: 4200,
				Status: paymentpb.PaymentStatus_PAYMENT_STATUS_SUCCEEDED, PaymentMethod: "card"},
			{Uid: "p-2", CustomerUid: callerUid, Status: paymentpb.PaymentStatus_PAYMENT_STATUS_FAILED,
				ErrorCode: "card_declined", ErrorMessage: "insufficient funds"},
		},
	}}
	adapter := newPaymentAdapter(t, srv)
	ctx := httpx.WithCaller(requestContext(t), httpx.Caller{UserUid: callerUid})

	total, got, err := adapter.ListPayments(ctx, callerUid, pageOf(5, 1, "created_at,asc"))
	if err != nil {
		t.Fatalf("list payments: %v", err)
	}

	if total != 3 || len(got) != 2 {
		t.Fatalf("total = %d, payments = %d, want 3 and 2", total, len(got))
	}
	if p := got[0]; p.Uid != "p-1" || p.PaymentUid != "px-1" || p.CustomerUid != callerUid || p.Currency != "USD" || p.Amount != 4200 ||
		p.Status != "PAYMENT_STATUS_SUCCEEDED" || p.PaymentMethod == nil || *p.PaymentMethod != "card" || p.ErrorCode != nil {
		t.Errorf("first payment = %+v", p)
	}
	if p := got[1]; p.PaymentMethod != nil || p.ErrorCode == nil || *p.ErrorCode != "card_declined" || p.ErrorMessage == nil || *p.ErrorMessage != "insufficient funds" {
		t.Errorf("second payment = %+v", p)
	}
	req := srv.gotList
	if req.GetCustomerUid() != callerUid || req.GetPageSize() != 5 || req.GetPageToken() != "1" || req.GetOrderBy() != "created_at,asc" {
		t.Errorf("request = %v", req)
	}
	assertActor(t, srv.gotMD, []string{callerUid}, nil)
}

func TestPaymentClientFindPayment_adminCaller_paymentMappedAndActorSentWithAdminFlag(t *testing.T) {
	srv := &fakePaymentServer{find: &paymentpb.FindPaymentResponse{Payment: &paymentpb.Payment{
		Uid: targetUid, CustomerUid: "u-9", Currency: "KRW", Amount: 15000, Status: paymentpb.PaymentStatus_PAYMENT_STATUS_SUCCEEDED,
	}}}
	adapter := newPaymentAdapter(t, srv)
	ctx := httpx.WithCaller(requestContext(t), httpx.Caller{UserUid: callerUid, IsAdmin: true})

	got, err := adapter.FindPayment(ctx, targetUid)
	if err != nil {
		t.Fatalf("find payment: %v", err)
	}

	if got == nil || got.Uid != targetUid || got.CustomerUid != "u-9" || got.Currency != "KRW" || got.Amount != 15000 {
		t.Errorf("got = %+v", got)
	}
	if srv.gotFind.GetUid() != targetUid {
		t.Errorf("requested uid = %q, want %q", srv.gotFind.GetUid(), targetUid)
	}
	assertActor(t, srv.gotMD, []string{callerUid}, []string{actorAdminMetadataValue})
}

func TestPaymentClientFindPayment_responseWithoutPayment_nilWithoutError(t *testing.T) {
	adapter := newPaymentAdapter(t, &fakePaymentServer{find: &paymentpb.FindPaymentResponse{}})

	got, err := adapter.FindPayment(requestContext(t), targetUid)
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil", got, err)
	}
}

func TestPaymentClient_upstreamStatus_codeAndMessagePreserved(t *testing.T) {
	for _, code := range upstreamCodes {
		t.Run(code.String(), func(t *testing.T) {
			adapter := newPaymentAdapter(t, &fakePaymentServer{err: status.Error(code, "payment says no")})

			found, err := adapter.FindPayment(requestContext(t), targetUid)
			assertGrpcCode(t, err, code, "payment says no")
			if found != nil {
				t.Errorf("found = %+v, want nil", found)
			}

			total, listed, err := adapter.ListPayments(requestContext(t), callerUid, pageOf(10, 0, ""))
			assertGrpcCode(t, err, code, "payment says no")
			if total != 0 || listed != nil {
				t.Errorf("list = %d/%v, want zero values", total, listed)
			}
		})
	}
}

type fakeUserServer struct {
	userpb.UnimplementedUserServiceServer
	err     error
	list    *userpb.ListUsersResponse
	find    *userpb.FindUserResponse
	gotList *userpb.ListUsersRequest
	gotFind *userpb.FindUserRequest
	gotMD   metadata.MD
}

func (s *fakeUserServer) ListUsers(ctx context.Context, in *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	s.gotList, s.gotMD = in, incoming(ctx)
	if s.err != nil {
		return nil, s.err
	}
	return s.list, nil
}

func (s *fakeUserServer) FindUser(ctx context.Context, in *userpb.FindUserRequest) (*userpb.FindUserResponse, error) {
	s.gotFind, s.gotMD = in, incoming(ctx)
	if s.err != nil {
		return nil, s.err
	}
	return s.find, nil
}

func newUserAdapter(t *testing.T, srv *fakeUserServer) *userClient {
	t.Helper()
	conn := dialInProcess(t, func(s *grpc.Server) { userpb.RegisterUserServiceServer(s, srv) })
	adapter, ok := NewUserClient(&UserGrpcClient{Client: userpb.NewUserServiceClient(conn)}).(*userClient)
	if !ok {
		t.Fatal("NewUserClient did not build the grpc adapter")
	}
	return adapter
}

func TestUserClientListUsers_bearerInContext_requestBuiltPageMappedAndTokenForwarded(t *testing.T) {
	srv := &fakeUserServer{list: &userpb.ListUsersResponse{
		TotalSize: 2,
		Users: []*userpb.User{
			{Uid: "u-1", Email: "a@example.com", FirstName: "A", LastName: "B", Phone: "+1-555", Roles: []string{"USER", "STAFF"}},
			{Uid: "u-2", Email: "c@example.com"},
		},
	}}
	adapter := newUserAdapter(t, srv)
	ctx := httpx.WithBearerToken(requestContext(t), "Bearer abc.def.ghi")

	total, got, err := adapter.ListUsers(ctx, "a@example.com", "A", "+1-555", pageOf(50, 3, "email,asc"))
	if err != nil {
		t.Fatalf("list users: %v", err)
	}

	if total != 2 || len(got) != 2 {
		t.Fatalf("total = %d, users = %d, want 2 and 2", total, len(got))
	}
	if u := got[0]; u.Uid != "u-1" || u.Email != "a@example.com" || u.FirstName != "A" || u.LastName != "B" || u.Phone == nil || *u.Phone != "+1-555" ||
		!reflect.DeepEqual(u.Roles, []constant.AuthRole{constant.AUTHROLE_USER, constant.AUTHROLE_STAFF}) {
		t.Errorf("first user = %+v", u)
	}
	if got[1].Phone != nil || len(got[1].Roles) != 0 {
		t.Errorf("second user = %+v, want no phone and no roles", got[1])
	}
	req := srv.gotList
	if req.GetEmail() != "a@example.com" || req.GetName() != "A" || req.GetPhone() != "+1-555" ||
		req.GetPageSize() != 50 || req.GetPageToken() != "3" || req.GetOrderBy() != "email,asc" {
		t.Errorf("request = %v", req)
	}
	if auth := srv.gotMD.Get(userAuthorizationMetadataKey); !reflect.DeepEqual(auth, []string{"Bearer abc.def.ghi"}) {
		t.Errorf("authorization metadata = %v, want the caller's bearer token", auth)
	}
}

func TestUserClientFindUser_bearerInContext_userMappedAndTokenForwarded(t *testing.T) {
	srv := &fakeUserServer{find: &userpb.FindUserResponse{User: &userpb.User{Uid: targetUid, Email: "a@example.com", Roles: []string{"MANAGER"}}}}
	adapter := newUserAdapter(t, srv)
	ctx := httpx.WithBearerToken(requestContext(t), "Bearer abc.def.ghi")

	got, err := adapter.FindUser(ctx, targetUid)
	if err != nil {
		t.Fatalf("find user: %v", err)
	}

	if got == nil || got.Uid != targetUid || got.Email != "a@example.com" || len(got.Roles) != 1 {
		t.Errorf("got = %+v", got)
	}
	if srv.gotFind.GetUid() != targetUid {
		t.Errorf("requested uid = %q, want %q", srv.gotFind.GetUid(), targetUid)
	}
	if auth := srv.gotMD.Get(userAuthorizationMetadataKey); !reflect.DeepEqual(auth, []string{"Bearer abc.def.ghi"}) {
		t.Errorf("authorization metadata = %v, want the caller's bearer token", auth)
	}
}

func TestUserClientFindUser_noBearerAndResponseWithoutUser_nilAndNoAuthorizationSent(t *testing.T) {
	srv := &fakeUserServer{find: &userpb.FindUserResponse{}}
	adapter := newUserAdapter(t, srv)

	got, err := adapter.FindUser(requestContext(t), targetUid)
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil", got, err)
	}
	if auth := srv.gotMD.Get(userAuthorizationMetadataKey); len(auth) != 0 {
		t.Errorf("authorization metadata = %v, want none", auth)
	}
}

func TestUserClient_upstreamStatus_codeAndMessagePreserved(t *testing.T) {
	for _, code := range upstreamCodes {
		t.Run(code.String(), func(t *testing.T) {
			adapter := newUserAdapter(t, &fakeUserServer{err: status.Error(code, "identity says no")})

			found, err := adapter.FindUser(requestContext(t), targetUid)
			assertGrpcCode(t, err, code, "identity says no")
			if found != nil {
				t.Errorf("found = %+v, want nil", found)
			}

			total, listed, err := adapter.ListUsers(requestContext(t), "", "", "", pageOf(10, 0, ""))
			assertGrpcCode(t, err, code, "identity says no")
			if total != 0 || listed != nil {
				t.Errorf("list = %d/%v, want zero values", total, listed)
			}
		})
	}
}

func TestGrpcClients_configuredServers_eachClientTargetsItsOwnServerAndCloses(t *testing.T) {
	cfg := &config.GrpcConfig{
		OrderServer:       &sharedconfig.GrpcServerConfig{Host: "order.bookstore", Port: "9091"},
		PaymentServer:     &sharedconfig.GrpcServerConfig{Host: "payment.bookstore", Port: "9092"},
		UserServer:        &sharedconfig.GrpcServerConfig{Host: "identity.bookstore", Port: "9093"},
		FulfillmentServer: &sharedconfig.GrpcServerConfig{Host: "fulfillment.bookstore", Port: "9094"},
	}
	telemetry := &support.TelemetryManager{
		TraceProvider: sdktrace.NewTracerProvider(),
		MeterProvider: sdkmetric.NewMeterProvider(),
	}

	order, err := NewOrderGrpcClient(cfg, telemetry)
	if err != nil {
		t.Fatalf("order client: %v", err)
	}
	payment, err := NewPaymentGrpcClient(cfg, telemetry)
	if err != nil {
		t.Fatalf("payment client: %v", err)
	}
	user, err := NewUserGrpcClient(cfg, telemetry)
	if err != nil {
		t.Fatalf("user client: %v", err)
	}
	delivery, err := NewDeliveryGrpcClient(cfg, telemetry)
	if err != nil {
		t.Fatalf("delivery client: %v", err)
	}

	targets := map[string]string{
		"order":    order.conn.Target(),
		"payment":  payment.conn.Target(),
		"user":     user.conn.Target(),
		"delivery": delivery.conn.Target(),
	}
	want := map[string]string{
		"order":    "order.bookstore:9091",
		"payment":  "payment.bookstore:9092",
		"user":     "identity.bookstore:9093",
		"delivery": "fulfillment.bookstore:9094",
	}
	if !reflect.DeepEqual(targets, want) {
		t.Errorf("targets = %v, want %v", targets, want)
	}
	if order.Client == nil || payment.Client == nil || user.Client == nil || delivery.Client == nil {
		t.Error("a grpc client was built without its service stub")
	}
	for name, closeFn := range map[string]func() error{
		"order": order.Close, "payment": payment.Close, "user": user.Close, "delivery": delivery.Close,
	} {
		if err := closeFn(); err != nil {
			t.Errorf("close %s client: %v", name, err)
		}
	}
}
