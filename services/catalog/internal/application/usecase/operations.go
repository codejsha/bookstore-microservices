package usecase

import (
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/aggregate/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/condition"
)

type (
	AuthorCommand       interface{ isAuthorCommand() }
	AuthorCreateCommand struct{ Ent entity.AuthorEntity }
	AuthorUpdateCommand struct{ Ent entity.AuthorEntity }
	AuthorDeleteCommand struct{ Id int64 }
)

func (AuthorCreateCommand) isAuthorCommand() {}
func (AuthorUpdateCommand) isAuthorCommand() {}
func (AuthorDeleteCommand) isAuthorCommand() {}

type (
	AuthorQuery        interface{ isAuthorQuery() }
	AuthorFindAllQuery struct{ Cond condition.AuthorCondition }
	AuthorFindQuery    struct{ Id int64 }
)

func (AuthorFindAllQuery) isAuthorQuery() {}
func (AuthorFindQuery) isAuthorQuery()    {}

type (
	BookCommand       interface{ isBookCommand() }
	BookCreateCommand struct{ Ent entity.BookEntity }
	BookUpdateCommand struct{ Ent entity.BookEntity }
	BookDeleteCommand struct{ Id int64 }
)

func (BookCreateCommand) isBookCommand() {}
func (BookUpdateCommand) isBookCommand() {}
func (BookDeleteCommand) isBookCommand() {}

type (
	BookQuery        interface{ isBookQuery() }
	BookFindAllQuery struct{ Cond condition.BookCondition }
	BookFindQuery    struct{ Id int64 }
)

func (BookFindAllQuery) isBookQuery() {}
func (BookFindQuery) isBookQuery()    {}

type (
	PublisherCommand       interface{ isPublisherCommand() }
	PublisherCreateCommand struct{ Ent entity.PublisherEntity }
	PublisherUpdateCommand struct{ Ent entity.PublisherEntity }
	PublisherDeleteCommand struct{ Id int64 }
)

func (PublisherCreateCommand) isPublisherCommand() {}
func (PublisherUpdateCommand) isPublisherCommand() {}
func (PublisherDeleteCommand) isPublisherCommand() {}

type (
	PublisherQuery        interface{ isPublisherQuery() }
	PublisherFindAllQuery struct{ Cond condition.PublisherCondition }
	PublisherFindQuery    struct{ Id int64 }
)

func (PublisherFindAllQuery) isPublisherQuery() {}
func (PublisherFindQuery) isPublisherQuery()    {}

type (
	CategoryQuery        interface{ isCategoryQuery() }
	CategoryFindAllQuery struct{ Cond condition.CategoryCondition }
	CategoryFindQuery    struct{ Id int64 }
)

func (CategoryFindAllQuery) isCategoryQuery() {}
func (CategoryFindQuery) isCategoryQuery()    {}
