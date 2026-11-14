package pgsql

import (
	"errors"
	"math"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
)

func newPointRepo(t *testing.T) (*pointRepository, *gorm.DB) {
	t.Helper()
	db := newSchemaDB(t)
	return &pointRepository{q: queryOf(db), db: db}, db
}

func seedPoint(t *testing.T, db *gorm.DB, uid, userUid string, balance int32) *entity.CustomerPointEntity {
	t.Helper()
	row := &entity.CustomerPointEntity{Uid: uid, UserId: 77, UserUid: userUid, Balance: balance, CreatedAt: baseTime, Version: 1}
	mustCreate(t, db, row)
	return row
}

func loadHistory(t *testing.T, db *gorm.DB, userUid string) []entity.CustomerPointHistoryEntity {
	t.Helper()
	var rows []entity.CustomerPointHistoryEntity
	if err := db.Where("user_uid = ?", userUid).Order("id").Find(&rows).Error; err != nil {
		t.Fatalf("load history: %v", err)
	}
	return rows
}

func storedBalance(t *testing.T, db *gorm.DB, userUid string) int32 {
	t.Helper()
	var row entity.CustomerPointEntity
	if err := db.Where("user_uid = ?", userUid).Take(&row).Error; err != nil {
		t.Fatalf("load balance: %v", err)
	}
	return row.Balance
}

func TestPointRepositoryFindByUserUid_existingBalance_everyFieldMapped(t *testing.T) {
	r, db := newPointRepo(t)
	seeded := seedPoint(t, db, "p-1", userA, 120)
	seedPoint(t, db, "p-2", userB, 5)

	got, err := r.FindByUserUid(bg(), userA)
	if err != nil {
		t.Fatalf("find by user: %v", err)
	}
	want := &repo.PointResult{Id: seeded.Id, Uid: "p-1", UserId: 77, UserUid: userA, Balance: 120}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got = %+v, want %+v", got, want)
	}
	byId, err := r.FindOne(bg(), seeded.Id)
	if err != nil {
		t.Fatalf("find one: %v", err)
	}
	if !reflect.DeepEqual(byId, want) {
		t.Fatalf("by id = %+v, want %+v", byId, want)
	}
}

func TestPointRepositoryFindByUserUid_userWithoutBalance_nilWithoutError(t *testing.T) {
	r, _ := newPointRepo(t)

	got, err := r.FindByUserUid(bg(), userA)
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil", got, err)
	}
}

func TestPointRepositoryFindOne_unknownId_recordNotFound(t *testing.T) {
	r, _ := newPointRepo(t)

	if _, err := r.FindOne(bg(), 404); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("err = %v, want record not found", err)
	}
}

func TestPointRepositoryEnsureBalance_userWithoutBalance_zeroBalanceCreated(t *testing.T) {
	r, db := newPointRepo(t)

	if err := r.EnsureBalance(bg(), userA); err != nil {
		t.Fatalf("ensure balance: %v", err)
	}

	got, err := r.FindByUserUid(bg(), userA)
	if err != nil {
		t.Fatalf("find by user: %v", err)
	}
	if got == nil || got.Balance != 0 || got.Uid == "" || got.UserUid != userA {
		t.Fatalf("got = %+v, want a zero balance for the user", got)
	}
	if len(loadHistory(t, db, userA)) != 0 {
		t.Error("opening a balance must not write a history entry")
	}
}

func TestPointRepositoryEnsureBalance_userWithBalance_balanceAndIdentityKept(t *testing.T) {
	r, db := newPointRepo(t)
	seedPoint(t, db, "p-1", userA, 300)

	for i := 0; i < 2; i++ {
		if err := r.EnsureBalance(bg(), userA); err != nil {
			t.Fatalf("ensure balance: %v", err)
		}
	}

	got, err := r.FindByUserUid(bg(), userA)
	if err != nil {
		t.Fatalf("find by user: %v", err)
	}
	if got.Balance != 300 || got.Uid != "p-1" {
		t.Fatalf("got = %+v, want the existing balance untouched", got)
	}
	var count int64
	if err := db.Model(&entity.CustomerPointEntity{}).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("rows = %d, want 1", count)
	}
}

