package pgsql

import (
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

const (
	billingCategoryUid  = "018f0000-0000-7000-8000-0000000000d1"
	shippingCategoryUid = "018f0000-0000-7000-8000-0000000000d2"
)

func newFaqRepo(t *testing.T) (*faqRepository, *gorm.DB) {
	t.Helper()
	db := newSchemaDB(t)
	return &faqRepository{q: queryOf(db), db: db}, db
}

type faqSeed struct {
	uid        string
	question   string
	answer     string
	views      int64
	published  bool
	categoryId *int64
	createdAt  time.Time
}

func seedFaq(t *testing.T, db *gorm.DB, s faqSeed) *entity.CustomerFaqEntity {
	t.Helper()
	if s.createdAt.IsZero() {
		s.createdAt = baseTime
	}
	row := &entity.CustomerFaqEntity{
		Uid: s.uid, CategoryId: s.categoryId, Question: s.question, Answer: s.answer,
		ViewCount: s.views, Published: s.published, CreatedAt: s.createdAt, Version: 1,
	}
	mustCreate(t, db, row)
	return row
}

func faqQuestions(rows []*repo.FaqResult) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.Question
	}
	return out
}

func loadFaq(t *testing.T, db *gorm.DB, uid string) entity.CustomerFaqEntity {
	t.Helper()
	var got entity.CustomerFaqEntity
	if err := db.Unscoped().Where("uid = ?", uid).Take(&got).Error; err != nil {
		t.Fatalf("load faq %s: %v", uid, err)
	}
	return got
}

func TestFaqRepositoryFindAll_defaultOption_unpublishedExcluded(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "live", answer: "a", published: true})
	seedFaq(t, db, faqSeed{uid: "f-2", question: "draft", answer: "a", published: false})

	total, got, err := r.FindAll(bg(), option.NewFaqQueryOption())
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 1 || !reflect.DeepEqual(faqQuestions(got), []string{"live"}) {
		t.Fatalf("total = %d, questions = %v, want only the published faq", total, faqQuestions(got))
	}
}

func TestFaqRepositoryFindAll_publishedOnlyDisabled_draftsIncluded(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "live", answer: "a", published: true, views: 5})
	seedFaq(t, db, faqSeed{uid: "f-2", question: "draft", answer: "a", published: false, views: 1})

	total, got, err := r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqPublishedOnly(false)))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 2 || !reflect.DeepEqual(faqQuestions(got), []string{"live", "draft"}) {
		t.Fatalf("total = %d, questions = %v, want both faqs", total, faqQuestions(got))
	}
}

func TestFaqRepositoryFindAll_softDeletedFaq_excluded(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "kept", answer: "a", published: true})
	gone := seedFaq(t, db, faqSeed{uid: "f-2", question: "gone", answer: "a", published: true})
	softDelete(t, db, gone)

	total, got, err := r.FindAll(bg(), option.NewFaqQueryOption())
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 1 || !reflect.DeepEqual(faqQuestions(got), []string{"kept"}) {
		t.Fatalf("total = %d, questions = %v, want only the live faq", total, faqQuestions(got))
	}
}

func TestFaqRepositoryFindAll_categoryFilter_onlyThatCategoryWithItsUid(t *testing.T) {
	r, db := newFaqRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	shipping := seedCategory(t, db, shippingCategoryUid, "Shipping", nil)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "refunds", answer: "a", published: true, categoryId: &billing.Id})
	seedFaq(t, db, faqSeed{uid: "f-2", question: "tracking", answer: "a", published: true, categoryId: &shipping.Id})
	seedFaq(t, db, faqSeed{uid: "f-3", question: "general", answer: "a", published: true})

	total, got, err := r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqCategoryUid(billingCategoryUid)))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 1 || len(got) != 1 || got[0].Question != "refunds" {
		t.Fatalf("total = %d, questions = %v, want only the billing faq", total, faqQuestions(got))
	}
	if got[0].CategoryUid == nil || *got[0].CategoryUid != billingCategoryUid {
		t.Errorf("category uid = %v, want %s", got[0].CategoryUid, billingCategoryUid)
	}
}

func TestFaqRepositoryFindAll_uncategorizedAndCategorizedMixed_categoryUidOnlyWhereSet(t *testing.T) {
	r, db := newFaqRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "a-refunds", answer: "a", published: true, categoryId: &billing.Id})
	seedFaq(t, db, faqSeed{uid: "f-2", question: "b-general", answer: "a", published: true})
	seedFaq(t, db, faqSeed{uid: "f-3", question: "c-orphan", answer: "a", published: true, categoryId: i64p(999)})

	_, got, err := r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqPage(pageOf(10, 0, "question,asc"))))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	if got[0].CategoryUid == nil || *got[0].CategoryUid != billingCategoryUid {
		t.Errorf("categorized faq uid = %v, want %s", got[0].CategoryUid, billingCategoryUid)
	}
	if got[1].CategoryUid != nil || got[2].CategoryUid != nil {
		t.Errorf("category uids = %v/%v, want nil for an uncategorized faq and a dangling category id", got[1].CategoryUid, got[2].CategoryUid)
	}
}

