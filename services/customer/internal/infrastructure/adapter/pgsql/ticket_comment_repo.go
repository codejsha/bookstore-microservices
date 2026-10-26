package pgsql

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/codejsha/shared-library-go/pkg/database"

	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/entity"
	genrepo "github.com/codejsha/bookstore-microservices/customer/generated/infrastructure/port/repo"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/repo"
)

var _ repo.TicketCommentRepo = (*ticketCommentRepository)(nil)

type ticketCommentRepository struct {
	q *dao.Query
	genrepo.CustomerTicketCommentRepo
}

func NewTicketCommentRepository(dataSource *database.DataSource) repo.TicketCommentRepo {
	db := dataSource.DB()
	return &ticketCommentRepository{
		q:                         dao.Use(db),
		CustomerTicketCommentRepo: genrepo.NewCustomerTicketCommentRepo(db),
	}
}

func (r *ticketCommentRepository) FindAllByTicketId(ctx context.Context, ticketId int64, includeInternal bool) ([]*repo.TicketCommentResult, error) {
	c := r.q.CustomerTicketCommentEntity
	q := c.WithContext(ctx).Where(c.TicketId.Eq(ticketId))
	if !includeInternal {
		q = q.Where(c.Internal.Is(false))
	}
	entities, err := q.Order(c.CreatedAt, c.Id).Find()
	if err != nil {
		return nil, err
	}
	results := make([]*repo.TicketCommentResult, len(entities))
	for i, e := range entities {
		results[i] = toTicketCommentResult(e)
	}
	return results, nil
}

func (r *ticketCommentRepository) Insert(ctx context.Context, comment repo.TicketCommentCreate) (*repo.TicketCommentResult, error) {
	now := time.Now()
	ent := &entity.CustomerTicketCommentEntity{
		Uid:        uuid.Must(uuid.NewV7()).String(),
		TicketId:   comment.TicketId,
		AuthorUid:  comment.AuthorUid,
		AuthorRole: comment.AuthorRole,
		Body:       comment.Body,
		Internal:   comment.Internal,
		CreatedAt:  now,
		UpdatedAt:  &now,
		Version:    1,
	}
	if err := r.q.CustomerTicketCommentEntity.WithContext(ctx).Create(ent); err != nil {
		return nil, err
	}
	return toTicketCommentResult(ent), nil
}

func toTicketCommentResult(e *entity.CustomerTicketCommentEntity) *repo.TicketCommentResult {
	return &repo.TicketCommentResult{
		Id:         e.Id,
		Uid:        e.Uid,
		TicketId:   e.TicketId,
		AuthorUid:  e.AuthorUid,
		AuthorRole: e.AuthorRole,
		Body:       e.Body,
		Internal:   e.Internal,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}
