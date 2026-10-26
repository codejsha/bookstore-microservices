package restcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

const (
	ownerUid    = "018f0000-0000-7000-8000-0000000000c1"
	otherUid    = "018f0000-0000-7000-8000-0000000000c2"
	staffUid    = "018f0000-0000-7000-8000-0000000000c3"
	ticketUid   = "018f0000-0000-7000-8000-0000000000e1"
	missingUid  = "018f0000-0000-7000-8000-0000000000e9"
	categoryUid = "018f0000-0000-7000-8000-0000000000d1"
	faqUid      = "018f0000-0000-7000-8000-0000000000f1"
)

type stubSupportUseCase struct {
	usecase.SupportUseCase
	ticket          *aggregate.TicketAggregate
	faq             *aggregate.FaqAggregate
	createdTicket   *command.TicketCreateCommand
	updatedTicket   *command.TicketUpdateCommand
	statusCmd       *command.TicketStatusCommand
	addedComment    *command.CommentAddCommand
	listOpt         *option.TicketQueryOption
	faqListOpt      *option.FaqQueryOption
	includeInternal *bool
	createErr       error
}

func (s *stubSupportUseCase) FindTicket(_ context.Context, uid string) (*aggregate.TicketAggregate, error) {
	if s.ticket == nil || s.ticket.Uid != uid {
		return nil, fmt.Errorf("find ticket %s: %w", uid, usecase.ErrNotFound)
	}
	return s.ticket, nil
}

func (s *stubSupportUseCase) FindAllTickets(_ context.Context, opt option.TicketQueryOption) (int64, []*aggregate.TicketAggregate, error) {
	s.listOpt = &opt
	return 0, []*aggregate.TicketAggregate{}, nil
}

func (s *stubSupportUseCase) CreateTicket(_ context.Context, cmd command.TicketCreateCommand) (*aggregate.TicketAggregate, error) {
	s.createdTicket = &cmd
	if s.createErr != nil {
		return nil, s.createErr
	}
	return &aggregate.TicketAggregate{Uid: ticketUid, CustomerUid: cmd.CustomerUid}, nil
}

func (s *stubSupportUseCase) UpdateTicket(_ context.Context, _ string, cmd command.TicketUpdateCommand) (*aggregate.TicketAggregate, error) {
	s.updatedTicket = &cmd
	return s.ticket, nil
}

func (s *stubSupportUseCase) UpdateTicketStatus(_ context.Context, uid string, cmd command.TicketStatusCommand) (*aggregate.TicketAggregate, error) {
	s.statusCmd = &cmd
	if s.ticket == nil || s.ticket.Uid != uid {
		return nil, fmt.Errorf("update ticket status %s: %w", uid, usecase.ErrNotFound)
	}
	return s.ticket, nil
}

func (s *stubSupportUseCase) AddComment(_ context.Context, cmd command.CommentAddCommand) (*aggregate.TicketCommentAggregate, error) {
	s.addedComment = &cmd
	return &aggregate.TicketCommentAggregate{Uid: "018f0000-0000-7000-8000-0000000000a9", TicketUid: cmd.TicketUid}, nil
}

func (s *stubSupportUseCase) FindComments(_ context.Context, _ string, includeInternal bool) ([]*aggregate.TicketCommentAggregate, error) {
	s.includeInternal = &includeInternal
	return []*aggregate.TicketCommentAggregate{}, nil
}

func (s *stubSupportUseCase) CreateCategory(_ context.Context, cmd command.CategoryCreateCommand) (*aggregate.TicketCategoryAggregate, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}
	return &aggregate.TicketCategoryAggregate{Uid: categoryUid, Name: cmd.Name}, nil
}

