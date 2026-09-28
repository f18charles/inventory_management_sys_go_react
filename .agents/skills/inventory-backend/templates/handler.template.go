// internal/sale/handler.go
//
// Handlers live INSIDE the owning domain package. They bind/validate the
// request, call ONE service method, and translate the result to the standard
// envelope via internal/utils/response. No db calls, no business logic, and no
// Gin types ever cross into the service layer.

package sale

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"i_m_s/internal/middleware"
	"i_m_s/internal/models"
	"i_m_s/internal/utils/response"
)

type Handler struct {
	service *Service // same package — the domain's service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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
// @Security     BearerAuth
// @Param        request body createSaleRequest true "Sale details"
// @Success      201 {object} object{data=models.Sale}
// @Failure      400 {object} object{error=object{code=string,message=string}}
// @Failure      409 {object} object{error=object{code=string,message=string}} "insufficient inventory"
// @Router       /sales [post]
func (h *Handler) CreateSale(c *gin.Context) {
	var req createSaleRequest
	// BindJSON writes the standard INVALID_REQUEST envelope and returns false.
	if !response.BindJSON(c, &req) {
		return
	}

	items := make([]SaleItemInput, 0, len(req.Items))
	for _, i := range req.Items {
		items = append(items, SaleItemInput{ProductID: uuid.MustParse(i.ProductID), Quantity: i.Quantity})
	}

	sale, err := h.service.CreateSale(c.Request.Context(), CreateSaleInput{
		CustomerID: uuid.MustParse(req.CustomerID),
		UserID:     currentUserID(c),
		Items:      items,
	})
	if err != nil {
		// RespondError is the single shared mapping from domain sentinels to
		// HTTP status + envelope — never hand-roll a switch per handler.
		response.RespondError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, sale)
}

// ListSales godoc — example of the pagination envelope.
// @Summary      List sales
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        page_size query int false "Items per page" default(20)
// @Success      200 {object} object{data=[]models.Sale,meta=response.PaginationMeta}
// @Router       /sales [get]
func (h *Handler) ListSales(c *gin.Context) {
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 20)

	sales, totalItems, err := h.service.ListSales(c.Request.Context(), page, pageSize)
	if err != nil {
		response.RespondError(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, sales, response.NewPaginationMeta(page, pageSize, totalItems))
}

// currentUserID reads the id attached by JWTAuthMiddleware. Query/param
// parsing helpers (queryInt, uuid parsing) stay in the handler layer.
func currentUserID(c *gin.Context) uuid.UUID {
	id, _ := uuid.Parse(c.GetString(middleware.ContextUserID))
	return id
}

func queryInt(c *gin.Context, key string, fallback int) int {
	// strconv.Atoi + fallback on error/absent.
	return fallback
}
