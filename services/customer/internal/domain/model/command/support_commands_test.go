package command

import (
	"errors"
	"testing"

	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
)

const (
	testCustomerUid = "018f0000-0000-7000-8000-0000000000c1"
	testCategoryUid = "018f0000-0000-7000-8000-0000000000d1"
	testTicketUid   = "018f0000-0000-7000-8000-0000000000e1"
)

func validTicketCreate() TicketCreateCommand {
	return TicketCreateCommand{
		CustomerUid: testCustomerUid,
		Subject:     "Order missing",
		Description: "My order never arrived",
		Priority:    constant.TICKET_PRIORITY_MEDIUM,
	}
}

func TestTicketCreateCommand_ValidInput_Passes(t *testing.T) {
	cmd := validTicketCreate()
	cmd.Subject = longText(255)
	cmd.Description = longText(16000)
	cmd.CategoryUid = ptr(testCategoryUid)
	if err := cmd.Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
}

func TestTicketCreateCommand_BlankOrOverlongFields_ReturnsErrInvalidCommand(t *testing.T) {
	cases := map[string]func(*TicketCreateCommand){
		"blankSubject":        func(c *TicketCreateCommand) { c.Subject = "   " },
		"overlongSubject":     func(c *TicketCreateCommand) { c.Subject = longText(256) },
		"blankDescription":    func(c *TicketCreateCommand) { c.Description = "\t" },
		"overlongDescription": func(c *TicketCreateCommand) { c.Description = longText(16001) },
		"nonUuidCustomer":     func(c *TicketCreateCommand) { c.CustomerUid = "7" },
		"nonUuidCategory":     func(c *TicketCreateCommand) { c.CategoryUid = ptr("cat") },
		"unknownPriority":     func(c *TicketCreateCommand) { c.Priority = "SOMEDAY" },
		"emptyPriority":       func(c *TicketCreateCommand) { c.Priority = "" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			cmd := validTicketCreate()
			mutate(&cmd)
			if err := cmd.Validate(); !errors.Is(err, ErrInvalidCommand) {
				t.Fatalf("want ErrInvalidCommand, got %v", err)
			}
		})
	}
}

func TestTicketUpdateCommand_EmptyUpdate_Passes(t *testing.T) {
	if err := (TicketUpdateCommand{}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
}

func TestTicketUpdateCommand_InvalidOptionalFields_ReturnsErrInvalidCommand(t *testing.T) {
	bad := constant.TicketPriority("LATER")
	cases := map[string]TicketUpdateCommand{
		"blankSubject":    {Subject: ptr(" ")},
		"overlongSubject": {Subject: ptr(longText(256))},
		"badPriority":     {Priority: &bad},
		"nonUuidAssignee": {AssigneeUid: ptr("agent-1")},
		"nonUuidCategory": {CategoryUid: ptr("cat")},
	}
	for name, cmd := range cases {
		t.Run(name, func(t *testing.T) {
			if err := cmd.Validate(); !errors.Is(err, ErrInvalidCommand) {
				t.Fatalf("want ErrInvalidCommand, got %v", err)
			}
		})
	}
}

func TestTicketStatusCommand_UnknownStatus_ReturnsErrInvalidCommand(t *testing.T) {
	if err := (TicketStatusCommand{Status: "DONE"}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand, got %v", err)
	}
	if err := (TicketStatusCommand{Status: constant.TICKET_STATUS_CLOSED}).Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
}

func TestCommentAddCommand_BlankBodyOrBadAuthor_ReturnsErrInvalidCommand(t *testing.T) {
	valid := CommentAddCommand{TicketUid: testTicketUid, AuthorUid: testCustomerUid, AuthorRole: constant.AUTHOR_ROLE_CUSTOMER, Body: "hello"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	cases := map[string]func(*CommentAddCommand){
		"blankBody":     func(c *CommentAddCommand) { c.Body = "  " },
		"overlongBody":  func(c *CommentAddCommand) { c.Body = longText(16001) },
		"badRole":       func(c *CommentAddCommand) { c.AuthorRole = "BOT" },
		"nonUuidAuthor": func(c *CommentAddCommand) { c.AuthorUid = "7" },
		"nonUuidTicket": func(c *CommentAddCommand) { c.TicketUid = "t" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			cmd := valid
			mutate(&cmd)
			if err := cmd.Validate(); !errors.Is(err, ErrInvalidCommand) {
				t.Fatalf("want ErrInvalidCommand, got %v", err)
			}
		})
	}
}

func TestCategoryCommands_NameAndDescriptionBounds_Enforced(t *testing.T) {
	if err := (CategoryCreateCommand{Name: longText(100), Description: ptr("")}).Validate(); err != nil {
		t.Fatalf("want nil for name at cap and empty description, got %v", err)
	}
	if err := (CategoryCreateCommand{Name: longText(101)}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand for overlong name, got %v", err)
	}
	if err := (CategoryCreateCommand{Name: "Billing", Description: ptr(longText(501))}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand for overlong description, got %v", err)
	}
	if err := (CategoryCreateCommand{Name: "Billing", ParentUid: ptr("root")}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand for non-uuid parent, got %v", err)
	}
	if err := (CategoryUpdateCommand{}).Validate(); err != nil {
		t.Fatalf("want nil for empty update, got %v", err)
	}
	if err := (CategoryUpdateCommand{Name: ptr(" ")}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand for blank name, got %v", err)
	}
}

func TestFaqCommands_QuestionAndAnswerBounds_Enforced(t *testing.T) {
	valid := FaqCreateCommand{Question: longText(500), Answer: "  padded answer  "}
	if err := valid.Validate(); err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if valid.Answer != "  padded answer  " {
		t.Fatalf("answer whitespace must be preserved, got %q", valid.Answer)
	}
	if err := (FaqCreateCommand{Question: longText(501), Answer: "a"}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand for overlong question, got %v", err)
	}
	if err := (FaqCreateCommand{Question: "q", Answer: longText(16001)}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand for overlong answer, got %v", err)
	}
	if err := (FaqCreateCommand{Question: "q", Answer: " "}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand for blank answer, got %v", err)
	}
	if err := (FaqUpdateCommand{}).Validate(); err != nil {
		t.Fatalf("want nil for empty update, got %v", err)
	}
	if err := (FaqUpdateCommand{CategoryUid: ptr("cat")}).Validate(); !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("want ErrInvalidCommand for non-uuid category, got %v", err)
	}
}
