package pgsql

import (
	"testing"

	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/support/utils"
)

func TestToWorkResult_WhenAuxiliaryRowsPresent_PopulatesAuthorsAndSubjects(t *testing.T) {
	desc := "epic"
	cover := utils.ToJsonString([]string{"c1", "c2"})
	ent := &entity.WorkEntity{
		Id:               42,
		Uid:              "w-uid",
		Title:            "Hobbit",
		Description:      &desc,
		CoverUid:         cover,
		FirstPublishDate: nil,
		OlKey:            nil,
	}
	authors := []workAuthorRow{{WorkId: 42, AuthorUid: "a-1", AuthorName: "Tolkien"}}
	subjects := []workSubjectRow{
		{WorkId: 42, SubjectUid: "s-1", SubjectName: "Fantasy"},
		{WorkId: 42, SubjectUid: "s-2", SubjectName: "Adventure"},
	}

	got := toWorkResult(ent, authors, subjects)
	if got.Id != 42 || got.Uid != "w-uid" || got.Title != "Hobbit" {
		t.Errorf("scalar fields = %+v", got)
	}
	if got.Description == nil || *got.Description != "epic" {
		t.Errorf("Description = %v", got.Description)
	}
	if len(got.CoverUids) != 2 || got.CoverUids[0] != "c1" {
		t.Errorf("CoverUids = %v", got.CoverUids)
	}
	if len(got.AuthorUids) != 1 || got.AuthorUids[0] != "a-1" || got.AuthorNames[0] != "Tolkien" {
		t.Errorf("Authors fan-out = %v / %v", got.AuthorUids, got.AuthorNames)
	}
	if len(got.SubjectUids) != 2 || got.SubjectNames[1] != "Adventure" {
		t.Errorf("Subjects fan-out = %v / %v", got.SubjectUids, got.SubjectNames)
	}
}

func TestToWorkResult_WhenAuxiliaryRowsEmpty_ReturnsEmptyCollections(t *testing.T) {
	ent := &entity.WorkEntity{Id: 1, Uid: "u", Title: "T"}
	got := toWorkResult(ent, nil, nil)
	if len(got.AuthorUids) != 0 || len(got.AuthorNames) != 0 {
		t.Errorf("AuthorUids/Names = %v/%v, want empty", got.AuthorUids, got.AuthorNames)
	}
	if len(got.SubjectUids) != 0 || len(got.SubjectNames) != 0 {
		t.Errorf("SubjectUids/Names = %v/%v, want empty", got.SubjectUids, got.SubjectNames)
	}
	if got.CoverUids != nil {
		t.Errorf("CoverUids = %v, want nil", got.CoverUids)
	}
}

func TestToAuthorResult_WhenRowHasEveryColumn_ReturnsResult(t *testing.T) {
	bio := "born 1920"
	birth := "1920-01-02"
	photo := utils.ToJsonString([]string{"p1"})
	alt := utils.ToJsonString([]string{"AA", "BB"})
	got := toAuthorResult(&entity.AuthorEntity{
		Id:            7,
		Uid:           "a-1",
		Name:          "Asimov",
		Bio:           &bio,
		BirthDate:     &birth,
		PhotoUid:      photo,
		AlternateName: alt,
	})
	if got.Id != 7 || got.Name != "Asimov" || got.Bio == nil || *got.Bio != bio {
		t.Errorf("scalar fields = %+v", got)
	}
	if len(got.PhotoUids) != 1 || got.PhotoUids[0] != "p1" {
		t.Errorf("PhotoUids = %v", got.PhotoUids)
	}
	if len(got.AlternateNames) != 2 || got.AlternateNames[1] != "BB" {
		t.Errorf("AlternateNames = %v", got.AlternateNames)
	}
}

func TestToPublisherResult_WhenRowHasEveryColumn_ReturnsResult(t *testing.T) {
	addr := "Mars"
	ol := "OL/123"
	got := toPublisherResult(&entity.PublisherEntity{
		Id: 9, Uid: "p-1", Name: "Acme", Address: &addr, OlKey: &ol,
	})
	if got.Id != 9 || got.Uid != "p-1" || got.Name != "Acme" {
		t.Errorf("scalar = %+v", got)
	}
	if got.Address == nil || *got.Address != "Mars" {
		t.Errorf("Address = %v", got.Address)
	}
	if got.OlKey == nil || *got.OlKey != "OL/123" {
		t.Errorf("OlKey = %v", got.OlKey)
	}
}

func TestToSubjectResult_WhenRowHasEveryColumn_ReturnsResult(t *testing.T) {
	got := toSubjectResult(&entity.SubjectEntity{Id: 5, Uid: "s-1", Name: "Fantasy"})
	if got.Id != 5 || got.Uid != "s-1" || got.Name != "Fantasy" {
		t.Errorf("got = %+v", got)
	}
}

func TestToEditionResult_WhenRowHasEveryColumn_ReturnsResult(t *testing.T) {
	isbn10 := "0-7475-3269-9"
	isbn13 := "978-0-7475-3269-9"
	pages := int32(223)
	cover := utils.ToJsonString([]string{"c1"})
	lang := utils.ToJsonString([]string{"en", "ko"})
	pubUid := "p-1"
	pubName := "Acme"
	row := &editionRow{
		Id:            11,
		Uid:           "e-1",
		Title:         "1st",
		Isbn10:        &isbn10,
		Isbn13:        &isbn13,
		NumberOfPage:  &pages,
		CoverUid:      cover,
		Language:      lang,
		WorkUid:       "w-1",
		WorkTitle:     "W",
		PublisherUid:  &pubUid,
		PublisherName: &pubName,
	}
	got := toEditionResult(row)
	if got.Id != 11 || got.WorkUid != "w-1" || got.WorkTitle != "W" {
		t.Errorf("scalar = %+v", got)
	}
	if got.Isbn10 == nil || *got.Isbn10 != isbn10 {
		t.Errorf("Isbn10 = %v", got.Isbn10)
	}
	if len(got.CoverUids) != 1 || got.CoverUids[0] != "c1" {
		t.Errorf("CoverUids = %v", got.CoverUids)
	}
	if len(got.Languages) != 2 || got.Languages[1] != "ko" {
		t.Errorf("Languages = %v", got.Languages)
	}
	if got.PublisherUid == nil || *got.PublisherUid != "p-1" {
		t.Errorf("PublisherUid = %v", got.PublisherUid)
	}
}
