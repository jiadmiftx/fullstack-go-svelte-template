package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post represents an announcement/posting in the system
// Scope: IsOrgWide = true → visible to all units in organization
//
//	IsOrgWide = false → visible only to the specific unit
type Post struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UnitId    uuid.UUID `gorm:"type:uuid;not null;index"`
	AuthorId  uuid.UUID `gorm:"type:uuid;not null;index"`
	IsOrgWide bool      `gorm:"default:false"` // true = org level, false = unit level

	Title    string `gorm:"type:varchar(500)"`
	Content  string `gorm:"type:text"`
	PostType string `gorm:"type:varchar(20);not null;default:'text'"` // text, photo, poll, link

	// For photo type
	ImageURL string `gorm:"type:varchar(500)"`

	// For link type
	LinkURL     string `gorm:"type:varchar(500)"`
	LinkTitle   string `gorm:"type:varchar(255)"`
	LinkPreview string `gorm:"type:text"` // description

	IsPinned     bool `gorm:"default:false"`
	IsImportant  bool `gorm:"default:false"`
	CommentCount int  `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	Unit     *Unit            `gorm:"foreignKey:UnitId"`
	Author   *User            `gorm:"foreignKey:AuthorId"`
	Comments []PostComment    `gorm:"foreignKey:PostId"`
	Options  []PostPollOption `gorm:"foreignKey:PostId"`
}

func (Post) TableName() string { return "posts" }

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if p.Id == uuid.Nil {
		p.Id = uuid.New()
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return
}

func (p *Post) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return
}

// PostComment represents a comment on a post
// ParentId = null means top-level comment
// ParentId = uuid means reply to another comment (single nesting only)
type PostComment struct {
	Id         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	PostId     uuid.UUID  `gorm:"type:uuid;not null;index"`
	ParentId   *uuid.UUID `gorm:"type:uuid;index"` // null = top-level, not null = reply
	AuthorId   uuid.UUID  `gorm:"type:uuid;not null;index"`
	Content    string     `gorm:"type:text;not null"`
	ReplyCount int        `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	Post    *Post         `gorm:"foreignKey:PostId"`
	Parent  *PostComment  `gorm:"foreignKey:ParentId"`
	Author  *User         `gorm:"foreignKey:AuthorId"`
	Replies []PostComment `gorm:"foreignKey:ParentId"`
}

func (PostComment) TableName() string { return "post_comments" }

func (c *PostComment) BeforeCreate(tx *gorm.DB) (err error) {
	if c.Id == uuid.Nil {
		c.Id = uuid.New()
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *PostComment) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}

// PostPollOption represents an option in a poll post
type PostPollOption struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostId    uuid.UUID `gorm:"type:uuid;not null;index"`
	Text      string    `gorm:"type:varchar(255);not null"`
	VoteCount int       `gorm:"default:0"`
	Urutan    int       `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Relations
	Post  *Post          `gorm:"foreignKey:PostId"`
	Votes []PostPollVote `gorm:"foreignKey:OptionId"`
}

func (PostPollOption) TableName() string { return "post_poll_options" }

func (o *PostPollOption) BeforeCreate(tx *gorm.DB) (err error) {
	if o.Id == uuid.Nil {
		o.Id = uuid.New()
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return
}

// PostPollVote represents a user's vote on a poll
type PostPollVote struct {
	Id       uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostId   uuid.UUID `gorm:"type:uuid;not null;index"`
	OptionId uuid.UUID `gorm:"type:uuid;not null;index"`
	UserId   uuid.UUID `gorm:"type:uuid;not null;index"`

	CreatedAt time.Time

	// Relations
	Post   *Post           `gorm:"foreignKey:PostId"`
	Option *PostPollOption `gorm:"foreignKey:OptionId"`
	User   *User           `gorm:"foreignKey:UserId"`
}

func (PostPollVote) TableName() string { return "post_poll_votes" }

func (v *PostPollVote) BeforeCreate(tx *gorm.DB) (err error) {
	if v.Id == uuid.Nil {
		v.Id = uuid.New()
	}
	v.CreatedAt = time.Now()
	return
}

// Post type constants
const (
	PostTypeText  = "text"
	PostTypePhoto = "photo"
	PostTypePoll  = "poll"
	PostTypeLink  = "link"
)
