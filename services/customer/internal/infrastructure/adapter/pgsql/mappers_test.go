package pgsql

import (
	"testing"
	"time"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
)

func TestToPointResult(t *testing.T) {
	got := toPointResult(&entity.CustomerPointEntity{Id: 1, Uid: "p-1", UserId: 7, Balance: 1500})
	if got.Id != 1 || got.Uid != "p-1" || got.UserId != 7 || got.Balance != 1500 {
		t.Errorf("got = %+v", got)
	}
}

func TestToPointHistoryResult(t *testing.T) {
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

func TestToReviewResult(t *testing.T) {
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

func TestToWishlistResult(t *testing.T) {
	now := time.Now()
	got := toWishlistResult(&entity.CustomerWishlistEntity{Id: 1, Uid: "w-1", UserId: 5, BookId: 9, CreatedAt: now})
	if got.Id != 1 || got.UserId != 5 || got.BookId != 9 {
		t.Errorf("got = %+v", got)
	}
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v", got.CreatedAt)
	}
}
