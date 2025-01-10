package main

import (
	"hysteria2-panel/config"
	"hysteria2-panel/database"
	"hysteria2-panel/handlers"
	"hysteria2-panel/middleware"
	"hysteria2-panel/services"

	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
	Config *config.Config
}

func main() {
	server := &Server{
		Router: gin.Default(),
		Config: loadConfig(),
	}

	// 初始化数据库
	initDB(server)

	// 设置路由
	setupRoutes(server)

	// 启动服务器
	server.Router.Run(server.Config.Listen)
}

func loadConfig() *config.Config {
	config, err := config.LoadConfig("configs/config.json")
	if err != nil {
		panic(err)
	}
	return config
}

func initDB(server *Server) {
	db, err := database.InitDB(server.Config)
	if err != nil {
		panic(err)
	}
	server.DB = db
}

func setupRoutes(server *Server) {
	// 创建服务实例
	userService := services.NewUserService(server.DB)
	userManager := services.NewUserManagerService(server.DB)
	configManager := services.NewConfigManagerService(server.DB)
	hy2Service := services.NewHysteria2Service(
		"configs/hysteria2",
		server.Config.TLSCertPath,
		server.Config.TLSKeyPath,
	)
	trafficService := services.NewTrafficService(server.DB)
	nodeService := services.NewNodeService(server.DB)
	nodeHandler := handlers.NewNodeHandler(nodeService)
	settingService := services.NewSettingService(server.DB)
	mailService := services.NewMailService(settingService)
	certService := services.NewCertService(settingService, "certs")
	planService := services.NewPlanService(server.DB)
	planHandler := handlers.NewPlanHandler(planService)
	notificationService := services.NewNotificationService(server.DB, mailService)
	paymentService := services.NewPaymentService(server.DB, settingService, planService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// 创建处理器
	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userManager)
	configHandler := handlers.NewConfigHandler(configManager)
	hy2Handler := handlers.NewHysteria2Handler(configManager, hy2Service)
	trafficHandler := handlers.NewTrafficHandler(trafficService)

	// 用户认证相关路由
	auth := server.Router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	// 需要认证的API路由
	api := server.Router.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		// 用户管理
		api.GET("/users", userHandler.GetUsers)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)

		// 配置管理
		api.GET("/configs/:id", configHandler.GetUserConfig)
		api.PUT("/configs/:id", configHandler.UpdateUserConfig)

		// 添加 Hysteria2 配置相关路由
		api.GET("/hy2/server/:id", hy2Handler.GenerateServerConfig)
		api.GET("/hy2/client/:id", hy2Handler.GetClientConfig)

		// 添加流量统计相关路由
		api.POST("/traffic/record", trafficHandler.RecordTraffic)
		api.GET("/traffic/check/:id", trafficHandler.CheckTrafficLimit)

		// 添加节点管理相关路由
		api.POST("/nodes", nodeHandler.CreateNode)
		api.GET("/nodes", nodeHandler.GetNodes)
		api.POST("/nodes/:id/status", nodeHandler.UpdateNodeStatus)
		api.GET("/nodes/:id/status", nodeHandler.GetNodeStatus)

		// 添加系统设置相关路由
		api.GET("/settings/tls", settingHandler.GetTLSConfig)
		api.PUT("/settings/tls", settingHandler.UpdateTLSConfig)
		api.GET("/settings/smtp", settingHandler.GetSMTPConfig)
		api.PUT("/settings/smtp", settingHandler.UpdateSMTPConfig)
		api.GET("/settings/announcement", settingHandler.GetAnnouncement)
		api.PUT("/settings/announcement", settingHandler.UpdateAnnouncement)

		// 添加套餐管理相关路由
		api.POST("/plans", planHandler.CreatePlan)
		api.GET("/plans", planHandler.GetPlans)
		api.PUT("/plans/:id", planHandler.UpdatePlan)
		api.POST("/plans/:id/subscribe", planHandler.Subscribe)
		api.POST("/plans/:id/order", planHandler.CreateOrder)

		// 添加支付相关路由
		api.POST("/payments", paymentHandler.CreatePayment)
		api.GET("/payments/status", paymentHandler.QueryPaymentStatus)
	}

	// 支付回调接口（不需要认证）
	server.Router.POST("/api/callback/:method", paymentHandler.HandleCallback)

	// 初始化证书
	if err := certService.ObtainCert(); err != nil {
		log.Printf("初始化证书失败: %v", err)
	}

	// 启动定时任务
	go func() {
		// 证书更新检查
		certTicker := time.NewTicker(24 * time.Hour)
		// 到期提醒检查
		expirationTicker := time.NewTicker(12 * time.Hour)
		// 流量提醒检查
		trafficTicker := time.NewTicker(6 * time.Hour)

		for {
			select {
			case <-certTicker.C:
				if err := certService.RenewCert(); err != nil {
					log.Printf("更新证书失败: %v", err)
				}
			case <-expirationTicker.C:
				if err := notificationService.SendExpirationNotices(); err != nil {
					log.Printf("发送到期提醒失败: %v", err)
				}
			case <-trafficTicker.C:
				if err := notificationService.SendTrafficNotices(); err != nil {
					log.Printf("发送流量提醒失败: %v", err)
				}
			}
		}
	}()
}
