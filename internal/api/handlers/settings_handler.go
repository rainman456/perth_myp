package handlers

import (
	"api-customer-merchant/internal/services/settings"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SettingsHandler struct {
	settingsService *settings.SettingsService
	logger          *zap.Logger
}

func NewSettingsHandler(settingsService *settings.SettingsService, logger *zap.Logger) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		logger:          logger,
	}
}

// GetSettings retrieves the global settings
// @Summary Get marketplace settings
// @Description Retrieves global marketplace settings including fees, tax, and shipping
// @Tags Settings
// @Produce json
// @Success 200 {object} models.Settings
// @Failure 500 {object} object{error=string}
// @Router /settings [get]
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	ctx := c.Request.Context()

	settings, err := h.settingsService.GetSettings(ctx)
	if err != nil {
		h.logger.Error("Failed to get settings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}