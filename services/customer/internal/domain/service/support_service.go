package service

import (
	"context"
	"fmt"
	"time"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
)

var _ usecase.SupportUseCase = (*supportService)(nil)

type supportService struct {
	ticketRepo   repo.TicketRepo
	commentRepo  repo.TicketCommentRepo
	categoryRepo repo.TicketCategoryRepo
	faqRepo      repo.FaqRepo
}

func NewSupportService(
	ticketRepo repo.TicketRepo,
	commentRepo repo.TicketCommentRepo,
	categoryRepo repo.TicketCategoryRepo,
	faqRepo repo.FaqRepo,
) usecase.SupportUseCase {
	return &supportService{
		ticketRepo:   ticketRepo,
		commentRepo:  commentRepo,
		categoryRepo: categoryRepo,
		faqRepo:      faqRepo,
	}
}

func (s *supportService) resolveCategoryId(ctx context.Context, categoryUid *string) (*int64, error) {
	if categoryUid == nil {
		return nil, nil
	}
	id, err := s.categoryRepo.FindIdByUid(ctx, *categoryUid)
	if err != nil {
		return nil, fmt.Errorf("resolve category: %w", err)
	}
	if id == nil {
		return nil, fmt.Errorf("category %s: %w", *categoryUid, repo.ErrUnknownReference)
	}
	return id, nil
}

// ─── Tickets ────────────────────────────────────────────────────────

func (s *supportService) CreateTicket(ctx context.Context, cmd command.TicketCreateCommand) (*aggregate.TicketAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	categoryId, err := s.resolveCategoryId(ctx, cmd.CategoryUid)
	if err != nil {
		return nil, err
	}
	result, err := s.ticketRepo.Insert(ctx, repo.TicketCreate{
		CustomerUid: cmd.CustomerUid,
		CategoryId:  categoryId,
		Subject:     cmd.Subject,
		Description: cmd.Description,
		Status:      string(constant.TICKET_STATUS_OPEN),
		Priority:    string(cmd.Priority),
	})
	if err != nil {
		return nil, fmt.Errorf("create ticket: %w", err)
	}
	return toTicketAggregate(result), nil
}

func (s *supportService) FindTicket(ctx context.Context, uid string) (*aggregate.TicketAggregate, error) {
	result, err := s.ticketRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("find ticket: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("find ticket %s: %w", uid, usecase.ErrNotFound)
	}
	return toTicketAggregate(result), nil
}

func (s *supportService) FindAllTickets(ctx context.Context, opt option.TicketQueryOption) (int64, []*aggregate.TicketAggregate, error) {
	total, results, err := s.ticketRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, fmt.Errorf("list tickets: %w", err)
	}
	out := make([]*aggregate.TicketAggregate, len(results))
	for i, r := range results {
		out[i] = toTicketAggregate(r)
	}
	return total, out, nil
}

