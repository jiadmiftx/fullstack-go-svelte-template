package post_controller

import (
	"net/http"

	"sekolah-madrasah/app/use_case/post_use_case"
	"sekolah-madrasah/pkg/gin_utils"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostController struct {
	useCase post_use_case.PostUseCase
}

func NewPostController(useCase post_use_case.PostUseCase) *PostController {
	return &PostController{useCase: useCase}
}

// ========== Posts ==========

// GetPosts godoc
// @Summary Get list of posts
// @Tags posts
func (c *PostController) GetPosts(ctx *gin.Context) {
	paginate := &paginate_utils.PaginateData{}
	queryParams := make(map[string]interface{})
	for k, v := range ctx.Request.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}
	paginate_utils.CheckPaginateFromMap(queryParams, paginate)

	filter := post_use_case.PostFilterDTO{}

	if unitId := ctx.Query("unit_id"); unitId != "" {
		if id, err := uuid.Parse(unitId); err == nil {
			filter.UnitId = &id
		}
	}

	if postType := ctx.Query("post_type"); postType != "" {
		filter.PostType = &postType
	}

	// Filter by org-wide or unit-specific
	if isOrgWide := ctx.Query("is_org_wide"); isOrgWide != "" {
		val := isOrgWide == "true"
		filter.IsOrgWide = &val
	}

	posts, code, err := c.useCase.GetPosts(ctx.Request.Context(), filter, paginate)
	if err != nil {
		ctx.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(code, gin_utils.DataWithPaginateResponse{
		DataResponse: gin_utils.DataResponse{Message: "success", Data: posts},
		Paginate:     paginate,
	})
}

// GetPost godoc
// @Summary Get single post with comments
// @Tags posts
func (c *PostController) GetPost(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid id"})
		return
	}

	userIdVal, _ := ctx.Get("user_id")
	userId, _ := userIdVal.(uuid.UUID)

	filter := post_use_case.PostFilterDTO{Id: &id}
	post, code, err := c.useCase.GetPost(ctx.Request.Context(), filter, userId)
	if err != nil {
		ctx.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(code, gin_utils.DataResponse{Message: "success", Data: post})
}

// CreatePost godoc
// @Summary Create a new post
// @Tags posts
func (c *PostController) CreatePost(ctx *gin.Context) {
	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	userIdVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "user not authenticated"})
		return
	}
	authorId := userIdVal.(uuid.UUID)

	dto := post_use_case.CreatePostDTO{
		UnitId:      req.UnitId,
		IsOrgWide:   req.IsOrgWide,
		Title:       req.Title,
		Content:     req.Content,
		PostType:    req.PostType,
		ImageURL:    req.ImageURL,
		LinkURL:     req.LinkURL,
		LinkTitle:   req.LinkTitle,
		LinkPreview: req.LinkPreview,
		IsPinned:    req.IsPinned,
		IsImportant: req.IsImportant,
		PollOptions: req.PollOptions,
	}

	created, code, err := c.useCase.CreatePost(ctx.Request.Context(), dto, authorId)
	if err != nil {
		ctx.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(code, gin_utils.DataResponse{Message: "post created successfully", Data: created})
}

// UpdatePost godoc
// @Summary Update a post
// @Tags posts
func (c *PostController) UpdatePost(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid id"})
		return
	}

	var req UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	filter := post_use_case.PostFilterDTO{Id: &id}
	dto := post_use_case.PostDTO{
		Title:       req.Title,
		Content:     req.Content,
		ImageURL:    req.ImageURL,
		LinkURL:     req.LinkURL,
		LinkTitle:   req.LinkTitle,
		LinkPreview: req.LinkPreview,
		IsPinned:    req.IsPinned,
		IsImportant: req.IsImportant,
	}

	code, err := c.useCase.UpdatePost(ctx.Request.Context(), filter, dto)
	if err != nil {
		ctx.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "post updated successfully"})
}

// DeletePost godoc
// @Summary Delete a post
// @Tags posts
func (c *PostController) DeletePost(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid id"})
		return
	}

	filter := post_use_case.PostFilterDTO{Id: &id}
	code, err := c.useCase.DeletePost(ctx.Request.Context(), filter)
	if err != nil {
		ctx.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "post deleted successfully"})
}

