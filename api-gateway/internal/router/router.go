package router

import (
	"api-gateway/config"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"api-gateway/pkg/constants"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func SetupRouter(
	gatewayHandler *handler.GatewayHandler,
	cfg *config.Config,
	logger *zap.Logger,
	redisClient *redis.Client,
) *gin.Engine {
	router := gin.New()
	// Add recovery middleware
	router.Use(gin.Recovery())
	// CRITICAL: Custom CORS middleware MUST be first
	router.Use(middleware.CORSMiddleware(&cfg.CORS, logger))
	// Strip any client-supplied internal headers before auth
	router.Use(middleware.StripInternalHeaders(logger))
	// Skip logging OPTIONS requests (CORS preflight) to reduce noise
	router.Use(middleware.SkipOptionsLoggingMiddleware(logger))
	// Request logging middleware
	router.Use(middleware.RequestLoggingMiddleware(logger))
	router.Use(middleware.ErrorLoggingMiddleware(logger))
	// Rate limiting middleware
	router.Use(middleware.RateLimitMiddleware(&cfg.RateLimit, logger))

	router.GET("/swagger-ui/*any", func(c *gin.Context) {
		path := c.Param("any")
		if path == "/" || path == "/index.html" || path == "" {
			c.File("docs/swagger-index.html")
			return
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("gateway"))(c)
	})

	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger-ui/")
	})

	// Swagger UI configuration endpoint (returns gateway-scoped spec URLs)
	// The gateway will serve service specs at /specs/{short}.json which
	// are fetched server-side and rewritten so that runtime calls go to the gateway.
	router.GET("/swagger-config", func(c *gin.Context) {
		urls := make([]gin.H, 0, len(cfg.Services))
		for key := range cfg.Services {
			short := strings.TrimSuffix(key, "_service")
			short = strings.ReplaceAll(short, "_", "-")
			// Point UI to gateway's spec proxy endpoint
			specURL := fmt.Sprintf("/specs/%s.json", short)
			name := strings.Title(short) + " Service"
			urls = append(urls, gin.H{"url": specURL, "name": name})
		}
		c.JSON(http.StatusOK, gin.H{"urls": urls})
	})

	// Gateway-side spec proxy endpoint. Returns the service's swagger.json
	// with `servers` rewritten to point at the gateway so the UI and Try-it-out
	// calls go through the gateway.
	// Single wildcard route to handle both `/specs/service` and `/specs/service.json`
	router.GET("/specs/*service", gatewayHandler.SpecProxy)

	router.POST("/auth/session", gatewayHandler.SetSessionCookie)
	router.GET("/health", gatewayHandler.HealthCheck)
	router.GET("/api/gateway/health", gatewayHandler.HealthCheck)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			products := v1.Group("/products")
			{
				products.GET("", gatewayHandler.ProxyRequest)
				products.GET("/:id", gatewayHandler.ProxyRequest)
				products.GET("/search", gatewayHandler.ProxyRequest)
				products.GET("/:id/items", gatewayHandler.ProxyRequest)
				products.GET("/:id/items/:item_id", gatewayHandler.ProxyRequest)
				products.GET("/:id/variations", gatewayHandler.ProxyRequest)
				products.POST("", gatewayHandler.ProxyRequest)

				// Protected: RoleCookieRouter → Auth → Session
				protected := products.Group("")
				protected.Use(
					middleware.RoleCookieRouter(logger),
					middleware.AuthMiddleware(&cfg.JWT, logger),
					middleware.SessionMiddleware(logger, redisClient),
				)
				{
					protected.PUT("/:id", gatewayHandler.ProxyRequest)
					protected.PATCH("/:id", gatewayHandler.ProxyRequest)
					protected.PATCH("/:id/inventory", gatewayHandler.ProxyRequest)
					protected.DELETE("/:id", gatewayHandler.ProxyRequest)
					protected.POST("/:id/items", gatewayHandler.ProxyRequest)
					protected.PUT("/:id/items/:item_id", gatewayHandler.ProxyRequest)
					protected.DELETE("/:id/items/:item_id", gatewayHandler.ProxyRequest)
				}
			}
			categories := v1.Group("/categories")
			{
				categories.GET("", gatewayHandler.ProxyRequest)
				categories.GET("/:id", gatewayHandler.ProxyRequest)
				categories.GET("/slug/:slug", gatewayHandler.ProxyRequest)
				categories.GET("/:id/children", gatewayHandler.ProxyRequest)
				categories.GET("/:id/products", gatewayHandler.ProxyRequest)

				// Protected: RoleCookieRouter → Auth → Session
				protected := categories.Group("")
				protected.Use(
					middleware.RoleCookieRouter(logger),
					middleware.AuthMiddleware(&cfg.JWT, logger),
					middleware.SessionMiddleware(logger, redisClient),
				)
				{
					protected.POST("", middleware.RoleMiddleware(logger, constants.RoleAdmin, constants.RoleBuyer), gatewayHandler.ProxyRequest)
					protected.PUT(":id", gatewayHandler.ProxyRequest)
					protected.PATCH(":id/active", gatewayHandler.ProxyRequest)
					protected.DELETE(":id", gatewayHandler.ProxyRequest)
				}
				admin := categories.Group("/admin")
				admin.Use(
					middleware.RoleCookieRouter(logger),
					middleware.AuthMiddleware(&cfg.JWT, logger),
					middleware.SessionMiddleware(logger, redisClient),
					middleware.RoleMiddleware(logger, constants.RoleAdmin),
				)
				{
					admin.GET("", gatewayHandler.ProxyRequest)
					admin.GET("/parents", gatewayHandler.ProxyRequest)
					admin.POST("/parent", gatewayHandler.ProxyRequest)
				}
			}
			search := v1.Group("/search")
			{
				search.GET("", gatewayHandler.ProxyRequest)
			}
			cart := v1.Group("/cart")
			cart.Use(
				middleware.RoleCookieRouter(logger),
				middleware.AuthMiddleware(&cfg.JWT, logger),
				middleware.SessionMiddleware(logger, redisClient),
			)
			{
				cart.GET("", gatewayHandler.ProxyRequest)
				cart.DELETE("", gatewayHandler.ProxyRequest)
				cart.POST("/items", gatewayHandler.ProxyRequest)
				cart.PUT("/items/:product_item_id", gatewayHandler.ProxyRequest)
				cart.DELETE("/items/:product_item_id", gatewayHandler.ProxyRequest)
				cart.POST("/items/:product_item_id/toggle", gatewayHandler.ProxyRequest)
				cart.PATCH("/items/:product_item_id/selection", gatewayHandler.ProxyRequest)
				cart.PATCH("/selection", gatewayHandler.ProxyRequest)
				cart.PATCH("/shops/:shop_id/selection", gatewayHandler.ProxyRequest)
				cart.DELETE("/selected", gatewayHandler.ProxyRequest)
				cart.POST("/validate", gatewayHandler.ProxyRequest)
			}
			orders := v1.Group("/orders")
			orders.Use(
				middleware.RoleCookieRouter(logger),
				middleware.AuthMiddleware(&cfg.JWT, logger),
				middleware.SessionMiddleware(logger, redisClient),
			)
			{
				orders.POST("", gatewayHandler.ProxyRequest)
				orders.GET("", gatewayHandler.ProxyRequest)
				orders.GET("/:id", gatewayHandler.ProxyRequest)
				orders.GET("/number/:order_number", gatewayHandler.ProxyRequest)
			}

			// Auth endpoints: RoleCookieRouterLax (validate header, don't require cookie)
			auth := v1.Group("/auth")
			auth.Use(middleware.RoleCookieRouterLax(logger))
			{
				auth.POST("/register", gatewayHandler.ProxyRequest)
				auth.POST("/login", gatewayHandler.ProxyRequest)
				auth.POST("/refresh", gatewayHandler.ProxyRequest)
			}

			// Logout: requires full auth pipeline
			authProtected := v1.Group("/auth")
			authProtected.Use(
				middleware.RoleCookieRouter(logger),
				middleware.AuthMiddleware(&cfg.JWT, logger),
				middleware.SessionMiddleware(logger, redisClient),
			)
			{
				authProtected.POST("/logout", gatewayHandler.ProxyRequest)
			}

			// Protected identity routes
			protectedIdentity := v1.Group("")
			protectedIdentity.Use(
				middleware.RoleCookieRouter(logger),
				middleware.AuthMiddleware(&cfg.JWT, logger),
				middleware.SessionMiddleware(logger, redisClient),
			)
			{
				users := protectedIdentity.Group("/users")
				{
					users.GET("/profile", gatewayHandler.ProxyRequest)
					users.PUT("/profile", gatewayHandler.ProxyRequest)
					users.PUT("/password", gatewayHandler.ProxyRequest)
				}

				addresses := protectedIdentity.Group("/addresses")
				{
					addresses.GET("", gatewayHandler.ProxyRequest)
					addresses.POST("", gatewayHandler.ProxyRequest)
					addresses.GET("/:id", gatewayHandler.ProxyRequest)
					addresses.PUT("/:id", gatewayHandler.ProxyRequest)
					addresses.DELETE("/:id", gatewayHandler.ProxyRequest)
					addresses.PUT("/:id/default", gatewayHandler.ProxyRequest)
				}
			}
		}
	}

	return router
}
