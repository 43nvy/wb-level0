package handler

import (
	"github.com/43nvy/wb_l0/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	order := router.Group("/message")
	{
		order.POST("/", h.createOrder)
		order.GET("/", h.getTenOrders)
		order.GET("/:id", h.cacheMiddleware, h.getOrderById)
	}

	router.GET("/sub", h.subscribeClient)
	router.GET("/unsub", h.unsubscribeClient)

	return router
}

func (h *Handler) InitCache() {
	go h.populateCache()
}
