package pgsql

import (
	"errors"
	"reflect"
	"testing"

	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
)

const refundsCategoryUid = "018f0000-0000-7000-8000-0000000000d3"

func newTicketCategoryRepo(t *testing.T) (*ticketCategoryRepository, *gorm.DB) {
	t.Helper()
	db := newSchemaDB(t)
	return &ticketCategoryRepository{q: queryOf(db), db: db}, db
}

func categoryNames(rows []*repo.TicketCategoryResult) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.Name
	}
	return out
}

func loadCategory(t *testing.T, db *gorm.DB, uid string) entity.CustomerTicketCategoryEntity {
	t.Helper()
	var got entity.CustomerTicketCategoryEntity
	if err := db.Unscoped().Where("uid = ?", uid).Take(&got).Error; err != nil {
		t.Fatalf("load category %s: %v", uid, err)
	}
	return got
}

func TestTicketCategoryRepositoryFindAll_nestedCategories_creationOrderWithParentUid(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	seedCategory(t, db, shippingCategoryUid, "Shipping", nil)
	seedCategory(t, db, refundsCategoryUid, "Refunds", &billing.Id)

	got, err := r.FindAll(bg())
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if want := []string{"Billing", "Shipping", "Refunds"}; !reflect.DeepEqual(categoryNames(got), want) {
		t.Fatalf("names = %v, want %v", categoryNames(got), want)
	}
	if got[0].ParentUid != nil || got[1].ParentUid != nil {
		t.Errorf("root parent uids = %v/%v, want nil", got[0].ParentUid, got[1].ParentUid)
	}
	if got[2].ParentUid == nil || *got[2].ParentUid != billingCategoryUid {
		t.Errorf("child parent uid = %v, want %s", got[2].ParentUid, billingCategoryUid)
	}
}

func TestTicketCategoryRepositoryFindAll_parentSoftDeleted_parentHiddenChildKeepsParentUid(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	seedCategory(t, db, refundsCategoryUid, "Refunds", &billing.Id)
	softDelete(t, db, billing)

	got, err := r.FindAll(bg())
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if !reflect.DeepEqual(categoryNames(got), []string{"Refunds"}) {
		t.Fatalf("names = %v, want [Refunds]", categoryNames(got))
	}
	if got[0].ParentUid == nil || *got[0].ParentUid != billingCategoryUid {
		t.Errorf("parent uid = %v, want the retired parent uid", got[0].ParentUid)
	}
}

func TestTicketCategoryRepositoryFindByUid_childCategory_everyFieldMapped(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	description := "money back"
	mustCreate(t, db, &entity.CustomerTicketCategoryEntity{
		Uid: refundsCategoryUid, Name: "Refunds", Description: &description, ParentId: &billing.Id, CreatedAt: baseTime, Version: 1,
	})

	got, err := r.FindByUid(bg(), refundsCategoryUid)
	if err != nil {
		t.Fatalf("find by uid: %v", err)
	}
	if got == nil || got.Uid != refundsCategoryUid || got.Name != "Refunds" || got.Description == nil || *got.Description != description {
		t.Fatalf("got = %+v", got)
	}
	if got.ParentId == nil || *got.ParentId != billing.Id || got.ParentUid == nil || *got.ParentUid != billingCategoryUid {
		t.Errorf("parent = %v/%v, want %d/%s", got.ParentId, got.ParentUid, billing.Id, billingCategoryUid)
	}
}

func TestTicketCategoryRepositoryFindByName_exactName_matchesOnlyThatCategory(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	seedCategory(t, db, billingCategoryUid, "Billing", nil)
	seedCategory(t, db, shippingCategoryUid, "Billing Disputes", nil)

	got, err := r.FindByName(bg(), "Billing")
	if err != nil {
		t.Fatalf("find by name: %v", err)
	}
	if got == nil || got.Uid != billingCategoryUid {
		t.Fatalf("got = %+v, want the Billing category", got)
	}
}

func TestTicketCategoryRepositoryFind_unknownOrSoftDeleted_nilWithoutError(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	gone := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	softDelete(t, db, gone)

	byUid, err := r.FindByUid(bg(), billingCategoryUid)
	if err != nil || byUid != nil {
		t.Errorf("soft-deleted by uid: got = %+v, err = %v, want nil, nil", byUid, err)
	}
	byName, err := r.FindByName(bg(), "Billing")
	if err != nil || byName != nil {
		t.Errorf("soft-deleted by name: got = %+v, err = %v, want nil, nil", byName, err)
	}
	unknown, err := r.FindByUid(bg(), absentUid)
	if err != nil || unknown != nil {
		t.Errorf("unknown uid: got = %+v, err = %v, want nil, nil", unknown, err)
	}
	id, err := r.FindIdByUid(bg(), billingCategoryUid)
	if err != nil || id != nil {
		t.Errorf("soft-deleted id lookup: id = %v, err = %v, want nil, nil", id, err)
	}
}