func (s *stubSupportUseCase) FindAllCategories(_ context.Context) ([]*aggregate.TicketCategoryAggregate, error) {
	return []*aggregate.TicketCategoryAggregate{{Uid: categoryUid, Name: "Billing", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
}

func (s *stubSupportUseCase) FindCategory(_ context.Context, uid string) (*aggregate.TicketCategoryAggregate, error) {
	if uid != categoryUid {
		return nil, fmt.Errorf("find category %s: %w", uid, usecase.ErrNotFound)
	}
	return &aggregate.TicketCategoryAggregate{Uid: categoryUid, Name: "Billing"}, nil
}

func (s *stubSupportUseCase) DeleteCategory(_ context.Context, uid string) error {
	if uid != categoryUid {
		return fmt.Errorf("delete category %s: %w", uid, usecase.ErrNotFound)
	}
	return nil
}

func (s *stubSupportUseCase) CreateFaq(_ context.Context, cmd command.FaqCreateCommand) (*aggregate.FaqAggregate, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}
	return &aggregate.FaqAggregate{Uid: faqUid, Question: cmd.Question, Published: cmd.Published}, nil
}

func (s *stubSupportUseCase) FindFaq(_ context.Context, uid string, _ bool) (*aggregate.FaqAggregate, error) {
	if s.faq == nil || s.faq.Uid != uid {
		return nil, fmt.Errorf("find faq %s: %w", uid, usecase.ErrNotFound)
	}
	return s.faq, nil
}

func (s *stubSupportUseCase) FindAllFaqs(_ context.Context, opt option.FaqQueryOption) (int64, []*aggregate.FaqAggregate, error) {
	s.faqListOpt = &opt
	return 0, []*aggregate.FaqAggregate{}, nil
}

func (s *stubSupportUseCase) DeleteFaq(_ context.Context, uid string) error {
	if uid != faqUid {
		return fmt.Errorf("delete faq %s: %w", uid, usecase.ErrNotFound)
	}
	return nil
}

// testCallerMiddleware mirrors what the principal middleware in the support package
// derives from the mesh headers; the package itself cannot be imported here without
// an import cycle, and the manager-only route gate it adds is covered by its own tests.
func testCallerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := c.GetHeader("X-User-Id")
		if sub == "" {
			c.Next()
			return
		}
		roles := strings.Split(c.GetHeader("X-User-Roles"), ",")
		isAdmin := false
		for _, r := range roles {
			if r == "STAFF" || r == "SYSTEM" {
				isAdmin = true
			}
		}
		c.Request = c.Request.WithContext(httpx.WithCaller(c.Request.Context(), httpx.Caller{UserUid: sub, IsAdmin: isAdmin}))
		c.Next()
	}
}

func newSupportRouter(uc usecase.SupportUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(httpx.GinResponseMapping())
	r.Use(testCallerMiddleware())
	openapi.RegisterTicketApiGinRoutes(r, NewTicketController(uc))
	openapi.RegisterCommentApiGinRoutes(r, NewCommentController(uc))
	openapi.RegisterCategoryApiGinRoutes(r, NewCategoryController(uc))
	openapi.RegisterFaqApiGinRoutes(r, NewFaqController(uc))
	return r
}

type principal struct {
	sub   string
	roles string
}

var (
	anonymous = principal{}
	owner     = principal{sub: ownerUid, roles: "USER"}
	stranger  = principal{sub: otherUid, roles: "USER"}
	staff     = principal{sub: staffUid, roles: "STAFF,USER"}
	manager   = principal{sub: staffUid, roles: "MANAGER,USER"}
	system    = principal{sub: staffUid, roles: "SYSTEM"}
	brokenSub = principal{sub: "7", roles: "USER"}
)

