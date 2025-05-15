package controller

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"my-web-template/internal/core/appcontext"
	"my-web-template/internal/entity/request"
	"my-web-template/internal/result"
	"my-web-template/internal/service"
)

type UserController struct {
	base        *AppBaseController
	userService *service.UserService
	logger      *zap.SugaredLogger
}

func NewUserController(logger *zap.SugaredLogger, base *AppBaseController, userService *service.UserService) *UserController {
	return &UserController{
		logger:      logger,
		base:        base,
		userService: userService,
	}
}

func (u *UserController) Register(ctx *fiber.Ctx) error {
	query := &request.RegisterRequest{}
	if err := u.base.parseAndValidateBody(ctx, query); err != nil {
		return ctx.JSON(err.ToAppResult())
	}

	user, err := u.userService.SaveUser(query.Username, query.Email, query.Password)
	if err != nil {
		return ctx.JSON(err.ToAppResult())
	}

	return ctx.JSON(result.NewSuccessResult(user))
}

func (u *UserController) GetUserInfoByUsername(ctx *fiber.Ctx) error {
	query := &request.GetUserInfoRequest{}
	if err := u.base.parseAndValidateQuery(ctx, query); err != nil {
		return ctx.JSON(err.ToAppResult())
	}

	// 测试获取全局变量和当前类中的日志
	appcontext.Get().Logger.Infof("GetUserInfoByUsername: %s", query.Username)
	u.logger.Infof("GetUserInfoByUsername: %s", query.Username)

	user, err := u.userService.GetUserByUsername(query.Username)
	if err != nil {
		return ctx.JSON(err.ToAppResult())
	}

	return ctx.JSON(result.NewSuccessResult(user))
}

func (u *UserController) SetupRouter(router fiber.Router, permissionMiddleware fiber.Handler) {
	userAPI := router.Group("/user")
	userAPI.Post("/v1/register", u.Register)
	userAPI.Get("/v1/info", permissionMiddleware, u.GetUserInfoByUsername)
}
