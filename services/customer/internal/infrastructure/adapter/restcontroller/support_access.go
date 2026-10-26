package restcontroller

import (
	"context"

	"github.com/google/uuid"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

type supportCaller struct {
	uid   string
	staff bool
}

func (c supportCaller) owns(ownerUid string) bool {
	owner, err := uuid.Parse(ownerUid)
	if err != nil {
		return false
	}
	return c.uid == owner.String()
}

func resolveSupportCaller(ctx context.Context) (supportCaller, error) {
	caller := httpx.CallerFromContext(ctx)
	if caller.UserUid == "" {
		return supportCaller{}, httpx.MapUnauthenticated(ctx, httpx.ErrUnauthenticated)
	}
	parsed, err := uuid.Parse(caller.UserUid)
	if err != nil {
		return supportCaller{}, httpx.MapForbidden(ctx, httpx.ErrForbidden)
	}
	return supportCaller{uid: parsed.String(), staff: caller.IsAdmin}, nil
}

func mapSupportError(ctx context.Context, err error) error {
	return httpx.MapBusinessError(ctx,
		httpx.MapUnknownReference(ctx,
			httpx.MapCategoryNameExists(ctx,
				httpx.MapForbidden(ctx,
					httpx.MapNotFound(ctx, err)))))
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toTicketFindResponse(a *aggregate.TicketAggregate) openapi.TicketFindResponse {
	return openapi.TicketFindResponse{
		Uid:         a.Uid,
		CustomerUid: a.CustomerUid,
		CategoryUid: a.CategoryUid,
		AssigneeUid: a.AssigneeUid,
		Subject:     a.Subject,
		Description: a.Description,
		Status:      openapi.TicketStatus(a.Status),
		Priority:    openapi.TicketPriority(a.Priority),
		ResolvedAt:  a.ResolvedAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func toCommentFindResponse(a *aggregate.TicketCommentAggregate) openapi.CommentFindResponse {
	return openapi.CommentFindResponse{
		Uid:        a.Uid,
		TicketUid:  a.TicketUid,
		AuthorUid:  a.AuthorUid,
		AuthorRole: openapi.AuthorRole(a.AuthorRole),
		Body:       a.Body,
		Internal:   a.Internal,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

func toCategoryFindResponse(a *aggregate.TicketCategoryAggregate) openapi.CategoryFindResponse {
	return openapi.CategoryFindResponse{
		Uid:         a.Uid,
		Name:        a.Name,
		Description: a.Description,
		ParentUid:   a.ParentUid,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func toFaqFindResponse(a *aggregate.FaqAggregate) openapi.FaqFindResponse {
	return openapi.FaqFindResponse{
		Uid:         a.Uid,
		CategoryUid: a.CategoryUid,
		Question:    a.Question,
		Answer:      a.Answer,
		ViewCount:   a.ViewCount,
		Published:   a.Published,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}