func call(t *testing.T, r *gin.Engine, p principal, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *strings.Reader
	if body != "" {
		reader = strings.NewReader(body)
	} else {
		reader = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Authorization", "Bearer test-token")
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if p.sub != "" {
		req.Header.Set("X-User-Id", p.sub)
		req.Header.Set("X-User-Roles", p.roles)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ownedTicket() *aggregate.TicketAggregate {
	now := time.Now()
	return &aggregate.TicketAggregate{
		Uid: ticketUid, CustomerUid: ownerUid, Subject: "s", Description: "d",
		Status: constant.TICKET_STATUS_OPEN, Priority: constant.TICKET_PRIORITY_LOW, CreatedAt: now, UpdatedAt: now,
	}
}

func problemDetail(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct {
		Detail string `json:"detail"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("body %q is not problem json: %v", w.Body.String(), err)
	}
	return body.Detail
}

// ─── tickets: authorization matrix ──────────────────────────────────────────

func TestTickets_Anonymous_Returns401(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{ticket: ownedTicket()})
	if w := call(t, r, anonymous, http.MethodGet, "/api/v1/tickets", ""); w.Code != http.StatusUnauthorized {
		t.Errorf("list = %d, want 401", w.Code)
	}
	if w := call(t, r, anonymous, http.MethodGet, "/api/v1/tickets/"+ticketUid, ""); w.Code != http.StatusUnauthorized {
		t.Errorf("read = %d, want 401", w.Code)
	}
}

func TestTicketsGetAll_NonUuidSubject_Returns403(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{})
	if w := call(t, r, brokenSub, http.MethodGet, "/api/v1/tickets", ""); w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403", w.Code)
	}
}

func TestTicketsGetAll_Customer_PinsCustomerFilterToCaller(t *testing.T) {
	uc := &stubSupportUseCase{}
	r := newSupportRouter(uc)
	w := call(t, r, owner, http.MethodGet, "/api/v1/tickets?customer_uid="+otherUid+"&status=OPEN", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.listOpt.CustomerUid() == nil || *uc.listOpt.CustomerUid() != ownerUid {
		t.Errorf("customer filter = %v, want caller", uc.listOpt.CustomerUid())
	}
	if uc.listOpt.Status() == nil || *uc.listOpt.Status() != constant.TICKET_STATUS_OPEN {
		t.Errorf("status filter = %v, want OPEN", uc.listOpt.Status())
	}
}

func TestTicketsGetAll_Staff_KeepsRequestedFilters(t *testing.T) {
	uc := &stubSupportUseCase{}
	r := newSupportRouter(uc)
	w := call(t, r, staff, http.MethodGet, "/api/v1/tickets?customer_uid="+otherUid+"&priority=HIGH&size=5&page=1", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.listOpt.CustomerUid() == nil || *uc.listOpt.CustomerUid() != otherUid {
		t.Errorf("customer filter = %v, want requested", uc.listOpt.CustomerUid())
	}
	if offset, limit := uc.listOpt.Page().OffsetLimit(); offset != 5 || limit != 5 {
		t.Errorf("page = (%d, %d), want (5, 5)", offset, limit)
	}
}

func TestTicketsCreate_Customer_OverridesCustomerUidWithCaller(t *testing.T) {
	uc := &stubSupportUseCase{}
	r := newSupportRouter(uc)
	body := `{"customer_uid":"` + otherUid + `","subject":"Help","description":"Please"}`
	w := call(t, r, owner, http.MethodPost, "/api/v1/tickets", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", w.Code, w.Body.String())
	}
	if uc.createdTicket.CustomerUid != ownerUid {
		t.Errorf("customer uid = %s, want caller", uc.createdTicket.CustomerUid)
	}
	if uc.createdTicket.Priority != constant.TICKET_PRIORITY_MEDIUM {
		t.Errorf("priority = %s, want MEDIUM default", uc.createdTicket.Priority)
	}
	if loc := w.Header().Get("Location"); loc != "/api/v1/tickets/"+ticketUid {
		t.Errorf("Location = %q", loc)
	}
}

func TestTicketsCreate_Staff_KeepsBodyCustomerUid(t *testing.T) {
	uc := &stubSupportUseCase{}
	r := newSupportRouter(uc)
	body := `{"customer_uid":"` + otherUid + `","subject":"Help","description":"Please","priority":"URGENT"}`
	if w := call(t, r, staff, http.MethodPost, "/api/v1/tickets", body); w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", w.Code, w.Body.String())
	}
	if uc.createdTicket.CustomerUid != otherUid || uc.createdTicket.Priority != constant.TICKET_PRIORITY_URGENT {
		t.Errorf("cmd = %+v", uc.createdTicket)
	}
}

func TestTicketsCreate_SchemaViolations_Return400(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{})
	cases := map[string]string{
		"missingSubject":  `{"customer_uid":"` + ownerUid + `","description":"Please"}`,
		"overlongSubject": `{"customer_uid":"` + ownerUid + `","subject":"` + strings.Repeat("x", 256) + `","description":"Please"}`,
		"badPriority":     `{"customer_uid":"` + ownerUid + `","subject":"s","description":"d","priority":"SOMEDAY"}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			if w := call(t, r, staff, http.MethodPost, "/api/v1/tickets", body); w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestTicketsCreate_UnknownCategory_Returns400ProblemDetails(t *testing.T) {
	uc := &stubSupportUseCase{createErr: fmt.Errorf("category x: %w", repo.ErrUnknownReference)}
	r := newSupportRouter(uc)
	body := `{"customer_uid":"` + ownerUid + `","subject":"s","description":"d","category_uid":"` + categoryUid + `"}`
	w := call(t, r, owner, http.MethodPost, "/api/v1/tickets", body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400: %s", w.Code, w.Body.String())
	}
	if got := problemDetail(t, w); got != "referenced category not found" {
		t.Errorf("detail = %q", got)
	}
}

func TestTicketsRead_OwnerStaffStranger_200_200_403(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{ticket: ownedTicket()})
	if w := call(t, r, owner, http.MethodGet, "/api/v1/tickets/"+ticketUid, ""); w.Code != http.StatusOK {
		t.Errorf("owner = %d, want 200: %s", w.Code, w.Body.String())
	}
	if w := call(t, r, staff, http.MethodGet, "/api/v1/tickets/"+ticketUid, ""); w.Code != http.StatusOK {
		t.Errorf("staff = %d, want 200", w.Code)
	}
	w := call(t, r, stranger, http.MethodGet, "/api/v1/tickets/"+ticketUid, "")
	if w.Code != http.StatusForbidden {
		t.Errorf("stranger = %d, want 403", w.Code)
	}
	if got := problemDetail(t, w); got != "forbidden" {
		t.Errorf("detail = %q", got)
	}
}

func TestTicketsRead_Missing_Returns404BeforeOwnership(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{ticket: ownedTicket()})
	w := call(t, r, stranger, http.MethodGet, "/api/v1/tickets/"+missingUid, "")
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404: %s", w.Code, w.Body.String())
	}
}

func TestTicketsUpdate_CustomerSettingAssignee_Returns403(t *testing.T) {
	uc := &stubSupportUseCase{ticket: ownedTicket()}
	r := newSupportRouter(uc)
	w := call(t, r, owner, http.MethodPatch, "/api/v1/tickets/"+ticketUid, `{"assignee_uid":"`+staffUid+`"}`)
	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403: %s", w.Code, w.Body.String())
	}
	if uc.updatedTicket != nil {
		t.Errorf("use case must not be called")
	}
}

func TestTicketsUpdate_OwnerChangingSubject_Returns200(t *testing.T) {
	uc := &stubSupportUseCase{ticket: ownedTicket()}
	r := newSupportRouter(uc)
	w := call(t, r, owner, http.MethodPatch, "/api/v1/tickets/"+ticketUid, `{"subject":"New","priority":"HIGH"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.updatedTicket.Subject == nil || *uc.updatedTicket.Subject != "New" || uc.updatedTicket.Priority == nil || *uc.updatedTicket.Priority != constant.TICKET_PRIORITY_HIGH {
		t.Errorf("cmd = %+v", uc.updatedTicket)
	}
}

func TestTicketsUpdate_StrangerOrMissing_403Or404(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{ticket: ownedTicket()})
	if w := call(t, r, stranger, http.MethodPatch, "/api/v1/tickets/"+ticketUid, `{"subject":"New"}`); w.Code != http.StatusForbidden {
		t.Errorf("stranger = %d, want 403", w.Code)
	}
	if w := call(t, r, owner, http.MethodPatch, "/api/v1/tickets/"+missingUid, `{"subject":"New"}`); w.Code != http.StatusNotFound {
		t.Errorf("missing = %d, want 404", w.Code)
	}
}

func TestTicketsUpdateStatus_RoleMatrix(t *testing.T) {
	uc := &stubSupportUseCase{ticket: ownedTicket()}
	r := newSupportRouter(uc)
	body := `{"status":"RESOLVED"}`
	if w := call(t, r, owner, http.MethodPatch, "/api/v1/tickets/"+ticketUid+"/status", body); w.Code != http.StatusForbidden {
		t.Errorf("customer = %d, want 403", w.Code)
	}
	if uc.statusCmd != nil {
		t.Errorf("use case must not be called for a customer")
	}
	if w := call(t, r, staff, http.MethodPatch, "/api/v1/tickets/"+ticketUid+"/status", body); w.Code != http.StatusOK {
		t.Errorf("staff = %d, want 200: %s", w.Code, w.Body.String())
	}
	if w := call(t, r, system, http.MethodPatch, "/api/v1/tickets/"+ticketUid+"/status", body); w.Code != http.StatusOK {
		t.Errorf("system = %d, want 200", w.Code)
	}
	if w := call(t, r, staff, http.MethodPatch, "/api/v1/tickets/"+missingUid+"/status", body); w.Code != http.StatusNotFound {
		t.Errorf("missing = %d, want 404", w.Code)
	}
	if w := call(t, r, staff, http.MethodPatch, "/api/v1/tickets/"+ticketUid+"/status", `{"status":"DONE"}`); w.Code != http.StatusBadRequest {
		t.Errorf("bad status = %d, want 400", w.Code)
	}
}

// ─── comments ───────────────────────────────────────────────────────────────

func TestCommentsCreate_Customer_AuthoredAsCustomerFromCaller(t *testing.T) {
	uc := &stubSupportUseCase{ticket: ownedTicket()}
	r := newSupportRouter(uc)
	w := call(t, r, owner, http.MethodPost, "/api/v1/tickets/"+ticketUid+"/comments", `{"body":"hello","author_uid":"`+otherUid+`","author_role":"AGENT"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", w.Code, w.Body.String())
	}
	if uc.addedComment.AuthorUid != ownerUid || uc.addedComment.AuthorRole != constant.AUTHOR_ROLE_CUSTOMER || uc.addedComment.Internal {
		t.Errorf("cmd = %+v", uc.addedComment)
	}
	if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/api/v1/tickets/"+ticketUid+"/comments/") {
		t.Errorf("Location = %q", loc)
	}
}

func TestCommentsCreate_Staff_AuthoredAsAgentAndMayBeInternal(t *testing.T) {
	uc := &stubSupportUseCase{ticket: ownedTicket()}
	r := newSupportRouter(uc)
	if w := call(t, r, staff, http.MethodPost, "/api/v1/tickets/"+ticketUid+"/comments", `{"body":"note","internal":true}`); w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", w.Code, w.Body.String())
	}
	if uc.addedComment.AuthorUid != staffUid || uc.addedComment.AuthorRole != constant.AUTHOR_ROLE_AGENT || !uc.addedComment.Internal {
		t.Errorf("cmd = %+v", uc.addedComment)
	}
}

