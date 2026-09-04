package restcontroller

import (
	"reflect"
	"testing"
	"time"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
)

func TestToCustomerFindResponse_WhenAggregateFullyPopulated_MapsEveryField(t *testing.T) {
	phone := "+1-555"
	a := &aggregate.CustomerAggregate{
		Uid:       "u-1",
		Email:     "a@b.com",
		FirstName: "A",
		LastName:  "B",
		Phone:     &phone,
		Roles:     []constant.AuthRole{constant.AUTHROLE_ORDER, constant.AUTHROLE_VIEW},
	}
	got := toCustomerFindResponse(a)
	if got.Uid != "u-1" {
		t.Errorf("Uid = %q", got.Uid)
	}
	if got.Email == nil || *got.Email != "a@b.com" {
		t.Errorf("Email = %v", got.Email)
	}
	if got.Phone == nil || *got.Phone != phone {
		t.Errorf("Phone = %v", got.Phone)
	}
	if got.Roles == nil || len(*got.Roles) != 2 {
		t.Fatalf("Roles = %v", got.Roles)
	}
	if (*got.Roles)[0] != openapi.AUTHROLE_ORDER || (*got.Roles)[1] != openapi.AUTHROLE_VIEW {
		t.Errorf("Roles = %+v", *got.Roles)
	}
}

func TestToCustomerUpdateResponse_WhenAggregateGiven_MapsEveryField(t *testing.T) {
	a := &aggregate.CustomerAggregate{
		Uid:       "u-1",
		Email:     "a@b.com",
		FirstName: "A",
		LastName:  "B",
		Roles:     []constant.AuthRole{constant.AUTHROLE_PROFILE},
	}
	got := toCustomerUpdateResponse(a)
	if got.Uid != "u-1" || *got.Email != "a@b.com" {
		t.Errorf("got = %+v", got)
	}
	if got.Roles == nil || len(*got.Roles) != 1 || (*got.Roles)[0] != openapi.AUTHROLE_PROFILE {
		t.Errorf("Roles = %v", got.Roles)
	}
}

func TestAuthRoles_WhenConvertedBothWays_RoundTripUnchanged(t *testing.T) {
	in := []openapi.AuthRole{openapi.AUTHROLE_ORDER, openapi.AUTHROLE_PROFILE}
	domain := toAuthRoles(&in)
	if !reflect.DeepEqual(domain, []constant.AuthRole{constant.AUTHROLE_ORDER, constant.AUTHROLE_PROFILE}) {
		t.Errorf("toAuthRoles = %v", domain)
	}
	rest := authRolesToRest(domain)
	if !reflect.DeepEqual(rest, in) {
		t.Errorf("authRolesToRest = %v, want %v", rest, in)
	}
}

func TestToAuthRoles_WhenRolesNil_ReturnsNil(t *testing.T) {
	if got := toAuthRoles(nil); got != nil {
		t.Errorf("nil -> %v, want nil", got)
	}
}

func TestToOrderFindResponse_WhenAggregateFullyPopulated_MapsEveryField(t *testing.T) {
	a := &aggregate.OrderAggregate{
		Uid:         "o-1",
		UserUid:     "u-1",
		TotalAmount: 99.99,
		Status:      aggregate.ORDERSTATUS_SHIPPING,
		OrderItems: []aggregate.OrderItem{
			{BookUid: "b-1", Quantity: 2},
		},
	}
	got := toOrderFindResponse(a)
	if got.Uid != "o-1" || *got.UserUid != "u-1" {
		t.Errorf("got = %+v", got)
	}
	if got.OrderItems == nil || len(*got.OrderItems) != 1 || (*got.OrderItems)[0].BookUid != "b-1" {
		t.Errorf("OrderItems = %v", got.OrderItems)
	}
	if got.TotalPrice == nil || *got.TotalPrice != 99.99 {
		t.Errorf("TotalPrice = %v", got.TotalPrice)
	}
	if got.Status == nil || *got.Status != openapi.ORDERSTATUS_SHIPPING {
		t.Errorf("Status = %v", got.Status)
	}
}

