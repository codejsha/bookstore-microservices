package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

const (
	testCustomerUid = "018f0000-0000-7000-8000-0000000000c1"
	testCategoryUid = "018f0000-0000-7000-8000-0000000000d1"
	testParentUid   = "018f0000-0000-7000-8000-0000000000d2"
	testTicketUid   = "018f0000-0000-7000-8000-0000000000e1"
	testFaqUid      = "018f0000-0000-7000-8000-0000000000f1"
)

func strp(s string) *string { return &s }

// ─── repo stubs ─────────────────────────────────────────────────────────────

type stubTicketRepo struct {
	repo.TicketRepo
	rows     map[string]*repo.TicketResult
	inserted []repo.TicketCreate
	listOpt  *option.TicketQueryOption
}

func (s *stubTicketRepo) FindAll(_ context.Context, opt option.TicketQueryOption) (int64, []*repo.TicketResult, error) {
	s.listOpt = &opt
	out := make([]*repo.TicketResult, 0, len(s.rows))
	for _, r := range s.rows {
		out = append(out, r)
	}
	return int64(len(out)), out, nil
}

func (s *stubTicketRepo) FindByUid(_ context.Context, uid string) (*repo.TicketResult, error) {
	return s.rows[uid], nil
}

func (s *stubTicketRepo) Insert(_ context.Context, t repo.TicketCreate) (*repo.TicketResult, error) {
	s.inserted = append(s.inserted, t)
	row := &repo.TicketResult{
		Id: int64(len(s.inserted)), Uid: testTicketUid, CustomerUid: t.CustomerUid, CategoryId: t.CategoryId,
		Subject: t.Subject, Description: t.Description, Status: t.Status, Priority: t.Priority, CreatedAt: time.Now(),
	}
	if t.CategoryId != nil {
		row.CategoryUid = strp(testCategoryUid)
	}
	s.rows[row.Uid] = row
	return row, nil
}

func (s *stubTicketRepo) UpdateByUid(_ context.Context, uid string, apply func(*repo.TicketResult)) (*repo.TicketResult, error) {
	row, ok := s.rows[uid]
	if !ok {
		return nil, nil
	}
	apply(row)
	return row, nil
}

type stubCommentRepo struct {
	repo.TicketCommentRepo
	inserted        []repo.TicketCommentCreate
	includeInternal *bool
}

func (s *stubCommentRepo) FindAllByTicketId(_ context.Context, _ int64, includeInternal bool) ([]*repo.TicketCommentResult, error) {
	s.includeInternal = &includeInternal
	return []*repo.TicketCommentResult{}, nil
}

func (s *stubCommentRepo) Insert(_ context.Context, c repo.TicketCommentCreate) (*repo.TicketCommentResult, error) {
	s.inserted = append(s.inserted, c)
	return &repo.TicketCommentResult{Id: 1, Uid: "cm-1", TicketId: c.TicketId, AuthorUid: c.AuthorUid, AuthorRole: c.AuthorRole, Body: c.Body, Internal: c.Internal, CreatedAt: time.Now()}, nil
}

type stubCategoryRepo struct {
	repo.TicketCategoryRepo
	byUid    map[string]*repo.TicketCategoryResult
	byName   map[string]*repo.TicketCategoryResult
	inserted []repo.TicketCategoryCreate
	deleted  []string
}

func (s *stubCategoryRepo) FindIdByUid(_ context.Context, uid string) (*int64, error) {
	if c, ok := s.byUid[uid]; ok {
		return &c.Id, nil
	}
	return nil, nil
}

func (s *stubCategoryRepo) FindByUid(_ context.Context, uid string) (*repo.TicketCategoryResult, error) {
	return s.byUid[uid], nil
}

func (s *stubCategoryRepo) FindByName(_ context.Context, name string) (*repo.TicketCategoryResult, error) {
	return s.byName[name], nil
}

func (s *stubCategoryRepo) Insert(_ context.Context, c repo.TicketCategoryCreate) (*repo.TicketCategoryResult, error) {
	s.inserted = append(s.inserted, c)
	return &repo.TicketCategoryResult{Id: 99, Uid: "cat-new", Name: c.Name, Description: c.Description, ParentId: c.ParentId, CreatedAt: time.Now()}, nil
}