func TestCommentsCreate_CustomerInternalOrStranger_403(t *testing.T) {
	uc := &stubSupportUseCase{ticket: ownedTicket()}
	r := newSupportRouter(uc)
	if w := call(t, r, owner, http.MethodPost, "/api/v1/tickets/"+ticketUid+"/comments", `{"body":"note","internal":true}`); w.Code != http.StatusForbidden {
		t.Errorf("internal by customer = %d, want 403", w.Code)
	}
	if w := call(t, r, stranger, http.MethodPost, "/api/v1/tickets/"+ticketUid+"/comments", `{"body":"note"}`); w.Code != http.StatusForbidden {
		t.Errorf("stranger = %d, want 403", w.Code)
	}
	if uc.addedComment != nil {
		t.Errorf("use case must not be called")
	}
}

func TestCommentsCreate_MissingTicket_404AndBlankBody_400(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{ticket: ownedTicket()})
	if w := call(t, r, owner, http.MethodPost, "/api/v1/tickets/"+missingUid+"/comments", `{"body":"note"}`); w.Code != http.StatusNotFound {
		t.Errorf("missing = %d, want 404", w.Code)
	}
	if w := call(t, r, owner, http.MethodPost, "/api/v1/tickets/"+ticketUid+"/comments", `{"body":""}`); w.Code != http.StatusBadRequest {
		t.Errorf("empty body = %d, want 400", w.Code)
	}
}

