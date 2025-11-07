package controller

import (
	_ "github.com/Oleska1601/WBURLShortener/docs"

	"github.com/wb-go/wbf/ginext"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupRouter() {
	ginMode := ""
	engine := ginext.New(ginMode)
	engine.Use(CORSMiddleware())
	engine.Static("/static", "./web")

	notifyGroup := engine.Group("/api/")
	{
		notifyGroup.POST("/shorten", s.CreateShortURLHandler)
		notifyGroup.GET("/s/:short_url", s.RedirectShortURLHandler)
		notifyGroup.GET("/analytics/:short_url", s.GetAnalyticsHandler)
	}
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.Srv.Handler = engine
}
