package pgsql

import (
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

const (
	agentA = "018f0000-0000-7000-8000-0000000000c1"
	agentB = "018f0000-0000-7000-8000-0000000000c2"
)

func newTicketRepo(t *testing.T) (*ticketRepository, *gorm.DB) {
	t.Helper()
	db := newSchemaDB(t)
	return &ticketRepository{q: queryOf(db), db: db}, db
}

type ticketSeed struct {
	uid        string
	customer   string
	assignee   *string
	categoryId *int64
	subject    string
	status     string
	priority   string
	createdAt  time.Time
}

func seedTicket(t *testing.T, db *gorm.DB, s ticketSeed) *entity.CustomerTicketEntity {
	t.Helper()
	if s.createdAt.IsZero() {
		s.createdAt = baseTime
	}
	if s.status == "" {
		s.status = "OPEN"
	}
	if s.priority == "" {
		s.priority = "LOW"
	}
	if s.customer == "" {
		s.customer = userA
	}
	row := &entity.CustomerTicketEntity{
		Uid: s.uid, CustomerUid: s.customer, AssigneeUid: s.assignee, CategoryId: s.categoryId,
		Subject: s.subject, Description: "d", Status: s.status, Priority: s.priority,
		CreatedAt: s.createdAt, Version: 1,
	}
	mustCreate(t, db, row)
	return row
}

func ticketSubjects(rows []*repo.TicketResult) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.Subject
	}
	return out
}

func loadTicket(t *testing.T, db *gorm.DB, uid string) entity.CustomerTicketEntity {
	t.Helper()
	var got entity.CustomerTicketEntity
	if err := db.Unscoped().Where("uid = ?", uid).Take(&got).Error; err != nil {
		t.Fatalf("load ticket %s: %v", uid, err)
	}
	return got
}

func seedTicketMatrix(t *testing.T, db *gorm.DB) {
	t.Helper()
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	a, b := agentA, agentB
	seedTicket(t, db, ticketSeed{uid: "t-1", subject: "s1", customer: userA, assignee: &a, status: "OPEN", priority: "HIGH", categoryId: &billing.Id, createdAt: hoursAfterBase(1)})
	seedTicket(t, db, ticketSeed{uid: "t-2", subject: "s2", customer: userA, assignee: &b, status: "CLOSED", priority: "LOW", createdAt: hoursAfterBase(2)})
	seedTicket(t, db, ticketSeed{uid: "t-3", subject: "s3", customer: userB, assignee: &a, status: "OPEN", priority: "LOW", createdAt: hoursAfterBase(3)})
	seedTicket(t, db, ticketSeed{uid: "t-4", subject: "s4", customer: userB, status: "IN_PROGRESS", priority: "URGENT", categoryId: &billing.Id, createdAt: hoursAfterBase(4)})
}

func TestTicketRepositoryFindAll_singleFilter_onlyMatchingTicketsNewestFirst(t *testing.T) {
	cases := []struct {
		name string
		opt  option.TicketQueryOptionFunc
		want []string
	}{
		{"customer", option.WithTicketCustomerUid(userA), []string{"s2", "s1"}},
		{"assignee", option.WithTicketAssigneeUid(agentA), []string{"s3", "s1"}},
		{"status", option.WithTicketStatus(constant.TICKET_STATUS_OPEN), []string{"s3", "s1"}},
		{"priority", option.WithTicketPriority(constant.TICKET_PRIORITY_LOW), []string{"s3", "s2"}},
		{"category", option.WithTicketCategoryUid(billingCategoryUid), []string{"s4", "s1"}},
		{"blankCustomerIgnored", option.WithTicketCustomerUid(""), []string{"s4", "s3", "s2", "s1"}},
		{"blankAssigneeIgnored", option.WithTicketAssigneeUid(""), []string{"s4", "s3", "s2", "s1"}},
		{"blankCategoryIgnored", option.WithTicketCategoryUid(""), []string{"s4", "s3", "s2", "s1"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, db := newTicketRepo(t)
			seedTicketMatrix(t, db)

			total, got, err := r.FindAll(bg(), option.NewTicketQueryOption(c.opt))
			if err != nil {
				t.Fatalf("find all: %v", err)
			}
			if total != int64(len(c.want)) || !reflect.DeepEqual(ticketSubjects(got), c.want) {
				t.Fatalf("total = %d, subjects = %v, want %v", total, ticketSubjects(got), c.want)
			}
		})
	}
}

