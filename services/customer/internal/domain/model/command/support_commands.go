package command

import "github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"

const (
	ticketSubjectMaxLen       = 255
	ticketDescriptionMaxLen   = 16000
	commentBodyMaxLen         = 16000
	categoryNameMaxLen        = 100
	categoryDescriptionMaxLen = 500
	faqQuestionMaxLen         = 500
	faqAnswerMaxLen           = 16000
)

type TicketCreateCommand struct {
	CustomerUid string
	Subject     string
	Description string
	Priority    constant.TicketPriority
	CategoryUid *string
}

type TicketUpdateCommand struct {
	Subject     *string
	Description *string
	Priority    *constant.TicketPriority
	CategoryUid *string
	AssigneeUid *string
}

type TicketStatusCommand struct {
	Status constant.TicketStatus
}

type CommentAddCommand struct {
	TicketUid  string
	AuthorUid  string
	AuthorRole constant.AuthorRole
	Body       string
	Internal   bool
}

type CategoryCreateCommand struct {
	Name        string
	Description *string
	ParentUid   *string
}

type CategoryUpdateCommand struct {
	Name        *string
	Description *string
	ParentUid   *string
}

type FaqCreateCommand struct {
	Question    string
	Answer      string
	CategoryUid *string
	Published   bool
}

type FaqUpdateCommand struct {
	Question    *string
	Answer      *string
	CategoryUid *string
	Published   *bool
}

func (c TicketCreateCommand) Validate() error {
	if err := requireUUID("customer_uid", c.CustomerUid); err != nil {
		return err
	}
	if err := requireNonBlank("subject", c.Subject); err != nil {
		return err
	}
	if err := requireMaxLen("subject", c.Subject, ticketSubjectMaxLen); err != nil {
		return err
	}
	if err := requireNonBlank("description", c.Description); err != nil {
		return err
	}
	if err := requireMaxLen("description", c.Description, ticketDescriptionMaxLen); err != nil {
		return err
	}
	if !c.Priority.IsValid() {
		return invalidf("priority must be one of LOW, MEDIUM, HIGH, URGENT: got %q", c.Priority)
	}
	return optionalUUID("category_uid", c.CategoryUid)
}

func (c TicketUpdateCommand) Validate() error {
	if err := optionalNonBlank("subject", c.Subject); err != nil {
		return err
	}
	if err := optionalMaxLen("subject", c.Subject, ticketSubjectMaxLen); err != nil {
		return err
	}
	if err := optionalNonBlank("description", c.Description); err != nil {
		return err
	}
	if err := optionalMaxLen("description", c.Description, ticketDescriptionMaxLen); err != nil {
		return err
	}
	if c.Priority != nil && !c.Priority.IsValid() {
		return invalidf("priority must be one of LOW, MEDIUM, HIGH, URGENT: got %q", *c.Priority)
	}
	if err := optionalUUID("category_uid", c.CategoryUid); err != nil {
		return err
	}
	return optionalUUID("assignee_uid", c.AssigneeUid)
}

func (c TicketStatusCommand) Validate() error {
	if !c.Status.IsValid() {
		return invalidf("status must be one of OPEN, IN_PROGRESS, RESOLVED, CLOSED: got %q", c.Status)
	}
	return nil
}

func (c CommentAddCommand) Validate() error {
	if err := requireUUID("ticket_uid", c.TicketUid); err != nil {
		return err
	}
	if err := requireUUID("author_uid", c.AuthorUid); err != nil {
		return err
	}
	if !c.AuthorRole.IsValid() {
		return invalidf("author_role must be one of CUSTOMER, AGENT: got %q", c.AuthorRole)
	}
	if err := requireNonBlank("body", c.Body); err != nil {
		return err
	}
	return requireMaxLen("body", c.Body, commentBodyMaxLen)
}

func (c CategoryCreateCommand) Validate() error {
	if err := requireNonBlank("name", c.Name); err != nil {
		return err
	}
	if err := requireMaxLen("name", c.Name, categoryNameMaxLen); err != nil {
		return err
	}
	if err := optionalMaxLen("description", c.Description, categoryDescriptionMaxLen); err != nil {
		return err
	}
	return optionalUUID("parent_uid", c.ParentUid)
}

func (c CategoryUpdateCommand) Validate() error {
	if err := optionalNonBlank("name", c.Name); err != nil {
		return err
	}
	if err := optionalMaxLen("name", c.Name, categoryNameMaxLen); err != nil {
		return err
	}
	if err := optionalMaxLen("description", c.Description, categoryDescriptionMaxLen); err != nil {
		return err
	}
	return optionalUUID("parent_uid", c.ParentUid)
}

func (c FaqCreateCommand) Validate() error {
	if err := requireNonBlank("question", c.Question); err != nil {
		return err
	}
	if err := requireMaxLen("question", c.Question, faqQuestionMaxLen); err != nil {
		return err
	}
	if err := requireNonBlank("answer", c.Answer); err != nil {
		return err
	}
	if err := requireMaxLen("answer", c.Answer, faqAnswerMaxLen); err != nil {
		return err
	}
	return optionalUUID("category_uid", c.CategoryUid)
}

func (c FaqUpdateCommand) Validate() error {
	if err := optionalNonBlank("question", c.Question); err != nil {
		return err
	}
	if err := optionalMaxLen("question", c.Question, faqQuestionMaxLen); err != nil {
		return err
	}
	if err := optionalNonBlank("answer", c.Answer); err != nil {
		return err
	}
	if err := optionalMaxLen("answer", c.Answer, faqAnswerMaxLen); err != nil {
		return err
	}
	return optionalUUID("category_uid", c.CategoryUid)
}
