package repo

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

type WorkRepo interface {
	FindAll(ctx context.Context, opt option.WorkQueryOption) (int64, []*WorkResult, error)
	FindOne(ctx context.Context, id int64) (*WorkResult, error)
	FindByUid(ctx context.Context, uid string) (*WorkResult, error)
	Create(ctx context.Context, cmd command.WorkCreateCommand) (int64, error)
	Update(ctx context.Context, id int64, cmd command.WorkUpdateCommand) error
}

type WorkSearchRepo interface {
	FindAll(ctx context.Context, opt option.WorkQueryOption) (int64, []*WorkSearchResult, error)
	Search(ctx context.Context, opt option.WorkSearchOption) (int64, []*WorkSearchResult, error)
}

type EditionRepo interface {
	FindAll(ctx context.Context, opt option.EditionQueryOption) (int64, []*EditionResult, error)
	FindOne(ctx context.Context, id int64) (*EditionResult, error)
	FindByUid(ctx context.Context, uid string) (*EditionResult, error)
	Create(ctx context.Context, cmd command.EditionCreateCommand, workId int64, publisherId *int64) (int64, error)
	Update(ctx context.Context, id int64, cmd command.EditionUpdateCommand, workId *int64, publisherId *int64) error
}

type AuthorRepo interface {
	FindAll(ctx context.Context, opt option.AuthorQueryOption) (int64, []*AuthorResult, error)
	FindOne(ctx context.Context, id int64) (*AuthorResult, error)
	FindByUid(ctx context.Context, uid string) (*AuthorResult, error)
	Create(ctx context.Context, cmd command.AuthorCreateCommand) (int64, error)
	Update(ctx context.Context, id int64, cmd command.AuthorUpdateCommand) error
}

type PublisherRepo interface {
	FindAll(ctx context.Context, opt option.PublisherQueryOption) (int64, []*PublisherResult, error)
	FindOne(ctx context.Context, id int64) (*PublisherResult, error)
	FindByUid(ctx context.Context, uid string) (*PublisherResult, error)
	Create(ctx context.Context, cmd command.PublisherCreateCommand) (int64, error)
	Update(ctx context.Context, id int64, cmd command.PublisherUpdateCommand) error
}

type SubjectRepo interface {
	FindAll(ctx context.Context, opt option.SubjectQueryOption) (int64, []*SubjectResult, error)
	FindOne(ctx context.Context, id int64) (*SubjectResult, error)
	FindByUid(ctx context.Context, uid string) (*SubjectResult, error)
	FindOrCreateByName(ctx context.Context, name string) (int64, error)
	Create(ctx context.Context, cmd command.SubjectCreateCommand) (int64, error)
	Update(ctx context.Context, id int64, cmd command.SubjectUpdateCommand) error
}