func TestFaqRepositoryFindAll_categorySoftDeleted_filterStillMatchesItsFaqs(t *testing.T) {
	r, db := newFaqRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "refunds", answer: "a", published: true, categoryId: &billing.Id})
	softDelete(t, db, billing)

	total, got, err := r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqCategoryUid(billingCategoryUid)))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Fatalf("total = %d, len = %d, want the faq of the retired category", total, len(got))
	}
	if got[0].CategoryUid == nil || *got[0].CategoryUid != billingCategoryUid {
		t.Errorf("category uid = %v, want the retired category uid", got[0].CategoryUid)
	}
}

func TestFaqRepositoryFindAll_unknownCategory_emptyPageWithoutError(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "refunds", answer: "a", published: true})

	total, got, err := r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqCategoryUid(absentUid)))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 0 || got == nil || len(got) != 0 {
		t.Fatalf("total = %d, items = %v, want an empty non-nil page", total, got)
	}
}

func TestFaqRepositoryFindAll_blankCategoryUid_filterIgnored(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "refunds", answer: "a", published: true})

	total, _, err := r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqCategoryUid("")))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 1 {
		t.Fatalf("total = %d, want 1", total)
	}
}

func TestFaqRepositoryFindAll_paddedMixedCaseQuery_matchesQuestionOrAnswer(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "How do Refunds work?", answer: "slowly", published: true, views: 3})
	seedFaq(t, db, faqSeed{uid: "f-2", question: "Late parcel", answer: "Ask for a REFUND", published: true, views: 2})
	seedFaq(t, db, faqSeed{uid: "f-3", question: "Opening hours", answer: "9 to 5", published: true, views: 1})

	total, got, err := r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqQuery("  reFund ")))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if want := []string{"How do Refunds work?", "Late parcel"}; total != 2 || !reflect.DeepEqual(faqQuestions(got), want) {
		t.Fatalf("total = %d, questions = %v, want %v", total, faqQuestions(got), want)
	}
}

func TestFaqRepositoryFindAll_whitespaceQuery_filterIgnored(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "q1", answer: "a", published: true})
	seedFaq(t, db, faqSeed{uid: "f-2", question: "q2", answer: "a", published: true})

	total, _, err := r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqQuery("   ")))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 2 {
		t.Fatalf("total = %d, want 2", total)
	}
}

func TestFaqRepositoryFindAll_sortOmitted_mostViewedFirst(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "few", answer: "a", published: true, views: 2})
	seedFaq(t, db, faqSeed{uid: "f-2", question: "many", answer: "a", published: true, views: 90})
	seedFaq(t, db, faqSeed{uid: "f-3", question: "some", answer: "a", published: true, views: 40})

	_, got, err := r.FindAll(bg(), option.NewFaqQueryOption())
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if want := []string{"many", "some", "few"}; !reflect.DeepEqual(faqQuestions(got), want) {
		t.Fatalf("questions = %v, want %v", faqQuestions(got), want)
	}
}

func TestFaqRepositoryFindAll_secondPageSortedByQuestion_skipsOneFullPage(t *testing.T) {
	r, db := newFaqRepo(t)
	for _, q := range []string{"e", "c", "a", "d", "b"} {
		seedFaq(t, db, faqSeed{uid: "f-" + q, question: q, answer: "a", published: true})
	}

	total, got, err := r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqPage(pageOf(2, 1, "question,asc"))))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}
	if want := []string{"c", "d"}; !reflect.DeepEqual(faqQuestions(got), want) {
		t.Fatalf("questions = %v, want %v", faqQuestions(got), want)
	}
}

func TestFaqRepositoryFindAll_sortByCreatedAtDesc_newestFirst(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "old", answer: "a", published: true, views: 9, createdAt: hoursAfterBase(1)})
	seedFaq(t, db, faqSeed{uid: "f-2", question: "new", answer: "a", published: true, views: 1, createdAt: hoursAfterBase(3)})
	seedFaq(t, db, faqSeed{uid: "f-3", question: "mid", answer: "a", published: true, views: 5, createdAt: hoursAfterBase(2)})

	_, got, err := r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqPage(pageOf(10, 0, "created_at,desc"))))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if want := []string{"new", "mid", "old"}; !reflect.DeepEqual(faqQuestions(got), want) {
		t.Fatalf("questions = %v, want %v", faqQuestions(got), want)
	}
}