func TestCommentsGetAll_IncludeInternalHonouredForStaffOnly(t *testing.T) {
	uc := &stubSupportUseCase{ticket: ownedTicket()}
	r := newSupportRouter(uc)
	if w := call(t, r, owner, http.MethodGet, "/api/v1/tickets/"+ticketUid+"/comments?include_internal=true", ""); w.Code != http.StatusOK {
		t.Fatalf("owner = %d, want 200: %s", w.Code, w.Body.String())
	}
	if uc.includeInternal == nil || *uc.includeInternal {
		t.Errorf("customer must never receive internal comments")
	}
	if w := call(t, r, staff, http.MethodGet, "/api/v1/tickets/"+ticketUid+"/comments?include_internal=true", ""); w.Code != http.StatusOK {
		t.Fatalf("staff = %d, want 200", w.Code)
	}
	if uc.includeInternal == nil || !*uc.includeInternal {
		t.Errorf("staff asking for internal comments must get them")
	}
	if w := call(t, r, staff, http.MethodGet, "/api/v1/tickets/"+ticketUid+"/comments", ""); w.Code != http.StatusOK {
		t.Fatalf("staff default = %d, want 200", w.Code)
	}
	if *uc.includeInternal {
		t.Errorf("include_internal must default to false")
	}
	if w := call(t, r, stranger, http.MethodGet, "/api/v1/tickets/"+ticketUid+"/comments", ""); w.Code != http.StatusForbidden {
		t.Errorf("stranger = %d, want 403", w.Code)
	}
}

