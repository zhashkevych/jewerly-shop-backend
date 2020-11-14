package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/service"
	"net/http"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init() *gin.Engine {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	// todo move to config
	config.AllowHeaders = append(config.AllowHeaders, "Access-Control-Request-Headers", "Authorization", "X-Forwarded-For",
		"Host", "User-Agent", "Accept")
	router.Use(cors.New(config))

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initPublicRoutes(router)
	h.initAPIRoutes(router)
	h.initAdminRoutes(router)

	return router
}

func (h *Handler) initPublicRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/admin/sign-in", h.adminSignIn)
	}

	payment := router.Group("/payment")
	{
		payment.POST("/callback", h.callback)
	}
}

func (h *Handler) initAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		products := api.Group("/products")
		{
			products.GET("", h.getAllProducts)
			products.GET("/:id", h.getProduct)
		}

		api.POST("/order", h.placeOrder)

		api.GET("/settings", h.getSettings)
	}
}

func (h *Handler) initAdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin", h.adminIdentity)
	{
		admin.POST("/products", h.createProduct)
		admin.GET("/products", h.getAllProducts)
		admin.GET("/products/:id", h.getProduct)
		admin.PUT("/products/:id", h.updateProduct)
		admin.DELETE("/products/:id", h.deleteProduct)

		admin.GET("/orders", h.getAllOrders)
		admin.GET("/orders/:id", h.getOrder)

		settings := admin.Group("/settings")
		{
			settings.GET("/homepage/images", h.getHomepageImages)
			settings.POST("/homepage/image", h.createHomepageImage)
			settings.PUT("/homepage/image/:id", h.updateHomepageImage)

			settings.GET("/text-blocks", h.getTextBlocks)
			settings.POST("/text-block", h.createTextBlock)
			settings.GET("/text-block/:id", h.getTextBlockById)
			settings.PUT("/text-block/:id", h.updateTextBlock)
		}

		admin.POST("/upload", h.uploadImage)
	}
}
