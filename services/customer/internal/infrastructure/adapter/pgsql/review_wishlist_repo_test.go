package pgsql

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

func newReviewRepo(t *testing.T) (*reviewRepository, *gorm.DB) {
	t.Helper()
	db := newSchemaDB(t)
	return &reviewRepository{q: queryOf(db)}, db
}

func seedReview(t *testing.T, db *gorm.DB, uid, userUid, bookUid string, rating int16, createdAt time.Time) *entity.CustomerReviewEntity {
	t.Helper()
	row := &entity.CustomerReviewEntity{Uid: uid, UserUid: userUid, BookUid: bookUid, Rating: rating, CreatedAt: createdAt, Version: 1}
	mustCreate(t, db, row)
	return row
}

func reviewUids(rows []*repo.ReviewResult) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.Uid
	}
	return out
}

func seedReviewMatrix(t *testing.T, db *gorm.DB) {
	t.Helper()
	seedReview(t, db, "r-1", userA, bookA, 5, hoursAfterBase(1))
	seedReview(t, db, "r-2", userA, bookB, 3, hoursAfterBase(2))
	seedReview(t, db, "r-3", userB, bookA, 5, hoursAfterBase(3))
	seedReview(t, db, "r-4", userB, bookB, 1, hoursAfterBase(4))
}

func TestReviewRepositoryFindAll_singleFilter_onlyMatchingReviewsNewestFirst(t *testing.T) {
	cases := []struct {
		name string
		opts []option.ReviewQueryOptionFunc
		want []string
	}{
		{"none", nil, []string{"r-4", "r-3", "r-2", "r-1"}},
		{"user", []option.ReviewQueryOptionFunc{option.WithReviewUserUid(userA)}, []string{"r-2", "r-1"}},
		{"book", []option.ReviewQueryOptionFunc{option.WithReviewBookUid(bookA)}, []string{"r-3", "r-1"}},
		{"rating", []option.ReviewQueryOptionFunc{option.WithReviewRating(5)}, []string{"r-3", "r-1"}},
		{"bookAndRating", []option.ReviewQueryOptionFunc{option.WithReviewBookUid(bookB), option.WithReviewRating(1)}, []string{"r-4"}},
		{"blankUserIgnored", []option.ReviewQueryOptionFunc{option.WithReviewUserUid("")}, []string{"r-4", "r-3", "r-2", "r-1"}},
		{"blankBookIgnored", []option.ReviewQueryOptionFunc{option.WithReviewBookUid("")}, []string{"r-4", "r-3", "r-2", "r-1"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, db := newReviewRepo(t)
			seedReviewMatrix(t, db)

			total, got, err := r.FindAll(bg(), option.NewReviewQueryOption(c.opts...))
			if err != nil {
				t.Fatalf("find all: %v", err)
			}
			if total != int64(len(c.want)) || !reflect.DeepEqual(reviewUids(got), c.want) {
				t.Fatalf("total = %d, uids = %v, want %v", total, reviewUids(got), c.want)
			}
		})
	}
}

func TestReviewRepositoryFindAll_secondPage_skipsOneFullPage(t *testing.T) {
	r, db := newReviewRepo(t)
	seedReviewMatrix(t, db)

	total, got, err := r.FindAll(bg(), option.NewReviewQueryOption(option.WithReviewPage(pageOf(3, 1, ""))))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 4 || !reflect.DeepEqual(reviewUids(got), []string{"r-1"}) {
		t.Fatalf("total = %d, uids = %v, want total 4 and [r-1]", total, reviewUids(got))
	}
}

func TestReviewRepositoryFindAll_softDeletedReview_excluded(t *testing.T) {
	r, db := newReviewRepo(t)
	seedReview(t, db, "r-1", userA, bookA, 5, hoursAfterBase(1))
	gone := seedReview(t, db, "r-2", userA, bookB, 4, hoursAfterBase(2))
	softDelete(t, db, gone)

	total, got, err := r.FindAll(bg(), option.NewReviewQueryOption(option.WithReviewUserUid(userA)))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 1 || !reflect.DeepEqual(reviewUids(got), []string{"r-1"}) {
		t.Fatalf("total = %d, uids = %v, want [r-1]", total, reviewUids(got))
	}
}

