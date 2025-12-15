package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"workpulse/internal/service"
)

type ConfigHandler struct {
	cfgService *service.ConfigService
	basePath   string
}

func NewConfigHandler(cfgService *service.ConfigService, basePath string) *ConfigHandler {
	return &ConfigHandler{cfgService: cfgService, basePath: basePath}
}

func (h *ConfigHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, h.cfgService.ClientConfig(h.basePath))
}