// ========== Comments ==========

// GetComments godoc
// @Summary Get comments for a post
// @Tags posts
func (c *PostController) GetComments(ctx *gin.Context) {
	postId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid post id"})
		return
	}

	paginate := &paginate_utils.PaginateData{}
	queryParams := make(map[string]interface{})
	for k, v := range ctx.Request.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}
	paginate_utils.CheckPaginateFromMap(queryParams, paginate)

	comments, code, err := c.useCase.GetComments(ctx.Request.Context(), postId, paginate)
	if err != nil {
		ctx.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(code, gin_utils.DataWithPaginateResponse{
		DataResponse: gin_utils.DataResponse{Message: "success", Data: comments},
		Paginate:     paginate,
	})
}

// CreateComment godoc
// @Summary Add a comment to a post
// @Tags posts
func (c *PostController) CreateComment(ctx *gin.Context) {
	postId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid post id"})
		return
	}

	var req CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	userIdVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "user not authenticated"})
		return
	}
	authorId := userIdVal.(uuid.UUID)

	dto := post_use_case.CreateCommentDTO{
		PostId:   postId,
		ParentId: req.ParentId,
		Content:  req.Content,
	}

	created, code, err := c.useCase.CreateComment(ctx.Request.Context(), dto, authorId)
	if err != nil {
		ctx.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(code, gin_utils.DataResponse{Message: "comment added successfully", Data: created})
}

// DeleteComment godoc
// @Summary Delete a comment
// @Tags posts
func (c *PostController) DeleteComment(ctx *gin.Context) {
	commentId, err := uuid.Parse(ctx.Param("commentId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid comment id"})
		return
	}

	userIdVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "user not authenticated"})
		return
	}
	userId := userIdVal.(uuid.UUID)

	code, err := c.useCase.DeleteComment(ctx.Request.Context(), commentId, userId)
	if err != nil {
		ctx.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "comment deleted successfully"})
}

// ========== Poll ==========

// VotePoll godoc
// @Summary Vote on a poll
// @Tags posts
func (c *PostController) VotePoll(ctx *gin.Context) {
	postId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid post id"})
		return
	}

	var req VotePollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	userIdVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin_utils.MessageResponse{Message: "user not authenticated"})
		return
	}
	userId := userIdVal.(uuid.UUID)

	dto := post_use_case.VotePollDTO{
		PostId:   postId,
		OptionId: req.OptionId,
	}

	code, err := c.useCase.VotePoll(ctx.Request.Context(), dto, userId)
	if err != nil {
		ctx.JSON(code, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "vote recorded successfully"})
}

// Helper
func splitByComma(s string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	current := ""
	for _, c := range s {
		if c == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// Request structs
type CreatePostRequest struct {
	UnitId    uuid.UUID `json:"unit_id" binding:"required"`
	IsOrgWide bool      `json:"is_org_wide"` // true = org level, false = unit level
	Title     string    `json:"title"`
	Content   string    `json:"content" binding:"required"`
	PostType  string    `json:"post_type" binding:"required"`

	ImageURL    string `json:"image_url"`
	LinkURL     string `json:"link_url"`
	LinkTitle   string `json:"link_title"`
	LinkPreview string `json:"link_preview"`

	IsPinned    bool `json:"is_pinned"`
	IsImportant bool `json:"is_important"`

	PollOptions []string `json:"poll_options"`
}

type UpdatePostRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	ImageURL    string `json:"image_url"`
	LinkURL     string `json:"link_url"`
	LinkTitle   string `json:"link_title"`
	LinkPreview string `json:"link_preview"`
	IsPinned    bool   `json:"is_pinned"`
	IsImportant bool   `json:"is_important"`
}

type CreateCommentRequest struct {
	ParentId *uuid.UUID `json:"parent_id"`
	Content  string     `json:"content" binding:"required"`
}

type VotePollRequest struct {
	OptionId uuid.UUID `json:"option_id" binding:"required"`
}