func TestReviewRepositoryInsert_firstReviewOfBook_persistedAndReadableByUidAndId(t *testing.T) {
	r, _ := newReviewRepo(t)
	title, content := "Great", "Loved it"

	got, err := r.Insert(bg(), repo.ReviewCreate{UserUid: userA, BookUid: bookA, Rating: 4, Title: &title, Content: &content})
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if got.Uid == "" || got.Id == 0 || got.UserUid != userA || got.BookUid != bookA || got.Rating != 4 {
		t.Fatalf("got = %+v", got)
	}
	byUid, err := r.FindByUid(bg(), got.Uid)
	if err != nil {
		t.Fatalf("find by uid: %v", err)
	}
	byId, err := r.FindOne(bg(), got.Id)
	if err != nil {
		t.Fatalf("find one: %v", err)
	}
	for _, found := range []*repo.ReviewResult{byUid, byId} {
		if found.Uid != got.Uid || found.Rating != 4 || found.Title == nil || *found.Title != title || found.Content == nil || *found.Content != content {
			t.Errorf("found = %+v, want the inserted review", found)
		}
	}
}

func TestReviewRepositoryInsert_sameUserAndBookAgain_reviewExistsAndNoSecondRow(t *testing.T) {
	r, db := newReviewRepo(t)
	seedReview(t, db, "r-1", userA, bookA, 5, baseTime)

	got, err := r.Insert(bg(), repo.ReviewCreate{UserUid: userA, BookUid: bookA, Rating: 1})
	if !errors.Is(err, repo.ErrReviewExists) {
		t.Fatalf("err = %v, want ErrReviewExists", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
	var count int64
	if err := db.Model(&entity.CustomerReviewEntity{}).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("rows = %d, want 1", count)
	}
}

func TestReviewRepositoryInsert_sameBookOtherUserOrOtherBook_accepted(t *testing.T) {
	r, db := newReviewRepo(t)
	seedReview(t, db, "r-1", userA, bookA, 5, baseTime)

	if _, err := r.Insert(bg(), repo.ReviewCreate{UserUid: userB, BookUid: bookA, Rating: 2}); err != nil {
		t.Errorf("other user, same book: %v", err)
	}
	if _, err := r.Insert(bg(), repo.ReviewCreate{UserUid: userA, BookUid: bookB, Rating: 2}); err != nil {
		t.Errorf("same user, other book: %v", err)
	}
}

func TestReviewRepositoryInsert_previousReviewSoftDeleted_accepted(t *testing.T) {
	r, db := newReviewRepo(t)
	gone := seedReview(t, db, "r-1", userA, bookA, 5, baseTime)
	softDelete(t, db, gone)

	got, err := r.Insert(bg(), repo.ReviewCreate{UserUid: userA, BookUid: bookA, Rating: 2})
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if got.Uid == "r-1" || got.Rating != 2 {
		t.Fatalf("got = %+v, want a fresh review", got)
	}
}

func TestReviewRepositoryInsert_concurrentDuplicateHitsUniqueIndex_reviewExists(t *testing.T) {
	r, db := newReviewRepo(t)
	failCreatesOn(t, db, entity.TableNameCustomerReviewEntity, uniqueViolation())

	_, err := r.Insert(bg(), repo.ReviewCreate{UserUid: userA, BookUid: bookA, Rating: 2})
	if !errors.Is(err, repo.ErrReviewExists) {
		t.Fatalf("err = %v, want ErrReviewExists", err)
	}
}

func TestReviewRepositoryInsert_otherWriteFailure_errorPassedThrough(t *testing.T) {
	r, db := newReviewRepo(t)
	failCreatesOn(t, db, entity.TableNameCustomerReviewEntity, errInjected)

	_, err := r.Insert(bg(), repo.ReviewCreate{UserUid: userA, BookUid: bookA, Rating: 2})
	assertInjected(t, err)
	if errors.Is(err, repo.ErrReviewExists) {
		t.Error("a non-unique failure was reported as a duplicate review")
	}
}

func TestReviewRepositoryFind_unknownOrSoftDeleted_recordNotFound(t *testing.T) {
	r, db := newReviewRepo(t)
	gone := seedReview(t, db, "r-1", userA, bookA, 5, baseTime)
	softDelete(t, db, gone)

	if _, err := r.FindByUid(bg(), absentUid); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("unknown uid: err = %v, want record not found", err)
	}
	if _, err := r.FindByUid(bg(), "r-1"); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("soft-deleted uid: err = %v, want record not found", err)
	}
	if _, err := r.FindOne(bg(), gone.Id); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("soft-deleted id: err = %v, want record not found", err)
	}
}

