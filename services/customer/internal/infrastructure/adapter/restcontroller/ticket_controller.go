package restcontroller

import (
	"context"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

var _ openapi.TicketApi = (*ticketController)(nil)

type ticketController struct {
	supportUseCase usecase.SupportUseCase
}

func NewTicketController(supportUseCase usecase.SupportUseCase) openapi.TicketApi {
	return &ticketController{supportUseCase: supportUseCase}
}

func (c *ticketController) TicketsGetAll(
	ctx context.Context,
	_ openapi.TicketApiCredentials,
	customerUid *string,
	assigneeUid *string,
	categoryUid *string,
	status *openapi.TicketStatus,
	priority *openapi.TicketPriority,
	size *int32,
	page *int32,
	sort *string,
) (openapi.TicketsGetAllResponse, error) {
	caller, err := resolveSupportCaller(ctx)
	if err != nil {
		return nil, err
	}

	opts := []option.TicketQueryOptionFunc{
		option.WithTicketPage(pagination.NewPageOption(size, page, sort)),
	}
	if caller.staff {
		if customerUid != nil {
			opts = append(opts, option.WithTicketCustomerUid(*customerUid))
		}
	} else {
		opts = append(opts, option.WithTicketCustomerUid(caller.uid))
	}
	if assigneeUid != nil {
		opts = append(opts, option.WithTicketAssigneeUid(*assigneeUid))
	}
	if categoryUid != nil {
		opts = append(opts, option.WithTicketCategoryUid(*categoryUid))
	}
	if status != nil {
		opts = append(opts, option.WithTicketStatus(constant.TicketStatus(*status)))
	}
	if priority != nil {
		opts = append(opts, option.WithTicketPriority(constant.TicketPriority(*priority)))
	}

	total, tickets, err := c.supportUseCase.FindAllTickets(ctx, option.NewTicketQueryOption(opts...))
	if err != nil {
		return nil, err
	}

	items := make([]openapi.TicketFindResponse, len(tickets))
	for i, t := range tickets {
		items[i] = toTicketFindResponse(t)
	}
	return openapi.TicketsGetAllOk{Body: openapi.TicketFindAllResponse{Items: items, Total: total}}, nil
}

func (c *ticketController) TicketsCreate(
	ctx context.Context,
	_ openapi.TicketApiCredentials,
	req openapi.TicketCreateRequest,
) (openapi.TicketsCreateResponse, error) {
	caller, err := resolveSupportCaller(ctx)
	if err != nil {
		return nil, err
	}

	customerUid := caller.uid
	if caller.staff {
		customerUid = req.CustomerUid
	}
	priority := constant.TICKET_PRIORITY_MEDIUM
	if req.Priority != nil {
		priority = constant.TicketPriority(*req.Priority)
	}
	cmd := command.TicketCreateCommand{
		CustomerUid: customerUid,
		Subject:     req.Subject,
		Description: req.Description,
		Priority:    priority,
		CategoryUid: req.CategoryUid,
	}

	ticket, err := c.supportUseCase.CreateTicket(ctx, cmd)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.TicketsCreateCreated{Location: "/api/v1/tickets/" + ticket.Uid}, nil
}

func (c *ticketController) loadAccessibleTicket(ctx context.Context, caller supportCaller, uid string) (*aggregate.TicketAggregate, error) {
	ticket, err := c.supportUseCase.FindTicket(ctx, uid)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	if !caller.staff && !caller.owns(ticket.CustomerUid) {
		return nil, httpx.MapForbidden(ctx, httpx.ErrForbidden)
	}
	return ticket, nil
}

func (c *ticketController) TicketsRead(
	ctx context.Context,
	_ openapi.TicketApiCredentials,
	uid string,
) (openapi.TicketsReadResponse, error) {
	caller, err := resolveSupportCaller(ctx)
	if err != nil {
		return nil, err
	}
	ticket, err := c.loadAccessibleTicket(ctx, caller, uid)
	if err != nil {
		return nil, err
	}
	return openapi.TicketsReadOk{Body: toTicketFindResponse(ticket)}, nil
}

func (c *ticketController) TicketsUpdate(
	ctx context.Context,
	_ openapi.TicketApiCredentials,
	uid string,
	req openapi.TicketUpdateRequest,
) (openapi.TicketsUpdateResponse, error) {
	caller, err := resolveSupportCaller(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := c.loadAccessibleTicket(ctx, caller, uid); err != nil {
		return nil, err
	}
	if !caller.staff && req.AssigneeUid != nil {
		return nil, httpx.MapForbidden(ctx, httpx.ErrForbidden)
	}

	cmd := command.TicketUpdateCommand{
		Subject:     req.Subject,
		Description: req.Description,
		CategoryUid: req.CategoryUid,
		AssigneeUid: req.AssigneeUid,
	}
	if req.Priority != nil {
		priority := constant.TicketPriority(*req.Priority)
		cmd.Priority = &priority
	}

	ticket, err := c.supportUseCase.UpdateTicket(ctx, uid, cmd)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.TicketsUpdateOk{Body: toTicketFindResponse(ticket)}, nil
}

func (c *ticketController) TicketsUpdateStatus(
	ctx context.Context,
	_ openapi.TicketApiCredentials,
	uid string,
	req openapi.TicketUpdateStatusRequest,
) (openapi.TicketsUpdateStatusResponse, error) {
	caller, err := resolveSupportCaller(ctx)
	if err != nil {
		return nil, err
	}
	if !caller.staff {
		return nil, httpx.MapForbidden(ctx, httpx.ErrForbidden)
	}

	cmd := command.TicketStatusCommand{Status: constant.TicketStatus(req.Status)}
	ticket, err := c.supportUseCase.UpdateTicketStatus(ctx, uid, cmd)
	if err != nil {
		return nil, mapSupportError(ctx, err)
	}
	return openapi.TicketsUpdateStatusOk{Body: toTicketFindResponse(ticket)}, nil
}