func TestFaqRepositoryFindByUid_categorizedFaq_resultCarriesCategoryUid(t *testing.T) {
	r, db := newFaqRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "refunds", answer: "slowly", views: 7, published: true, categoryId: &billing.Id})

	got, err := r.FindByUid(bg(), "f-1")
	if err != nil {
		t.Fatalf("find by uid: %v", err)
	}
	if got == nil || got.Uid != "f-1" || got.Question != "refunds" || got.Answer != "slowly" || got.ViewCount != 7 || !got.Published {
		t.Fatalf("got = %+v", got)
	}
	if got.CategoryUid == nil || *got.CategoryUid != billingCategoryUid {
		t.Errorf("category uid = %v, want %s", got.CategoryUid, billingCategoryUid)
	}
}

func TestFaqRepositoryFindByUid_unknownUid_nilWithoutError(t *testing.T) {
	r, _ := newFaqRepo(t)

	got, err := r.FindByUid(bg(), absentUid)
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil", got, err)
	}
}

func TestFaqRepositoryFindByUid_softDeletedFaq_nilWithoutError(t *testing.T) {
	r, db := newFaqRepo(t)
	gone := seedFaq(t, db, faqSeed{uid: "f-1", question: "gone", answer: "a", published: true})
	softDelete(t, db, gone)

	got, err := r.FindByUid(bg(), "f-1")
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil", got, err)
	}
}

func TestFaqRepositoryInsert_categorizedDraft_persistedWithZeroViewsAndVersionOne(t *testing.T) {
	r, db := newFaqRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)

	got, err := r.Insert(bg(), repo.FaqCreate{CategoryId: &billing.Id, Question: "q", Answer: "a", Published: false})
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if got.Uid == "" || got.Id == 0 {
		t.Fatalf("identifiers = %d/%q, want both assigned", got.Id, got.Uid)
	}
	if got.CategoryUid == nil || *got.CategoryUid != billingCategoryUid {
		t.Errorf("category uid = %v, want %s", got.CategoryUid, billingCategoryUid)
	}
	stored := loadFaq(t, db, got.Uid)
	if stored.Question != "q" || stored.Answer != "a" || stored.Published || stored.ViewCount != 0 || stored.Version != 1 {
		t.Errorf("stored = %+v", stored)
	}
	if stored.UpdatedAt == nil || stored.CreatedAt.IsZero() {
		t.Errorf("timestamps = %v/%v, want both set by the application", stored.CreatedAt, stored.UpdatedAt)
	}
}

func TestFaqRepositoryInsert_twoFaqs_distinctUids(t *testing.T) {
	r, _ := newFaqRepo(t)

	first, err := r.Insert(bg(), repo.FaqCreate{Question: "q1", Answer: "a", Published: true})
	if err != nil {
		t.Fatalf("insert first: %v", err)
	}
	second, err := r.Insert(bg(), repo.FaqCreate{Question: "q2", Answer: "a", Published: true})
	if err != nil {
		t.Fatalf("insert second: %v", err)
	}
	if first.Uid == second.Uid {
		t.Fatalf("both faqs share uid %q", first.Uid)
	}
	if first.CategoryUid != nil {
		t.Errorf("category uid = %v, want nil for an uncategorized faq", first.CategoryUid)
	}
}

func TestFaqRepositoryUpdateByUid_existingFaq_fieldsAppliedAndVersionBumped(t *testing.T) {
	r, db := newFaqRepo(t)
	shipping := seedCategory(t, db, shippingCategoryUid, "Shipping", nil)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "old q", answer: "old a", views: 11, published: false})

	got, err := r.UpdateByUid(bg(), "f-1", func(f *repo.FaqResult) {
		f.Question = "new q"
		f.Answer = "new a"
		f.Published = true
		f.CategoryId = &shipping.Id
		f.ViewCount = 9999
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if got == nil || got.Question != "new q" || got.Answer != "new a" || !got.Published {
		t.Fatalf("got = %+v", got)
	}
	if got.CategoryUid == nil || *got.CategoryUid != shippingCategoryUid {
		t.Errorf("category uid = %v, want %s", got.CategoryUid, shippingCategoryUid)
	}
	stored := loadFaq(t, db, "f-1")
	if stored.Question != "new q" || stored.Answer != "new a" || !stored.Published || stored.Version != 2 {
		t.Errorf("stored = %+v", stored)
	}
	if stored.ViewCount != 11 {
		t.Errorf("view count = %d, want 11 untouched by an edit", stored.ViewCount)
	}
	if stored.UpdatedAt == nil || !stored.UpdatedAt.After(baseTime) {
		t.Errorf("updated_at = %v, want refreshed", stored.UpdatedAt)
	}
}

func TestFaqRepositoryUpdateByUid_unknownUid_nilAndApplySkipped(t *testing.T) {
	r, _ := newFaqRepo(t)
	applied := false

	got, err := r.UpdateByUid(bg(), absentUid, func(*repo.FaqResult) { applied = true })
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil", got, err)
	}
	if applied {
		t.Error("apply ran for a faq that does not exist")
	}
}

