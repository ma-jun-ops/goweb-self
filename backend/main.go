package main

import (
	"log"
	"time"
	"web/database"
	"web/handler"

	_ "web/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title          my API
// @version         1.0
// @description     这是后端 API 接口文档
// @host            localhost:8080
// @BasePath        /

func main() {
	database.ConnectDB()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/articles", handler.ListArticle)
		api.GET("/articles/:id", handler.GetArticle)
		api.POST("/articles", handler.CreateArticle)
		api.PUT("/articles/:id", handler.UpdateArticle)
		api.DELETE("/articles/:id", handler.DeleteArticle)
		api.POST("/files/upload", handler.UploadFile)
		api.GET("/files/:id/download", handler.DownloadFile)
		api.GET("/files/:id/image", handler.ImageFile)
	}

	log.Println("后端服务启动成功，监听地址: http://localhost:8080")
	r.Run()
}