func TestTicketRepositoryFindAll_combinedFilters_intersectionOnly(t *testing.T) {
	r, db := newTicketRepo(t)
	seedTicketMatrix(t, db)

	total, got, err := r.FindAll(bg(), option.NewTicketQueryOption(
		option.WithTicketCustomerUid(userB),
		option.WithTicketStatus(constant.TICKET_STATUS_OPEN),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 1 || !reflect.DeepEqual(ticketSubjects(got), []string{"s3"}) {
		t.Fatalf("total = %d, subjects = %v, want [s3]", total, ticketSubjects(got))
	}
}

func TestTicketRepositoryFindAll_categorizedTickets_resultsCarryCategoryUid(t *testing.T) {
	r, db := newTicketRepo(t)
	seedTicketMatrix(t, db)

	_, got, err := r.FindAll(bg(), option.NewTicketQueryOption(option.WithTicketPage(pageOf(10, 0, "subject,asc"))))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("len = %d, want 4", len(got))
	}
	if got[0].CategoryUid == nil || *got[0].CategoryUid != billingCategoryUid {
		t.Errorf("s1 category uid = %v, want %s", got[0].CategoryUid, billingCategoryUid)
	}
	if got[1].CategoryUid != nil {
		t.Errorf("s2 category uid = %v, want nil", got[1].CategoryUid)
	}
	if got[0].AssigneeUid == nil || *got[0].AssigneeUid != agentA || got[3].AssigneeUid != nil {
		t.Errorf("assignees = %v/%v, want %s/nil", got[0].AssigneeUid, got[3].AssigneeUid, agentA)
	}
}

func TestTicketRepositoryFindAll_unknownCategory_emptyPageWithoutError(t *testing.T) {
	r, db := newTicketRepo(t)
	seedTicketMatrix(t, db)

	total, got, err := r.FindAll(bg(), option.NewTicketQueryOption(option.WithTicketCategoryUid(absentUid)))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 0 || got == nil || len(got) != 0 {
		t.Fatalf("total = %d, items = %v, want an empty non-nil page", total, got)
	}
}

func TestTicketRepositoryFindAll_secondPageSortedBySubjectDesc_skipsOneFullPage(t *testing.T) {
	r, db := newTicketRepo(t)
	seedTicketMatrix(t, db)

	total, got, err := r.FindAll(bg(), option.NewTicketQueryOption(option.WithTicketPage(pageOf(3, 1, "subject,desc"))))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 4 || !reflect.DeepEqual(ticketSubjects(got), []string{"s1"}) {
		t.Fatalf("total = %d, subjects = %v, want total 4 and [s1]", total, ticketSubjects(got))
	}
}

func TestTicketRepositoryFindAll_softDeletedTicket_excluded(t *testing.T) {
	r, db := newTicketRepo(t)
	seedTicket(t, db, ticketSeed{uid: "t-1", subject: "kept"})
	gone := seedTicket(t, db, ticketSeed{uid: "t-2", subject: "gone"})
	softDelete(t, db, gone)

	total, got, err := r.FindAll(bg(), option.NewTicketQueryOption())
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 1 || !reflect.DeepEqual(ticketSubjects(got), []string{"kept"}) {
		t.Fatalf("total = %d, subjects = %v, want [kept]", total, ticketSubjects(got))
	}
}

func TestTicketRepositoryFindByUid_categorizedTicket_everyFieldMapped(t *testing.T) {
	r, db := newTicketRepo(t)
	seedTicketMatrix(t, db)

	got, err := r.FindByUid(bg(), "t-1")
	if err != nil {
		t.Fatalf("find by uid: %v", err)
	}
	if got == nil || got.Uid != "t-1" || got.CustomerUid != userA || got.Subject != "s1" || got.Description != "d" ||
		got.Status != "OPEN" || got.Priority != "HIGH" || !got.CreatedAt.Equal(hoursAfterBase(1)) {
		t.Fatalf("got = %+v", got)
	}
	if got.CategoryUid == nil || *got.CategoryUid != billingCategoryUid {
		t.Errorf("category uid = %v, want %s", got.CategoryUid, billingCategoryUid)
	}
}

func TestTicketRepositoryFindByUid_unknownOrSoftDeleted_nilWithoutError(t *testing.T) {
	r, db := newTicketRepo(t)
	gone := seedTicket(t, db, ticketSeed{uid: "t-1", subject: "gone"})
	softDelete(t, db, gone)

	for _, uid := range []string{absentUid, "t-1"} {
		got, err := r.FindByUid(bg(), uid)
		if err != nil || got != nil {
			t.Fatalf("uid %s: got = %+v, err = %v, want nil, nil", uid, got, err)
		}
	}
}

func TestTicketRepositoryInsert_categorizedTicket_persistedUnassignedWithVersionOne(t *testing.T) {
	r, db := newTicketRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)

	got, err := r.Insert(bg(), repo.TicketCreate{
		CustomerUid: userA, CategoryId: &billing.Id, Subject: "broken", Description: "it broke", Status: "OPEN", Priority: "MEDIUM",
	})
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if got.Uid == "" || got.Id == 0 {
		t.Fatalf("identifiers = %d/%q, want both assigned", got.Id, got.Uid)
	}
	if got.CategoryUid == nil || *got.CategoryUid != billingCategoryUid {
		t.Errorf("category uid = %v, want %s", got.CategoryUid, billingCategoryUid)
	}
	stored := loadTicket(t, db, got.Uid)
	if stored.CustomerUid != userA || stored.Subject != "broken" || stored.Description != "it broke" ||
		stored.Status != "OPEN" || stored.Priority != "MEDIUM" || stored.Version != 1 {
		t.Errorf("stored = %+v", stored)
	}
	if stored.AssigneeUid != nil || stored.ResolvedAt != nil {
		t.Errorf("assignee/resolved = %v/%v, want a new ticket unassigned and unresolved", stored.AssigneeUid, stored.ResolvedAt)
	}
}

