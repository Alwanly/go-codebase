package repository

import (
	"context"
	"fmt"

	"github.com/Alwanly/go-codebase/internal/example/schema"
	"github.com/Alwanly/go-codebase/model"
	"github.com/Alwanly/go-codebase/pkg/database"
	"github.com/Alwanly/go-codebase/pkg/redis"
	"github.com/Alwanly/go-codebase/pkg/utils"
)

const ContextName = "Internal.User.Repository"

type (
	Repository struct {
		DB    database.IDBService
		Redis redis.IRedisService
	}

	IRepository interface {
		Create(context.Context, *model.Book) error
		Get(context.Context, string) *model.Book
		List(context.Context, schema.RequestBookList) ([]model.Book, int64)
		Update(context.Context, *model.Book) error
		Delete(context.Context, string) error
	}
)

func NewRepository(r Repository) IRepository {
	return &Repository{
		DB:    r.DB,
		Redis: r.Redis,
	}
}

func (r *Repository) Create(ctx context.Context, book *model.Book) error {
	return r.DB.GetTransaction(ctx).Create(book).Error
}

func (r *Repository) Get(ctx context.Context, id string) *model.Book {
	var book model.Book
	r.DB.GetTransaction(ctx).Where("id = ?", id).First(&book)
	return &book
}

func (r *Repository) List(ctx context.Context, req schema.RequestBookList) ([]model.Book, int64) {
	var books []model.Book
	var total int64
	tx := r.DB.GetTransaction(ctx)

	tx.Model(&model.Book{}).Count(&total)

	offset := utils.CalculatePageSkip(req.Page, req.PageSize)
	tx.Offset(offset).Limit(req.PageSize)
	tx.Order(fmt.Sprintf("%s %s", req.SortBy, req.SortOrder))
	tx.Find(&books)

	return books, total
}

func (r *Repository) Update(ctx context.Context, book *model.Book) error {
	return r.DB.GetTransaction(ctx).Save(book).Error
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	return r.DB.GetTransaction(ctx).Where("id = ?", id).Delete(&model.Book{}).Error
}
