package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/Alwanly/go-codebase/config"
	"github.com/Alwanly/go-codebase/internal/example/repository"
	"github.com/Alwanly/go-codebase/internal/example/schema"
	"github.com/Alwanly/go-codebase/model"
	"github.com/Alwanly/go-codebase/pkg/contract"
	"github.com/Alwanly/go-codebase/pkg/wrapper"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const ContextName = "Internal.Book.Usecase"

type (
	UseCase struct {
		Config     *config.GlobalConfig
		Logger     *zap.Logger
		Repository repository.IRepository
	}

	IUseCase interface {
		Create(context.Context, *schema.RequestBookCreate) wrapper.JSONResult
		Get(context.Context, *schema.RequestBookGet) wrapper.JSONResult
		List(context.Context, *schema.RequestBookList) wrapper.JSONResult
		Update(context.Context, *schema.RequestBookUpdate) wrapper.JSONResult
		Delete(context.Context, *schema.RequestBookDelete) wrapper.JSONResult
	}
)

func NewUseCase(uc UseCase) IUseCase {
	return &UseCase{
		Config:     uc.Config,
		Logger:     uc.Logger,
		Repository: uc.Repository,
	}
}

func (u *UseCase) Create(ctx context.Context, req *schema.RequestBookCreate) wrapper.JSONResult {
	l := u.Logger.With(zap.String("usecase", "Create"))

	// Create a new book
	now := time.Now()
	book := &model.Book{
		ID:        uuid.New().String(),
		Title:     req.Title,
		Author:    req.Author,
		CreatedBy: req.AuthUserData.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.Repository.Create(ctx, book); err != nil {
		l.Error("failed to create a book", zap.Error(err))
		return wrapper.ResponseFailed(500, contract.StatusCodeInternalServerError, "Failed to create a book", nil)
	}

	l.Debug("book created", zap.String("id", book.ID))

	return wrapper.ResponseSuccess(http.StatusCreated, schema.ResponseBookCreate{ID: book.ID})
}

func (u *UseCase) Get(ctx context.Context, req *schema.RequestBookGet) wrapper.JSONResult {
	l := u.Logger.With(zap.String("usecase", "Get"))

	book := u.Repository.Get(ctx, req.ID)

	if book == nil {
		l.Error("book not found", zap.String("id", req.ID))
		return wrapper.ResponseFailed(http.StatusNotFound, contract.CreateStatusCode("0004"), "Book not found", nil)
	}

	return wrapper.ResponseSuccess(http.StatusOK, schema.ResponseBookGet{
		ID:     book.ID,
		Title:  book.Title,
		Author: book.Author,
	})
}

func (u *UseCase) List(ctx context.Context, req *schema.RequestBookList) wrapper.JSONResult {
	l := u.Logger.With(zap.String("usecase", "List"))

	filter := schema.RequestBookList{
		Page:      req.Page,
		PageSize:  req.PageSize,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,

		AuthUserData: req.AuthUserData,
	}
	books, total := u.Repository.List(ctx, filter)

	response := req.ToResponse(books)
	l.Debug("books listed", zap.Int64("total", total))
	return wrapper.ResponsePagination(req.Page, req.PageSize, len(books), int(total), response, nil)
}

func (u *UseCase) Update(ctx context.Context, req *schema.RequestBookUpdate) wrapper.JSONResult {
	l := u.Logger.With(zap.String("usecase", "Update"))

	book := u.Repository.Get(ctx, req.ID)
	if book == nil {
		l.Error("book not found", zap.String("id", req.ID))
		return wrapper.ResponseFailed(http.StatusNotFound, contract.CreateStatusCode("0004"), "Book not found", nil)
	}

	book.Title = req.Title
	book.Author = req.Author
	book.UpdatedAt = time.Now()

	if err := u.Repository.Update(ctx, book); err != nil {
		l.Error("failed to update a book", zap.Error(err))
		return wrapper.ResponseFailed(500, contract.StatusCodeInternalServerError, "Failed to update a book", nil)
	}

	l.Debug("book updated", zap.String("id", book.ID))

	return wrapper.ResponseSuccess(http.StatusOK, schema.ResponseBookUpdate{ID: book.ID})
}

func (u *UseCase) Delete(ctx context.Context, req *schema.RequestBookDelete) wrapper.JSONResult {
	l := u.Logger.With(zap.String("usecase", "Delete"))

	book := u.Repository.Get(ctx, req.ID)
	if book == nil {
		l.Error("book not found", zap.String("id", req.ID))
		return wrapper.ResponseFailed(http.StatusNotFound, contract.CreateStatusCode("0004"), "Book not found", nil)
	}

	if err := u.Repository.Delete(ctx, book.ID); err != nil {
		l.Error("failed to delete a book", zap.Error(err))
		return wrapper.ResponseFailed(500, contract.StatusCodeInternalServerError, "Failed to delete a book", nil)
	}

	l.Debug("book deleted", zap.String("id", book.ID))

	return wrapper.ResponseSuccess(http.StatusNoContent, schema.ResponseBookDelete{})
}