func TestTicketRepositoryUpdateByUid_existingTicket_fieldsAppliedAndVersionBumped(t *testing.T) {
	r, db := newTicketRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	seedTicket(t, db, ticketSeed{uid: "t-1", subject: "old", status: "OPEN", priority: "LOW"})
	resolvedAt := hoursAfterBase(5)
	assignee := agentA

	got, err := r.UpdateByUid(bg(), "t-1", func(tk *repo.TicketResult) {
		tk.Subject = "new"
		tk.Description = "new d"
		tk.Status = "RESOLVED"
		tk.Priority = "URGENT"
		tk.AssigneeUid = &assignee
		tk.CategoryId = &billing.Id
		tk.ResolvedAt = &resolvedAt
		tk.CustomerUid = userB
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if got == nil || got.Subject != "new" || got.Status != "RESOLVED" || got.Priority != "URGENT" {
		t.Fatalf("got = %+v", got)
	}
	if got.CategoryUid == nil || *got.CategoryUid != billingCategoryUid {
		t.Errorf("category uid = %v, want %s", got.CategoryUid, billingCategoryUid)
	}
	stored := loadTicket(t, db, "t-1")
	if stored.Subject != "new" || stored.Description != "new d" || stored.Status != "RESOLVED" || stored.Priority != "URGENT" || stored.Version != 2 {
		t.Errorf("stored = %+v", stored)
	}
	if stored.AssigneeUid == nil || *stored.AssigneeUid != agentA {
		t.Errorf("assignee = %v, want %s", stored.AssigneeUid, agentA)
	}
	if stored.ResolvedAt == nil || !stored.ResolvedAt.Equal(resolvedAt) {
		t.Errorf("resolved_at = %v, want %v", stored.ResolvedAt, resolvedAt)
	}
	if stored.CustomerUid != userA {
		t.Errorf("customer uid = %s, want the owner unchanged by an edit", stored.CustomerUid)
	}
	if stored.UpdatedAt == nil || !stored.UpdatedAt.After(baseTime) {
		t.Errorf("updated_at = %v, want refreshed", stored.UpdatedAt)
	}
}

func TestTicketRepositoryUpdateByUid_unknownUid_nilAndApplySkipped(t *testing.T) {
	r, _ := newTicketRepo(t)
	applied := false

	got, err := r.UpdateByUid(bg(), absentUid, func(*repo.TicketResult) { applied = true })
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil", got, err)
	}
	if applied {
		t.Error("apply ran for a ticket that does not exist")
	}
}

func TestTicketRepositoryUpdateByUid_saveFails_errorSurfacedAndRowUnchanged(t *testing.T) {
	r, db := newTicketRepo(t)
	seedTicket(t, db, ticketSeed{uid: "t-1", subject: "old"})
	failCreatesOn(t, db, entity.TableNameCustomerTicketEntity, errInjected)

	got, err := r.UpdateByUid(bg(), "t-1", func(tk *repo.TicketResult) { tk.Subject = "new" })
	assertInjected(t, err)
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
	if stored := loadTicket(t, db, "t-1"); stored.Subject != "old" || stored.Version != 1 {
		t.Errorf("stored = %+v, want the row untouched", stored)
	}
}

func TestTicketRepository_categoryLookupFails_errorSurfaced(t *testing.T) {
	r, db := newTicketRepo(t)
	seedTicketMatrix(t, db)
	billingId, err := categoryIdByUid(bg(), r.q, billingCategoryUid, false)
	if err != nil || billingId == nil {
		t.Fatalf("billing category id = %v, err = %v", billingId, err)
	}
	failQueriesOn(t, db, entity.TableNameCustomerTicketCategoryEntity)

	_, _, err = r.FindAll(bg(), option.NewTicketQueryOption())
	assertInjected(t, err)
	_, _, err = r.FindAll(bg(), option.NewTicketQueryOption(option.WithTicketCategoryUid(billingCategoryUid)))
	assertInjected(t, err)
	_, err = r.FindByUid(bg(), "t-1")
	assertInjected(t, err)
	_, err = r.Insert(bg(), repo.TicketCreate{CustomerUid: userA, CategoryId: billingId, Subject: "s", Description: "d", Status: "OPEN", Priority: "LOW"})
	assertInjected(t, err)
	_, err = r.UpdateByUid(bg(), "t-1", func(*repo.TicketResult) {})
	assertInjected(t, err)
}

func TestTicketRepository_databaseUnavailable_everyOperationFails(t *testing.T) {
	r, db := newTicketRepo(t)
	closeDatabase(t, db)

	_, _, err := r.FindAll(bg(), option.NewTicketQueryOption())
	assertFails(t, "FindAll", err)
	_, err = r.FindByUid(bg(), "t-1")
	assertFails(t, "FindByUid", err)
	_, err = r.Insert(bg(), repo.TicketCreate{CustomerUid: userA, Subject: "s", Description: "d", Status: "OPEN", Priority: "LOW"})
	assertFails(t, "Insert", err)
	_, err = r.UpdateByUid(bg(), "t-1", func(*repo.TicketResult) {})
	assertFails(t, "UpdateByUid", err)
}

func newTicketCommentRepo(t *testing.T) (*ticketCommentRepository, *gorm.DB) {
	t.Helper()
	db := newSchemaDB(t)
	return &ticketCommentRepository{q: queryOf(db)}, db
}

func seedComment(t *testing.T, db *gorm.DB, uid string, ticketId int64, body string, internal bool, createdAt time.Time) *entity.CustomerTicketCommentEntity {
	t.Helper()
	row := &entity.CustomerTicketCommentEntity{
		Uid: uid, TicketId: ticketId, AuthorUid: userA, AuthorRole: "CUSTOMER", Body: body, Internal: internal,
		CreatedAt: createdAt, Version: 1,
	}
	mustCreate(t, db, row)
	return row
}

func commentBodies(rows []*repo.TicketCommentResult) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.Body
	}
	return out
}