// ─── categories ─────────────────────────────────────────────────────────────

func TestCategories_ReadsAndWrites_MapStatusesAndLocation(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{})
	if w := call(t, r, owner, http.MethodGet, "/api/v1/categories", ""); w.Code != http.StatusOK {
		t.Errorf("list as customer = %d, want 200: %s", w.Code, w.Body.String())
	}
	if w := call(t, r, owner, http.MethodGet, "/api/v1/categories/"+categoryUid, ""); w.Code != http.StatusOK {
		t.Errorf("read as customer = %d, want 200", w.Code)
	}
	if w := call(t, r, owner, http.MethodGet, "/api/v1/categories/"+missingUid, ""); w.Code != http.StatusNotFound {
		t.Errorf("read missing = %d, want 404", w.Code)
	}
	body := `{"name":"Billing"}`
	if w := call(t, r, system, http.MethodPost, "/api/v1/categories", body); w.Code != http.StatusCreated {
		t.Errorf("create as SYSTEM = %d, want 201: %s", w.Code, w.Body.String())
	}
	w := call(t, r, manager, http.MethodPost, "/api/v1/categories", body)
	if w.Code != http.StatusCreated {
		t.Errorf("create as MANAGER = %d, want 201", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/api/v1/categories/"+categoryUid {
		t.Errorf("Location = %q", loc)
	}
	if w := call(t, r, manager, http.MethodDelete, "/api/v1/categories/"+categoryUid, ""); w.Code != http.StatusNoContent {
		t.Errorf("delete = %d, want 204", w.Code)
	}
	if w := call(t, r, manager, http.MethodDelete, "/api/v1/categories/"+missingUid, ""); w.Code != http.StatusNotFound {
		t.Errorf("delete missing = %d, want 404", w.Code)
	}
}