func TestPointRepositoryApplyPointChange_earnThenSpend_balanceAndHistoryRecorded(t *testing.T) {
	r, db := newPointRepo(t)
	seedPoint(t, db, "p-1", userA, 100)
	seedPoint(t, db, "p-2", userB, 100)
	reason := "order refund"

	earned, err := r.ApplyPointChange(bg(), userA, 50, "EARN", &reason)
	if err != nil {
		t.Fatalf("earn: %v", err)
	}
	spent, err := r.ApplyPointChange(bg(), userA, -120, "SPEND", nil)
	if err != nil {
		t.Fatalf("spend: %v", err)
	}

	if earned.Balance != 150 || spent.Balance != 30 || spent.Uid != "p-1" || spent.UserUid != userA {
		t.Fatalf("earned = %+v, spent = %+v, want balances 150 then 30", earned, spent)
	}
	if got := storedBalance(t, db, userA); got != 30 {
		t.Errorf("stored balance = %d, want 30", got)
	}
	if got := storedBalance(t, db, userB); got != 100 {
		t.Errorf("other user balance = %d, want 100", got)
	}
	history := loadHistory(t, db, userA)
	if len(history) != 2 {
		t.Fatalf("history entries = %d, want 2", len(history))
	}
	if history[0].ChangeType != "EARN" || history[0].Amount != 50 || history[0].Reason == nil || *history[0].Reason != reason || history[0].UserId != 77 {
		t.Errorf("earn entry = %+v", history[0])
	}
	if history[1].ChangeType != "SPEND" || history[1].Amount != -120 || history[1].Reason != nil {
		t.Errorf("spend entry = %+v", history[1])
	}
	if history[0].Uid == "" || history[0].Uid == history[1].Uid {
		t.Errorf("history uids = %q/%q, want distinct", history[0].Uid, history[1].Uid)
	}
}

func TestPointRepositoryApplyPointChange_spendWholeBalance_zeroBalance(t *testing.T) {
	r, db := newPointRepo(t)
	seedPoint(t, db, "p-1", userA, 100)

	got, err := r.ApplyPointChange(bg(), userA, -100, "SPEND", nil)
	if err != nil {
		t.Fatalf("spend: %v", err)
	}
	if got.Balance != 0 {
		t.Fatalf("balance = %d, want 0", got.Balance)
	}
}