func (s *supportService) UpdateTicket(ctx context.Context, uid string, cmd command.TicketUpdateCommand) (*aggregate.TicketAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	categoryId, err := s.resolveCategoryId(ctx, cmd.CategoryUid)
	if err != nil {
		return nil, err
	}
	result, err := s.ticketRepo.UpdateByUid(ctx, uid, func(t *repo.TicketResult) {
		if cmd.Subject != nil {
			t.Subject = *cmd.Subject
		}
		if cmd.Description != nil {
			t.Description = *cmd.Description
		}
		if cmd.Priority != nil {
			t.Priority = string(*cmd.Priority)
		}
		if categoryId != nil {
			t.CategoryId = categoryId
		}
		if cmd.AssigneeUid != nil {
			t.AssigneeUid = cmd.AssigneeUid
		}
	})
	if err != nil {
		return nil, fmt.Errorf("update ticket: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("update ticket %s: %w", uid, usecase.ErrNotFound)
	}
	return toTicketAggregate(result), nil
}

func (s *supportService) UpdateTicketStatus(ctx context.Context, uid string, cmd command.TicketStatusCommand) (*aggregate.TicketAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	result, err := s.ticketRepo.UpdateByUid(ctx, uid, func(t *repo.TicketResult) {
		t.Status = string(cmd.Status)
		if cmd.Status == constant.TICKET_STATUS_RESOLVED && t.ResolvedAt == nil {
			now := time.Now().UTC()
			t.ResolvedAt = &now
		}
	})
	if err != nil {
		return nil, fmt.Errorf("update ticket status: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("update ticket status %s: %w", uid, usecase.ErrNotFound)
	}
	return toTicketAggregate(result), nil
}

// ─── Ticket comments ────────────────────────────────────────────────

func (s *supportService) AddComment(ctx context.Context, cmd command.CommentAddCommand) (*aggregate.TicketCommentAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	ticket, err := s.ticketRepo.FindByUid(ctx, cmd.TicketUid)
	if err != nil {
		return nil, fmt.Errorf("find ticket: %w", err)
	}
	if ticket == nil {
		return nil, fmt.Errorf("add comment to ticket %s: %w", cmd.TicketUid, usecase.ErrNotFound)
	}
	result, err := s.commentRepo.Insert(ctx, repo.TicketCommentCreate{
		TicketId:   ticket.Id,
		AuthorUid:  cmd.AuthorUid,
		AuthorRole: string(cmd.AuthorRole),
		Body:       cmd.Body,
		Internal:   cmd.Internal,
	})
	if err != nil {
		return nil, fmt.Errorf("add comment: %w", err)
	}
	return toTicketCommentAggregate(ticket.Uid, result), nil
}

func (s *supportService) FindComments(ctx context.Context, ticketUid string, includeInternal bool) ([]*aggregate.TicketCommentAggregate, error) {
	ticket, err := s.ticketRepo.FindByUid(ctx, ticketUid)
	if err != nil {
		return nil, fmt.Errorf("find ticket: %w", err)
	}
	if ticket == nil {
		return nil, fmt.Errorf("list comments of ticket %s: %w", ticketUid, usecase.ErrNotFound)
	}
	results, err := s.commentRepo.FindAllByTicketId(ctx, ticket.Id, includeInternal)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	out := make([]*aggregate.TicketCommentAggregate, len(results))
	for i, r := range results {
		out[i] = toTicketCommentAggregate(ticket.Uid, r)
	}
	return out, nil
}

// ─── Ticket categories ──────────────────────────────────────────────

func (s *supportService) CreateCategory(ctx context.Context, cmd command.CategoryCreateCommand) (*aggregate.TicketCategoryAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.categoryRepo.FindByName(ctx, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("find category by name: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("create category %q: %w", cmd.Name, repo.ErrCategoryNameExists)
	}
	parentId, err := s.resolveCategoryId(ctx, cmd.ParentUid)
	if err != nil {
		return nil, err
	}
	result, err := s.categoryRepo.Insert(ctx, repo.TicketCategoryCreate{
		Name:        cmd.Name,
		Description: cmd.Description,
		ParentId:    parentId,
	})
	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	return toTicketCategoryAggregate(result), nil
}

func (s *supportService) FindCategory(ctx context.Context, uid string) (*aggregate.TicketCategoryAggregate, error) {
	result, err := s.categoryRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("find category: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("find category %s: %w", uid, usecase.ErrNotFound)
	}
	return toTicketCategoryAggregate(result), nil
}

func (s *supportService) FindAllCategories(ctx context.Context) ([]*aggregate.TicketCategoryAggregate, error) {
	results, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	out := make([]*aggregate.TicketCategoryAggregate, len(results))
	for i, r := range results {
		out[i] = toTicketCategoryAggregate(r)
	}
	return out, nil
}

func (s *supportService) UpdateCategory(ctx context.Context, uid string, cmd command.CategoryUpdateCommand) (*aggregate.TicketCategoryAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	current, err := s.categoryRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("find category: %w", err)
	}
	if current == nil {
		return nil, fmt.Errorf("update category %s: %w", uid, usecase.ErrNotFound)
	}
	if cmd.Name != nil && *cmd.Name != current.Name {
		existing, err := s.categoryRepo.FindByName(ctx, *cmd.Name)
		if err != nil {
			return nil, fmt.Errorf("find category by name: %w", err)
		}
		if existing != nil && existing.Uid != uid {
			return nil, fmt.Errorf("rename category to %q: %w", *cmd.Name, repo.ErrCategoryNameExists)
		}
	}
	parentId, err := s.resolveCategoryId(ctx, cmd.ParentUid)
	if err != nil {
		return nil, err
	}
	result, err := s.categoryRepo.UpdateByUid(ctx, uid, func(c *repo.TicketCategoryResult) {
		if cmd.Name != nil {
			c.Name = *cmd.Name
		}
		if cmd.Description != nil {
			c.Description = cmd.Description
		}
		if parentId != nil {
			c.ParentId = parentId
		}
	})
	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("update category %s: %w", uid, usecase.ErrNotFound)
	}
	return toTicketCategoryAggregate(result), nil
}

func (s *supportService) DeleteCategory(ctx context.Context, uid string) error {
	deleted, err := s.categoryRepo.SoftDeleteByUid(ctx, uid)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	if !deleted {
		return fmt.Errorf("delete category %s: %w", uid, usecase.ErrNotFound)
	}
	return nil
}

// ─── FAQs ───────────────────────────────────────────────────────────

func (s *supportService) CreateFaq(ctx context.Context, cmd command.FaqCreateCommand) (*aggregate.FaqAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	categoryId, err := s.resolveCategoryId(ctx, cmd.CategoryUid)
	if err != nil {
		return nil, err
	}
	result, err := s.faqRepo.Insert(ctx, repo.FaqCreate{
		CategoryId: categoryId,
		Question:   cmd.Question,
		Answer:     cmd.Answer,
		Published:  cmd.Published,
	})
	if err != nil {
		return nil, fmt.Errorf("create faq: %w", err)
	}
	return toFaqAggregate(result), nil
}

func (s *supportService) FindFaq(ctx context.Context, uid string, incrementView bool) (*aggregate.FaqAggregate, error) {
	result, err := s.faqRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("find faq: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("find faq %s: %w", uid, usecase.ErrNotFound)
	}
	if incrementView && result.Published {
		if err := s.faqRepo.IncrementViewCount(ctx, uid); err == nil {
			result.ViewCount++
		}
	}
	return toFaqAggregate(result), nil
}

func (s *supportService) FindAllFaqs(ctx context.Context, opt option.FaqQueryOption) (int64, []*aggregate.FaqAggregate, error) {
	total, results, err := s.faqRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, fmt.Errorf("list faqs: %w", err)
	}
	out := make([]*aggregate.FaqAggregate, len(results))
	for i, r := range results {
		out[i] = toFaqAggregate(r)
	}
	return total, out, nil
}

