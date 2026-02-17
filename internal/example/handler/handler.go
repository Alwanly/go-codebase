package handler

import (
	"github.com/Alwanly/go-codebase/internal/example/repository"
	"github.com/Alwanly/go-codebase/internal/example/schema"
	"github.com/Alwanly/go-codebase/internal/example/usecase"
	"github.com/Alwanly/go-codebase/pkg/binding"
	"github.com/Alwanly/go-codebase/pkg/deps"
	"github.com/Alwanly/go-codebase/pkg/logger"
	"github.com/Alwanly/go-codebase/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const ContextName = "Internal.Book.Handler"

type (
	Handler struct {
		Logger    *zap.Logger
		Validator validator.IValidatorService
		UseCase   usecase.IUseCase
	}
)

func NewHandler(d *deps.App) *Handler {
	repository := repository.NewRepository(repository.Repository{
		DB:    d.DB,
		Redis: d.Redis,
	})
	usecase := usecase.NewUseCase(usecase.UseCase{
		Config:     d.Config,
		Logger:     d.Logger,
		Repository: repository,
	})
	handler := &Handler{
		Logger:    d.Logger,
		Validator: d.Validator,
		UseCase:   usecase,
	}

	e := d.Fiber.Group("/books/v1", d.Auth.JwtAuth())
	e.Post("/", handler.Create)
	e.Get("/", handler.List)
	e.Get("/:id", handler.Get)
	e.Put("/:id", handler.Update)
	e.Delete("/:id", handler.Delete)
	return handler
}

// Create creates a new book.
func (h *Handler) Create(c *fiber.Ctx) error {
	l := logger.WithID(h.Logger, ContextName, "Create")

	// bind model
	model := &schema.RequestBookCreate{}
	if err := binding.BindModel(l, c, model, binding.BindFromBody()); err != nil {
		perr := err.(*binding.ModelBindingError)
		return c.Status(perr.Code).JSON(perr.ResponseBody)
	}

	// validate model
	if err := validator.ValidateModel(l, h.Validator, model); err != nil {
		perr := err.(*validator.ModelValidationError)
		return c.Status(perr.Code).JSON(perr.ResponseBody)
	}

	// create a new book
	response := h.UseCase.Create(c.UserContext(), model)
	return c.Status(response.Code).JSON(response)
}

// List returns a list of books.
func (h *Handler) List(c *fiber.Ctx) error {
	l := logger.WithID(h.Logger, ContextName, "List")

	// bind model
	model := &schema.RequestBookList{
		Page:      1,
		PageSize:  10,
		SortBy:    "title",
		SortOrder: "desc",
	}
	if err := binding.BindModel(l, c, model, binding.BindFromQuery()); err != nil {
		perr := err.(*binding.ModelBindingError)
		return c.Status(perr.Code).JSON(perr.ResponseBody)
	}

	// validate model
	if err := validator.ValidateModel(l, h.Validator, model); err != nil {
		perr := err.(*validator.ModelValidationError)
		return c.Status(perr.Code).JSON(perr.ResponseBody)
	}

	// get list of books
	response := h.UseCase.List(c.UserContext(), model)
	return c.Status(response.Code).JSON(response)
}

// Get returns a book by ID.
func (h *Handler) Get(c *fiber.Ctx) error {
	l := logger.WithID(h.Logger, ContextName, "Get")

	// bind model
	model := &schema.RequestBookGet{}
	if err := binding.BindModel(l, c, model, binding.BindFromParams()); err != nil {
		perr := err.(*binding.ModelBindingError)
		return c.Status(perr.Code).JSON(perr.ResponseBody)
	}

	// validate model
	if err := validator.ValidateModel(l, h.Validator, model); err != nil {
		perr := err.(*validator.ModelValidationError)
		return c.Status(perr.Code).JSON(perr.ResponseBody)
	}

	// get book by ID
	response := h.UseCase.Get(c.UserContext(), model)
	return c.Status(response.Code).JSON(response)
}

// Update updates a book by ID.
func (h *Handler) Update(c *fiber.Ctx) error {
	l := logger.WithID(h.Logger, ContextName, "Update")

	// bind model
	model := &schema.RequestBookUpdate{}
	if err := binding.BindModel(l, c, model, binding.BindFromParams(), binding.BindFromBody()); err != nil {
		perr := err.(*binding.ModelBindingError)
		return c.Status(perr.Code).JSON(perr.ResponseBody)
	}

	// validate model
	if err := validator.ValidateModel(l, h.Validator, model); err != nil {
		perr := err.(*validator.ModelValidationError)
		return c.Status(perr.Code).JSON(perr.ResponseBody)
	}

	// update book by ID
	response := h.UseCase.Update(c.UserContext(), model)
	return c.Status(response.Code).JSON(response)
}

// Delete deletes a book by ID.
func (h *Handler) Delete(c *fiber.Ctx) error {
	l := logger.WithID(h.Logger, ContextName, "Delete")

	// bind model
	model := &schema.RequestBookDelete{}
	if err := binding.BindModel(l, c, model, binding.BindFromParams()); err != nil {
		perr := err.(*binding.ModelBindingError)
		return c.Status(perr.Code).JSON(perr.ResponseBody)
	}

	// validate model
	if err := validator.ValidateModel(l, h.Validator, model); err != nil {
		perr := err.(*validator.ModelValidationError)
		return c.Status(perr.Code).JSON(perr.ResponseBody)
	}

	// delete book by ID
	response := h.UseCase.Delete(c.UserContext(), model)
	return c.Status(response.Code).JSON(response)
}
