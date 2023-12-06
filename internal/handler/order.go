package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/43nvy/wb_l0"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createOrder(ctx *gin.Context) {
	var input wb_l0.Order
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.WbOrder.CreateOrder(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.services.NotifyNewOrder(id)

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

var (
	cache      = make(map[int]wb_l0.Order)
	cacheMutex sync.Mutex
)

func (h *Handler) cacheMiddleware(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if order, ok := cache[id]; ok {
		fmt.Println("Cache hit!")
		ctx.JSON(http.StatusOK, order)
		return
	}

	fmt.Println("Cache miss!")

	ctx.Next()

	order, exists := ctx.Get("order")
	if !exists {
		return
	}

	cachedOrder, ok := order.(wb_l0.Order)
	if !ok {
		return
	}

	cache[id] = cachedOrder
	fmt.Println("Cached order!")

	ctx.JSON(http.StatusOK, cachedOrder)
}

func (h *Handler) getOrderById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	order, err := h.services.WbOrder.GetOrder(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.Set("order", order)

	ctx.JSON(http.StatusOK, order)
}

func (h *Handler) getTenOrders(ctx *gin.Context) {
	orders, err := h.services.WbOrder.GetTenOrders()
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (h *Handler) populateCache() {
	orders, err := h.services.WbOrder.GetTenOrders()
	if err != nil {
		fmt.Println("Failed to populate cache:", err)
		return
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	for _, order := range orders {
		cache[order.ID] = order
	}
}
