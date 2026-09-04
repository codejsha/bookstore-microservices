package command

import (
	"errors"
	"strings"
	"testing"
)

func ptr[T any](v T) *T { return &v }

const (
	authorUid    = "0197f7c0-1c1a-7000-8000-000000000001"
	altAuthorUid = "0197f7c0-1c1a-7000-8000-000000000002"
	workUid      = "0197f7c0-1c1a-7000-8000-0000000000a1"
	publisherUid = "0197f7c0-1c1a-7000-8000-0000000000b1"
)

func TestWorkCreateCommand_Validate(t *testing.T) {
	valid := WorkCreateCommand{Title: "Dune", AuthorUids: []string{authorUid}}
	tests := []struct {
		name    string
		mutate  func(c *WorkCreateCommand)
		wantErr bool
	}{
		{"whenAllFieldsValid_returnsNil", func(c *WorkCreateCommand) {}, false},
		{"whenTitleBlank_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.Title = "   " }, true},
		{"whenAuthorUidsEmpty_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.AuthorUids = nil }, true},
		{"whenAuthorUidBlank_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.AuthorUids = []string{authorUid, " "} }, true},
		{"whenAuthorUidNotUuid_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.AuthorUids = []string{"a1"} }, true},
		{"whenAuthorUidsDuplicated_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.AuthorUids = []string{authorUid, authorUid} }, true},
		{"whenAuthorUidsDistinct_returnsNil", func(c *WorkCreateCommand) { c.AuthorUids = []string{authorUid, altAuthorUid} }, false},
		{"whenSubjectNamesDuplicated_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.SubjectNames = []string{"sf", "sf"} }, true},
		{"whenTitleOverLimit_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.Title = strings.Repeat("a", 513) }, true},
		{"whenTitleAtLimit_returnsNil", func(c *WorkCreateCommand) { c.Title = strings.Repeat("a", 512) }, false},
		{"whenFirstPublishDateOverLimit_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.FirstPublishDate = ptr(strings.Repeat("1", 33)) }, true},
		{"whenOlKeyOverLimit_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.OlKey = ptr(strings.Repeat("k", 65)) }, true},
		{"whenOlKeyEmpty_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.OlKey = ptr("") }, true},
		{"whenSubjectNameOverLimit_returnsErrInvalidCommand", func(c *WorkCreateCommand) { c.SubjectNames = []string{strings.Repeat("s", 256)} }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := valid
			tt.mutate(&cmd)
			err := cmd.Validate()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidCommand) {
					t.Fatalf("want ErrInvalidCommand, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("want nil, got %v", err)
			}
		})
	}
}

// A partial update carries only the fields the caller set, so an empty command is
// valid and only a field that is present but blank is rejected.
func TestWorkUpdateCommand_WhenTitleBlank_ReturnsErrInvalidCommand(t *testing.T) {
	if err := (WorkUpdateCommand{}).Validate(); err != nil {
		t.Fatalf("empty update must be valid, got %v", err)
	}
	if err := (WorkUpdateCommand{Title: ptr("  ")}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
}

func TestEditionCreateCommand_Validate(t *testing.T) {
	valid := EditionCreateCommand{Title: "Dune", WorkUid: workUid}
	tests := []struct {
		name    string
		mutate  func(c *EditionCreateCommand)
		wantErr bool
	}{
		{"whenIsbnOmitted_returnsNil", func(c *EditionCreateCommand) {}, false},
		{"whenIsbn10Valid_returnsNil", func(c *EditionCreateCommand) { c.Isbn10 = ptr("0441172717") }, false},
		{"whenIsbn10EndsWithX_returnsNil", func(c *EditionCreateCommand) { c.Isbn10 = ptr("097522980X") }, false},
		{"whenIsbn10ChecksumWrong_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.Isbn10 = ptr("0441172718") }, true},
		{"whenIsbn10LengthWrong_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.Isbn10 = ptr("044117271") }, true},
		{"whenIsbn13Valid_returnsNil", func(c *EditionCreateCommand) { c.Isbn13 = ptr("9780441172719") }, false},
		{"whenIsbn13ChecksumWrong_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.Isbn13 = ptr("9780441172718") }, true},
		{"whenIsbn13PrefixWrong_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.Isbn13 = ptr("1234567890128") }, true},
		{"whenIsbn13HasNonDigit_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.Isbn13 = ptr("97804411727AB") }, true},
		{"whenNumberOfPagesZero_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.NumberOfPages = ptr(int32(0)) }, true},
		{"whenNumberOfPagesNegative_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.NumberOfPages = ptr(int32(-5)) }, true},
		{"whenWorkUidBlank_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.WorkUid = " " }, true},
		{"whenWorkUidNotUuid_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.WorkUid = "w1" }, true},
		{"whenPublisherUidNotUuid_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.PublisherUid = ptr("p1") }, true},
		{"whenPublisherUidValid_returnsNil", func(c *EditionCreateCommand) { c.PublisherUid = ptr(publisherUid) }, false},
		{"whenTitleOverLimit_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.Title = strings.Repeat("a", 513) }, true},
		{"whenPublishDateOverLimit_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.PublishDate = ptr(strings.Repeat("1", 33)) }, true},
		{"whenPhysicalFormatOverLimit_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.PhysicalFormat = ptr(strings.Repeat("f", 65)) }, true},
		{"whenOlKeyOverLimit_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.OlKey = ptr(strings.Repeat("k", 65)) }, true},
		{"whenOlKeyEmpty_returnsErrInvalidCommand", func(c *EditionCreateCommand) { c.OlKey = ptr("") }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := valid
			tt.mutate(&cmd)
			err := cmd.Validate()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidCommand) {
					t.Fatalf("want ErrInvalidCommand, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("want nil, got %v", err)
			}
		})
	}
}

func TestAuthorCreateCommand_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cmd     AuthorCreateCommand
		wantErr bool
	}{
		{"whenBirthBeforeDeath_returnsNil", AuthorCreateCommand{Name: "Frank Herbert", BirthDate: ptr("1920-10-08"), DeathDate: ptr("1986-02-11")}, false},
		{"whenDeathBeforeBirth_returnsErrInvalidCommand", AuthorCreateCommand{Name: "Frank Herbert", BirthDate: ptr("1986-02-11"), DeathDate: ptr("1920-10-08")}, true},
		{"whenDatesNotParseable_returnsNil", AuthorCreateCommand{Name: "Homer", BirthDate: ptr("8th century BC"), DeathDate: ptr("unknown")}, false},
		{"whenNameBlank_returnsErrInvalidCommand", AuthorCreateCommand{Name: ""}, true},
		{"whenNameOverLimit_returnsErrInvalidCommand", AuthorCreateCommand{Name: strings.Repeat("n", 256)}, true},
		{"whenBirthDateOverLimit_returnsErrInvalidCommand", AuthorCreateCommand{Name: "Homer", BirthDate: ptr(strings.Repeat("1", 33))}, true},
		{"whenDeathDateOverLimit_returnsErrInvalidCommand", AuthorCreateCommand{Name: "Homer", DeathDate: ptr(strings.Repeat("1", 33))}, true},
		{"whenOlKeyOverLimit_returnsErrInvalidCommand", AuthorCreateCommand{Name: "Homer", OlKey: ptr(strings.Repeat("k", 65))}, true},
		{"whenOlKeyEmpty_returnsErrInvalidCommand", AuthorCreateCommand{Name: "Homer", OlKey: ptr("")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidCommand) {
					t.Fatalf("want ErrInvalidCommand, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("want nil, got %v", err)
			}
		})
	}
}

// Publisher and subject share one rule set, so both commands are exercised here:
// blank names, over-long names and addresses, and an empty ol key on update.
func TestPublisherAndSubjectCommands_WhenNameOrOlKeyInvalid_ReturnErrInvalidCommand(t *testing.T) {
	if err := (PublisherCreateCommand{Name: " "}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (SubjectCreateCommand{Name: "sf"}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if err := (SubjectUpdateCommand{Name: ptr("")}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (SubjectCreateCommand{Name: strings.Repeat("s", 256)}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (PublisherCreateCommand{Name: strings.Repeat("p", 256)}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (PublisherCreateCommand{Name: "Ace", Address: ptr(strings.Repeat("a", 256))}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (PublisherUpdateCommand{OlKey: ptr("")}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
}
