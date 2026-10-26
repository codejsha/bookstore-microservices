package pgsql

import (
	"reflect"
	"testing"
	"time"

	"gorm.io/gen/field"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
)

func strp(s string) *string { return &s }

func testSortColumns() (map[string]sortColumn, field.Expr) {
	created := field.NewTime("t", "created_at")
	subject := field.NewString("t", "subject")
	return map[string]sortColumn{
		"created_at": {asc: created, desc: created.Desc()},
		"subject":    {asc: subject, desc: subject.Desc()},
	}, created.Desc()
}

func TestOrderExprs_UnknownOrEmptySort_FallsBackToDefault(t *testing.T) {
	columns, fallback := testSortColumns()
	for name, sorts := range map[string][]*string{
		"empty":   {},
		"nil":     {nil},
		"unknown": {strp("status,desc")},
	} {
		t.Run(name, func(t *testing.T) {
			got := orderExprs(sorts, columns, fallback)
			if len(got) != 1 || !reflect.DeepEqual(got[0], fallback) {
				t.Errorf("got %v, want the fallback only", got)
			}
		})
	}
}

func TestOrderExprs_KnownSorts_KeepOrderAndDirection(t *testing.T) {
	columns, fallback := testSortColumns()
	got := orderExprs([]*string{strp("subject,DESC"), strp(" created_at , asc ")}, columns, fallback)
	if len(got) != 2 {
		t.Fatalf("got %d exprs, want 2", len(got))
	}
	if !reflect.DeepEqual(got[0], columns["subject"].desc) {
		t.Errorf("first expr must be subject desc")
	}
	if !reflect.DeepEqual(got[1], columns["created_at"].asc) {
		t.Errorf("second expr must be created_at asc")
	}
}

func TestCollectCategoryIds_NilAndDuplicates_Dropped(t *testing.T) {
	a, b := int64(1), int64(2)
	got := collectCategoryIds([]*int64{&a, nil, &b, &a})
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Errorf("got %v, want [1 2]", got)
	}
}

func TestLookupCategoryUid_MissingIdOrNil_ReturnsNil(t *testing.T) {
	uids := map[int64]string{1: "cat-1"}
	one := int64(1)
	two := int64(2)
	if got := lookupCategoryUid(uids, &one); got == nil || *got != "cat-1" {
		t.Errorf("got %v, want cat-1", got)
	}
	if got := lookupCategoryUid(uids, &two); got != nil {
		t.Errorf("got %v, want nil for unknown id", got)
	}
	if got := lookupCategoryUid(uids, nil); got != nil {
		t.Errorf("got %v, want nil for nil id", got)
	}
}

func TestToTicketResult_WhenRowGiven_CarriesCategoryUidAndTimestamps(t *testing.T) {
	now := time.Now()
	categoryId := int64(7)
	got := toTicketResult(&entity.CustomerTicketEntity{
		Id: 3, Uid: "t-1", CustomerUid: "c-1", CategoryId: &categoryId, Subject: "s", Description: "d",
		Status: "OPEN", Priority: "LOW", CreatedAt: now, UpdatedAt: &now,
	}, strp("cat-7"))
	if got.Id != 3 || got.Status != "OPEN" || got.CategoryUid == nil || *got.CategoryUid != "cat-7" {
		t.Errorf("got = %+v", got)
	}
	if got.UpdatedAt == nil || !got.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v", got.UpdatedAt)
	}
}

func TestToFaqResult_WhenRowGiven_ReturnsResult(t *testing.T) {
	now := time.Now()
	got := toFaqResult(&entity.CustomerFaqEntity{Id: 2, Uid: "f-1", Question: "q", Answer: "a", ViewCount: 9, Published: true, CreatedAt: now}, nil)
	if got.Id != 2 || got.ViewCount != 9 || !got.Published || got.CategoryUid != nil {
		t.Errorf("got = %+v", got)
	}
}

func TestToTicketCategoryResult_WhenRowGiven_ReturnsResult(t *testing.T) {
	now := time.Now()
	parentId := int64(1)
	got := toTicketCategoryResult(&entity.CustomerTicketCategoryEntity{Id: 2, Uid: "cat-2", Name: "Billing", ParentId: &parentId, CreatedAt: now}, strp("cat-1"))
	if got.Name != "Billing" || got.ParentUid == nil || *got.ParentUid != "cat-1" {
		t.Errorf("got = %+v", got)
	}
}

func TestToTicketCommentResult_WhenRowGiven_ReturnsResult(t *testing.T) {
	now := time.Now()
	got := toTicketCommentResult(&entity.CustomerTicketCommentEntity{Id: 4, Uid: "cm-1", TicketId: 3, AuthorUid: "a-1", AuthorRole: "AGENT", Body: "b", Internal: true, CreatedAt: now})
	if got.TicketId != 3 || got.AuthorRole != "AGENT" || !got.Internal {
		t.Errorf("got = %+v", got)
	}
}