func TestFaqRepositoryUpdateByUid_saveFails_errorSurfacedAndRowUnchanged(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "old q", answer: "a", published: true})
	failCreatesOn(t, db, entity.TableNameCustomerFaqEntity, errInjected)

	got, err := r.UpdateByUid(bg(), "f-1", func(f *repo.FaqResult) { f.Question = "new q" })
	assertInjected(t, err)
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
	if stored := loadFaq(t, db, "f-1"); stored.Question != "old q" || stored.Version != 1 {
		t.Errorf("stored = %+v, want the row untouched", stored)
	}
}

func TestFaqRepositorySoftDeleteByUid_existingFaq_hiddenButRowKept(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "q", answer: "a", published: true})

	deleted, err := r.SoftDeleteByUid(bg(), "f-1")
	if err != nil || !deleted {
		t.Fatalf("deleted = %v, err = %v, want true, nil", deleted, err)
	}
	if stored := loadFaq(t, db, "f-1"); !stored.DeletedAt.Valid {
		t.Error("deleted_at not set; the row must be soft deleted, not removed")
	}
	again, err := r.SoftDeleteByUid(bg(), "f-1")
	if err != nil || again {
		t.Fatalf("second delete = %v, err = %v, want false, nil", again, err)
	}
}

func TestFaqRepositorySoftDeleteByUid_unknownUid_false(t *testing.T) {
	r, _ := newFaqRepo(t)

	deleted, err := r.SoftDeleteByUid(bg(), absentUid)
	if err != nil || deleted {
		t.Fatalf("deleted = %v, err = %v, want false, nil", deleted, err)
	}
}

func TestFaqRepositoryIncrementViewCount_existingFaq_onlyThatFaqIncremented(t *testing.T) {
	r, db := newFaqRepo(t)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "q1", answer: "a", published: true, views: 4})
	seedFaq(t, db, faqSeed{uid: "f-2", question: "q2", answer: "a", published: true, views: 4})

	for i := 0; i < 2; i++ {
		if err := r.IncrementViewCount(bg(), "f-1"); err != nil {
			t.Fatalf("increment: %v", err)
		}
	}
	if got := loadFaq(t, db, "f-1").ViewCount; got != 6 {
		t.Errorf("view count = %d, want 6", got)
	}
	if got := loadFaq(t, db, "f-2").ViewCount; got != 4 {
		t.Errorf("other faq view count = %d, want 4", got)
	}
}

func TestFaqRepositoryIncrementViewCount_unknownUid_noError(t *testing.T) {
	r, _ := newFaqRepo(t)

	if err := r.IncrementViewCount(bg(), absentUid); err != nil {
		t.Fatalf("increment: %v", err)
	}
}

func TestFaqRepository_categoryLookupFails_errorSurfaced(t *testing.T) {
	r, db := newFaqRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	seedFaq(t, db, faqSeed{uid: "f-1", question: "q", answer: "a", published: true, categoryId: &billing.Id})
	failQueriesOn(t, db, entity.TableNameCustomerTicketCategoryEntity)

	_, _, err := r.FindAll(bg(), option.NewFaqQueryOption())
	assertInjected(t, err)
	_, _, err = r.FindAll(bg(), option.NewFaqQueryOption(option.WithFaqCategoryUid(billingCategoryUid)))
	assertInjected(t, err)
	_, err = r.FindByUid(bg(), "f-1")
	assertInjected(t, err)
	_, err = r.Insert(bg(), repo.FaqCreate{CategoryId: &billing.Id, Question: "q2", Answer: "a"})
	assertInjected(t, err)
	_, err = r.UpdateByUid(bg(), "f-1", func(*repo.FaqResult) {})
	assertInjected(t, err)
}

func TestFaqRepository_databaseUnavailable_everyOperationFails(t *testing.T) {
	r, db := newFaqRepo(t)
	closeDatabase(t, db)

	_, _, err := r.FindAll(bg(), option.NewFaqQueryOption())
	assertFails(t, "FindAll", err)
	_, err = r.FindByUid(bg(), "f-1")
	assertFails(t, "FindByUid", err)
	_, err = r.Insert(bg(), repo.FaqCreate{Question: "q", Answer: "a"})
	assertFails(t, "Insert", err)
	_, err = r.UpdateByUid(bg(), "f-1", func(*repo.FaqResult) {})
	assertFails(t, "UpdateByUid", err)
	_, err = r.SoftDeleteByUid(bg(), "f-1")
	assertFails(t, "SoftDeleteByUid", err)
	assertFails(t, "IncrementViewCount", r.IncrementViewCount(bg(), "f-1"))
}