func TestCategoriesCreate_DuplicateName_Returns409(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{createErr: fmt.Errorf("create category: %w", repo.ErrCategoryNameExists)})
	w := call(t, r, manager, http.MethodPost, "/api/v1/categories", `{"name":"Billing"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409: %s", w.Code, w.Body.String())
	}
	if got := problemDetail(t, w); got != "category name already exists" {
		t.Errorf("detail = %q", got)
	}
}

func TestCategoriesCreate_OverlongName_Returns400(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{})
	if w := call(t, r, manager, http.MethodPost, "/api/v1/categories", `{"name":"`+strings.Repeat("n", 101)+`"}`); w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

// ─── faqs ───────────────────────────────────────────────────────────────────

func TestFaqsGetAll_CustomerForcedToPublishedOnly(t *testing.T) {
	uc := &stubSupportUseCase{}
	r := newSupportRouter(uc)
	if w := call(t, r, owner, http.MethodGet, "/api/v1/faqs?published_only=false", ""); w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	if !uc.faqListOpt.PublishedOnly() {
		t.Errorf("customer must only see published faqs")
	}
	if w := call(t, r, staff, http.MethodGet, "/api/v1/faqs?published_only=false&q=refund", ""); w.Code != http.StatusOK {
		t.Fatalf("staff = %d, want 200", w.Code)
	}
	if uc.faqListOpt.PublishedOnly() || uc.faqListOpt.Query() == nil || *uc.faqListOpt.Query() != "refund" {
		t.Errorf("staff option = %+v", uc.faqListOpt)
	}
	if w := call(t, r, staff, http.MethodGet, "/api/v1/faqs", ""); w.Code != http.StatusOK {
		t.Fatalf("staff default = %d, want 200", w.Code)
	}
	if !uc.faqListOpt.PublishedOnly() {
		t.Errorf("published_only must default to true")
	}
}

func TestFaqsRead_UnpublishedHiddenFromCustomers(t *testing.T) {
	draft := &aggregate.FaqAggregate{Uid: faqUid, Question: "q", Answer: "a", Published: false}
	r := newSupportRouter(&stubSupportUseCase{faq: draft})
	if w := call(t, r, owner, http.MethodGet, "/api/v1/faqs/"+faqUid, ""); w.Code != http.StatusNotFound {
		t.Errorf("customer = %d, want 404", w.Code)
	}
	if w := call(t, r, staff, http.MethodGet, "/api/v1/faqs/"+faqUid, ""); w.Code != http.StatusOK {
		t.Errorf("staff = %d, want 200: %s", w.Code, w.Body.String())
	}
	if w := call(t, r, staff, http.MethodGet, "/api/v1/faqs/"+missingUid, ""); w.Code != http.StatusNotFound {
		t.Errorf("missing = %d, want 404", w.Code)
	}
}

func TestFaqs_Writes_MapStatusesAndLocation(t *testing.T) {
	r := newSupportRouter(&stubSupportUseCase{})
	body := `{"question":"Q?","answer":"A.","published":true}`
	w := call(t, r, system, http.MethodPost, "/api/v1/faqs", body)
	if w.Code != http.StatusCreated {
		t.Errorf("create as SYSTEM = %d, want 201: %s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); loc != "/api/v1/faqs/"+faqUid {
		t.Errorf("Location = %q", loc)
	}
	if w := call(t, r, manager, http.MethodDelete, "/api/v1/faqs/"+faqUid, ""); w.Code != http.StatusNoContent {
		t.Errorf("delete = %d, want 204", w.Code)
	}
	if w := call(t, r, manager, http.MethodDelete, "/api/v1/faqs/"+missingUid, ""); w.Code != http.StatusNotFound {
		t.Errorf("delete missing = %d, want 404", w.Code)
	}
}