func (s *supportService) UpdateFaq(ctx context.Context, uid string, cmd command.FaqUpdateCommand) (*aggregate.FaqAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	categoryId, err := s.resolveCategoryId(ctx, cmd.CategoryUid)
	if err != nil {
		return nil, err
	}
	result, err := s.faqRepo.UpdateByUid(ctx, uid, func(f *repo.FaqResult) {
		if cmd.Question != nil {
			f.Question = *cmd.Question
		}
		if cmd.Answer != nil {
			f.Answer = *cmd.Answer
		}
		if categoryId != nil {
			f.CategoryId = categoryId
		}
		if cmd.Published != nil {
			f.Published = *cmd.Published
		}
	})
	if err != nil {
		return nil, fmt.Errorf("update faq: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("update faq %s: %w", uid, usecase.ErrNotFound)
	}
	return toFaqAggregate(result), nil
}

func (s *supportService) DeleteFaq(ctx context.Context, uid string) error {
	deleted, err := s.faqRepo.SoftDeleteByUid(ctx, uid)
	if err != nil {
		return fmt.Errorf("delete faq: %w", err)
	}
	if !deleted {
		return fmt.Errorf("delete faq %s: %w", uid, usecase.ErrNotFound)
	}
	return nil
}

// ─── mapping helpers ────────────────────────────────────────────────

func updatedOrCreated(createdAt time.Time, updatedAt *time.Time) time.Time {
	if updatedAt != nil {
		return *updatedAt
	}
	return createdAt
}

func toTicketAggregate(r *repo.TicketResult) *aggregate.TicketAggregate {
	return &aggregate.TicketAggregate{
		Uid:         r.Uid,
		CustomerUid: r.CustomerUid,
		CategoryUid: r.CategoryUid,
		AssigneeUid: r.AssigneeUid,
		Subject:     r.Subject,
		Description: r.Description,
		Status:      constant.TicketStatus(r.Status),
		Priority:    constant.TicketPriority(r.Priority),
		ResolvedAt:  r.ResolvedAt,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   updatedOrCreated(r.CreatedAt, r.UpdatedAt),
	}
}

func toTicketCommentAggregate(ticketUid string, r *repo.TicketCommentResult) *aggregate.TicketCommentAggregate {
	return &aggregate.TicketCommentAggregate{
		Uid:        r.Uid,
		TicketUid:  ticketUid,
		AuthorUid:  r.AuthorUid,
		AuthorRole: constant.AuthorRole(r.AuthorRole),
		Body:       r.Body,
		Internal:   r.Internal,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  updatedOrCreated(r.CreatedAt, r.UpdatedAt),
	}
}

func toTicketCategoryAggregate(r *repo.TicketCategoryResult) *aggregate.TicketCategoryAggregate {
	return &aggregate.TicketCategoryAggregate{
		Uid:         r.Uid,
		Name:        r.Name,
		Description: r.Description,
		ParentUid:   r.ParentUid,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   updatedOrCreated(r.CreatedAt, r.UpdatedAt),
	}
}

func toFaqAggregate(r *repo.FaqResult) *aggregate.FaqAggregate {
	return &aggregate.FaqAggregate{
		Uid:         r.Uid,
		CategoryUid: r.CategoryUid,
		Question:    r.Question,
		Answer:      r.Answer,
		ViewCount:   r.ViewCount,
		Published:   r.Published,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   updatedOrCreated(r.CreatedAt, r.UpdatedAt),
	}
}
