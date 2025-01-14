package handler

import (
	"strconv"

	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreateBusinessCategory godoc
// @Router /business-category [post]
// @Summary Create a new business-category
// @Description Create a new business-category
// @Security BearerAuth
// @Tags business-category
// @Accept  json
// @Produce  json
// @Param business-category body entity.BusinessCategory true "BusinessCategory object"
// @Success 201 {object} entity.BusinessCategory
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateBusinessCategory(ctx *gin.Context) {
	var (
		body entity.BusinessCategory
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	businessCategory, err := h.UseCase.BusinessCategoryRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating business-category") {
		return
	}

	ctx.JSON(201, businessCategory)
}

// GetBusinessCategory godoc
// @Router /business-category/{id} [get]
// @Summary Get a business-category by ID
// @Description Get a business-category by ID
// @Security BearerAuth
// @Tags business-category
// @Accept  json
// @Produce  json
// @Param id path string true "BusinessCategory ID"
// @Param category_id path string true "Category ID"
// @Success 200 {object} entity.BusinessCategory
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBusinessCategory(ctx *gin.Context) {
	var (
		req entity.BusinessCategorySingleRequest
	)

	req.ID = ctx.Param("id")
	req.Name = ctx.Param("name")

	businessCategory, err := h.UseCase.BusinessCategoryRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting business-category") {
		return
	}

	ctx.JSON(200, businessCategory)
}

// GetBusinessCategorys godoc
// @Router /business-category/list [get]
// @Summary Get a list of users
// @Description Get a list of users
// @Security BearerAuth
// @Tags business-category
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.BusinessCategoryList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBusinessCategories(ctx *gin.Context) {
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
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	users, err := h.UseCase.BusinessCategoryRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting users") {
		return
	}

	ctx.JSON(200, users)
}

// UpdateBusinessCategory godoc
// @Router /business-category [put]
// @Summary Update a business-category
// @Description Update a business-category
// @Security BearerAuth
// @Tags business-category
// @Accept  json
// @Produce  json
// @Param business-category body entity.BusinessCategory true "BusinessCategory object"
// @Success 200 {object} entity.BusinessCategory
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateBusinessCategory(ctx *gin.Context) {
	var (
		body entity.BusinessCategory
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_role") == "super_admin" {
		body.ID = ctx.GetHeader("sub")
	}

	businessCategory, err := h.UseCase.BusinessCategoryRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating business-category") {
		return
	}

	ctx.JSON(200, businessCategory)
}

// DeleteBusinessCategory godoc
// @Router /business-category/{id} [delete]
// @Summary Delete a business-category
// @Description Delete a business-category
// @Security BearerAuth
// @Tags business-category
// @Accept  json
// @Produce  json
// @Param id path string true "BusinessCategory ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteBusinessCategory(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	if ctx.GetHeader("user_role") == "super_admin" {
		req.ID = ctx.GetHeader("sub")
	}

	err := h.UseCase.BusinessCategoryRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting business-category") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "BusinessCategory deleted successfully",
	})
}
