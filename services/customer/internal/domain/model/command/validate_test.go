package command

import (
	"errors"
	"testing"

	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
)

func ptr[T any](v T) *T { return &v }

const (
	testUserUid  = "018f0000-0000-7000-8000-0000000000a1"
	testBookUid1 = "018f0000-0000-7000-8000-0000000000b1"
	testBookUid2 = "018f0000-0000-7000-8000-0000000000b2"
)

func longText(n int) string {
	out := make([]rune, n)
	for i := range out {
		out[i] = 'x'
	}
	return string(out)
}

func TestCustomerUpdateCommand_Validate(t *testing.T) {
	valid := CustomerUpdateCommand{Uid: "u-1"}
	tests := []struct {
		name    string
		mutate  func(c *CustomerUpdateCommand)
		wantErr bool
	}{
		{"onlyUidSet_noError", func(c *CustomerUpdateCommand) {}, false},
		{"everyFieldSet_noError", func(c *CustomerUpdateCommand) {
			c.Email = ptr("a@b.co")
			c.FirstName = ptr("Jin")
			c.Roles = []constant.AuthRole{constant.AUTHROLE_STAFF, constant.AUTHROLE_MANAGE}
		}, false},
		{"uidBlank_errInvalidCommand", func(c *CustomerUpdateCommand) { c.Uid = " " }, true},
		{"emailMalformed_errInvalidCommand", func(c *CustomerUpdateCommand) { c.Email = ptr("not-an-email") }, true},
		{"passwordBlank_errInvalidCommand", func(c *CustomerUpdateCommand) { c.Password = ptr("  ") }, true},
		{"phoneBlank_errInvalidCommand", func(c *CustomerUpdateCommand) { c.Phone = ptr("") }, true},
		{"roleUnknown_errInvalidCommand", func(c *CustomerUpdateCommand) { c.Roles = []constant.AuthRole{constant.AUTHROLE_UNKNOWN} }, true},
		{"roleNotAKnownMember_errInvalidCommand", func(c *CustomerUpdateCommand) { c.Roles = []constant.AuthRole{99} }, true},
		{"rolesDuplicated_errInvalidCommand", func(c *CustomerUpdateCommand) {
			c.Roles = []constant.AuthRole{constant.AUTHROLE_USER, constant.AUTHROLE_USER}
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := valid
			tt.mutate(&cmd)
			err := cmd.Validate()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidCommand) {
					t.Fatalf("want ErrInvalidCommand, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("want nil, got %v", err)
			}
		})
	}
}

