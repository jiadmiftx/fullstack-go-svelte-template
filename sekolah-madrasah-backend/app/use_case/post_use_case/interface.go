package post_use_case

import (
	"context"

	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
)

type PostUseCase interface {
	// Posts
	GetPost(ctx context.Context, filter PostFilterDTO, userId uuid.UUID) (PostDTO, int, error)
	GetPosts(ctx context.Context, filter PostFilterDTO, paginate *paginate_utils.PaginateData) ([]PostDTO, int, error)
	CreatePost(ctx context.Context, dto CreatePostDTO, authorId uuid.UUID) (PostDTO, int, error)
	UpdatePost(ctx context.Context, filter PostFilterDTO, dto PostDTO) (int, error)
	DeletePost(ctx context.Context, filter PostFilterDTO) (int, error)

	// Comments
	GetComments(ctx context.Context, postId uuid.UUID, paginate *paginate_utils.PaginateData) ([]PostCommentDTO, int, error)
	CreateComment(ctx context.Context, dto CreateCommentDTO, authorId uuid.UUID) (PostCommentDTO, int, error)
	DeleteComment(ctx context.Context, commentId uuid.UUID, userId uuid.UUID) (int, error)

	// Poll
	VotePoll(ctx context.Context, dto VotePollDTO, userId uuid.UUID) (int, error)
}
