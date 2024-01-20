package sdump

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type URLEndpointMetadata struct{}

type URLEndpoint struct {
	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()" json:"id,omitempty"`
	Reference string    `json:"reference,omitempty"`
	IsActive  bool      `json:"is_active,omitempty"`

	Metadata URLEndpointMetadata `json:"metadata,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `bun:"table:urls"`
}

func NewURLEndpoint() *URLEndpoint {
	return &URLEndpoint{
		Reference: xid.New().String(),
		IsActive:  true,
	}
}

type URLRepository interface {
	Create(context.Context, *URLEndpoint) error
}
