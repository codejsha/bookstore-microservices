package pgsql

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
)

func TestToPointResult_WhenRowGiven_ReturnsResult(t *testing.T) {
	got := toPointResult(&entity.CustomerPointEntity{Id: 1, Uid: "p-1", UserId: 7, Balance: 1500})
	if got.Id != 1 || got.Uid != "p-1" || got.UserId != 7 || got.Balance != 1500 {
		t.Errorf("got = %+v", got)
	}
}

func TestToPointHistoryResult_WhenRowGiven_ReturnsResult(t *testing.T) {
	now := time.Now()
	reason := "earned via review"
	got := toPointHistoryResult(&entity.CustomerPointHistoryEntity{
		Id: 9, Uid: "ph-1", UserId: 1, ChangeType: "EARN", Amount: 50, Reason: &reason, CreatedAt: now,
	})
	if got.Id != 9 || got.ChangeType != "EARN" || got.Amount != 50 {
		t.Errorf("got = %+v", got)
	}
	if got.Reason == nil || *got.Reason != reason {
		t.Errorf("Reason = %v", got.Reason)
	}
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v", got.CreatedAt)
	}
}

func TestToReviewResult_WhenRowGiven_ReturnsResult(t *testing.T) {
	now := time.Now()
	updated := now.Add(time.Hour)
	title := "Good"
	content := "Liked it"
	got := toReviewResult(&entity.CustomerReviewEntity{
		Id: 3, Uid: "rv-1", UserId: 1, BookId: 2,
		Rating: 4, Title: &title, Content: &content, CreatedAt: now, UpdatedAt: &updated,
	})
	if got.Id != 3 || got.Rating != 4 {
		t.Errorf("got = %+v", got)
	}
	if got.Title == nil || *got.Title != title {
		t.Errorf("Title = %v", got.Title)
	}
	if got.UpdatedAt == nil || !got.UpdatedAt.Equal(updated) {
		t.Errorf("UpdatedAt = %v", got.UpdatedAt)
	}
}

func TestToWishlistResult_WhenRowGiven_ReturnsResult(t *testing.T) {
	now := time.Now()
	got := toWishlistResult(&entity.CustomerWishlistEntity{Id: 1, Uid: "w-1", UserId: 5, BookId: 9, CreatedAt: now})
	if got.Id != 1 || got.UserId != 5 || got.BookId != 9 {
		t.Errorf("got = %+v", got)
	}
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v", got.CreatedAt)
	}
}

func TestNextBalance(t *testing.T) {
	tests := []struct {
		name    string
		current int32
		delta   int32
		want    int32
		wantErr error
	}{
		{"whenDeltaPositive_returnsSum", 100, 50, 150, nil},
		{"whenDeltaNegative_returnsRemainder", 100, -40, 60, nil},
		{"whenDeltaDrainsBalance_returnsZero", 100, -100, 0, nil},
		{"whenDeltaExceedsBalance_returnsErrInsufficientPoints", 10, -50, 0, repo.ErrInsufficientPoints},
		{"whenSumExceedsInt32_returnsErrPointBalanceOverflow", math.MaxInt32 - 5, 10, 0, repo.ErrPointBalanceOverflow},
		{"whenSumHitsInt32Ceiling_returnsMaxInt32", math.MaxInt32 - 5, 5, math.MaxInt32, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nextBalance(tt.current, tt.delta)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("err = %v, want nil", err)
			}
			if got != tt.want {
				t.Errorf("got = %d, want %d", got, tt.want)
			}
		})
	}
}