// Earn and spend share one rule set, so both commands run here: a non-positive or
// over-cap amount, a blank or non-uuid user uid, and a blank or over-cap reason are
// rejected, while the values sitting exactly on the caps pass.
func TestPointCommands_WhenAmountOrReasonOutOfBounds_ReturnErrInvalidCommand(t *testing.T) {
	if err := (PointEarnCommand{UserUid: testUserUid, Amount: 10}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if err := (PointEarnCommand{UserUid: testUserUid, Amount: 0}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (PointSpendCommand{UserUid: " ", Amount: 10}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (PointSpendCommand{UserUid: "not-a-uuid", Amount: 10}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (PointSpendCommand{UserUid: testUserUid, Amount: -5}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (PointEarnCommand{UserUid: testUserUid, Amount: 10, Reason: ptr(" ")}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (PointEarnCommand{UserUid: testUserUid, Amount: pointAmountMaxPerOp}).Validate(); err != nil {
		t.Fatalf("want nil at the per-operation cap, got %v", err)
	}
	if err := (PointEarnCommand{UserUid: testUserUid, Amount: pointAmountMaxPerOp + 1}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (PointSpendCommand{UserUid: testUserUid, Amount: pointAmountMaxPerOp + 1}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (PointEarnCommand{UserUid: testUserUid, Amount: 10, Reason: ptr(longText(pointReasonMaxLen))}).Validate(); err != nil {
		t.Fatalf("want nil at the reason cap, got %v", err)
	}
	if err := (PointEarnCommand{UserUid: testUserUid, Amount: 10, Reason: ptr(longText(pointReasonMaxLen + 1))}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
}

func TestWishlistAddCommand_Validate(t *testing.T) {
	tests := []struct {
		name     string
		userUid  string
		bookUids []string
		wantErr  bool
	}{
		{"whenUserAndBooksValid_returnsNil", testUserUid, []string{testBookUid1, testBookUid2}, false},
		{"whenUserUidBlank_returnsErrInvalidCommand", " ", []string{testBookUid1}, true},
		{"whenUserUidNotUuid_returnsErrInvalidCommand", "u-1", []string{testBookUid1}, true},
		{"whenBookUidsNil_returnsErrInvalidCommand", testUserUid, nil, true},
		{"whenBookUidsEmpty_returnsErrInvalidCommand", testUserUid, []string{}, true},
		{"whenBookUidBlank_returnsErrInvalidCommand", testUserUid, []string{testBookUid1, "  "}, true},
		{"whenBookUidNotUuid_returnsErrInvalidCommand", testUserUid, []string{"b-1"}, true},
		{"whenBookUidsDuplicated_returnsErrInvalidCommand", testUserUid, []string{testBookUid1, testBookUid1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			add := WishlistAddCommand{UserUid: tt.userUid, BookUids: tt.bookUids}
			remove := WishlistRemoveCommand{UserUid: tt.userUid, BookUids: tt.bookUids}
			for _, err := range []error{add.Validate(), remove.Validate()} {
				if tt.wantErr {
					if !errors.Is(err, ErrInvalidCommand) {
						t.Fatalf("want ErrInvalidCommand, got %v", err)
					}
				} else if err != nil {
					t.Fatalf("want nil, got %v", err)
				}
			}
		})
	}
}

func TestReviewWriteCommand_Validate(t *testing.T) {
	valid := ReviewWriteCommand{UserUid: testUserUid, BookUid: testBookUid1, Rating: 3}
	tests := []struct {
		name    string
		mutate  func(c *ReviewWriteCommand)
		wantErr bool
	}{
		{"whenEveryFieldValid_returnsNil", func(c *ReviewWriteCommand) {}, false},
		{"whenRatingZero_returnsErrInvalidCommand", func(c *ReviewWriteCommand) { c.Rating = 0 }, true},
		{"whenRatingAboveMax_returnsErrInvalidCommand", func(c *ReviewWriteCommand) { c.Rating = 6 }, true},
		{"whenBookUidBlank_returnsErrInvalidCommand", func(c *ReviewWriteCommand) { c.BookUid = " " }, true},
		{"whenBookUidNotUuid_returnsErrInvalidCommand", func(c *ReviewWriteCommand) { c.BookUid = "b-1" }, true},
		{"whenUserUidNotUuid_returnsErrInvalidCommand", func(c *ReviewWriteCommand) { c.UserUid = "u-1" }, true},
		{"whenTitleBlank_returnsErrInvalidCommand", func(c *ReviewWriteCommand) { c.Title = ptr("  ") }, true},
		{"whenTitleAtCap_returnsNil", func(c *ReviewWriteCommand) { c.Title = ptr(longText(reviewTitleMaxLen)) }, false},
		{"whenTitleOverCap_returnsErrInvalidCommand", func(c *ReviewWriteCommand) { c.Title = ptr(longText(reviewTitleMaxLen + 1)) }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := valid
			tt.mutate(&cmd)
			err := cmd.Validate()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidCommand) {
					t.Fatalf("want ErrInvalidCommand, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("want nil, got %v", err)
			}
		})
	}

	if err := (ReviewEditCommand{}).Validate(); err != nil {
		t.Fatalf("empty edit must be valid, got %v", err)
	}
	if err := (ReviewEditCommand{Rating: ptr(int32(0))}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (ReviewEditCommand{Content: ptr("")}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (ReviewEditCommand{Title: ptr(longText(reviewTitleMaxLen))}).Validate(); err != nil {
		t.Fatalf("want nil at the title cap, got %v", err)
	}
	if err := (ReviewEditCommand{Title: ptr(longText(reviewTitleMaxLen + 1))}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
}