func TestReviewRepository_databaseUnavailable_everyOperationFails(t *testing.T) {
	r, db := newReviewRepo(t)
	closeDatabase(t, db)

	_, _, err := r.FindAll(bg(), option.NewReviewQueryOption())
	assertFails(t, "FindAll", err)
	_, err = r.Insert(bg(), repo.ReviewCreate{UserUid: userA, BookUid: bookA, Rating: 2})
	assertFails(t, "Insert", err)
	if errors.Is(err, repo.ErrReviewExists) {
		t.Error("a lookup failure was reported as a duplicate review")
	}
}

func TestIsUniqueViolation_errorVariants_onlySqlstate23505(t *testing.T) {
	if !isUniqueViolation(uniqueViolation()) {
		t.Error("23505 not recognised as a unique violation")
	}
	if !isUniqueViolation(errors.Join(errors.New("create"), uniqueViolation())) {
		t.Error("wrapped 23505 not recognised as a unique violation")
	}
	if isUniqueViolation(errInjected) || isUniqueViolation(nil) {
		t.Error("a non-postgres error was treated as a unique violation")
	}
}

func newWishlistRepo(t *testing.T) (*wishlistRepository, *gorm.DB) {
	t.Helper()
	db := newSchemaDB(t)
	return &wishlistRepository{q: queryOf(db)}, db
}

func seedWish(t *testing.T, db *gorm.DB, uid, userUid, bookUid string, createdAt time.Time) *entity.CustomerWishlistEntity {
	t.Helper()
	row := &entity.CustomerWishlistEntity{Uid: uid, UserUid: userUid, BookUid: bookUid, CreatedAt: createdAt, Version: 1}
	mustCreate(t, db, row)
	return row
}

func wishBooks(rows []*repo.WishlistResult) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.BookUid
	}
	return out
}

func wishFilter(userUid *string, bookUids *[]string, size, page int32) option.WishlistQueryOption {
	o := option.WishlistQueryOption{}
	return option.NewWishlistQueryOption(o.WithUserUid(userUid), o.WithBookUids(bookUids), o.WithPage(pageOf(size, page, "")))
}

func TestWishlistRepositoryFindAll_userAndBookFilters_onlyMatchingEntriesNewestFirst(t *testing.T) {
	r, db := newWishlistRepo(t)
	seedWish(t, db, "w-1", userA, bookA, hoursAfterBase(1))
	seedWish(t, db, "w-2", userA, bookB, hoursAfterBase(2))
	seedWish(t, db, "w-3", userA, bookC, hoursAfterBase(3))
	seedWish(t, db, "w-4", userB, bookA, hoursAfterBase(4))
	user, empty := userA, ""
	books := []string{bookA, bookC}
	noBooks := []string{}

	cases := []struct {
		name string
		opt  option.WishlistQueryOption
		want []string
	}{
		{"user", wishFilter(&user, nil, 10, 0), []string{bookC, bookB, bookA}},
		{"userAndBooks", wishFilter(&user, &books, 10, 0), []string{bookC, bookA}},
		{"booksAcrossUsers", wishFilter(nil, &books, 10, 0), []string{bookA, bookC, bookA}},
		{"blankUserAndEmptyBooksIgnored", wishFilter(&empty, &noBooks, 10, 0), []string{bookA, bookC, bookB, bookA}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			total, got, err := r.FindAll(bg(), c.opt)
			if err != nil {
				t.Fatalf("find all: %v", err)
			}
			if total != int64(len(c.want)) || !reflect.DeepEqual(wishBooks(got), c.want) {
				t.Fatalf("total = %d, books = %v, want %v", total, wishBooks(got), c.want)
			}
		})
	}

	total, got, err := r.FindAll(bg(), wishFilter(&user, nil, 2, 1))
	if err != nil {
		t.Fatalf("find second page: %v", err)
	}
	if total != 3 || !reflect.DeepEqual(wishBooks(got), []string{bookA}) {
		t.Fatalf("second page: total = %d, books = %v, want total 3 and [%s]", total, wishBooks(got), bookA)
	}
}

