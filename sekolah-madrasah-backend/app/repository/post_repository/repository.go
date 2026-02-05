package post_repository

import (
	"context"
	"net/http"
	"time"

	"sekolah-madrasah/app/repository/common"
	"sekolah-madrasah/database/schemas"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

// ========== Posts ==========

func (r *postRepository) applyPostFilter(query *gorm.DB, filter PostFilter) *gorm.DB {
	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.UnitId != nil {
		query = query.Where("unit_id = ?", *filter.UnitId)
	}
	if filter.AuthorId != nil {
		query = query.Where("author_id = ?", *filter.AuthorId)
	}
	if filter.PostType != nil {
		query = query.Where("post_type = ?", *filter.PostType)
	}
	if filter.IsPinned != nil {
		query = query.Where("is_pinned = ?", *filter.IsPinned)
	}
	if filter.IsImportant != nil {
		query = query.Where("is_important = ?", *filter.IsImportant)
	}

	// Visibility filter: org-wide or unit-specific
	if filter.IsOrgWide != nil {
		query = query.Where("is_org_wide = ?", *filter.IsOrgWide)
	}

	return query
}

func (r *postRepository) postToModel(schema schemas.Post) Post {
	authorName := ""
	if schema.Author != nil {
		authorName = schema.Author.FullName
	}

	return Post{
		Id:           schema.Id,
		UnitId:       schema.UnitId,
		AuthorId:     schema.AuthorId,
		AuthorName:   authorName,
		IsOrgWide:    schema.IsOrgWide,
		Title:        schema.Title,
		Content:      schema.Content,
		PostType:     schema.PostType,
		ImageURL:     schema.ImageURL,
		LinkURL:      schema.LinkURL,
		LinkTitle:    schema.LinkTitle,
		LinkPreview:  schema.LinkPreview,
		IsPinned:     schema.IsPinned,
		IsImportant:  schema.IsImportant,
		CommentCount: schema.CommentCount,
		CreatedAt:    schema.CreatedAt,
		UpdatedAt:    schema.UpdatedAt,
	}
}

func (r *postRepository) GetPost(ctx context.Context, filter PostFilter) (Post, int, error) {
	var schema schemas.Post
	query := r.db.WithContext(ctx).Model(&schemas.Post{}).Preload("Author")
	query = r.applyPostFilter(query, filter)

	if err := query.First(&schema).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return Post{}, code, err
	}

	return r.postToModel(schema), http.StatusOK, nil
}

func (r *postRepository) GetPosts(ctx context.Context, filter PostFilter, paginate *paginate_utils.PaginateData) ([]Post, int, error) {
	var schemaList []schemas.Post
	query := r.db.WithContext(ctx).Model(&schemas.Post{}).Preload("Author")
	query = r.applyPostFilter(query, filter)

	if err := common.CountTotal(query, paginate); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	query = common.ApplyPagination(query, paginate)
	query = query.Order("is_pinned DESC, created_at DESC")

	if err := query.Find(&schemaList).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return nil, code, err
	}

	posts := make([]Post, len(schemaList))
	for i, s := range schemaList {
		posts[i] = r.postToModel(s)
	}

	return posts, http.StatusOK, nil
}

func (r *postRepository) CreatePost(ctx context.Context, post Post) (Post, int, error) {
	schema := schemas.Post{
		Id:          post.Id,
		UnitId:      post.UnitId,
		AuthorId:    post.AuthorId,
		IsOrgWide:   post.IsOrgWide,
		Title:       post.Title,
		Content:     post.Content,
		PostType:    post.PostType,
		ImageURL:    post.ImageURL,
		LinkURL:     post.LinkURL,
		LinkTitle:   post.LinkTitle,
		LinkPreview: post.LinkPreview,
		IsPinned:    post.IsPinned,
		IsImportant: post.IsImportant,
	}

	if schema.Id == uuid.Nil {
		schema.Id = uuid.New()
	}
	schema.CreatedAt = time.Now()
	schema.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return Post{}, http.StatusInternalServerError, err
	}

	post.Id = schema.Id
	post.CreatedAt = schema.CreatedAt
	post.UpdatedAt = schema.UpdatedAt
	return post, http.StatusCreated, nil
}

