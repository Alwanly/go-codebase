# Contributing Guide

Thank you for considering contributing to this Go Codebase project! This guide will help you understand the project structure and how to add new features following our architecture patterns.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Project Architecture](#project-architecture)
4. [Adding New Features](#adding-new-features)
5. [Code Style Guidelines](#code-style-guidelines)
6. [Testing](#testing)
7. [Pull Request Process](#pull-request-process)

## Code of Conduct

- Be respectful and inclusive
- Write clean, maintainable code
- Follow the established coding standards
- Document your code appropriately
- Write tests for new features

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/go-codebase.git`
3. Create a feature branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit with clear messages
7. Push to your fork
8. Submit a pull request

## Project Architecture

This project follows Clean Architecture and Domain-Driven Design (DDD) principles:

### Layers

```
┌─────────────────────────────────────────┐
│           Handler Layer                  │  HTTP handlers (presentation)
├─────────────────────────────────────────┤
│           UseCase Layer                  │  Business logic
├─────────────────────────────────────────┤
│         Repository Layer                 │  Data access
├─────────────────────────────────────────┤
│           Model Layer                    │  Database models
└─────────────────────────────────────────┘
```

### Dependency Flow

- **Handler** → **UseCase** → **Repository** → **Model**
- Higher layers can depend on lower layers
- Lower layers should NOT depend on higher layers

## Adding New Features

### Example: Adding a "Product" Feature

#### Step 1: Create Domain Structure

Create the following directory structure:

```
internal/
└── product/
    ├── handler/
    │   └── handler.go
    ├── repository/
    │   └── repository.go
    ├── usecase/
    │   └── usecase.go
    └── schema/
        ├── request.go
        ├── response.go
        └── constant.go
```

#### Step 2: Define the Model

Create `model/product.go`:

```go
package model

import "time"

type Product struct {
	ID          string    `gorm:"primaryKey;column:id;type:varchar(255);not null"`
	Name        string    `gorm:"column:name;type:varchar(255);not null"`
	Description string    `gorm:"column:description;type:text"`
	Price       float64   `gorm:"column:price;type:decimal(10,2);not null"`
	Stock       int       `gorm:"column:stock;type:integer;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	CreatedBy   string    `gorm:"column:created_by;type:varchar(255);not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	UpdatedBy   string    `gorm:"column:updated_by;type:varchar(255);not null"`
}

func (Product) TableName() string {
	return "products"
}

type Products []Product
```

#### Step 3: Define Schemas

Create `internal/product/schema/request.go`:

```go
package schema

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	Stock       int     `json:"stock" validate:"omitempty,gte=0"`
}
```

Create `internal/product/schema/response.go`:

```go
package schema

import "time"

type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
```

#### Step 4: Create Repository Interface

Create `internal/product/repository/repository.go`:

```go
package repository

import (
	"context"

	"github.com/Alwanly/go-codebase/model"
	"github.com/Alwanly/go-codebase/pkg/database"
	"github.com/Alwanly/go-codebase/pkg/redis"
)

type IRepository interface {
	Create(ctx context.Context, product *model.Product) error
	FindByID(ctx context.Context, id string) (*model.Product, error)
	FindAll(ctx context.Context, limit, offset int) (model.Products, int64, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id string) error
}

type Repository struct {
	DB    *database.DBService
	Redis *redis.Service
}

func NewRepository(r Repository) IRepository {
	return &r
}

func (r *Repository) Create(ctx context.Context, product *model.Product) error {
	return r.DB.Gorm.WithContext(ctx).Create(product).Error
}

func (r *Repository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	err := r.DB.Gorm.WithContext(ctx).Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *Repository) FindAll(ctx context.Context, limit, offset int) (model.Products, int64, error) {
	var products model.Products
	var count int64

	err := r.DB.Gorm.WithContext(ctx).Model(&model.Product{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.DB.Gorm.WithContext(ctx).Limit(limit).Offset(offset).Find(&products).Error
	return products, count, err
}

func (r *Repository) Update(ctx context.Context, product *model.Product) error {
	return r.DB.Gorm.WithContext(ctx).Save(product).Error
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	return r.DB.Gorm.WithContext(ctx).Delete(&model.Product{}, "id = ?", id).Error
}
```

#### Step 5: Create UseCase

Create `internal/product/usecase/usecase.go`:

```go
package usecase

import (
	"context"
	"time"

	"github.com/Alwanly/go-codebase/config"
	"github.com/Alwanly/go-codebase/internal/product/repository"
	"github.com/Alwanly/go-codebase/internal/product/schema"
	"github.com/Alwanly/go-codebase/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IUseCase interface {
	Create(ctx context.Context, req *schema.CreateProductRequest, userID string) (*schema.ProductResponse, error)
	GetByID(ctx context.Context, id string) (*schema.ProductResponse, error)
	GetAll(ctx context.Context, page, limit int) ([]*schema.ProductResponse, int64, error)
	Update(ctx context.Context, id string, req *schema.UpdateProductRequest, userID string) error
	Delete(ctx context.Context, id string) error
}

type UseCase struct {
	Config     *config.GlobalConfig
	Logger     *zap.Logger
	Repository repository.IRepository
}

func NewUseCase(u UseCase) IUseCase {
	return &u
}

func (u *UseCase) Create(ctx context.Context, req *schema.CreateProductRequest, userID string) (*schema.ProductResponse, error) {
	product := &model.Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
		UpdatedAt:   time.Now(),
		UpdatedBy:   userID,
	}

	if err := u.Repository.Create(ctx, product); err != nil {
		u.Logger.Error("Failed to create product", zap.Error(err))
		return nil, err
	}

	return u.toResponse(product), nil
}

func (u *UseCase) GetByID(ctx context.Context, id string) (*schema.ProductResponse, error) {
	product, err := u.Repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u.toResponse(product), nil
}

func (u *UseCase) GetAll(ctx context.Context, page, limit int) ([]*schema.ProductResponse, int64, error) {
	offset := (page - 1) * limit
	products, count, err := u.Repository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*schema.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = u.toResponse(&product)
	}

	return responses, count, nil
}

func (u *UseCase) Update(ctx context.Context, id string, req *schema.UpdateProductRequest, userID string) error {
	product, err := u.Repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	product.UpdatedAt = time.Now()
	product.UpdatedBy = userID

	return u.Repository.Update(ctx, product)
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.Repository.Delete(ctx, id)
}

func (u *UseCase) toResponse(product *model.Product) *schema.ProductResponse {
	return &schema.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
```

#### Step 6: Create Handler

Create `internal/product/handler/handler.go`:

```go
package handler

import (
	"strconv"

	"github.com/Alwanly/go-codebase/internal/product/repository"
	"github.com/Alwanly/go-codebase/internal/product/schema"
	"github.com/Alwanly/go-codebase/internal/product/usecase"
	"github.com/Alwanly/go-codebase/pkg/binding"
	"github.com/Alwanly/go-codebase/pkg/contract"
	"github.com/Alwanly/go-codebase/pkg/deps"
	"github.com/Alwanly/go-codebase/pkg/wrapper"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	Logger  *zap.Logger
	UseCase usecase.IUseCase
}

func NewHandler(d *deps.App) *Handler {
	repo := repository.NewRepository(repository.Repository{
		DB:    d.DB,
		Redis: d.Redis,
	})
	uc := usecase.NewUseCase(usecase.UseCase{
		Config:     d.Config,
		Logger:     d.Logger,
		Repository: repo,
	})

	handler := &Handler{
		Logger:  d.Logger,
		UseCase: uc,
	}

	// Register routes
	e := d.Fiber.Group("/products/v1", d.Auth.JwtAuth())
	e.Post("/", handler.Create)
	e.Get("/", handler.List)
	e.Get("/:id", handler.Get)
	e.Put("/:id", handler.Update)
	e.Delete("/:id", handler.Delete)

	return handler
}

// Create godoc
// @Summary Create product
// @Tags products
// @Accept json
// @Produce json
// @Param request body schema.CreateProductRequest true "Product data"
// @Success 201 {object} wrapper.JSONResult{data=schema.ProductResponse}
// @Router /products/v1/ [post]
// @Security Bearer
func (h *Handler) Create(c *fiber.Ctx) error {
	var req schema.CreateProductRequest
	if err := binding.Bind(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(wrapper.ResponseFailed(
			fiber.StatusBadRequest,
			contract.StatusCodeBadRequest,
			"Invalid request",
			nil,
		))
	}

	// Get user ID from context (set by JWT middleware)
	userID := c.Locals("user_id").(string)

	result, err := h.UseCase.Create(c.Context(), &req, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(wrapper.ResponseFailed(
			fiber.StatusInternalServerError,
			contract.StatusCodeInternalServerError,
			"Failed to create product",
			nil,
		))
	}

	return c.Status(fiber.StatusCreated).JSON(wrapper.ResponseSuccess(fiber.StatusCreated, result))
}

// List godoc
// @Summary List products
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} wrapper.JSONResult{data=[]schema.ProductResponse}
// @Router /products/v1/ [get]
// @Security Bearer
func (h *Handler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	products, count, err := h.UseCase.GetAll(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(wrapper.ResponseFailed(
			fiber.StatusInternalServerError,
			contract.StatusCodeInternalServerError,
			"Failed to fetch products",
			nil,
		))
	}

	return c.JSON(wrapper.ResponsePagination(page, limit, len(products), int(count), products, nil))
}

// Get godoc
// @Summary Get product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} wrapper.JSONResult{data=schema.ProductResponse}
// @Router /products/v1/{id} [get]
// @Security Bearer
func (h *Handler) Get(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := h.UseCase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(wrapper.ResponseFailed(
			fiber.StatusNotFound,
			contract.StatusCodeNotFound,
			"Product not found",
			nil,
		))
	}

	return c.JSON(wrapper.ResponseSuccess(fiber.StatusOK, product))
}

// Update godoc
// @Summary Update product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body schema.UpdateProductRequest true "Product data"
// @Success 200 {object} wrapper.JSONResult
// @Router /products/v1/{id} [put]
// @Security Bearer
func (h *Handler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req schema.UpdateProductRequest

	if err := binding.Bind(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(wrapper.ResponseFailed(
			fiber.StatusBadRequest,
			contract.StatusCodeBadRequest,
			"Invalid request",
			nil,
		))
	}

	userID := c.Locals("user_id").(string)

	if err := h.UseCase.Update(c.Context(), id, &req, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(wrapper.ResponseFailed(
			fiber.StatusInternalServerError,
			contract.StatusCodeInternalServerError,
			"Failed to update product",
			nil,
		))
	}

	return c.JSON(wrapper.ResponseSuccess(fiber.StatusOK, nil))
}

// Delete godoc
// @Summary Delete product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} wrapper.JSONResult
// @Router /products/v1/{id} [delete]
// @Security Bearer
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.UseCase.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(wrapper.ResponseFailed(
			fiber.StatusInternalServerError,
			contract.StatusCodeInternalServerError,
			"Failed to delete product",
			nil,
		))
	}

	return c.JSON(wrapper.ResponseSuccess(fiber.StatusOK, nil))
}
```

#### Step 7: Register Handler in Bootstrap

Edit `cmd/main/bootstrap.go`:

```go
import (
	// ... other imports
	product_handler "github.com/Alwanly/go-codebase/internal/product/handler"
)

// In Bootstrap function:
book_handler.NewHandler(inst)
product_handler.NewHandler(inst) // Add this line
```

#### Step 8: Generate Documentation

```bash
make docs
```

## Code Style Guidelines

### Follow Go Standards

- Use `gofmt` for formatting
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Follow [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `product`, `user`)
- **Files**: lowercase with underscores (e.g., `product_handler.go`)
- **Types**: PascalCase (e.g., `ProductService`)
- **Functions**: camelCase for private, PascalCase for exported
- **Constants**: PascalCase (e.g., `MaxRetries`)

### Documentation

- Document all exported types and functions
- Use godoc format
- Include Swagger annotations for API endpoints

```go
// ProductService handles product-related operations.
// It implements business logic for product management.
type ProductService struct {
    repo IProductRepository
}

// Create creates a new product.
// It validates the input and returns the created product or an error.
func (s *ProductService) Create(ctx context.Context, req *CreateProductRequest) (*Product, error) {
    // Implementation
}
```

### Error Handling

- Always check errors
- Use descriptive error messages
- Log errors with context

```go
if err := repo.Create(ctx, product); err != nil {
    logger.Error("Failed to create product", 
        zap.Error(err),
        zap.String("product_id", product.ID),
    )
    return nil, fmt.Errorf("create product: %w", err)
}
```

## Testing

### Unit Tests

Write unit tests for all business logic:

```go
func TestProductUseCase_Create(t *testing.T) {
    // Setup
    mockRepo := new(MockProductRepository)
    uc := NewProductUseCase(mockRepo)

    // Test
    req := &CreateProductRequest{
        Name:  "Test Product",
        Price: 99.99,
    }

    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    result, err := uc.Create(context.Background(), req, "user123")

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Test Product", result.Name)
    mockRepo.AssertExpectations(t)
}
```

### Run Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run specific package
go test ./internal/product/...
```

## Pull Request Process

1. **Update tests**: Ensure all tests pass
2. **Update documentation**: Update README if needed
3. **Generate mocks**: Run `make mock` if interfaces changed
4. **Generate docs**: Run `make docs` if API changed
5. **Lint code**: Run `make lint`
6. **Format code**: Run `make format`
7. **Write clear commit messages**:
   ```
   feat: add product management endpoints
   
   - Add Product model
   - Implement CRUD operations
   - Add unit tests
   - Generate Swagger docs
   ```

8. **Fill PR template**:
   - Description of changes
   - Related issues
   - Testing done
   - Screenshots (if UI changes)

9. **Wait for review**: Address feedback from maintainers

## Questions?

- Check existing [issues](https://github.com/Alwanly/go-codebase/issues)
- Read the [documentation](../README.md)
- Ask in discussions

Thank you for contributing! 🎉
