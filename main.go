package main

import (
	"github.com/ShynaliOtep/paydal-api/controllers"
	_ "github.com/ShynaliOtep/paydal-api/docs"
	"github.com/ShynaliOtep/paydal-api/middlewares"
	"github.com/ShynaliOtep/paydal-api/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"time"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BearerAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	models.ConnectDataBase()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET, POST, PUT, PATCH, DELETE, OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	public := r.Group("/api")
	// тіркелген кезде 100 бонус беріледі
	public.POST("/register", controllers.ByPhoneRegister)
	public.POST("/login", controllers.Login)
	// Обработчик для обновления access токена с использованием refresh токена
	public.GET("/refresh-token", controllers.RefreshTokenHandler)

	protected := r.Group("/api/p/")
	protected.Use(middlewares.JwtAuthMiddleware())

	protected.GET("/user", controllers.CurrentUser)
	protected.PUT("/user/avatar/upload", controllers.UserUploadAvatar)
	protected.GET("/cities", controllers.GetCities)

	protected.POST("/volunteer/contract/create", controllers.VolunteerContractCreate)
	// волонтер тіркелген кезде 1000 бонус беріледі
	protected.POST("/police/volunteer/approve", controllers.VolunteerContractApprove)
	protected.GET("/volunteer/contract/status", controllers.VolunteerContractStatus)

	protected.GET("/violation/categories", controllers.GetViolationCategories)
	protected.POST("/violation/create", controllers.ViolationApplicationCreate)
	protected.POST("/violation/upload/media", controllers.ViolationUploadMedia)

	protected.GET("/violation/:id", controllers.GetViolation)
	public.GET("/violations", controllers.GetViolations)
	protected.GET("/police/violations", controllers.GetViolationsForPolice)
	protected.GET("/police/violation/:id", controllers.GetViolationForPolice)
	protected.PUT("/violation/reject", controllers.ViolationReject)

	protected.POST("/police/protocol/create", controllers.ProtocolCreate)

	public.GET("/penalty/search", controllers.ProtocolSearch)

	public.GET("/article/search", controllers.ArticleSearch)

	protected.GET("/volunteer/violation/applications", controllers.GerVolunteerViolations)
	protected.GET("/police/protocols", controllers.GetPoliceProtocols)

	//url := ginSwagger.URL("http://localhost:8080/docs/swagger.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve uploaded files
	r.Static("/uploads", "./uploads")

	r.Run(":8080")

}
