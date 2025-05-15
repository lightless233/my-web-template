package appcontext

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"my-web-template/internal/config"
	"xorm.io/xorm"
)

var globalAppContext *Context

// Context 应用上下文，只包含核心组件
type Context struct {
	AppConfig *config.AppConfig
	DBEngine  *xorm.Engine
	WebApp    *fiber.App
	Logger    *zap.SugaredLogger
}

// Initialize 初始化全局应用上下文。
// 此函数应在应用启动时由 main 函数调用。
func Initialize(cfg *config.AppConfig, db *xorm.Engine, webApp *fiber.App, logger *zap.SugaredLogger) {
	if globalAppContext != nil {
		logger.Warn("ApplicationContext already initialized")
		return
	}
	globalAppContext = &Context{
		AppConfig: cfg,
		DBEngine:  db,
		WebApp:    webApp,
		Logger:    logger,
	}
}

// Get 获取全局 application context 变量，如果未初始化，则会 panic。
func Get() *Context {
	if globalAppContext == nil {
		panic("ApplicationContext not initialized. Call appcontext.Initialize() first.")
	}
	return globalAppContext
}
