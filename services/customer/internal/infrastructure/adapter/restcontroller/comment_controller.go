package restcontroller

import (
	"context"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

var _ openapi.CommentApi = (*commentController)(nil)

type commentController struct {
	supportUseCase usecase.SupportUseCase
}

func NewCommentController(supportUseCase usecase.SupportUseCase) openapi.CommentApi {
	return &commentController{supportUseCase: supportUseCase}
}

func (c *commentController) authorizeTicket(ctx context.Context, caller supportCaller, ticketUid string) error {
	ticket, err := c.supportUseCase.FindTicket(ctx, ticketUid)
	if err != nil {
		return mapSupportError(ctx, err)
	}
	if !caller.staff && !caller.owns(ticket.CustomerUid) {
		return httpx.MapForbidden(ctx, httpx.ErrForbidden)
	}
	return nil
}

func (c *commentController) CommentsGetAll(
	ctx context.Context,
	_ openapi.CommentApiCredentials,
	ticketUid string,
	includeInternal *bool,
) (openapi.CommentsGetAllResponse, error) {
	caller, err := resolveSupportCaller(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.authorizeTicket(ctx, caller, ticketUid); err != nil {
		return nil, err
	}

	withInternal := includeInternal != nil && *includeInternal && caller.staff
	comments, err := c.supportUseCase.FindComments(ctx, ticketUid, withInternal)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}

	items := make([]openapi.CommentFindResponse, len(comments))
	for i, cm := range comments {
		items[i] = toCommentFindResponse(cm)
	}
	return openapi.CommentsGetAllOk{Body: openapi.CommentFindAllResponse{Items: items}}, nil
}

func (c *commentController) CommentsCreate(
	ctx context.Context,
	_ openapi.CommentApiCredentials,
	ticketUid string,
	req openapi.CommentCreateRequest,
) (openapi.CommentsCreateResponse, error) {
	caller, err := resolveSupportCaller(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.authorizeTicket(ctx, caller, ticketUid); err != nil {
		return nil, err
	}

	internal := req.Internal != nil && *req.Internal
	if internal && !caller.staff {
		return nil, httpx.MapForbidden(ctx, httpx.ErrForbidden)
	}
	authorRole := constant.AUTHOR_ROLE_CUSTOMER
	if caller.staff {
		authorRole = constant.AUTHOR_ROLE_AGENT
	}

	cmd := command.CommentAddCommand{
		TicketUid:  ticketUid,
		AuthorUid:  caller.uid,
		AuthorRole: authorRole,
		Body:       req.Body,
		Internal:   internal,
	}
	comment, err := c.supportUseCase.AddComment(ctx, cmd)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.CommentsCreateCreated{Location: "/api/v1/tickets/" + ticketUid + "/comments/" + comment.Uid}, nil
}
