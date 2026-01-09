package handler

import (
	"net/http"
	"product-service/internal/domain"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// VariationHandler handles variation-related HTTP requests
type VariationHandler struct {
	variationRepo    domain.VariationRepository
	variationOptRepo domain.VariationOptionRepository
	logger           *zap.Logger
}

// NewVariationHandler creates a new variation handler
func NewVariationHandler(
	variationRepo domain.VariationRepository,
	variationOptRepo domain.VariationOptionRepository,
	logger *zap.Logger,
) *VariationHandler {
	return &VariationHandler{
		variationRepo:    variationRepo,
		variationOptRepo: variationOptRepo,
		logger:           logger,
	}
}

// VariationWithOptions represents a variation with its options
type VariationWithOptions struct {
	ID      uint                     `json:"id"`
	Name    string                   `json:"name"` // "Màu Sắc", "Kích Thước"
	Options []domain.VariationOption `json:"options"`
}

// GetProductVariations godoc
// @Summary Get all variations for a product
// @Description Get variations (Color, Size, etc.) with their options for product detail page
// @Tags variations
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {array} VariationWithOptions
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id}/variations [get]
func (h *VariationHandler) GetProductVariations(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	// Get all variations for product
	variations, err := h.variationRepo.GetByProductID(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get variations"})
		return
	}

	// Build response with options
	var response []VariationWithOptions
	for _, v := range variations {
		options, err := h.variationOptRepo.GetByVariationID(v.ID)
		if err != nil {
			h.logger.Error("Failed to get variation options",
				zap.Uint("variation_id", v.ID),
				zap.Error(err))
			continue
		}

		// Convert []*VariationOption to []VariationOption
		optionsList := make([]domain.VariationOption, len(options))
		for i, opt := range options {
			optionsList[i] = *opt
		}

		response = append(response, VariationWithOptions{
			ID:      v.ID,
			Name:    v.Name,
			Options: optionsList,
		})
	}

	c.JSON(http.StatusOK, response)
}
