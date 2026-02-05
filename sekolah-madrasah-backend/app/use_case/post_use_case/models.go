package post_use_case

import (
	"time"

	"github.com/google/uuid"
)

// Post DTOs
type PostDTO struct {
	Id         uuid.UUID `json:"id"`
	UnitId     uuid.UUID `json:"unit_id"`
	AuthorId   uuid.UUID `json:"author_id"`
	AuthorName string    `json:"author_name"`
	IsOrgWide  bool      `json:"is_org_wide"` // true = org level, false = unit level

	Title    string `json:"title"`
	Content  string `json:"content"`
	PostType string `json:"post_type"` // text, photo, poll, link

	ImageURL    string `json:"image_url,omitempty"`
	LinkURL     string `json:"link_url,omitempty"`
	LinkTitle   string `json:"link_title,omitempty"`
	LinkPreview string `json:"link_preview,omitempty"`

	IsPinned     bool `json:"is_pinned"`
	IsImportant  bool `json:"is_important"`
	CommentCount int  `json:"comment_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Options  []PostPollOptionDTO `json:"options,omitempty"`
	Comments []PostCommentDTO    `json:"comments,omitempty"`
}

type PostFilterDTO struct {
	Id          *uuid.UUID
	UnitId      *uuid.UUID
	OrgId       *uuid.UUID
	AuthorId    *uuid.UUID
	PostType    *string
	IsPinned    *bool
	IsImportant *bool
	IsOrgWide   *bool
}

type PostCommentDTO struct {
	Id         uuid.UUID  `json:"id"`
	PostId     uuid.UUID  `json:"post_id"`
	ParentId   *uuid.UUID `json:"parent_id,omitempty"`
	AuthorId   uuid.UUID  `json:"author_id"`
	AuthorName string     `json:"author_name"`
	Content    string     `json:"content"`
	ReplyCount int        `json:"reply_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Replies []PostCommentDTO `json:"replies,omitempty"`
}

type PostPollOptionDTO struct {
	Id        uuid.UUID `json:"id"`
	PostId    uuid.UUID `json:"post_id"`
	Text      string    `json:"text"`
	VoteCount int       `json:"vote_count"`
	Urutan    int       `json:"urutan"`
	HasVoted  bool      `json:"has_voted"` // user has voted for this option
}

// Create Request DTOs
type CreatePostDTO struct {
	UnitId    uuid.UUID `json:"unit_id" validate:"required"`
	IsOrgWide bool      `json:"is_org_wide"` // true = org level, false = unit level
	Title     string    `json:"title"`
	Content   string    `json:"content" validate:"required"`
	PostType  string    `json:"post_type" validate:"required,oneof=text photo poll link"`

	ImageURL    string `json:"image_url"`
	LinkURL     string `json:"link_url"`
	LinkTitle   string `json:"link_title"`
	LinkPreview string `json:"link_preview"`

	IsPinned    bool `json:"is_pinned"`
	IsImportant bool `json:"is_important"`

	// For poll type
	PollOptions []string `json:"poll_options"`
}

type CreateCommentDTO struct {
	PostId   uuid.UUID  `json:"post_id" validate:"required"`
	ParentId *uuid.UUID `json:"parent_id"`
	Content  string     `json:"content" validate:"required"`
}

type VotePollDTO struct {
	PostId   uuid.UUID `json:"post_id" validate:"required"`
	OptionId uuid.UUID `json:"option_id" validate:"required"`
}