func (s *stubCategoryRepo) UpdateByUid(_ context.Context, uid string, apply func(*repo.TicketCategoryResult)) (*repo.TicketCategoryResult, error) {
	row, ok := s.byUid[uid]
	if !ok {
		return nil, nil
	}
	apply(row)
	return row, nil
}

func (s *stubCategoryRepo) SoftDeleteByUid(_ context.Context, uid string) (bool, error) {
	if _, ok := s.byUid[uid]; !ok {
		return false, nil
	}
	s.deleted = append(s.deleted, uid)
	return true, nil
}

type stubFaqRepo struct {
	repo.FaqRepo
	rows         map[string]*repo.FaqResult
	increments   []string
	incrementErr error
}

func (s *stubFaqRepo) FindByUid(_ context.Context, uid string) (*repo.FaqResult, error) {
	return s.rows[uid], nil
}

func (s *stubFaqRepo) IncrementViewCount(_ context.Context, uid string) error {
	s.increments = append(s.increments, uid)
	return s.incrementErr
}

func (s *stubFaqRepo) Insert(_ context.Context, f repo.FaqCreate) (*repo.FaqResult, error) {
	return &repo.FaqResult{Id: 1, Uid: testFaqUid, CategoryId: f.CategoryId, Question: f.Question, Answer: f.Answer, Published: f.Published, CreatedAt: time.Now()}, nil
}

func (s *stubFaqRepo) UpdateByUid(_ context.Context, uid string, apply func(*repo.FaqResult)) (*repo.FaqResult, error) {
	row, ok := s.rows[uid]
	if !ok {
		return nil, nil
	}
	apply(row)
	return row, nil
}

func (s *stubFaqRepo) SoftDeleteByUid(_ context.Context, uid string) (bool, error) {
	_, ok := s.rows[uid]
	return ok, nil
}

func newSupportFixture() (*supportService, *stubTicketRepo, *stubCommentRepo, *stubCategoryRepo, *stubFaqRepo) {
	category := &repo.TicketCategoryResult{Id: 10, Uid: testCategoryUid, Name: "Billing", CreatedAt: time.Now()}
	parent := &repo.TicketCategoryResult{Id: 11, Uid: testParentUid, Name: "Root", CreatedAt: time.Now()}
	tickets := &stubTicketRepo{rows: map[string]*repo.TicketResult{}}
	comments := &stubCommentRepo{}
	categories := &stubCategoryRepo{
		byUid:  map[string]*repo.TicketCategoryResult{testCategoryUid: category, testParentUid: parent},
		byName: map[string]*repo.TicketCategoryResult{"Billing": category, "Root": parent},
	}
	faqs := &stubFaqRepo{rows: map[string]*repo.FaqResult{}}
	svc := &supportService{ticketRepo: tickets, commentRepo: comments, categoryRepo: categories, faqRepo: faqs}
	return svc, tickets, comments, categories, faqs
}

func seedTicket(tickets *stubTicketRepo) *repo.TicketResult {
	row := &repo.TicketResult{Id: 5, Uid: testTicketUid, CustomerUid: testCustomerUid, Subject: "s", Description: "d", Status: "OPEN", Priority: "LOW", CreatedAt: time.Now()}
	tickets.rows[row.Uid] = row
	return row
}

// ─── tickets ────────────────────────────────────────────────────────────────

func TestCreateTicket_WithoutCategory_OpensWithGivenPriority(t *testing.T) {
	svc, tickets, _, _, _ := newSupportFixture()
	got, err := svc.CreateTicket(context.Background(), command.TicketCreateCommand{
		CustomerUid: testCustomerUid, Subject: "s", Description: "d", Priority: constant.TICKET_PRIORITY_HIGH,
	})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if got.Status != constant.TICKET_STATUS_OPEN || got.Priority != constant.TICKET_PRIORITY_HIGH {
		t.Errorf("got status=%s priority=%s", got.Status, got.Priority)
	}
	if got.CategoryUid != nil || tickets.inserted[0].CategoryId != nil {
		t.Errorf("category must stay unset")
	}
	if got.UpdatedAt.IsZero() {
		t.Errorf("updated_at must fall back to created_at")
	}
}