func TestTicketCommentRepositoryFindAllByTicketId_internalExcluded_publicCommentsOldestFirst(t *testing.T) {
	r, db := newTicketCommentRepo(t)
	seedComment(t, db, "c-1", 1, "second", false, hoursAfterBase(2))
	seedComment(t, db, "c-2", 1, "staff note", true, hoursAfterBase(1))
	seedComment(t, db, "c-3", 1, "first", false, hoursAfterBase(1))
	seedComment(t, db, "c-4", 2, "other ticket", false, hoursAfterBase(1))

	got, err := r.FindAllByTicketId(bg(), 1, false)
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if want := []string{"first", "second"}; !reflect.DeepEqual(commentBodies(got), want) {
		t.Fatalf("bodies = %v, want %v", commentBodies(got), want)
	}
}

func TestTicketCommentRepositoryFindAllByTicketId_internalIncluded_sameInstantOrderedByInsertion(t *testing.T) {
	r, db := newTicketCommentRepo(t)
	seedComment(t, db, "c-1", 1, "second", false, hoursAfterBase(2))
	seedComment(t, db, "c-2", 1, "staff note", true, hoursAfterBase(1))
	seedComment(t, db, "c-3", 1, "first", false, hoursAfterBase(1))
	gone := seedComment(t, db, "c-4", 1, "retracted", false, hoursAfterBase(3))
	softDelete(t, db, gone)

	got, err := r.FindAllByTicketId(bg(), 1, true)
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if want := []string{"staff note", "first", "second"}; !reflect.DeepEqual(commentBodies(got), want) {
		t.Fatalf("bodies = %v, want %v", commentBodies(got), want)
	}
	if !got[0].Internal || got[1].Internal {
		t.Errorf("internal flags = %v/%v, want true/false", got[0].Internal, got[1].Internal)
	}
}