func TestTicketCategoryRepositoryFindIdByUid_liveCategory_internalId(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	seedCategory(t, db, billingCategoryUid, "Billing", nil)
	shipping := seedCategory(t, db, shippingCategoryUid, "Shipping", nil)

	id, err := r.FindIdByUid(bg(), shippingCategoryUid)
	if err != nil {
		t.Fatalf("find id: %v", err)
	}
	if id == nil || *id != shipping.Id {
		t.Fatalf("id = %v, want %d", id, shipping.Id)
	}
}

func TestTicketCategoryRepositoryInsert_childCategory_persistedWithParentUid(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	description := "money back"

	got, err := r.Insert(bg(), repo.TicketCategoryCreate{Name: "Refunds", Description: &description, ParentId: &billing.Id})
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if got.Uid == "" || got.Id == 0 || got.Name != "Refunds" {
		t.Fatalf("got = %+v", got)
	}
	if got.ParentUid == nil || *got.ParentUid != billingCategoryUid {
		t.Errorf("parent uid = %v, want %s", got.ParentUid, billingCategoryUid)
	}
	stored := loadCategory(t, db, got.Uid)
	if stored.Name != "Refunds" || stored.Description == nil || *stored.Description != description || stored.Version != 1 || stored.UpdatedAt == nil {
		t.Errorf("stored = %+v", stored)
	}
}

func TestTicketCategoryRepositoryInsert_nameOfSoftDeletedCategory_accepted(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	gone := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	softDelete(t, db, gone)

	got, err := r.Insert(bg(), repo.TicketCategoryCreate{Name: "Billing"})
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if got.Uid == billingCategoryUid || got.ParentUid != nil {
		t.Fatalf("got = %+v, want a new root category", got)
	}
}

