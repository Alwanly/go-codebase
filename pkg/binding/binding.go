package binding

import (
	"net/http"
	"reflect"

	"github.com/Alwanly/go-codebase/pkg/contract"
	"github.com/Alwanly/go-codebase/pkg/logger"
	"github.com/Alwanly/go-codebase/pkg/middleware"
	"github.com/Alwanly/go-codebase/pkg/wrapper"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const ContextName = "Binding"

type (
	Binder struct {
		l   *zap.Logger
		ctx *fiber.Ctx
		m   interface{}
	}

	Source func(*Binder) error

	ModelBindingError struct {
		Code         int
		ResponseBody wrapper.JSONResult
	}
)

func (e *ModelBindingError) Error() string {
	return "Failed to bind request body"
}

func BindFromBody() Source {
	return func(b *Binder) error {
		if err := b.ctx.BodyParser(b.m); err != nil {
			b.l.Debug("Error when binding from body", zap.Error(err))
			return err
		}

		return nil
	}
}

func BindFromQuery() Source {
	return func(b *Binder) error {
		if err := b.ctx.QueryParser(b.m); err != nil {
			b.l.Debug("Error when binding from query string", zap.Error(err))
			return err
		}

		return nil
	}
}

func BindFromParams() Source {
	return func(b *Binder) error {
		if err := b.ctx.ParamsParser(b.m); err != nil {
			b.l.Debug("Error when binding from path params", zap.Error(err))
			return err
		}

		return nil
	}
}

func BindFromHeaders() Source {
	return func(b *Binder) error {
		if err := b.ctx.ReqHeaderParser(b.m); err != nil {
			b.l.Debug("Error when binding from request headers", zap.Error(err))
			return err
		}

		return nil
	}
}

func BindModel(log *zap.Logger, c *fiber.Ctx, m interface{}, sources ...Source) error {
	// create local logger
	l := logger.WithID(log, ContextName, "BindModel")

	// create binder instance
	binder := &Binder{
		l:   l,
		ctx: c,
		m:   m,
	}

	// process data binding
	for _, source := range sources {
		// execute binding
		if err := source(binder); err != nil {
			result := wrapper.ResponseFailed(http.StatusBadRequest, contract.StatusCodeBindingFailed, contract.ErrorValidatePayload, nil)
			return &ModelBindingError{
				Code:         result.Code,
				ResponseBody: result,
			}
		}
	}

	// check if the target has AuthUserData field and set it
	if authUser, ok := c.Locals(middleware.LocalTokenKey).(*middleware.AuthUserData); ok {
		dataField := reflect.Indirect(reflect.ValueOf(m)).FieldByName("AuthUserData")
		if dataField.IsValid() && dataField.CanSet() {
			dataField.Set(reflect.ValueOf(authUser))
		}
	}

	return nil
}
