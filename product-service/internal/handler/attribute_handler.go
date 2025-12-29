package handler

import (
	"net/http"
	"product-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AttributeHandler handles HTTP requests for EAV attribute operations
type AttributeHandler struct {
	attributeService *service.AttributeService
	logger           *zap.Logger
}

// NewAttributeHandler creates a new attribute handler
func NewAttributeHandler(attributeService *service.AttributeService, logger *zap.Logger) *AttributeHandler {
	return &AttributeHandler{
		attributeService: attributeService,
		logger:           logger,
	}
}

// CreateCategoryAttribute godoc
// @Summary Create a category attribute
// @Description Create a new attribute for a category (e.g. RAM, Màn hình for Điện thoại category)
// @Tags attributes
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Param attribute body service.CreateCategoryAttributeRequest true "Attribute info"
// @Success 201 {object} domain.CategoryAttribute
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{category_id}/attributes [post]
func (h *AttributeHandler) CreateCategoryAttribute(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
		return
	}

	var req service.CreateCategoryAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set category_id from path
	req.CategoryID = uint(categoryID)

	attr, err := h.attributeService.CreateCategoryAttribute(&req)
	if err != nil {
		h.logger.Error("failed to create category attribute", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attr)
}

// GetCategoryAttributes godoc
// @Summary Get category attributes
// @Description Get all attributes for a category
// @Tags attributes
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {array} domain.CategoryAttribute
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{category_id}/attributes [get]
func (h *AttributeHandler) GetCategoryAttributes(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
		return
	}

	attrs, err := h.attributeService.GetCategoryAttributes(uint(categoryID))
	if err != nil {
		h.logger.Error("failed to get category attributes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get attributes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attributes": attrs,
		"count":      len(attrs),
	})
}

// SetProductAttributes godoc
// @Summary Set product attributes
// @Description Set attribute values for a product (replaces all existing values)
// @Tags attributes
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param attributes body service.SetProductAttributesRequest true "Attributes map"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{product_id}/attributes [post]
func (h *AttributeHandler) SetProductAttributes(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	var req service.SetProductAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.attributeService.SetProductAttributes(uint(productID), &req); err != nil {
		h.logger.Error("failed to set product attributes", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product attributes set successfully"})
}

// GetProductAttributes godoc
// @Summary Get product attributes
// @Description Get all attribute values for a product
// @Tags attributes
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{product_id}/attributes [get]
func (h *AttributeHandler) GetProductAttributes(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	attrs, err := h.attributeService.GetProductAttributes(uint(productID))
	if err != nil {
		h.logger.Error("failed to get product attributes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get attributes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attributes": attrs,
	})
}

// DeleteCategoryAttribute godoc
// @Summary Delete category attribute
// @Description Delete a category attribute
// @Tags attributes
// @Produce json
// @Param category_id path int true "Category ID"
// @Param attr_id path int true "Attribute ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{category_id}/attributes/{attr_id} [delete]
func (h *AttributeHandler) DeleteCategoryAttribute(c *gin.Context) {
	attrID, err := strconv.ParseUint(c.Param("attr_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attr_id"})
		return
	}

	if err := h.attributeService.DeleteCategoryAttribute(uint(attrID)); err != nil {
		h.logger.Error("failed to delete category attribute", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete attribute"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category attribute deleted successfully"})
}