func TestTicketCommentRepositoryFindAllByTicketId_ticketWithoutComments_emptySlice(t *testing.T) {
	r, _ := newTicketCommentRepo(t)

	got, err := r.FindAllByTicketId(bg(), 42, true)
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Fatalf("got = %v, want an empty non-nil slice", got)
	}
}

func TestTicketCommentRepositoryInsert_internalStaffComment_persistedAndListed(t *testing.T) {
	r, _ := newTicketCommentRepo(t)

	got, err := r.Insert(bg(), repo.TicketCommentCreate{TicketId: 7, AuthorUid: agentA, AuthorRole: "STAFF", Body: "checking", Internal: true})
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if got.Uid == "" || got.Id == 0 || got.TicketId != 7 || got.AuthorUid != agentA || got.AuthorRole != "STAFF" || got.Body != "checking" || !got.Internal {
		t.Fatalf("got = %+v", got)
	}
	if got.UpdatedAt == nil || got.CreatedAt.IsZero() {
		t.Errorf("timestamps = %v/%v, want both set by the application", got.CreatedAt, got.UpdatedAt)
	}
	public, err := r.FindAllByTicketId(bg(), 7, false)
	if err != nil {
		t.Fatalf("find public: %v", err)
	}
	all, err := r.FindAllByTicketId(bg(), 7, true)
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if len(public) != 0 || len(all) != 1 || all[0].Uid != got.Uid {
		t.Fatalf("public = %d, all = %d, want the internal comment visible to staff only", len(public), len(all))
	}
}

func TestTicketCommentRepository_databaseUnavailable_everyOperationFails(t *testing.T) {
	r, db := newTicketCommentRepo(t)
	closeDatabase(t, db)

	_, err := r.FindAllByTicketId(bg(), 1, true)
	assertFails(t, "FindAllByTicketId", err)
	_, err = r.Insert(bg(), repo.TicketCommentCreate{TicketId: 1, AuthorUid: userA, AuthorRole: "CUSTOMER", Body: "b"})
	assertFails(t, "Insert", err)
}
