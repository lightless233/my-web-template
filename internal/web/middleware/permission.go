package middleware

import (
	"github.com/gofiber/fiber/v2"
	"my-web-template/internal/service"
)

func PermissionMiddleware(userService service.UserServiceInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// 调用服务层的方法
		// userService.GetUserByUsername("admin")

		// 如果想中断当前请求
		// return c.JSON(fiber.Map{"message": "权限不足"})

		// 把数据放到 ctx 中
		// c.Locals("user", userService.GetCurrentUser(c))

		// 继续执行
		return c.Next()
	}
}
