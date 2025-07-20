package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go.bdev.cn/shortener/internal/handlers"
)

func NewRouter() *gin.Engine {
	g := gin.Default()

	// swagger api docs
	// g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// pprof router 性能分析路由
	// 默认关闭，开发环境下可以打开
	// 访问方式: HOST/debug/pprof
	// 通过 HOST/debug/pprof/profile 生成profile
	// 查看分析图 go tool pprof -http=:5000 profile
	// see: https://github.com/gin-contrib/pprof
	// pprof.Register(g)

	// favicon.ico
	g.GET("/favicon.ico", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
	})

	account := handlers.Handle.AccountHandler
	user := handlers.Handle.UserHandler
	shortener := handlers.Handle.ShortenHandler
	history := handlers.Handle.HistoryHandler

	//apiV1 := g.Group("/api/v1")
	apiV1 := g.Group("/api")

	// PING
	apiV1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	apiV1.POST("/account/login", account.Login)
	apiV1.Use(authMiddleware())
	{
		apiV1.POST("/shortens", shortener.ShortenAdd)
		apiV1.GET("/shortens", shortener.ShortenList)
		apiV1.DELETE("/shortens", shortener.ShortenDeleteAll)
		apiV1.GET("/shortens/:code", shortener.ShortenFind)
		apiV1.PUT("/shortens/:code", shortener.ShortenUpdate)
		apiV1.DELETE("/shortens/:code", shortener.ShortenDelete)

		apiV1.GET("/histories", history.HistoryList)
		apiV1.DELETE("/histories", history.HistoryDeleteAll)

		apiV1.POST("/account/logout", account.Logout)
		apiV1.GET("/users/current", user.Current)
	}

	// 短链接跳转路由
	g.GET("/:code", shortener.ShortenRedirect)
	g.HEAD("/:code", shortener.ShortenRedirect)

	return g
}
