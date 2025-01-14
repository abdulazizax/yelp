package handler

import (
	"net/http"
	"strconv"

	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreateReview godoc
// @Router /review [post]
// @Summary Create a new review
// @Description Create a new review
// @Security BearerAuth
// @Tags review
// @Accept  json
// @Produce  json
// @Param review body entity.Review true "Review object"
// @Success 201 {object} entity.Review
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateReview(ctx *gin.Context) {
	var (
		body entity.Review
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	body.UserID = ctx.GetHeader("sub")

	review, err := h.UseCase.ReviewRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating review") {
		return
	}

	review.Attachments, err = h.UseCase.ReviewAttachmentRepo.MultipleUpsert(ctx, entity.ReviewAttachmentMultipleInsertRequest{
		ReviewId:    review.ID,
		Attachments: body.Attachments,
	})
	if h.HandleDbError(ctx, err, "Error creating review") {
		return
	}

	ctx.JSON(201, review)
}

// GetReview godoc
// @Router /review/{id} [get]
// @Summary Get a review by ID
// @Description Get a review by ID
// @Security BearerAuth
// @Tags review
// @Accept  json
// @Produce  json
// @Param id path string true "Review ID"
// @Success 200 {object} entity.Review
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetReview(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	review, err := h.UseCase.ReviewRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting review") {
		return
	}

	ctx.JSON(200, review)
}

// GetReviews godoc
// @Router /review/list [get]
// @Summary Get a list of users
// @Description Get a list of users
// @Security BearerAuth
// @Tags review
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.ReviewList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetReviews(ctx *gin.Context) {
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
			Column: "rating",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "comment",
			Type:   "search",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	users, err := h.UseCase.ReviewRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting users") {
		return
	}

	ctx.JSON(200, users)
}

// UpdateReview godoc
// @Router /review [put]
// @Summary Update a review
// @Description Update a review
// @Security BearerAuth
// @Tags review
// @Accept  json
// @Produce  json
// @Param review body entity.Review true "Review object"
// @Success 200 {object} entity.Review
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateReview(ctx *gin.Context) {
	var (
		body entity.Review
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if body.UserID != ctx.GetHeader("sub") || ctx.GetHeader("user_type") != "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "You have no access to the comment", http.StatusForbidden)
		return
	}

	review, err := h.UseCase.ReviewRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating review") {
		return
	}

	review.Attachments, err = h.UseCase.ReviewAttachmentRepo.MultipleUpsert(ctx, entity.ReviewAttachmentMultipleInsertRequest{
		ReviewId:    review.ID,
		Attachments: body.Attachments,
	})
	if h.HandleDbError(ctx, err, "Error upserting review attachments") {
		return
	}

	ctx.JSON(200, review)
}

// DeleteReview godoc
// @Router /review/{id} [delete]
// @Summary Delete a review
// @Description Delete a review
// @Security BearerAuth
// @Tags review
// @Accept  json
// @Produce  json
// @Param id path string true "Review ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteReview(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	body, err := h.UseCase.ReviewRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting review") {
		return
	}

	if body.UserID != ctx.GetHeader("sub") || ctx.GetHeader("user_type") != "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "You have no access to the comment", http.StatusForbidden)
		return
	}

	err = h.UseCase.ReviewRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting review") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Review deleted successfully",
	})
}
