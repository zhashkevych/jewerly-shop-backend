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
	config.AllowHeaders = append(config.AllowHeaders, "Access-Control-Request-Headers", "Authorization")
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
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/admin/sign-in", h.adminSignIn)
	}

	payment := router.Group("/payment")
	{
		payment.POST("/callback", h.callback)
	}
}

func (h *Handler) initAPIRoutes(router *gin.Engine) {
	//api := router.Group("/api", h.userIdentity)
	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.GET("/profile", h.getUserProfile)
			user.GET("/orders", h.getUserOrders)
		}

		products := api.Group("/products")
		{
			products.GET("", h.getAllProducts)
			products.GET("/:id", h.getProduct)
		}

		api.POST("/order", h.placeOrder)
	}
}

func (h *Handler) initAdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin", h.adminIdentity)
	{
		// product routes
		admin.POST("/products", h.createProduct)
		admin.GET("/products", h.getAllProducts)
		admin.GET("/products/:id", h.getProduct)
		admin.PUT("/products/:id", h.updateProduct)
		admin.DELETE("/products/:id", h.deleteProduct)
		// orders routes
		admin.GET("/orders", h.getAllOrders)
		admin.GET("/orders/:id", h.getOrder)
		// product images
		admin.POST("/upload", h.uploadImage)
	}
}
