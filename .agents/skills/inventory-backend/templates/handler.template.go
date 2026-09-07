// internal/handlers/sale_handler.go
//
// Handlers: bind/validate, call ONE service method, translate the result to
// the standard envelope. No db calls, no business logic, no gin-specific
// leakage into the service layer.

package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inventory/internal/models"
	"inventory/internal/services"
	"inventory/internal/utils/response"
)

type SaleHandler struct {
	saleService *services.SaleService
}

func NewSaleHandler(saleService *services.SaleService) *SaleHandler {
	return &SaleHandler{saleService: saleService}
}

type createSaleRequest struct {
	CustomerID string `json:"customer_id" binding:"required,uuid"`
	Items      []struct {
		ProductID string `json:"product_id" binding:"required,uuid"`
		Quantity  int    `json:"quantity" binding:"required,gt=0"`
	} `json:"items" binding:"required,min=1"`
}

// CreateSale godoc
// @Summary      Create a sale
// @Description  Creates a sale, verifies stock, and decreases inventory in one transaction.
// @Tags         sales
// @Accept       json
// @Produce      json
// @Param        request body createSaleRequest true "Sale details"
// @Success      201 {object} object{data=models.Sale}
// @Failure      400 {object} object{error=object{code=string,message=string}}
// @Failure      409 {object} object{error=object{code=string,message=string}} "insufficient inventory"
// @Router       /api/v1/sales [post]
func (h *SaleHandler) CreateSale(c *gin.Context) {
	var req createSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	userID := c.GetString("user_id") // set by auth middleware

	items := make([]services.SaleItemInput, 0, len(req.Items))
	for _, i := range req.Items {
		items = append(items, services.SaleItemInput{ProductID: i.ProductID, Quantity: i.Quantity})
	}

	sale, err := h.saleService.CreateSale(c.Request.Context(), services.CreateSaleInput{
		CustomerID: req.CustomerID,
		UserID:     userID,
		Items:      items,
	})
	if err != nil {
		respondServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, sale)
}

// ListSales godoc — example of the pagination envelope.
// @Summary      List sales
// @Tags         sales
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        page_size query int false "Items per page" default(20)
// @Success      200 {object} object{data=[]models.Sale,meta=response.PaginationMeta}
// @Router       /api/v1/sales [get]
func (h *SaleHandler) ListSales(c *gin.Context) {
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 20)

	sales, totalItems, err := h.saleService.ListSales(c.Request.Context(), page, pageSize)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, sales, response.NewPaginationMeta(page, pageSize, totalItems))
}

func queryInt(c *gin.Context, key string, fallback int) int {
	// Implementation detail (strconv.Atoi + fallback on error) omitted —
	// the point being illustrated is where page/page_size parsing belongs.
	return fallback
}

// respondServiceError maps sentinel errors from the service layer to HTTP
// status codes in ONE place, so handlers don't each hand-roll a switch.
func respondServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "resource not found")
	case errors.Is(err, models.ErrInsufficientInventory):
		response.Error(c, http.StatusConflict, "INSUFFICIENT_INVENTORY", err.Error())
	case errors.Is(err, models.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	case errors.Is(err, models.ErrForbidden):
		response.Error(c, http.StatusForbidden, "FORBIDDEN", "forbidden")
	default:
		// Never leak the raw error (SQL text, GORM internals) to the client —
		// this is unrelated to the logger's dev/prod switch, which only
		// affects what gets WRITTEN TO LOGS, never what's returned to callers.
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
	}
}