func TestWishlistRepositoryFindAllByUserUid_sameInstantEntries_newestFirstThenLatestInserted(t *testing.T) {
	r, db := newWishlistRepo(t)
	seedWish(t, db, "w-1", userA, bookA, hoursAfterBase(1))
	seedWish(t, db, "w-2", userA, bookB, hoursAfterBase(2))
	seedWish(t, db, "w-3", userA, bookC, hoursAfterBase(2))
	seedWish(t, db, "w-4", userB, bookA, hoursAfterBase(9))
	gone := seedWish(t, db, "w-5", userA, "018f0000-0000-7000-8000-0000000000b9", hoursAfterBase(8))
	softDelete(t, db, gone)

	got, err := r.FindAllByUserUid(bg(), userA)
	if err != nil {
		t.Fatalf("find all by user: %v", err)
	}
	if want := []string{bookC, bookB, bookA}; !reflect.DeepEqual(wishBooks(got), want) {
		t.Fatalf("books = %v, want %v", wishBooks(got), want)
	}
}

func TestWishlistRepositoryAddBooks_newAndAlreadyWishedBooks_onlyNewOnesAdded(t *testing.T) {
	r, db := newWishlistRepo(t)
	seedWish(t, db, "w-1", userA, bookA, hoursAfterBase(1))
	seedWish(t, db, "w-2", userB, bookB, hoursAfterBase(1))

	if err := r.AddBooks(bg(), userA, []string{bookA, bookB, "", bookB}); err != nil {
		t.Fatalf("add books: %v", err)
	}

	got, err := r.FindAllByUserUid(bg(), userA)
	if err != nil {
		t.Fatalf("find all by user: %v", err)
	}
	if want := []string{bookB, bookA}; !reflect.DeepEqual(wishBooks(got), want) {
		t.Fatalf("books = %v, want %v", wishBooks(got), want)
	}
	if got[1].Uid != "w-1" {
		t.Errorf("existing entry uid = %q, want w-1 left untouched", got[1].Uid)
	}
	added, err := r.FindOne(bg(), got[0].Id)
	if err != nil {
		t.Fatalf("find one: %v", err)
	}
	if added.Uid == "" || added.UserUid != userA || added.BookUid != bookB || added.CreatedAt.IsZero() {
		t.Errorf("added = %+v", added)
	}
	other, err := r.FindAllByUserUid(bg(), userB)
	if err != nil {
		t.Fatalf("find other user: %v", err)
	}
	if len(other) != 1 {
		t.Errorf("other user entries = %d, want 1", len(other))
	}
}

func TestWishlistRepositoryAddBooks_bookRemovedEarlier_addedAgain(t *testing.T) {
	r, db := newWishlistRepo(t)
	gone := seedWish(t, db, "w-1", userA, bookA, hoursAfterBase(1))
	softDelete(t, db, gone)

	if err := r.AddBooks(bg(), userA, []string{bookA}); err != nil {
		t.Fatalf("add books: %v", err)
	}

	got, err := r.FindAllByUserUid(bg(), userA)
	if err != nil {
		t.Fatalf("find all by user: %v", err)
	}
	if len(got) != 1 || got[0].BookUid != bookA || got[0].Uid == "w-1" {
		t.Fatalf("got = %+v, want one fresh entry for the book", got)
	}
}

func TestWishlistRepositoryAddBooks_nothingToAdd_noRowsWritten(t *testing.T) {
	r, db := newWishlistRepo(t)

	for _, books := range [][]string{nil, {}, {"", ""}} {
		if err := r.AddBooks(bg(), userA, books); err != nil {
			t.Fatalf("add %v: %v", books, err)
		}
	}
	var count int64
	if err := db.Model(&entity.CustomerWishlistEntity{}).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 0 {
		t.Errorf("rows = %d, want 0", count)
	}
}

func TestWishlistRepositoryFindOne_unknownId_recordNotFound(t *testing.T) {
	r, _ := newWishlistRepo(t)

	if _, err := r.FindOne(bg(), 404); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("err = %v, want record not found", err)
	}
}

func TestWishlistRepository_databaseUnavailable_everyOperationFails(t *testing.T) {
	r, db := newWishlistRepo(t)
	closeDatabase(t, db)

	_, _, err := r.FindAll(bg(), option.NewWishlistQueryOption())
	assertFails(t, "FindAll", err)
	_, err = r.FindAllByUserUid(bg(), userA)
	assertFails(t, "FindAllByUserUid", err)
	assertFails(t, "AddBooks", r.AddBooks(bg(), userA, []string{bookA}))
}
