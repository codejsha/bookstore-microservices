package protostub

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/orderpb"
	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/paymentpb"
	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/userpb"
)

type captureOrderServiceClient struct {
	orderpb.OrderServiceClient
	got *orderpb.ListOrdersRequest
}

func (c *captureOrderServiceClient) ListOrders(_ context.Context, in *orderpb.ListOrdersRequest, _ ...grpc.CallOption) (*orderpb.ListOrdersResponse, error) {
	c.got = in
	return &orderpb.ListOrdersResponse{}, nil
}

type capturePaymentServiceClient struct {
	paymentpb.PaymentServiceClient
	got *paymentpb.ListPaymentsRequest
}

func (c *capturePaymentServiceClient) ListPayments(_ context.Context, in *paymentpb.ListPaymentsRequest, _ ...grpc.CallOption) (*paymentpb.ListPaymentsResponse, error) {
	c.got = in
	return &paymentpb.ListPaymentsResponse{}, nil
}

type captureUserServiceClient struct {
	userpb.UserServiceClient
	got *userpb.ListUsersRequest
}

func (c *captureUserServiceClient) ListUsers(_ context.Context, in *userpb.ListUsersRequest, _ ...grpc.CallOption) (*userpb.ListUsersResponse, error) {
	c.got = in
	return &userpb.ListUsersResponse{}, nil
}

func pageOf(size, page int32, sort string) pagination.PageOption {
	return pagination.NewPageOption(&size, &page, &sort)
}

func TestToPageRequest_PageOptionVariants_NormalizedRequest(t *testing.T) {
	cases := []struct {
		name string
		in   pagination.PageOption
		want pageRequest
	}{
		{"unset_defaultSizeNoTokenNoOrder", pagination.PageOption{}, pageRequest{pageSize: 20}},
		{"firstPage_noToken", pageOf(10, 0, "created_at,desc"), pageRequest{pageSize: 10, orderBy: "created_at,desc"}},
		{"thirdPage_tokenIsIndex", pageOf(10, 2, ""), pageRequest{pageSize: 10, pageToken: "2"}},
		{"negativePage_noToken", pageOf(10, -3, ""), pageRequest{pageSize: 10}},
		{"oversized_cappedAtMax", pageOf(500, 1, ""), pageRequest{pageSize: 100, pageToken: "1"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := toPageRequest(tc.in); got != tc.want {
				t.Errorf("toPageRequest = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestOrderClientListOrders_PageAndSortGiven_ForwardsPageTokenAndOrderBy(t *testing.T) {
	stub := &captureOrderServiceClient{}
	c := &orderClient{client: stub}
	if _, _, err := c.ListOrders(context.Background(), "u-1", nil, pageOf(5, 3, "created_at,asc")); err != nil {
		t.Fatalf("ListOrders err: %v", err)
	}
	if stub.got.GetUserUid() != "u-1" || stub.got.GetPageSize() != 5 || stub.got.GetPageToken() != "3" || stub.got.GetOrderBy() != "created_at,asc" {
		t.Errorf("request = %+v, want user u-1 size 5 token 3 order_by created_at,asc", stub.got)
	}
}

func TestPaymentClientListPayments_PageAndSortGiven_ForwardsPageTokenAndOrderBy(t *testing.T) {
	stub := &capturePaymentServiceClient{}
	c := &paymentClient{client: stub}
	if _, _, err := c.ListPayments(context.Background(), "u-1", pageOf(7, 1, "amount,desc")); err != nil {
		t.Fatalf("ListPayments err: %v", err)
	}
	if stub.got.GetCustomerUid() != "u-1" || stub.got.GetPageSize() != 7 || stub.got.GetPageToken() != "1" || stub.got.GetOrderBy() != "amount,desc" {
		t.Errorf("request = %+v, want customer u-1 size 7 token 1 order_by amount,desc", stub.got)
	}
}

func TestUserClientListUsers_PageAndSortGiven_ForwardsPageTokenAndOrderBy(t *testing.T) {
	stub := &captureUserServiceClient{}
	c := &userClient{client: stub}
	if _, _, err := c.ListUsers(context.Background(), "a@b.com", "", "", pageOf(15, 4, "email,asc")); err != nil {
		t.Fatalf("ListUsers err: %v", err)
	}
	if stub.got.GetEmail() != "a@b.com" || stub.got.GetPageSize() != 15 || stub.got.GetPageToken() != "4" || stub.got.GetOrderBy() != "email,asc" {
		t.Errorf("request = %+v, want email a@b.com size 15 token 4 order_by email,asc", stub.got)
	}
}

func TestUserClientListUsers_FirstPageNoSort_OmitsPageTokenAndOrderBy(t *testing.T) {
	stub := &captureUserServiceClient{}
	c := &userClient{client: stub}
	if _, _, err := c.ListUsers(context.Background(), "", "", "", pagination.PageOption{}); err != nil {
		t.Fatalf("ListUsers err: %v", err)
	}
	if stub.got.GetPageSize() != 20 || stub.got.GetPageToken() != "" || stub.got.GetOrderBy() != "" {
		t.Errorf("request = %+v, want size 20 and no token/order_by", stub.got)
	}
}
