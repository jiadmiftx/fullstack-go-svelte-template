package post_repository

import (
	"context"

	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
)

type PostRepository interface {
	// Posts
	GetPost(ctx context.Context, filter PostFilter) (Post, int, error)
	GetPosts(ctx context.Context, filter PostFilter, paginate *paginate_utils.PaginateData) ([]Post, int, error)
	CreatePost(ctx context.Context, post Post) (Post, int, error)
	UpdatePost(ctx context.Context, filter PostFilter, post Post) (int, error)
	DeletePost(ctx context.Context, filter PostFilter) (int, error)
	IncrementCommentCount(ctx context.Context, postId uuid.UUID, delta int) error

	// Comments
	GetComment(ctx context.Context, filter PostCommentFilter) (PostComment, int, error)
	GetComments(ctx context.Context, filter PostCommentFilter, paginate *paginate_utils.PaginateData) ([]PostComment, int, error)
	CreateComment(ctx context.Context, comment PostComment) (PostComment, int, error)
	DeleteComment(ctx context.Context, filter PostCommentFilter) (int, error)
	IncrementReplyCount(ctx context.Context, commentId uuid.UUID, delta int) error

	// Poll Options
	GetPollOptions(ctx context.Context, filter PostPollOptionFilter) ([]PostPollOption, int, error)
	CreatePollOptions(ctx context.Context, options []PostPollOption) (int, error)
	IncrementVoteCount(ctx context.Context, optionId uuid.UUID, delta int) error

	// Poll Votes
	GetPollVote(ctx context.Context, filter PostPollVoteFilter) (PostPollVote, int, error)
	CreatePollVote(ctx context.Context, vote PostPollVote) (PostPollVote, int, error)
	DeletePollVote(ctx context.Context, filter PostPollVoteFilter) (int, error)
	GetUserVoteForPost(ctx context.Context, postId, userId uuid.UUID) (*PostPollVote, error)
}