func TestTicketCategoryRepositoryInsert_uniqueViolation_categoryNameExists(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	failCreatesOn(t, db, entity.TableNameCustomerTicketCategoryEntity, uniqueViolation())

	got, err := r.Insert(bg(), repo.TicketCategoryCreate{Name: "Billing"})
	if !errors.Is(err, repo.ErrCategoryNameExists) {
		t.Fatalf("err = %v, want ErrCategoryNameExists", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
}

func TestTicketCategoryRepositoryInsert_otherWriteFailure_errorPassedThrough(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	failCreatesOn(t, db, entity.TableNameCustomerTicketCategoryEntity, errInjected)

	_, err := r.Insert(bg(), repo.TicketCategoryCreate{Name: "Billing"})
	assertInjected(t, err)
	if errors.Is(err, repo.ErrCategoryNameExists) {
		t.Error("a non-unique failure was reported as a duplicate name")
	}
}

func TestTicketCategoryRepositoryUpdateByUid_existingCategory_fieldsAppliedAndVersionBumped(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
	seedCategory(t, db, refundsCategoryUid, "Refund", nil)
	description := "money back"

	got, err := r.UpdateByUid(bg(), refundsCategoryUid, func(c *repo.TicketCategoryResult) {
		c.Name = "Refunds"
		c.Description = &description
		c.ParentId = &billing.Id
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if got == nil || got.Name != "Refunds" || got.ParentUid == nil || *got.ParentUid != billingCategoryUid {
		t.Fatalf("got = %+v", got)
	}
	stored := loadCategory(t, db, refundsCategoryUid)
	if stored.Name != "Refunds" || stored.Description == nil || *stored.Description != description || stored.Version != 2 {
		t.Errorf("stored = %+v", stored)
	}
	if stored.ParentId == nil || *stored.ParentId != billing.Id {
		t.Errorf("parent id = %v, want %d", stored.ParentId, billing.Id)
	}
	if stored.UpdatedAt == nil || !stored.UpdatedAt.After(baseTime) {
		t.Errorf("updated_at = %v, want refreshed", stored.UpdatedAt)
	}
}

func TestTicketCategoryRepositoryUpdateByUid_unknownUid_nilAndApplySkipped(t *testing.T) {
	r, _ := newTicketCategoryRepo(t)
	applied := false

	got, err := r.UpdateByUid(bg(), absentUid, func(*repo.TicketCategoryResult) { applied = true })
	if err != nil || got != nil {
		t.Fatalf("got = %+v, err = %v, want nil, nil", got, err)
	}
	if applied {
		t.Error("apply ran for a category that does not exist")
	}
}

func TestTicketCategoryRepositoryUpdateByUid_uniqueViolation_categoryNameExistsAndRowUnchanged(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	seedCategory(t, db, billingCategoryUid, "Billing", nil)
	failCreatesOn(t, db, entity.TableNameCustomerTicketCategoryEntity, uniqueViolation())

	got, err := r.UpdateByUid(bg(), billingCategoryUid, func(c *repo.TicketCategoryResult) { c.Name = "Shipping" })
	if !errors.Is(err, repo.ErrCategoryNameExists) {
		t.Fatalf("err = %v, want ErrCategoryNameExists", err)
	}
	if got != nil {
		t.Errorf("got = %+v, want nil", got)
	}
	if stored := loadCategory(t, db, billingCategoryUid); stored.Name != "Billing" || stored.Version != 1 {
		t.Errorf("stored = %+v, want the row untouched", stored)
	}
}

func TestTicketCategoryRepositoryUpdateByUid_otherWriteFailure_errorPassedThrough(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	seedCategory(t, db, billingCategoryUid, "Billing", nil)
	failCreatesOn(t, db, entity.TableNameCustomerTicketCategoryEntity, errInjected)

	_, err := r.UpdateByUid(bg(), billingCategoryUid, func(c *repo.TicketCategoryResult) { c.Name = "Shipping" })
	assertInjected(t, err)
}

func TestTicketCategoryRepositorySoftDeleteByUid_existingCategory_hiddenButRowKept(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	seedCategory(t, db, billingCategoryUid, "Billing", nil)

	deleted, err := r.SoftDeleteByUid(bg(), billingCategoryUid)
	if err != nil || !deleted {
		t.Fatalf("deleted = %v, err = %v, want true, nil", deleted, err)
	}
	if stored := loadCategory(t, db, billingCategoryUid); !stored.DeletedAt.Valid {
		t.Error("deleted_at not set; the row must be soft deleted, not removed")
	}
	again, err := r.SoftDeleteByUid(bg(), billingCategoryUid)
	if err != nil || again {
		t.Fatalf("second delete = %v, err = %v, want false, nil", again, err)
	}
	missing, err := r.SoftDeleteByUid(bg(), absentUid)
	if err != nil || missing {
		t.Fatalf("unknown uid delete = %v, err = %v, want false, nil", missing, err)
	}
}

func TestTicketCategoryRepository_parentLookupFails_errorSurfaced(t *testing.T) {
	cases := []struct {
		name           string
		allowedQueries int
		run            func(r *ticketCategoryRepository, parentId *int64) error
	}{
		{"findAll", 1, func(r *ticketCategoryRepository, _ *int64) error {
			_, err := r.FindAll(bg())
			return err
		}},
		{"findByUid", 1, func(r *ticketCategoryRepository, _ *int64) error {
			_, err := r.FindByUid(bg(), refundsCategoryUid)
			return err
		}},
		{"insert", 0, func(r *ticketCategoryRepository, parentId *int64) error {
			_, err := r.Insert(bg(), repo.TicketCategoryCreate{Name: "Chargebacks", ParentId: parentId})
			return err
		}},
		{"updateByUid", 1, func(r *ticketCategoryRepository, _ *int64) error {
			_, err := r.UpdateByUid(bg(), refundsCategoryUid, func(*repo.TicketCategoryResult) {})
			return err
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, db := newTicketCategoryRepo(t)
			billing := seedCategory(t, db, billingCategoryUid, "Billing", nil)
			seedCategory(t, db, refundsCategoryUid, "Refunds", &billing.Id)
			failQueriesAfter(t, db, entity.TableNameCustomerTicketCategoryEntity, c.allowedQueries)

			assertInjected(t, c.run(r, &billing.Id))
		})
	}
}

func TestTicketCategoryRepository_databaseUnavailable_everyOperationFails(t *testing.T) {
	r, db := newTicketCategoryRepo(t)
	closeDatabase(t, db)

	_, err := r.FindAll(bg())
	assertFails(t, "FindAll", err)
	_, err = r.FindByUid(bg(), billingCategoryUid)
	assertFails(t, "FindByUid", err)
	_, err = r.FindByName(bg(), "Billing")
	assertFails(t, "FindByName", err)
	_, err = r.FindIdByUid(bg(), billingCategoryUid)
	assertFails(t, "FindIdByUid", err)
	_, err = r.Insert(bg(), repo.TicketCategoryCreate{Name: "Billing"})
	assertFails(t, "Insert", err)
	_, err = r.UpdateByUid(bg(), billingCategoryUid, func(*repo.TicketCategoryResult) {})
	assertFails(t, "UpdateByUid", err)
	_, err = r.SoftDeleteByUid(bg(), billingCategoryUid)
	assertFails(t, "SoftDeleteByUid", err)
}
