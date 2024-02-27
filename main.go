package main

import (
	"funding-app/auth"
	"funding-app/campaign"
	"funding-app/handler"
	"funding-app/helper"
	"funding-app/payment"
	"funding-app/transaction"
	"funding-app/users"
	webHandler "funding-app/web/handler"
	webHelper "funding-app/web/helper"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/fundstartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	userRepository := users.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := users.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewUserCampaign(campaignService)
	transactionHandler := handler.NewTransaction(transactionService)

	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService, userService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	sessionWebHandler := webHandler.NewSessionHandler(userService)

	router := gin.Default()
	router.Use(cors.Default())

	router.HTMLRender = loadTemplates("./web/templates")

	router.Static("/avatar_images", "./avatar_images")
	router.Static("/campaign_images", "./campaign_images")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.LoginUser)
	api.POST("/email_checkers", userHandler.CheckEMailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/create_campaign", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/update_campaign/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetTransactionCampaign)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetTransactionUser)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.GET("/users", authAdminMiddleware(), userWebHandler.Index)
	router.GET("/users/new", authAdminMiddleware(), userWebHandler.NewUser)
	router.POST("/users", authAdminMiddleware(), userWebHandler.CreateUser)
	router.GET("/users/edit/:id", authAdminMiddleware(), userWebHandler.GetUserById)
	router.POST("/users/update/:id", authAdminMiddleware(), userWebHandler.UpdateUser)
	router.GET("/users/avatar/:id", authAdminMiddleware(), userWebHandler.UploadAvatar)
	router.POST("/users/avatar/:id", authAdminMiddleware(), userWebHandler.CreateAvatar)
	router.GET("/campaigns", authAdminMiddleware(), campaignWebHandler.Index)
	router.GET("/campaign/new", authAdminMiddleware(), campaignWebHandler.NewCampaign)
	router.POST("/campaigns", authAdminMiddleware(), campaignWebHandler.CreateNewCampaign)
	router.GET("/campaign/images/:id", authAdminMiddleware(), campaignWebHandler.CampaignNewImage)
	router.POST("/campaign/images/:id", authAdminMiddleware(), campaignWebHandler.CreateCampaignNewImage)
	router.GET("/campaign/edit/:id", authAdminMiddleware(), campaignWebHandler.GetCampaignById)
	router.POST("/campaign/update/:id", authAdminMiddleware(), campaignWebHandler.UpdateCampaign)
	router.GET("/campaign/show/:id", authAdminMiddleware(), campaignWebHandler.ShowCampaignDetail)
	router.GET("/transactions", authAdminMiddleware(), transactionWebHandler.Index)
	router.GET("/login", sessionWebHandler.New)
	router.POST("/session", sessionWebHandler.CreateSession)
	router.GET("/logout", sessionWebHandler.Destroy)
	_ = router.Run(helper.GoDotEnvVariable("PORT"))
}

func authMiddleware(authService auth.Service, service users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			//error *
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		payload, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "failed", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userID := int(payload["user_id"].(float64))
		user, err := service.GetUserById(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "failed", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("currentUser", user)
	}
}

func authAdminMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		session, _ := webHelper.CookieStore.Get(context.Request, "login-admin-session")
		if session.Values["nameUser"] == nil {
			context.Redirect(http.StatusFound, "/login")
			return
		}
	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
