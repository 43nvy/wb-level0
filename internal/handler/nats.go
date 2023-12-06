package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) subscribeClient(ctx *gin.Context) {
	client_id := ctx.Query("client_id")
	if client_id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, "client_id is required")
		return
	}

	err := h.services.SubscribeClient(client_id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "failed to subscribe client")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "subscribed successfully"})
}

func (h *Handler) unsubscribeClient(ctx *gin.Context) {
	client_id := ctx.Query("client_id")
	if client_id == "" {
		newErrorResponse(ctx, http.StatusBadRequest, "client_id is required")
		return
	}

	h.services.UnsubscribeClient(client_id)

	ctx.JSON(http.StatusOK, gin.H{"message": "unsubscribed successfully"})
}