func TestToOrderStatusRest_WhenDomainStatusGiven_ReturnsRestStatus(t *testing.T) {
	cases := []struct {
		in   aggregate.OrderStatus
		want openapi.OrderStatus
	}{
		{aggregate.ORDERSTATUS_PENDING, openapi.ORDERSTATUS_PENDING},
		{aggregate.ORDERSTATUS_PAID, openapi.ORDERSTATUS_PAID},
		{aggregate.ORDERSTATUS_SHIPPING, openapi.ORDERSTATUS_SHIPPING},
		{aggregate.ORDERSTATUS_COMPLETED, openapi.ORDERSTATUS_COMPLETED},
		{aggregate.ORDERSTATUS_CANCELLED, openapi.ORDERSTATUS_CANCELLED},
		{aggregate.ORDERSTATUS_UNKNOWN, openapi.ORDERSTATUS_UNKNOWN},
	}
	for _, c := range cases {
		if got := toOrderStatusRest(c.in); got != c.want {
			t.Errorf("toOrderStatusRest(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToPaymentFindResponse_WhenAggregateFullyPopulated_MapsEveryField(t *testing.T) {
	errMsg := "card declined"
	a := &aggregate.PaymentAggregate{
		PaymentUid:   "px-1",
		CustomerUid:  "u-1",
		Currency:     "USD",
		Amount:       50,
		Status:       "PAYMENT_STATUS_FAILED",
		ErrorMessage: &errMsg,
	}
	got := toPaymentFindResponse(a)
	if got.Uid != "px-1" || *got.UserUid != "u-1" || *got.Currency != "USD" || *got.Amount != 50 {
		t.Errorf("got = %+v", got)
	}
	if got.Status == nil || *got.Status != "PAYMENT_STATUS_FAILED" {
		t.Errorf("Status = %v", got.Status)
	}
	if got.Reason == nil || *got.Reason != errMsg {
		t.Errorf("Reason = %v", got.Reason)
	}
}

func TestToPointHistoryResponse_WhenEntryGiven_MapsEveryField(t *testing.T) {
	now := time.Now()
	reason := "signup bonus"
	got := toPointHistoryResponse(&aggregate.PointHistoryEntry{
		Uid:        "ph-1",
		ChangeType: aggregate.POINTCHANGE_EARN,
		Amount:     100,
		Reason:     &reason,
		CreatedAt:  now,
	})
	if got.Uid != "ph-1" || got.Amount != 100 {
		t.Errorf("got = %+v", got)
	}
	if got.ChangeType != openapi.PointChangeType("EARN") {
		t.Errorf("ChangeType = %v", got.ChangeType)
	}
	if got.Reason == nil || *got.Reason != reason {
		t.Errorf("Reason = %v", got.Reason)
	}
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v", got.CreatedAt)
	}
}

func TestToReviewFindResponse_WhenAggregateFullyPopulated_MapsEveryField(t *testing.T) {
	now := time.Now()
	updated := now.Add(time.Hour)
	title := "Loved it"
	content := "Great"
	got := toReviewFindResponse(&aggregate.ReviewAggregate{
		Uid:       "rv-1",
		UserUid:   "u-1",
		BookUid:   "b-1",
		Rating:    5,
		Title:     &title,
		Content:   &content,
		CreatedAt: now,
		UpdatedAt: &updated,
	})
	if got.Uid != "rv-1" || got.Rating != 5 {
		t.Errorf("got = %+v", got)
	}
	if got.Title == nil || *got.Title != "Loved it" {
		t.Errorf("Title = %v", got.Title)
	}
	if got.UpdatedAt == nil || !got.UpdatedAt.Equal(updated) {
		t.Errorf("UpdatedAt = %v", got.UpdatedAt)
	}
}

func TestToWishlistResponse_WhenAggregateGiven_MapsBookUids(t *testing.T) {
	got := toWishlistResponse(&aggregate.WishlistAggregate{
		UserUid:  "u-1",
		BookUids: []string{"b-1", "b-2"},
	})
	if !reflect.DeepEqual(got.BookUids, []string{"b-1", "b-2"}) {
		t.Errorf("BookUids = %v", got.BookUids)
	}
}
