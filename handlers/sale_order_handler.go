package handlers

import (
	"fmt"
	"strconv"
	"time"

	"interview-user/models"
	"interview-user/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SaleOrderHandler struct {
	DB *gorm.DB
}

func NewSaleOrderHandler(db *gorm.DB) *SaleOrderHandler {
	return &SaleOrderHandler{DB: db}
}

type CreateSaleOrderRequest struct {
	CustomerName string                       `json:"customer_name" binding:"required"`
	Notes        string                       `json:"notes"`
	Items        []CreateSaleOrderItemRequest `json:"items" binding:"required,min=1"`
}

type CreateSaleOrderItemRequest struct {
	ProductName string  `json:"product_name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" binding:"required,min=0"`
}

type UpdateSaleOrderRequest struct {
	CustomerName string                       `json:"customer_name"`
	Notes        string                       `json:"notes"`
	Items        []CreateSaleOrderItemRequest `json:"items"`
}

// GetAll returns all sale orders with pagination
func (h *SaleOrderHandler) GetAll(c *gin.Context) {
	pagination := utils.GetPagination(c)

	var total int64
	var orders []models.SaleOrder

	// Count total
	if err := h.DB.Model(&models.SaleOrder{}).Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count sale orders")
		return
	}

	// Get paginated data
	if err := h.DB.Preload("CreatedBy").Preload("SaleOrderItems").
		Order("created_at DESC").
		Limit(pagination.Limit).
		Offset(pagination.GetOffset()).
		Find(&orders).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch sale orders")
		return
	}

	utils.OKResponse(c, "Sale orders retrieved successfully", utils.PaginatedResponse{
		Items:      orders,
		TotalItems: total,
		TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		Page:       pagination.Page,
		Limit:      pagination.Limit,
	})
}

// GetByID returns a sale order by ID
func (h *SaleOrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid sale order ID")
		return
	}

	var order models.SaleOrder
	if err := h.DB.Preload("CreatedBy").Preload("SaleOrderItems").First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Sale order not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to fetch sale order")
		return
	}

	utils.OKResponse(c, "Sale order retrieved successfully", order)
}

// Create creates a new sale order
func (h *SaleOrderHandler) Create(c *gin.Context) {
	var req CreateSaleOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, _ := c.Get("user_id")

	// Calculate total amount
	var totalAmount float64
	var items []models.SaleOrderItem
	for _, item := range req.Items {
		subtotal := float64(item.Quantity) * item.UnitPrice
		totalAmount += subtotal
		items = append(items, models.SaleOrderItem{
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    subtotal,
		})
	}

	// Generate order number
	orderNumber := fmt.Sprintf("SO-%s-%d", time.Now().Format("20060102150405"), userID.(uint))

	order := models.SaleOrder{
		OrderNumber:    orderNumber,
		CustomerName:   req.CustomerName,
		TotalAmount:    totalAmount,
		Notes:          req.Notes,
		CreatedByID:    userID.(uint),
		SaleOrderItems: items,
	}

	if err := h.DB.Create(&order).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create sale order")
		return
	}

	// Reload with associations
	h.DB.Preload("CreatedBy").Preload("SaleOrderItems").First(&order, order.ID)

	utils.CreatedResponse(c, "Sale order created successfully", order)
}

// Update updates a sale order
func (h *SaleOrderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid sale order ID")
		return
	}

	var order models.SaleOrder
	if err := h.DB.Preload("SaleOrderItems").First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Sale order not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to fetch sale order")
		return
	}

	var req UpdateSaleOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Update fields
	order.CustomerName = req.CustomerName
	order.Notes = req.Notes

	// If items are provided, update them
	if len(req.Items) > 0 {
		// Delete existing items
		h.DB.Where("sale_order_id = ?", order.ID).Delete(&models.SaleOrderItem{})

		// Create new items
		var totalAmount float64
		var items []models.SaleOrderItem
		for _, item := range req.Items {
			subtotal := float64(item.Quantity) * item.UnitPrice
			totalAmount += subtotal
			items = append(items, models.SaleOrderItem{
				SaleOrderID: order.ID,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				Subtotal:    subtotal,
			})
		}
		order.TotalAmount = totalAmount
		order.SaleOrderItems = items

		if err := h.DB.Create(&items).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to update sale order items")
			return
		}
	}

	if err := h.DB.Save(&order).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update sale order")
		return
	}

	// Reload with associations
	h.DB.Preload("CreatedBy").Preload("SaleOrderItems").First(&order, order.ID)

	utils.OKResponse(c, "Sale order updated successfully", order)
}

// Delete soft deletes a sale order
func (h *SaleOrderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid sale order ID")
		return
	}

	var order models.SaleOrder
	if err := h.DB.First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Sale order not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to fetch sale order")
		return
	}

	// Soft delete the order and its items
	if err := h.DB.Delete(&order).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete sale order")
		return
	}

	utils.OKResponse(c, "Sale order deleted successfully", nil)
}
