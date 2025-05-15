package bootstrap

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mysql/v2"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/gofiber/storage/sqlite3"
	"go.uber.org/zap"
	"my-web-template/internal/config"
	"my-web-template/internal/core/appcontext"
	"my-web-template/internal/logging"
	"my-web-template/internal/model"
	"my-web-template/internal/repository"
	"my-web-template/internal/service"
	"my-web-template/internal/web/controller"
	"my-web-template/internal/web/middleware"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

// AppComponents 包含所有初始化和组装好的应用组件
// 方便在 bootstrap 包内部传递，或者如果 Run() 函数需要返回这些以便进行测试或进一步操作
type AppComponents struct {
	Config         *config.AppConfig
	Logger         *zap.SugaredLogger
	DBEngine       *xorm.Engine
	WebApp         *fiber.App
	SessionStore   *session.Store
	Validator      *validator.Validate
	UserRepo       repository.UserRepositoryInterface
	UserService    service.UserServiceInterface
	BaseController *controller.AppBaseController
	UserController *controller.UserController
}

// Run 函数负责整个应用的初始化、组装和启动
func Run() error {
	// 1. 解析命令行参数
	cfgFilePath, err := parseCliArgs()
	if err != nil {
		return fmt.Errorf("解析命令行参数失败：%+v", err)
	}

	// 2. 加载配置文件
	appConfig, err := config.LoadConfig(cfgFilePath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败：%+v", err)
	}

	// 3. 初始化日志
	if err := logging.InitLogger(appConfig.Debug, "app.log"); err != nil {
		return fmt.Errorf("初始化日志失败: %w", err)
	}
	logger := logging.Sugar()
	defer func() { _ = logger.Sync() }()
	logger.Info("配置文件加载完毕，日志系统初始化完成.")

	// 4. 连接数据库
	dbEngine, err := initDatabase(appConfig, true)
	if err != nil {
		logger.Errorf("连接数据库失败: %v", err)
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	logger.Info("数据库连接成功.")

	// 5. 初始化 fiber App
	webApp := fiber.New(fiber.Config{
		AppName:           "appname",
		BodyLimit:         10 * 1024 * 1024,
		ReadBufferSize:    10 * 1024 * 1024,
		EnablePrintRoutes: appConfig.Debug,
		//ErrorHandler: func(c *fiber.Ctx, err error) error {
		//	code := fiber.StatusInternalServerError
		//	var e *fiber.Error
		//	if errors.As(err, &e) {
		//		code = e.Code
		//	}
		//	logger.Errorf("全局错误处理器捕获错误 - Path: %s, Error: %v", c.Path(), err)
		//	return c.Status(code).JSON(fiber.Map{"error": "服务器内部错误", "details": err.Error()}) // 生产环境可隐藏 details
		//},
	})

	// 6. 初始化核心 appcontext
	appcontext.Initialize(appConfig, dbEngine, webApp, logger)

	// 7. 初始化其他组件：session、validate
	sessionStore, err := initAppSession(appConfig)
	if err != nil {
		return fmt.Errorf("初始化 session 失败: %w", err)
	}
	logger.Infof("session 初始化成功")
	validate := validator.New()
	logger.Infof("validate 初始化成功")

	// 8. 依赖注入、组装
	userRepo := repository.NewUserRepository(dbEngine, logger)
	userService := service.NewUserService(userRepo, logger)
	baseController := controller.NewAppBaseController(validate, sessionStore)
	userController := controller.NewUserController(logger, baseController, userService)
	logger.Debugf("依赖注入完成")

	// 9. 组装组件
	components := &AppComponents{
		Config:         appConfig,
		Logger:         logger,
		DBEngine:       dbEngine,
		WebApp:         webApp,
		SessionStore:   sessionStore,
		Validator:      validate,
		UserRepo:       userRepo,
		UserService:    userService,
		BaseController: baseController,
		UserController: userController,
	}

	// 10. 配置 web 和路由
	setupWebApp(components)

	// (可选) 如果有其他的需要跑在后台的任务，可以在这里添加

	// 11. 启动 Web 服务
	listenAddr := "127.0.0.1:3000"
	if strings.TrimSpace(appConfig.Web.ListenAddr) != "" {
		listenAddr = strings.TrimSpace(appConfig.Web.ListenAddr)
	}
	logger.Infof("启动 web 服务，监听地址: %s", listenAddr)
	return webApp.Listen(listenAddr)
}

// parseCliArgs 解析命令行参数
func parseCliArgs() (string, error) {
	cli := kingpin.New("<AppName>", "<AppHelp>")                                                 // 替换为你的应用名和帮助信息
	cfgFile := cli.Flag("config", "config file path").Short('c').Default("config.toml").String() // 改为 String(), 在后面检查文件是否存在
	cli.HelpFlag.Short('h')
	// 解析但不立即退出，允许主调函数处理错误
	_, err := cli.Parse(os.Args[1:])
	if err != nil {
		return "", err
	}

	// 检查文件是否存在
	if _, err := os.Stat(*cfgFile); os.IsNotExist(err) {
		return "", fmt.Errorf("配置文件 '%s' 不存在", *cfgFile)
	}
	return *cfgFile, nil
}

// initDatabase 初始化数据库
func initDatabase(appConfig *config.AppConfig, syncTable bool) (*xorm.Engine, error) {
	dsn := ""
	dbCfg := appConfig.Database

	// 构造 DSN
	if appConfig.Database.Driver == "sqlite3" {
		dsn = fmt.Sprintf("./%s", dbCfg.Database)
	} else if appConfig.Database.Driver == "mysql" {
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
			dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Database,
		)
	} else if appConfig.Database.Driver == "postgres" {
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbCfg.Host, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.Database,
		)
	} else {
		return nil, fmt.Errorf("不支持的数据库驱动: %s", appConfig.Database.Driver)
	}

	// 创建xorm引擎
	engine, err := xorm.NewEngine(dbCfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("创建 XORM 引擎失败，错误: %+v", err)
	}

	// 设置 XORM 日志 (可选，但推荐)
	// 如果你希望 XORM 的日志也使用 zap，你需要一个适配器
	// engine.SetLogger(NewXormZapLogger(logger)) // 假设你有这样的适配器
	// engine.SetLogger(xorm.NewSimpleLogger(os.Stdout)) // 或者使用 xorm 默认的简单日志

	// 测试数据库连接
	if err := engine.Ping(); err != nil {
		return nil, fmt.Errorf("ping 数据库失败: %w", err)
	}

	// 是否显示sql
	engine.ShowSQL(dbCfg.ShowSQL)

	// set name mapper
	engine.SetMapper(names.GonicMapper{})

	// 同步表结构
	if syncTable {
		err = engine.Sync(
			// TODO 指定要同步的表
			new(model.AppUserModel),
		)
		if err != nil {
			return nil, fmt.Errorf("同步表结构失败，错误: %+v", err)
		}
	}
	return engine, nil
}