func TestCreateTicket_WithCategoryUid_ResolvesCategoryId(t *testing.T) {
	svc, tickets, _, _, _ := newSupportFixture()
	got, err := svc.CreateTicket(context.Background(), command.TicketCreateCommand{
		CustomerUid: testCustomerUid, Subject: "s", Description: "d", Priority: constant.TICKET_PRIORITY_MEDIUM, CategoryUid: strp(testCategoryUid),
	})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if tickets.inserted[0].CategoryId == nil || *tickets.inserted[0].CategoryId != 10 {
		t.Errorf("CategoryId = %v, want 10", tickets.inserted[0].CategoryId)
	}
	if got.CategoryUid == nil || *got.CategoryUid != testCategoryUid {
		t.Errorf("CategoryUid = %v", got.CategoryUid)
	}
}

func TestCreateTicket_UnknownCategory_ReturnsErrUnknownReference(t *testing.T) {
	svc, tickets, _, _, _ := newSupportFixture()
	_, err := svc.CreateTicket(context.Background(), command.TicketCreateCommand{
		CustomerUid: testCustomerUid, Subject: "s", Description: "d", Priority: constant.TICKET_PRIORITY_MEDIUM, CategoryUid: strp(testParentUid[:35] + "9"),
	})
	if !errors.Is(err, repo.ErrUnknownReference) {
		t.Fatalf("want ErrUnknownReference, got %v", err)
	}
	if len(tickets.inserted) != 0 {
		t.Errorf("nothing must be inserted")
	}
}

