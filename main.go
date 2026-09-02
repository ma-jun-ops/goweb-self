package main

import (
	"net/http"
	_ "web/docs" // 这个包会在生成文档后自动创建

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
	r := gin.Default()

	// 添加 Swagger 路由（访问 http://localhost:8080/swagger/index.html 查看文档）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// @Summary      获取问候语
	// @Description  根据用户名返回问候语
	// @Param        name  query  string  false  "用户名"  default(World)
	// @Success      200   {object}  map[string]string
	// @Router       /hello [get]
	r.GET("/hello", func(c *gin.Context) {
		username := c.DefaultQuery("name", "World")
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello " + username,
			"status":  "success",
		})
	})

	// @Summary      创建用户
	// @Description  创建一个新用户
	// @Param        user  body  object  true  "用户信息"
	// @Success      200   {object}  map[string]string
	// @Router       /user [post]
	r.POST("/user", func(c *gin.Context) {
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":  "创建成功",
			"username": user.Username,
		})
	})

	r.Run()
}
