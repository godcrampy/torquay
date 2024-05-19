package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godcrampy/torquay/pkg/counter"
)

type Handler struct {
	Counter *counter.Counter
}

func NewHandler(c *counter.Counter) *Handler {
	return &Handler{Counter: c}
}

func (h *Handler) GetToken(c *gin.Context) {
	newValue, err := h.Counter.GetAndIncrement()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newValue})
}
