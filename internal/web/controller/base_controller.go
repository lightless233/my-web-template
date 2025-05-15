package controller

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"my-web-template/internal/constant"
	"my-web-template/internal/result"
)

type AppBaseController struct {
	validator    *validator.Validate
	sessionStore *session.Store
}

func NewAppBaseController(v *validator.Validate, s *session.Store) *AppBaseController {
	return &AppBaseController{
		validator:    v,
		sessionStore: s,
	}
}

func (c *AppBaseController) parseAndValidateBody(ctx *fiber.Ctx, request interface{}) result.AppError {
	if err := ctx.BodyParser(request); err != nil {
		return result.NewAppErrorFromError(constant.CodeParamError, err)
	}
	if err := c.validator.Struct(request); err != nil {
		return result.NewAppErrorFromError(constant.CodeParamError, err)
	}

	c.trimStringField(request)

	return nil
}

func (c *AppBaseController) parseAndValidateQuery(ctx *fiber.Ctx, request interface{}) result.AppError {
	if err := ctx.QueryParser(request); err != nil {
		return result.NewAppErrorFromError(constant.CodeParamError, err)
	}
	if err := c.validator.Struct(request); err != nil {
		return result.NewAppErrorFromError(constant.CodeParamError, err)
	}

	c.trimStringField(request)

	return nil
}

// trimStringField 通过反射，将结构体中的 string 字段去掉前后空格
func (c *AppBaseController) trimStringField(request interface{}) {
	v := reflect.ValueOf(request)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			if field.CanSet() {
				field.SetString(strings.TrimSpace(field.String()))
			}
		}
	}
}
