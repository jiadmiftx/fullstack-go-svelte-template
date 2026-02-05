package post_use_case

import (
	"context"
	"errors"
	"net/http"

	"sekolah-madrasah/app/repository/post_repository"
	"sekolah-madrasah/app/repository/user_repository"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
)

type postUseCase struct {
	postRepo post_repository.PostRepository
	userRepo user_repository.UserRepository
}

func NewPostUseCase(postRepo post_repository.PostRepository, userRepo user_repository.UserRepository) PostUseCase {
	return &postUseCase{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (u *postUseCase) toRepoFilter(dto PostFilterDTO) post_repository.PostFilter {
	return post_repository.PostFilter{
		Id:          dto.Id,
		UnitId:      dto.UnitId,
		OrgId:       dto.OrgId,
		AuthorId:    dto.AuthorId,
		PostType:    dto.PostType,
		IsPinned:    dto.IsPinned,
		IsImportant: dto.IsImportant,
		IsOrgWide:   dto.IsOrgWide,
	}
}

func (u *postUseCase) toDTO(repoModel post_repository.Post) PostDTO {
	return PostDTO{
		Id:           repoModel.Id,
		UnitId:       repoModel.UnitId,
		AuthorId:     repoModel.AuthorId,
		AuthorName:   repoModel.AuthorName,
		IsOrgWide:    repoModel.IsOrgWide,
		Title:        repoModel.Title,
		Content:      repoModel.Content,
		PostType:     repoModel.PostType,
		ImageURL:     repoModel.ImageURL,
		LinkURL:      repoModel.LinkURL,
		LinkTitle:    repoModel.LinkTitle,
		LinkPreview:  repoModel.LinkPreview,
		IsPinned:     repoModel.IsPinned,
		IsImportant:  repoModel.IsImportant,
		CommentCount: repoModel.CommentCount,
		CreatedAt:    repoModel.CreatedAt,
		UpdatedAt:    repoModel.UpdatedAt,
	}
}

func (u *postUseCase) commentToDTO(repoModel post_repository.PostComment) PostCommentDTO {
	dto := PostCommentDTO{
		Id:         repoModel.Id,
		PostId:     repoModel.PostId,
		ParentId:   repoModel.ParentId,
		AuthorId:   repoModel.AuthorId,
		AuthorName: repoModel.AuthorName,
		Content:    repoModel.Content,
		ReplyCount: repoModel.ReplyCount,
		CreatedAt:  repoModel.CreatedAt,
		UpdatedAt:  repoModel.UpdatedAt,
	}

	if len(repoModel.Replies) > 0 {
		dto.Replies = make([]PostCommentDTO, len(repoModel.Replies))
		for i, r := range repoModel.Replies {
			dto.Replies[i] = u.commentToDTO(r)
		}
	}

	return dto
}

func (u *postUseCase) optionToDTO(repoModel post_repository.PostPollOption) PostPollOptionDTO {
	return PostPollOptionDTO{
		Id:        repoModel.Id,
		PostId:    repoModel.PostId,
		Text:      repoModel.Text,
		VoteCount: repoModel.VoteCount,
		Urutan:    repoModel.Urutan,
		HasVoted:  repoModel.HasVoted,
	}
}

// ========== Posts ==========

func (u *postUseCase) GetPost(ctx context.Context, filter PostFilterDTO, userId uuid.UUID) (PostDTO, int, error) {
	post, code, err := u.postRepo.GetPost(ctx, u.toRepoFilter(filter))
	if err != nil {
		return PostDTO{}, code, err
	}

	dto := u.toDTO(post)

	// Load poll options if poll type
	if post.PostType == "poll" {
		options, _, err := u.postRepo.GetPollOptions(ctx, post_repository.PostPollOptionFilter{
			PostId: &post.Id,
		})
		if err == nil {
			// Check user vote
			userVote, _ := u.postRepo.GetUserVoteForPost(ctx, post.Id, userId)

			dto.Options = make([]PostPollOptionDTO, len(options))
			for i, opt := range options {
				dto.Options[i] = u.optionToDTO(opt)
				if userVote != nil && userVote.OptionId == opt.Id {
					dto.Options[i].HasVoted = true
				}
			}
		}
	}

	return dto, code, nil
}

func (u *postUseCase) GetPosts(ctx context.Context, filter PostFilterDTO, paginate *paginate_utils.PaginateData) ([]PostDTO, int, error) {
	posts, code, err := u.postRepo.GetPosts(ctx, u.toRepoFilter(filter), paginate)
	if err != nil {
		return nil, code, err
	}

	dtos := make([]PostDTO, len(posts))
	for i, p := range posts {
		dtos[i] = u.toDTO(p)
	}

	return dtos, code, nil
}

func (u *postUseCase) CreatePost(ctx context.Context, dto CreatePostDTO, authorId uuid.UUID) (PostDTO, int, error) {
	// Get author name
	authorName := ""
	user, _, err := u.userRepo.GetUser(ctx, user_repository.UserFilter{Id: &authorId})
	if err == nil {
		authorName = user.FullName
	}

	repoPost := post_repository.Post{
		UnitId:      dto.UnitId,
		AuthorId:    authorId,
		AuthorName:  authorName,
		IsOrgWide:   dto.IsOrgWide,
		Title:       dto.Title,
		Content:     dto.Content,
		PostType:    dto.PostType,
		ImageURL:    dto.ImageURL,
		LinkURL:     dto.LinkURL,
		LinkTitle:   dto.LinkTitle,
		LinkPreview: dto.LinkPreview,
		IsPinned:    dto.IsPinned,
		IsImportant: dto.IsImportant,
	}

	created, code, err := u.postRepo.CreatePost(ctx, repoPost)
	if err != nil {
		return PostDTO{}, code, err
	}

	// Create poll options if poll type
	if dto.PostType == "poll" && len(dto.PollOptions) > 0 {
		options := make([]post_repository.PostPollOption, len(dto.PollOptions))
		for i, optText := range dto.PollOptions {
			options[i] = post_repository.PostPollOption{
				PostId: created.Id,
				Text:   optText,
				Urutan: i + 1,
			}
		}
		u.postRepo.CreatePollOptions(ctx, options)
	}

	result := u.toDTO(created)
	result.AuthorName = authorName
	return result, code, nil
}

func (u *postUseCase) UpdatePost(ctx context.Context, filter PostFilterDTO, dto PostDTO) (int, error) {
	repoPost := post_repository.Post{
		Title:       dto.Title,
		Content:     dto.Content,
		ImageURL:    dto.ImageURL,
		LinkURL:     dto.LinkURL,
		LinkTitle:   dto.LinkTitle,
		LinkPreview: dto.LinkPreview,
		IsPinned:    dto.IsPinned,
		IsImportant: dto.IsImportant,
	}

	return u.postRepo.UpdatePost(ctx, u.toRepoFilter(filter), repoPost)
}

func (u *postUseCase) DeletePost(ctx context.Context, filter PostFilterDTO) (int, error) {
	return u.postRepo.DeletePost(ctx, u.toRepoFilter(filter))
}

// ========== Comments ==========

func (u *postUseCase) GetComments(ctx context.Context, postId uuid.UUID, paginate *paginate_utils.PaginateData) ([]PostCommentDTO, int, error) {
	comments, code, err := u.postRepo.GetComments(ctx, post_repository.PostCommentFilter{
		PostId: &postId,
	}, paginate)
	if err != nil {
		return nil, code, err
	}

	dtos := make([]PostCommentDTO, len(comments))
	for i, c := range comments {
		dtos[i] = u.commentToDTO(c)
	}

	return dtos, code, nil
}

func (u *postUseCase) CreateComment(ctx context.Context, dto CreateCommentDTO, authorId uuid.UUID) (PostCommentDTO, int, error) {
	// Get author name
	authorName := ""
	user, _, err := u.userRepo.GetUser(ctx, user_repository.UserFilter{Id: &authorId})
	if err == nil {
		authorName = user.FullName
	}

	repoComment := post_repository.PostComment{
		PostId:     dto.PostId,
		ParentId:   dto.ParentId,
		AuthorId:   authorId,
		AuthorName: authorName,
		Content:    dto.Content,
	}

	created, code, err := u.postRepo.CreateComment(ctx, repoComment)
	if err != nil {
		return PostCommentDTO{}, code, err
	}

	// Increment comment count on post
	u.postRepo.IncrementCommentCount(ctx, dto.PostId, 1)

	// If reply, increment reply count on parent
	if dto.ParentId != nil {
		u.postRepo.IncrementReplyCount(ctx, *dto.ParentId, 1)
	}

	result := u.commentToDTO(created)
	result.AuthorName = authorName
	return result, code, nil
}

func (u *postUseCase) DeleteComment(ctx context.Context, commentId uuid.UUID, userId uuid.UUID) (int, error) {
	// Get comment first
	comment, code, err := u.postRepo.GetComment(ctx, post_repository.PostCommentFilter{
		Id: &commentId,
	})
	if err != nil {
		return code, err
	}

	// Only author can delete
	if comment.AuthorId != userId {
		return http.StatusForbidden, errors.New("you can only delete your own comments")
	}

	// Delete the comment
	code, err = u.postRepo.DeleteComment(ctx, post_repository.PostCommentFilter{
		Id: &commentId,
	})
	if err != nil {
		return code, err
	}

	// Decrement counts
	u.postRepo.IncrementCommentCount(ctx, comment.PostId, -1)
	if comment.ParentId != nil {
		u.postRepo.IncrementReplyCount(ctx, *comment.ParentId, -1)
	}

	return code, nil
}

// ========== Poll ==========

func (u *postUseCase) VotePoll(ctx context.Context, dto VotePollDTO, userId uuid.UUID) (int, error) {
	// Check if already voted
	existingVote, _ := u.postRepo.GetUserVoteForPost(ctx, dto.PostId, userId)
	if existingVote != nil {
		// Remove old vote
		u.postRepo.DeletePollVote(ctx, post_repository.PostPollVoteFilter{
			Id: &existingVote.Id,
		})
		u.postRepo.IncrementVoteCount(ctx, existingVote.OptionId, -1)
	}

	// Create new vote
	vote := post_repository.PostPollVote{
		PostId:   dto.PostId,
		OptionId: dto.OptionId,
		UserId:   userId,
	}

	_, code, err := u.postRepo.CreatePollVote(ctx, vote)
	if err != nil {
		return code, err
	}

	// Increment vote count
	u.postRepo.IncrementVoteCount(ctx, dto.OptionId, 1)

	return http.StatusOK, nil
}