// initAppSession 初始化session
func initAppSession(appConfig *config.AppConfig) (*session.Store, error) {
	var storage fiber.Storage
	switch appConfig.Database.Driver {
	case "sqlite3":
		storage = sqlite3.New(sqlite3.Config{Database: appConfig.Database.Database})
	case "mysql":
		storage = mysql.New(mysql.Config{
			Host: appConfig.Database.Host, Port: appConfig.Database.Port,
			Database: appConfig.Database.Database, Username: appConfig.Database.Username,
			Password: appConfig.Database.Password, Table: "fiber_storage", // 建议指定表名
		})
	case "postgres":
		storage = postgres.New(postgres.Config{
			Host: appConfig.Database.Host, Port: appConfig.Database.Port,
			Database: appConfig.Database.Database, Username: appConfig.Database.Username,
			Password: appConfig.Database.Password, Table: "fiber_storage", // 建议指定表名
		})
	default:
		return nil, fmt.Errorf("session storage 初始化失败，不支持的数据库类型 %s", appConfig.Database.Driver)
	}

	sessionConfig := session.ConfigDefault
	sessionConfig.Expiration = 7 * 24 * time.Hour
	sessionConfig.Storage = storage
	// sessionConfig.KeyLookup = "cookie:your_app_session_id" // 最好指定一个唯一的 session key
	// sessionConfig.CookieSecure = !appConfig.Debug // HTTPS only in production
	// sessionConfig.CookieHTTPOnly = true
	// sessionConfig.CookieSameSite = "Lax"

	return session.New(sessionConfig), nil
}

func setupWebApp(components *AppComponents) {
	// 核心中间件
	components.WebApp.Use(recover.New(recover.Config{EnableStackTrace: components.Config.Debug}))
	components.WebApp.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// 设置 access log 中间件
	accessFile := path.Join(logging.GetExecPath(), "logs", "access.log")
	accessLogFile, err := os.OpenFile(accessFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("打开 access.log 失败，错误 %s", err))
	}
	components.WebApp.Use(fiberLogger.New(fiberLogger.Config{
		Output: io.MultiWriter(os.Stdout, accessLogFile),
		Format: "[${time}] ${ip}:${port} ${status} - ${latency} ${method} ${path} Error: ${error}\n",
	}))

	// 健康检查路由
	components.WebApp.Get("/status", func(ctx *fiber.Ctx) error { return ctx.SendString("ok") })

	// API 路由组
	apiGroup := components.WebApp.Group("/api")

	// 中间件
	permissionMW := middleware.PermissionMiddleware(components.UserService)
	// 如果需要给中间件动态传递参数，可以使用
	// loginCheckMiddleware := func(requireAdmin bool) func(ctx *fiber.Ctx) error {
	//		return middleware.LoginCheckMiddleware(a.baseController.SessionStore, userService, requireAdmin)
	//	}
	// 另外一种注册中间件方法
	// app.WebApp.Use("/api/user", permissionMiddleware)

	// 设置每个 controller 模块的路由
	components.UserController.SetupRouter(apiGroup, permissionMW)

	// TODO 设置前端项目
	//app.WebApp.Use("/", filesystem.New(filesystem.Config{
	//	Root:       http.FS(fe.Dist),
	//	PathPrefix: "dist",
	//	Browse:     true,
	//}))

	// TODO 如果是内嵌前端资源，还需要额外设置静态文件路由，确保前端页面刷新时不会 404
	//app.WebApp.Use(rewrite.New(rewrite.Config{
	//	Rules: map[string]string{
	//		"/user/*": "/",
	//		// ....
	//	},
	//}))
}
