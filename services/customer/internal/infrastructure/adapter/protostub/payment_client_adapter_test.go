package protostub

import (
	"testing"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/paymentpb"
)

func TestToPaymentAggregate(t *testing.T) {
	got := toPaymentAggregate(&paymentpb.Payment{
		Uid:           "p-1",
		PaymentUid:    "px-1",
		CustomerUid:   "u-1",
		Currency:      "USD",
		Amount:        4200,
		Status:        paymentpb.PaymentStatus_PAYMENT_STATUS_SUCCEEDED,
		PaymentMethod: "card",
		ErrorCode:     "",
		ErrorMessage:  "",
	})
	if got.Uid != "p-1" || got.CustomerUid != "u-1" || got.Currency != "USD" {
		t.Errorf("scalar fields = %+v", got)
	}
	if got.Amount != 4200 {
		t.Errorf("amount = %v, want 4200", got.Amount)
	}
	if got.Status != "PAYMENT_STATUS_SUCCEEDED" {
		t.Errorf("status = %q, want PAYMENT_STATUS_SUCCEEDED", got.Status)
	}
	if got.PaymentMethod == nil || *got.PaymentMethod != "card" {
		t.Errorf("paymentMethod = %v, want card", got.PaymentMethod)
	}
	if got.ErrorCode != nil || got.ErrorMessage != nil {
		t.Errorf("error fields should be nil for empty strings, got code=%v msg=%v", got.ErrorCode, got.ErrorMessage)
	}
}

func TestOptionalString(t *testing.T) {
	if got := optionalString(""); got != nil {
		t.Errorf("empty string -> %v, want nil", got)
	}
	if got := optionalString("x"); got == nil || *got != "x" {
		t.Errorf("\"x\" -> %v, want pointer to \"x\"", got)
	}
}