func TestCreateTicket_InvalidCommand_ReturnsErrInvalidCommandWithoutInsert(t *testing.T) {
	svc, tickets, _, _, _ := newSupportFixture()
	_, err := svc.CreateTicket(context.Background(), command.TicketCreateCommand{CustomerUid: testCustomerUid, Subject: " ", Description: "d", Priority: constant.TICKET_PRIORITY_LOW})
	if !errors.Is(err, command.ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if len(tickets.inserted) != 0 {
		t.Errorf("nothing must be inserted")
	}
}

func TestFindTicket_Missing_ReturnsErrNotFound(t *testing.T) {
	svc, _, _, _, _ := newSupportFixture()
	if _, err := svc.FindTicket(context.Background(), testTicketUid); !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestUpdateTicket_PartialCommand_ChangesOnlyGivenFields(t *testing.T) {
	svc, tickets, _, _, _ := newSupportFixture()
	seedTicket(tickets)
	priority := constant.TICKET_PRIORITY_URGENT
	got, err := svc.UpdateTicket(context.Background(), testTicketUid, command.TicketUpdateCommand{Priority: &priority, CategoryUid: strp(testCategoryUid)})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if got.Subject != "s" || got.Description != "d" || got.Priority != constant.TICKET_PRIORITY_URGENT {
		t.Errorf("got %+v", got)
	}
	if tickets.rows[testTicketUid].CategoryId == nil || *tickets.rows[testTicketUid].CategoryId != 10 {
		t.Errorf("CategoryId = %v, want 10", tickets.rows[testTicketUid].CategoryId)
	}
}

func TestUpdateTicket_Missing_ReturnsErrNotFound(t *testing.T) {
	svc, _, _, _, _ := newSupportFixture()
	if _, err := svc.UpdateTicket(context.Background(), testTicketUid, command.TicketUpdateCommand{Subject: strp("x")}); !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestUpdateTicketStatus_Resolved_SetsResolvedAtOnce(t *testing.T) {
	svc, tickets, _, _, _ := newSupportFixture()
	seedTicket(tickets)
	first, err := svc.UpdateTicketStatus(context.Background(), testTicketUid, command.TicketStatusCommand{Status: constant.TICKET_STATUS_RESOLVED})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if first.ResolvedAt == nil {
		t.Fatalf("resolved_at must be set on RESOLVED")
	}
	stamp := *first.ResolvedAt
	time.Sleep(time.Millisecond)
	second, err := svc.UpdateTicketStatus(context.Background(), testTicketUid, command.TicketStatusCommand{Status: constant.TICKET_STATUS_RESOLVED})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if !second.ResolvedAt.Equal(stamp) {
		t.Errorf("resolved_at must keep the first stamp: %v vs %v", second.ResolvedAt, stamp)
	}
	reopened, err := svc.UpdateTicketStatus(context.Background(), testTicketUid, command.TicketStatusCommand{Status: constant.TICKET_STATUS_IN_PROGRESS})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if reopened.Status != constant.TICKET_STATUS_IN_PROGRESS || reopened.ResolvedAt == nil {
		t.Errorf("reopen must keep resolved_at: %+v", reopened)
	}
}

func TestUpdateTicketStatus_InProgressOnFreshTicket_LeavesResolvedAtNil(t *testing.T) {
	svc, tickets, _, _, _ := newSupportFixture()
	seedTicket(tickets)
	got, err := svc.UpdateTicketStatus(context.Background(), testTicketUid, command.TicketStatusCommand{Status: constant.TICKET_STATUS_IN_PROGRESS})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if got.ResolvedAt != nil {
		t.Errorf("resolved_at must stay nil")
	}
}

func TestUpdateTicketStatus_NonUtcLocalZone_StampsResolvedAtInUtc(t *testing.T) {
	local := time.Local
	time.Local = time.FixedZone("KST", 9*60*60)
	t.Cleanup(func() { time.Local = local })
	svc, tickets, _, _, _ := newSupportFixture()
	seedTicket(tickets)
	got, err := svc.UpdateTicketStatus(context.Background(), testTicketUid, command.TicketStatusCommand{Status: constant.TICKET_STATUS_RESOLVED})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if got.ResolvedAt == nil || got.ResolvedAt.Location() != time.UTC {
		t.Errorf("resolved_at must be UTC, got %v", got.ResolvedAt)
	}
}

// ─── comments ───────────────────────────────────────────────────────────────

func TestAddComment_TicketMissing_ReturnsErrNotFound(t *testing.T) {
	svc, _, comments, _, _ := newSupportFixture()
	_, err := svc.AddComment(context.Background(), command.CommentAddCommand{TicketUid: testTicketUid, AuthorUid: testCustomerUid, AuthorRole: constant.AUTHOR_ROLE_CUSTOMER, Body: "hi"})
	if !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
	if len(comments.inserted) != 0 {
		t.Errorf("nothing must be inserted")
	}
}

func TestAddComment_TicketPresent_UsesTicketIdAndUid(t *testing.T) {
	svc, tickets, comments, _, _ := newSupportFixture()
	seedTicket(tickets)
	got, err := svc.AddComment(context.Background(), command.CommentAddCommand{TicketUid: testTicketUid, AuthorUid: testCustomerUid, AuthorRole: constant.AUTHOR_ROLE_AGENT, Body: "hi", Internal: true})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if comments.inserted[0].TicketId != 5 || comments.inserted[0].AuthorRole != "AGENT" || !comments.inserted[0].Internal {
		t.Errorf("inserted = %+v", comments.inserted[0])
	}
	if got.TicketUid != testTicketUid || got.AuthorRole != constant.AUTHOR_ROLE_AGENT {
		t.Errorf("got = %+v", got)
	}
}

func TestFindComments_PassesIncludeInternalThrough(t *testing.T) {
	svc, tickets, comments, _, _ := newSupportFixture()
	seedTicket(tickets)
	if _, err := svc.FindComments(context.Background(), testTicketUid, false); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if comments.includeInternal == nil || *comments.includeInternal {
		t.Errorf("includeInternal must be false")
	}
	if _, err := svc.FindComments(context.Background(), "018f0000-0000-7000-8000-0000000000e9", true); !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("want ErrNotFound for missing ticket, got %v", err)
	}
}

// ─── categories ─────────────────────────────────────────────────────────────

func TestCreateCategory_DuplicateName_ReturnsErrCategoryNameExists(t *testing.T) {
	svc, _, _, categories, _ := newSupportFixture()
	_, err := svc.CreateCategory(context.Background(), command.CategoryCreateCommand{Name: "Billing"})
	if !errors.Is(err, repo.ErrCategoryNameExists) {
		t.Fatalf("want ErrCategoryNameExists, got %v", err)
	}
	if len(categories.inserted) != 0 {
		t.Errorf("nothing must be inserted")
	}
}

func TestCreateCategory_WithParent_ResolvesParentId(t *testing.T) {
	svc, _, _, categories, _ := newSupportFixture()
	got, err := svc.CreateCategory(context.Background(), command.CategoryCreateCommand{Name: "Refunds", ParentUid: strp(testParentUid)})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if categories.inserted[0].ParentId == nil || *categories.inserted[0].ParentId != 11 {
		t.Errorf("ParentId = %v, want 11", categories.inserted[0].ParentId)
	}
	if got.Name != "Refunds" {
		t.Errorf("got = %+v", got)
	}
}

func TestCreateCategory_UnknownParent_ReturnsErrUnknownReference(t *testing.T) {
	svc, _, _, _, _ := newSupportFixture()
	_, err := svc.CreateCategory(context.Background(), command.CategoryCreateCommand{Name: "Refunds", ParentUid: strp("018f0000-0000-7000-8000-0000000000d9")})
	if !errors.Is(err, repo.ErrUnknownReference) {
		t.Fatalf("want ErrUnknownReference, got %v", err)
	}
}

func TestUpdateCategory_RenameToTakenName_ReturnsErrCategoryNameExists(t *testing.T) {
	svc, _, _, _, _ := newSupportFixture()
	_, err := svc.UpdateCategory(context.Background(), testCategoryUid, command.CategoryUpdateCommand{Name: strp("Root")})
	if !errors.Is(err, repo.ErrCategoryNameExists) {
		t.Fatalf("want ErrCategoryNameExists, got %v", err)
	}
}

func TestUpdateCategory_SameNameAndDescription_Passes(t *testing.T) {
	svc, _, _, categories, _ := newSupportFixture()
	got, err := svc.UpdateCategory(context.Background(), testCategoryUid, command.CategoryUpdateCommand{Name: strp("Billing"), Description: strp("")})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if got.Description == nil || *got.Description != "" {
		t.Errorf("description must be the empty string, got %v", got.Description)
	}
	if categories.byUid[testCategoryUid].Name != "Billing" {
		t.Errorf("name changed unexpectedly")
	}
}

func TestUpdateCategory_Missing_ReturnsErrNotFound(t *testing.T) {
	svc, _, _, _, _ := newSupportFixture()
	if _, err := svc.UpdateCategory(context.Background(), "018f0000-0000-7000-8000-0000000000d9", command.CategoryUpdateCommand{}); !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestDeleteCategory_PresentOrMissing_ReportsAccordingly(t *testing.T) {
	svc, _, _, categories, _ := newSupportFixture()
	if err := svc.DeleteCategory(context.Background(), testCategoryUid); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if len(categories.deleted) != 1 {
		t.Errorf("deleted = %v", categories.deleted)
	}
	if err := svc.DeleteCategory(context.Background(), "018f0000-0000-7000-8000-0000000000d9"); !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

// ─── faqs ───────────────────────────────────────────────────────────────────

func TestCreateFaq_DefaultsToUnpublishedWithZeroViews(t *testing.T) {
	svc, _, _, _, _ := newSupportFixture()
	got, err := svc.CreateFaq(context.Background(), command.FaqCreateCommand{Question: "q", Answer: "a"})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if got.Published || got.ViewCount != 0 {
		t.Errorf("got = %+v", got)
	}
}

func TestCreateFaq_UnknownCategory_ReturnsErrUnknownReference(t *testing.T) {
	svc, _, _, _, _ := newSupportFixture()
	_, err := svc.CreateFaq(context.Background(), command.FaqCreateCommand{Question: "q", Answer: "a", CategoryUid: strp("018f0000-0000-7000-8000-0000000000d9")})
	if !errors.Is(err, repo.ErrUnknownReference) {
		t.Fatalf("want ErrUnknownReference, got %v", err)
	}
}

func TestFindFaq_PublishedAndIncrementRequested_IncrementsViewCount(t *testing.T) {
	svc, _, _, _, faqs := newSupportFixture()
	faqs.rows[testFaqUid] = &repo.FaqResult{Id: 1, Uid: testFaqUid, Question: "q", Answer: "a", ViewCount: 3, Published: true, CreatedAt: time.Now()}
	got, err := svc.FindFaq(context.Background(), testFaqUid, true)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if got.ViewCount != 4 || len(faqs.increments) != 1 {
		t.Errorf("view count = %d, increments = %v", got.ViewCount, faqs.increments)
	}
}

func TestFindFaq_UnpublishedOrIncrementDisabled_LeavesViewCount(t *testing.T) {
	svc, _, _, _, faqs := newSupportFixture()
	faqs.rows[testFaqUid] = &repo.FaqResult{Id: 1, Uid: testFaqUid, Question: "q", Answer: "a", ViewCount: 3, Published: false, CreatedAt: time.Now()}
	if _, err := svc.FindFaq(context.Background(), testFaqUid, true); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	faqs.rows[testFaqUid].Published = true
	if _, err := svc.FindFaq(context.Background(), testFaqUid, false); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if len(faqs.increments) != 0 {
		t.Errorf("increments = %v, want none", faqs.increments)
	}
}

func TestFindFaq_IncrementFails_StillReturnsFaq(t *testing.T) {
	svc, _, _, _, faqs := newSupportFixture()
	faqs.rows[testFaqUid] = &repo.FaqResult{Id: 1, Uid: testFaqUid, Question: "q", Answer: "a", ViewCount: 3, Published: true, CreatedAt: time.Now()}
	faqs.incrementErr = errors.New("db down")
	got, err := svc.FindFaq(context.Background(), testFaqUid, true)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if got.ViewCount != 3 {
		t.Errorf("view count = %d, want 3", got.ViewCount)
	}
}

func TestFindFaq_Missing_ReturnsErrNotFound(t *testing.T) {
	svc, _, _, _, _ := newSupportFixture()
	if _, err := svc.FindFaq(context.Background(), testFaqUid, true); !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestUpdateFaq_PartialCommand_ChangesOnlyGivenFields(t *testing.T) {
	svc, _, _, _, faqs := newSupportFixture()
	faqs.rows[testFaqUid] = &repo.FaqResult{Id: 1, Uid: testFaqUid, Question: "q", Answer: "a", Published: false, CreatedAt: time.Now()}
	published := true
	got, err := svc.UpdateFaq(context.Background(), testFaqUid, command.FaqUpdateCommand{Published: &published, Answer: strp("b")})
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if got.Question != "q" || got.Answer != "b" || !got.Published {
		t.Errorf("got = %+v", got)
	}
}

func TestDeleteFaq_Missing_ReturnsErrNotFound(t *testing.T) {
	svc, _, _, _, _ := newSupportFixture()
	if err := svc.DeleteFaq(context.Background(), testFaqUid); !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestFindAllTickets_PassesOptionThrough(t *testing.T) {
	svc, tickets, _, _, _ := newSupportFixture()
	seedTicket(tickets)
	opt := option.NewTicketQueryOption(option.WithTicketCustomerUid(testCustomerUid), option.WithTicketStatus(constant.TICKET_STATUS_OPEN))
	total, out, err := svc.FindAllTickets(context.Background(), opt)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if total != 1 || len(out) != 1 {
		t.Errorf("total=%d len=%d", total, len(out))
	}
	if tickets.listOpt == nil || tickets.listOpt.CustomerUid() == nil || *tickets.listOpt.CustomerUid() != testCustomerUid {
		t.Errorf("option not passed through: %+v", tickets.listOpt)
	}
}
