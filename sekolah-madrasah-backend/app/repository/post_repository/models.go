package post_repository

import (
	"time"

	"github.com/google/uuid"
)

// Post represents the domain model
type Post struct {
	Id         uuid.UUID
	UnitId     uuid.UUID
	AuthorId   uuid.UUID
	AuthorName string // populated from User
	IsOrgWide  bool   // true = org level, false = unit level

	Title    string
	Content  string
	PostType string // text, photo, poll, link

	ImageURL    string
	LinkURL     string
	LinkTitle   string
	LinkPreview string

	IsPinned     bool
	IsImportant  bool
	CommentCount int

	CreatedAt time.Time
	UpdatedAt time.Time

	// Relations (populated on read)
	Options  []PostPollOption
	Comments []PostComment
}

type PostComment struct {
	Id         uuid.UUID
	PostId     uuid.UUID
	ParentId   *uuid.UUID
	AuthorId   uuid.UUID
	AuthorName string
	Content    string
	ReplyCount int

	CreatedAt time.Time
	UpdatedAt time.Time

	Replies []PostComment // populated for top-level comments only
}

type PostPollOption struct {
	Id        uuid.UUID
	PostId    uuid.UUID
	Text      string
	VoteCount int
	Urutan    int
	HasVoted  bool // populated per user
}

type PostPollVote struct {
	Id       uuid.UUID
	PostId   uuid.UUID
	OptionId uuid.UUID
	UserId   uuid.UUID
}

// Filters
type PostFilter struct {
	Id          *uuid.UUID
	UnitId      *uuid.UUID
	OrgId       *uuid.UUID // for org-wide posts
	AuthorId    *uuid.UUID
	PostType    *string
	IsPinned    *bool
	IsImportant *bool
	IsOrgWide   *bool // filter by org-wide or unit-specific
}

type PostCommentFilter struct {
	Id       *uuid.UUID
	PostId   *uuid.UUID
	ParentId *uuid.UUID
	AuthorId *uuid.UUID
}

type PostPollOptionFilter struct {
	Id     *uuid.UUID
	PostId *uuid.UUID
}

type PostPollVoteFilter struct {
	Id       *uuid.UUID
	PostId   *uuid.UUID
	OptionId *uuid.UUID
	UserId   *uuid.UUID
}