func (r *postRepository) UpdatePost(ctx context.Context, filter PostFilter, post Post) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.Post{})
	query = r.applyPostFilter(query, filter)

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if post.Title != "" {
		updates["title"] = post.Title
	}
	if post.Content != "" {
		updates["content"] = post.Content
	}
	if post.ImageURL != "" {
		updates["image_url"] = post.ImageURL
	}
	if post.LinkURL != "" {
		updates["link_url"] = post.LinkURL
	}
	updates["is_pinned"] = post.IsPinned
	updates["is_important"] = post.IsImportant

	result := query.Updates(updates)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}

func (r *postRepository) DeletePost(ctx context.Context, filter PostFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.Post{})
	query = r.applyPostFilter(query, filter)

	result := query.Delete(&schemas.Post{})
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}

func (r *postRepository) IncrementCommentCount(ctx context.Context, postId uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).Model(&schemas.Post{}).
		Where("id = ?", postId).
		Update("comment_count", gorm.Expr("comment_count + ?", delta)).Error
}

// ========== Comments ==========

func (r *postRepository) applyCommentFilter(query *gorm.DB, filter PostCommentFilter) *gorm.DB {
	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.PostId != nil {
		query = query.Where("post_id = ?", *filter.PostId)
	}
	if filter.ParentId != nil {
		query = query.Where("parent_id = ?", *filter.ParentId)
	} else {
		query = query.Where("parent_id IS NULL") // default: top-level only
	}
	if filter.AuthorId != nil {
		query = query.Where("author_id = ?", *filter.AuthorId)
	}
	return query
}

func (r *postRepository) commentToModel(schema schemas.PostComment) PostComment {
	var parentId *uuid.UUID
	if schema.ParentId != nil {
		parentId = schema.ParentId
	}

	authorName := ""
	if schema.Author != nil {
		authorName = schema.Author.FullName
	}

	return PostComment{
		Id:         schema.Id,
		PostId:     schema.PostId,
		ParentId:   parentId,
		AuthorId:   schema.AuthorId,
		AuthorName: authorName,
		Content:    schema.Content,
		ReplyCount: schema.ReplyCount,
		CreatedAt:  schema.CreatedAt,
		UpdatedAt:  schema.UpdatedAt,
	}
}

func (r *postRepository) GetComment(ctx context.Context, filter PostCommentFilter) (PostComment, int, error) {
	var schema schemas.PostComment
	query := r.db.WithContext(ctx).Model(&schemas.PostComment{}).Preload("Author")

	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.PostId != nil {
		query = query.Where("post_id = ?", *filter.PostId)
	}

	if err := query.First(&schema).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return PostComment{}, code, err
	}

	return r.commentToModel(schema), http.StatusOK, nil
}

func (r *postRepository) GetComments(ctx context.Context, filter PostCommentFilter, paginate *paginate_utils.PaginateData) ([]PostComment, int, error) {
	var schemaList []schemas.PostComment
	query := r.db.WithContext(ctx).Model(&schemas.PostComment{}).Preload("Author").Preload("Replies.Author")
	query = r.applyCommentFilter(query, filter)

	if err := common.CountTotal(query, paginate); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	query = common.ApplyPagination(query, paginate)
	query = query.Order("created_at ASC")

	if err := query.Find(&schemaList).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return nil, code, err
	}

	comments := make([]PostComment, len(schemaList))
	for i, s := range schemaList {
		comments[i] = r.commentToModel(s)
		// Convert replies
		if len(s.Replies) > 0 {
			comments[i].Replies = make([]PostComment, len(s.Replies))
			for j, reply := range s.Replies {
				comments[i].Replies[j] = r.commentToModel(reply)
			}
		}
	}

	return comments, http.StatusOK, nil
}

func (r *postRepository) CreateComment(ctx context.Context, comment PostComment) (PostComment, int, error) {
	schema := schemas.PostComment{
		Id:       comment.Id,
		PostId:   comment.PostId,
		ParentId: comment.ParentId,
		AuthorId: comment.AuthorId,
		Content:  comment.Content,
	}

	if schema.Id == uuid.Nil {
		schema.Id = uuid.New()
	}
	schema.CreatedAt = time.Now()
	schema.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return PostComment{}, http.StatusInternalServerError, err
	}

	comment.Id = schema.Id
	comment.CreatedAt = schema.CreatedAt
	comment.UpdatedAt = schema.UpdatedAt
	return comment, http.StatusCreated, nil
}

