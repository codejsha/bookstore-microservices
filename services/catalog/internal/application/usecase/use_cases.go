package usecase

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

type CatalogUseCase interface {
	// ─── Work lifecycle ─────────────────────────────────────────────────
	SearchWorks(ctx context.Context, opt option.WorkQueryOption) (int64, []*aggregate.WorkAggregate, error)
	FullTextSearchWorks(ctx context.Context, opt option.WorkSearchOption) (int64, []*repo.WorkSearchResult, error)
	FindWork(ctx context.Context, uid string) (*aggregate.WorkAggregate, error)
	FindWorkWithRelated(ctx context.Context, uid string) (*aggregate.WorkAggregate, []*repo.WorkSearchResult, error)
	CreateWork(ctx context.Context, cmd command.WorkCreateCommand) (*aggregate.WorkAggregate, error)
	UpdateWork(ctx context.Context, uid string, cmd command.WorkUpdateCommand) (*aggregate.WorkAggregate, error)

	// ─── Edition lifecycle ──────────────────────────────────────────────
	SearchEditions(ctx context.Context, opt option.EditionQueryOption) (int64, []*aggregate.EditionAggregate, error)
	FindEdition(ctx context.Context, uid string) (*aggregate.EditionAggregate, error)
	CreateEdition(ctx context.Context, cmd command.EditionCreateCommand) (*aggregate.EditionAggregate, error)
	UpdateEdition(ctx context.Context, uid string, cmd command.EditionUpdateCommand) (*aggregate.EditionAggregate, error)

	// ─── Author management ──────────────────────────────────────────────
	FindAllAuthors(ctx context.Context, opt option.AuthorQueryOption) (int64, []*aggregate.AuthorAggregate, error)
	FindAuthor(ctx context.Context, uid string) (*aggregate.AuthorAggregate, error)
	CreateAuthor(ctx context.Context, cmd command.AuthorCreateCommand) (*aggregate.AuthorAggregate, error)
	UpdateAuthor(ctx context.Context, uid string, cmd command.AuthorUpdateCommand) (*aggregate.AuthorAggregate, error)

	// ─── Publisher management ───────────────────────────────────────────
	FindAllPublishers(ctx context.Context, opt option.PublisherQueryOption) (int64, []*aggregate.PublisherAggregate, error)
	FindPublisher(ctx context.Context, uid string) (*aggregate.PublisherAggregate, error)
	CreatePublisher(ctx context.Context, cmd command.PublisherCreateCommand) (*aggregate.PublisherAggregate, error)
	UpdatePublisher(ctx context.Context, uid string, cmd command.PublisherUpdateCommand) (*aggregate.PublisherAggregate, error)

	// ─── Subject management ─────────────────────────────────────────────
	FindAllSubjects(ctx context.Context, opt option.SubjectQueryOption) (int64, []*aggregate.SubjectAggregate, error)
	FindSubject(ctx context.Context, uid string) (*aggregate.SubjectAggregate, error)
	CreateSubject(ctx context.Context, cmd command.SubjectCreateCommand) (*aggregate.SubjectAggregate, error)
	UpdateSubject(ctx context.Context, uid string, cmd command.SubjectUpdateCommand) (*aggregate.SubjectAggregate, error)
}
