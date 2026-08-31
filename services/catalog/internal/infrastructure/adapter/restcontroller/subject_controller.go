package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/catalog/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/httpx"
)

var _ openapi.SubjectApi = (*subjectController)(nil)

type subjectController struct {
	catalogUseCase usecase.CatalogUseCase
}

func NewSubjectController(catalogUseCase usecase.CatalogUseCase) openapi.SubjectApi {
	return &subjectController{catalogUseCase: catalogUseCase}
}

func (c *subjectController) SubjectsGetAll(
	ctx context.Context,
	name *string,
	size *int32,
	page *int32,
	sort *string,
) (*openapi.SubjectFindAllResponse, error) {
	opt := option.NewSubjectQueryOption(
		option.SubjectQueryOption{}.WithName(name),
		option.SubjectQueryOption{}.WithPage(pagination.NewPageOption(size, page, sort)),
	)

	total, subjects, err := c.catalogUseCase.FindAllSubjects(ctx, opt)
	if err != nil {
		return nil, err
	}

	items := make([]openapi.SubjectItem, len(subjects))
	for i, s := range subjects {
		items[i] = toSubjectItem(s)
	}
	return &openapi.SubjectFindAllResponse{Items: items, Total: total}, nil
}

func (c *subjectController) SubjectsCreate(ctx context.Context, req openapi.SubjectCreateRequest) error {
	cmd := command.SubjectCreateCommand{
		Name: req.Name,
	}

	subject, err := c.catalogUseCase.CreateSubject(ctx, cmd)
	if err != nil {
		return httpx.MapConflict(ctx, httpx.MapBusinessError(ctx, err))
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "subject", subject.Uid, "created", logrus.Fields{"name": req.Name})
	return nil
}

func (c *subjectController) SubjectsRead(ctx context.Context, uid string) (*openapi.SubjectFindResponse, error) {
	subject, err := c.catalogUseCase.FindSubject(ctx, uid)
	if err != nil {
		return nil, httpx.MapNotFound(ctx, err)
	}

	return &openapi.SubjectFindResponse{
		Uid:  subject.Uid,
		Name: subject.Name,
	}, nil
}

func (c *subjectController) SubjectsUpdate(
	ctx context.Context,
	uid string,
	req openapi.SubjectUpdateRequest,
) (*openapi.SubjectUpdateResponse, error) {
	cmd := command.SubjectUpdateCommand{
		Name: req.Name,
	}

	subject, err := c.catalogUseCase.UpdateSubject(ctx, uid, cmd)
	if err != nil {
		return nil, httpx.MapConflict(ctx, httpx.MapBusinessError(ctx, httpx.MapNotFound(ctx, err)))
	}

	dispatchSideEffects(context.WithoutCancel(ctx), "subject", uid, "updated", logrus.Fields{})

	return &openapi.SubjectUpdateResponse{
		Uid:  subject.Uid,
		Name: subject.Name,
	}, nil
}
