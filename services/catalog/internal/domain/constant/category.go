package constant

import (
	"sort"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/categorypb"
)

type CategoryType string

// General
const (
	CategoryTypeArt         CategoryType = "Art"
	CategoryTypeBusiness    CategoryType = "Business"
	CategoryTypeChildren    CategoryType = "Children"
	CategoryTypeComic       CategoryType = "Comic"
	CategoryTypeCookbook    CategoryType = "Cookbook"
	CategoryTypeEducation   CategoryType = "Education"
	CategoryTypeEngineering CategoryType = "Engineering"
	CategoryTypeFiction     CategoryType = "Fiction"
	CategoryTypeHealth      CategoryType = "Health"
	CategoryTypeHistory     CategoryType = "History"
	CategoryTypeMagazine    CategoryType = "Magazine"
	CategoryTypeNovel       CategoryType = "Novel"
	CategoryTypePoetry      CategoryType = "Poetry"
	CategoryTypeScience     CategoryType = "Science"
	CategoryTypeSports      CategoryType = "Sports"
	CategoryTypeTechnology  CategoryType = "Technology"
	CategoryTypeTravel      CategoryType = "Travel"
)

// Computer Science
const (
	CategoryTypeAgile                CategoryType = "Agile"
	CategoryTypeSoftwareArchitecture CategoryType = "Software Architecture"
	CategoryTypeDataEngineering      CategoryType = "Data Engineering"
	CategoryTypeCodingPractices      CategoryType = "Coding Practices"
	CategoryTypeRefactoring          CategoryType = "Refactoring"
	CategoryTypeKubernetes           CategoryType = "Kubernetes"
)

type Category struct {
	Id   int64        `json:"id"`
	Name CategoryType `json:"name"`
}

func (c Category) GetAllCategories() []Category {
	items := []Category{
		{Id: 1, Name: CategoryTypeArt},
		{Id: 2, Name: CategoryTypeBusiness},
		{Id: 3, Name: CategoryTypeChildren},
		{Id: 4, Name: CategoryTypeComic},
		{Id: 5, Name: CategoryTypeCookbook},
		{Id: 6, Name: CategoryTypeEducation},
		{Id: 7, Name: CategoryTypeEngineering},
		{Id: 8, Name: CategoryTypeFiction},
		{Id: 9, Name: CategoryTypeHealth},
		{Id: 10, Name: CategoryTypeHistory},
		{Id: 11, Name: CategoryTypeMagazine},
		{Id: 12, Name: CategoryTypeNovel},
		{Id: 13, Name: CategoryTypePoetry},
		{Id: 14, Name: CategoryTypeScience},
		{Id: 15, Name: CategoryTypeSports},
		{Id: 16, Name: CategoryTypeTechnology},
		{Id: 17, Name: CategoryTypeTravel},
		{Id: 18, Name: CategoryTypeAgile},
		{Id: 19, Name: CategoryTypeSoftwareArchitecture},
		{Id: 20, Name: CategoryTypeDataEngineering},
		{Id: 21, Name: CategoryTypeCodingPractices},
		{Id: 22, Name: CategoryTypeRefactoring},
		{Id: 23, Name: CategoryTypeKubernetes},
	}
	return items
}

func (c Category) FindAllCategories(sortBy, orderBy string) []Category {
	items := c.GetAllCategories()
	if sortBy == "id" {
		sort.Slice(items, func(i, j int) bool {
			if orderBy == "asc" {
				return items[i].Id < items[j].Id
			}
			return items[i].Id > items[j].Id
		})
	}
	if sortBy == "name" {
		sort.Slice(items, func(i, j int) bool {
			if orderBy == "asc" {
				return items[i].Name < items[j].Name
			}
			return items[i].Name > items[j].Name
		})
	}
	return items
}

func (c Category) GetCategoryById(id int64) Category {
	items := c.GetAllCategories()
	for _, item := range items {
		if item.Id == id {
			return item
		}
	}
	return Category{}
}

func (c Category) ToCategoryFindWebResp() openapi.CategoryFindWebResp {
	resp := openapi.CategoryFindWebResp{
		Id:   c.Id,
		Name: string(c.Name),
	}
	return resp
}

func (c Category) ToCategoryFindProtoResp() *categorypb.CategoryFindProtoResp {
	resp := &categorypb.CategoryFindProtoResp{
		Id:   c.Id,
		Name: string(c.Name),
	}
	return resp
}