func TestPointRepositoryApplyPointChange_spendBeyondBalance_insufficientPointsAndNothingWritten(t *testing.T) {
	r, db := newPointRepo(t)
	seedPoint(t, db, "p-1", userA, 100)

	got, err := r.ApplyPointChange(bg(), userA, -101, "SPEND", nil)
	if !errors.Is(err, repo.ErrInsufficientPoints) {
		t.Fatalf("err = %v, want ErrInsufficientPoints", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
	if balance := storedBalance(t, db, userA); balance != 100 {
		t.Errorf("stored balance = %d, want 100", balance)
	}
	if len(loadHistory(t, db, userA)) != 0 {
		t.Error("a rejected spend left a history entry")
	}
}

func TestPointRepositoryApplyPointChange_earnBeyondLimit_balanceOverflowAndNothingWritten(t *testing.T) {
	r, db := newPointRepo(t)
	seedPoint(t, db, "p-1", userA, math.MaxInt32-5)

	_, err := r.ApplyPointChange(bg(), userA, 6, "EARN", nil)
	if !errors.Is(err, repo.ErrPointBalanceOverflow) {
		t.Fatalf("err = %v, want ErrPointBalanceOverflow", err)
	}
	if balance := storedBalance(t, db, userA); balance != math.MaxInt32-5 {
		t.Errorf("stored balance = %d, want it unchanged", balance)
	}
	if len(loadHistory(t, db, userA)) != 0 {
		t.Error("a rejected earn left a history entry")
	}
}

func TestPointRepositoryApplyPointChange_userWithoutBalance_nilWithoutErrorOrHistory(t *testing.T) {
	r, db := newPointRepo(t)

	got, err := r.ApplyPointChange(bg(), userA, 10, "EARN", nil)
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil", got, err)
	}
	if len(loadHistory(t, db, userA)) != 0 {
		t.Error("history written for a user without a balance")
	}
}

func TestPointRepositoryApplyPointChange_historyWriteFails_balanceRolledBack(t *testing.T) {
	r, db := newPointRepo(t)
	seedPoint(t, db, "p-1", userA, 100)
	failCreatesOn(t, db, entity.TableNameCustomerPointHistoryEntity, errInjected)

	got, err := r.ApplyPointChange(bg(), userA, 50, "EARN", nil)
	assertInjected(t, err)
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
	if balance := storedBalance(t, db, userA); balance != 100 {
		t.Errorf("stored balance = %d, want 100 after rollback", balance)
	}
}

func TestPointRepositoryApplyPointChange_balanceWriteFails_errorSurfacedWithoutHistory(t *testing.T) {
	r, db := newPointRepo(t)
	seedPoint(t, db, "p-1", userA, 100)
	failCreatesOn(t, db, entity.TableNameCustomerPointEntity, errInjected)

	_, err := r.ApplyPointChange(bg(), userA, 50, "EARN", nil)
	assertInjected(t, err)
	if len(loadHistory(t, db, userA)) != 0 {
		t.Error("history written although the balance update failed")
	}
}

func TestPointRepository_databaseUnavailable_everyOperationFails(t *testing.T) {
	r, db := newPointRepo(t)
	closeDatabase(t, db)

	_, err := r.FindByUserUid(bg(), userA)
	assertFails(t, "FindByUserUid", err)
	assertFails(t, "EnsureBalance", r.EnsureBalance(bg(), userA))
	_, err = r.ApplyPointChange(bg(), userA, 1, "EARN", nil)
	assertFails(t, "ApplyPointChange", err)
}

func newPointHistoryRepo(t *testing.T) (*pointHistoryRepository, *gorm.DB) {
	t.Helper()
	db := newSchemaDB(t)
	return &pointHistoryRepository{q: queryOf(db)}, db
}

func seedHistory(t *testing.T, db *gorm.DB, uid, userUid string, amount int32, createdAt time.Time) {
	t.Helper()
	mustCreate(t, db, &entity.CustomerPointHistoryEntity{
		Uid: uid, UserUid: userUid, ChangeType: "EARN", Amount: amount, CreatedAt: createdAt,
	})
}

func historyAmounts(rows []*repo.PointHistoryResult) []int32 {
	out := make([]int32, len(rows))
	for i, r := range rows {
		out[i] = r.Amount
	}
	return out
}

func TestPointHistoryRepositoryFindAll_twoUsers_onlyRequestedUserNewestFirst(t *testing.T) {
	r, db := newPointHistoryRepo(t)
	seedHistory(t, db, "h-1", userA, 10, hoursAfterBase(1))
	seedHistory(t, db, "h-2", userA, 30, hoursAfterBase(3))
	seedHistory(t, db, "h-3", userA, 20, hoursAfterBase(2))
	seedHistory(t, db, "h-4", userB, 99, hoursAfterBase(4))

	total, got, err := r.FindAll(bg(), userA, pageOf(10, 0, ""))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 3 || !reflect.DeepEqual(historyAmounts(got), []int32{30, 20, 10}) {
		t.Fatalf("total = %d, amounts = %v, want [30 20 10]", total, historyAmounts(got))
	}
	if got[0].Uid != "h-2" || got[0].ChangeType != "EARN" || !got[0].CreatedAt.Equal(hoursAfterBase(3)) {
		t.Errorf("first entry = %+v", got[0])
	}
}

func TestPointHistoryRepositoryFindAll_secondPage_skipsOneFullPage(t *testing.T) {
	r, db := newPointHistoryRepo(t)
	for i := 1; i <= 5; i++ {
		seedHistory(t, db, "h-"+string(rune('0'+i)), userA, int32(i), hoursAfterBase(i))
	}

	total, got, err := r.FindAll(bg(), userA, pageOf(2, 1, ""))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 5 || !reflect.DeepEqual(historyAmounts(got), []int32{3, 2}) {
		t.Fatalf("total = %d, amounts = %v, want total 5 and [3 2]", total, historyAmounts(got))
	}
}

func TestPointHistoryRepositoryFindAll_userWithoutHistory_emptyPage(t *testing.T) {
	r, _ := newPointHistoryRepo(t)

	total, got, err := r.FindAll(bg(), userA, pageOf(10, 0, ""))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 0 || got == nil || len(got) != 0 {
		t.Fatalf("total = %d, items = %v, want an empty non-nil page", total, got)
	}
}

func TestPointHistoryRepositoryFindAll_databaseUnavailable_errorSurfaced(t *testing.T) {
	r, db := newPointHistoryRepo(t)
	closeDatabase(t, db)

	_, _, err := r.FindAll(bg(), userA, pageOf(10, 0, ""))
	assertFails(t, "FindAll", err)
}
