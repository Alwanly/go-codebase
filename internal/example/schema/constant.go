package schema

import (
	"github.com/Alwanly/go-codebase/model"
	"github.com/Alwanly/go-codebase/pkg/middleware"
)

type RequestBookCreate struct {
	Title  string `json:"title" validate:"required,min=3,max=255"`
	Author string `json:"author" validate:"required"`

	AuthUserData *middleware.AuthUserData
}

type ResponseBookCreate struct {
	ID string `json:"id"`
}

type RequestBookGet struct {
	ID string `params:"id" validate:"required"`

	AuthUserData *middleware.AuthUserData
}

type ResponseBookGet struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type RequestBookList struct {
	Page      int    `query:"page" validate:"required,min=1"`
	PageSize  int    `query:"page_size" validate:"required,min=1,max=100"`
	SortBy    string `query:"sort_by" validate:"required,oneof=title author"`
	SortOrder string `query:"sort_order" validate:"required,oneof=asc desc"`

	AuthUserData *middleware.AuthUserData
}

type RequestBookUpdate struct {
	ID string `params:"id" validate:"required"`

	Title  string `json:"title" validate:"required,min=3,max=255"`
	Author string `json:"author" validate:"required"`

	AuthUserData *middleware.AuthUserData
}

type ResponseBookUpdate struct {
	ID string `json:"id"`
}

type RequestBookDelete struct {
	ID string `params:"id" validate:"required"`

	AuthUserData *middleware.AuthUserData
}

type ResponseBookDelete struct{}

func (r *RequestBookList) ToResponse(books []model.Book) []ResponseBookGet {
	responseBooks := make([]ResponseBookGet, len(books))
	for i, book := range books {
		responseBooks[i] = ResponseBookGet{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
		}
	}
	return responseBooks
}