func (r *postRepository) DeleteComment(ctx context.Context, filter PostCommentFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.PostComment{})
	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}

	result := query.Delete(&schemas.PostComment{})
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}

func (r *postRepository) IncrementReplyCount(ctx context.Context, commentId uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).Model(&schemas.PostComment{}).
		Where("id = ?", commentId).
		Update("reply_count", gorm.Expr("reply_count + ?", delta)).Error
}

// ========== Poll Options ==========

func (r *postRepository) GetPollOptions(ctx context.Context, filter PostPollOptionFilter) ([]PostPollOption, int, error) {
	var schemaList []schemas.PostPollOption
	query := r.db.WithContext(ctx).Model(&schemas.PostPollOption{})

	if filter.PostId != nil {
		query = query.Where("post_id = ?", *filter.PostId)
	}

	if err := query.Order("urutan ASC").Find(&schemaList).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return nil, code, err
	}

	options := make([]PostPollOption, len(schemaList))
	for i, s := range schemaList {
		options[i] = PostPollOption{
			Id:        s.Id,
			PostId:    s.PostId,
			Text:      s.Text,
			VoteCount: s.VoteCount,
			Urutan:    s.Urutan,
		}
	}

	return options, http.StatusOK, nil
}

func (r *postRepository) CreatePollOptions(ctx context.Context, options []PostPollOption) (int, error) {
	if len(options) == 0 {
		return http.StatusOK, nil
	}

	pollOptions := make([]schemas.PostPollOption, len(options))
	for i, opt := range options {
		id := opt.Id
		if id == uuid.Nil {
			id = uuid.New()
		}
		pollOptions[i] = schemas.PostPollOption{
			Id:     id,
			PostId: opt.PostId,
			Text:   opt.Text,
			Urutan: opt.Urutan,
		}
	}

	if err := r.db.WithContext(ctx).Create(&pollOptions).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (r *postRepository) IncrementVoteCount(ctx context.Context, optionId uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).Model(&schemas.PostPollOption{}).
		Where("id = ?", optionId).
		Update("vote_count", gorm.Expr("vote_count + ?", delta)).Error
}

// ========== Poll Votes ==========

func (r *postRepository) GetPollVote(ctx context.Context, filter PostPollVoteFilter) (PostPollVote, int, error) {
	var schema schemas.PostPollVote
	query := r.db.WithContext(ctx).Model(&schemas.PostPollVote{})

	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.PostId != nil {
		query = query.Where("post_id = ?", *filter.PostId)
	}
	if filter.UserId != nil {
		query = query.Where("user_id = ?", *filter.UserId)
	}

	if err := query.First(&schema).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return PostPollVote{}, code, err
	}

	return PostPollVote{
		Id:       schema.Id,
		PostId:   schema.PostId,
		OptionId: schema.OptionId,
		UserId:   schema.UserId,
	}, http.StatusOK, nil
}

func (r *postRepository) CreatePollVote(ctx context.Context, vote PostPollVote) (PostPollVote, int, error) {
	schema := schemas.PostPollVote{
		Id:       vote.Id,
		PostId:   vote.PostId,
		OptionId: vote.OptionId,
		UserId:   vote.UserId,
	}

	if schema.Id == uuid.Nil {
		schema.Id = uuid.New()
	}
	schema.CreatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return PostPollVote{}, http.StatusInternalServerError, err
	}

	vote.Id = schema.Id
	return vote, http.StatusCreated, nil
}

func (r *postRepository) DeletePollVote(ctx context.Context, filter PostPollVoteFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.PostPollVote{})

	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.PostId != nil {
		query = query.Where("post_id = ?", *filter.PostId)
	}
	if filter.UserId != nil {
		query = query.Where("user_id = ?", *filter.UserId)
	}

	result := query.Delete(&schemas.PostPollVote{})
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	return http.StatusOK, nil
}

func (r *postRepository) GetUserVoteForPost(ctx context.Context, postId, userId uuid.UUID) (*PostPollVote, error) {
	var schema schemas.PostPollVote
	err := r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postId, userId).
		First(&schema).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &PostPollVote{
		Id:       schema.Id,
		PostId:   schema.PostId,
		OptionId: schema.OptionId,
		UserId:   schema.UserId,
	}, nil
}
