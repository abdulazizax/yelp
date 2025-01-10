package handler

import (
	"strconv"

	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreateBusiness godoc
// @Router /business [post]
// @Summary Create a new business
// @Description Create a new business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param business body entity.Business true "Business object"
// @Success 201 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateBusiness(ctx *gin.Context) {
	var (
		body entity.Business
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	business, err := h.UseCase.BusinessRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating business") {
		return
	}

	ctx.JSON(201, business)
}

// GetBusiness godoc
// @Router /business/{id} [get]
// @Summary Get a business by ID
// @Description Get a business by ID
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Param owner_id path string true "Owner ID"
// @Param category_id path string true "Category ID"
// @Success 200 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBusiness(ctx *gin.Context) {
	var (
		req entity.BusinessSingleRequest
	)

	req.ID = ctx.Param("id")
	req.OwnerID = ctx.Param("owner_id")
	req.CategoryID = ctx.Param("category_id")

	business, err := h.UseCase.BusinessRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting business") {
		return
	}

	ctx.JSON(200, business)
}

// GetBusinesss godoc
// @Router /business/list [get]
// @Summary Get a list of users
// @Description Get a list of users
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.BusinessList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBusinesses(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "name",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "address",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "description",
			Type:   "search",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	users, err := h.UseCase.BusinessRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting users") {
		return
	}

	ctx.JSON(200, users)
}

// UpdateBusiness godoc
// @Router /business [put]
// @Summary Update a business
// @Description Update a business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param business body entity.Business true "Business object"
// @Success 200 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateBusiness(ctx *gin.Context) {
	var (
		body entity.Business
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_role") == "business_owner" {
		body.ID = ctx.GetHeader("sub")
	}

	business, err := h.UseCase.BusinessRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating business") {
		return
	}

	ctx.JSON(200, business)
}

// DeleteBusiness godoc
// @Router /business/{id} [delete]
// @Summary Delete a business
// @Description Delete a business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteBusiness(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	if ctx.GetHeader("user_role") == "business_owner" {
		req.ID = ctx.GetHeader("sub")
	}

	err := h.UseCase.BusinessRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting business") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Business deleted successfully",
	})
}
