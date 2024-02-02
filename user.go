package sdump

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	ErrPlanNotFound     = appError("plan does not exists")
	ErrUserNotFound     = appError("user not found")
	ErrCounterExhausted = appError("no more units left")
)

type Counter int64

func (c *Counter) Add() {
	(*c)++
}

func (c *Counter) Take() error {
	if *c <= 0 {
		return ErrCounterExhausted
	}

	(*c)--
	return nil
}

func (c *Counter) TakeN(n int64) error {
	if *c <= 0 {
		return c.Take()
	}

	(*c) -= Counter(n)
	return nil
}

type User struct {
	ID             uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()" json:"id,omitempty"`
	SSHFingerPrint string    `json:"ssh_finger_print,omitempty"`
	IsBanned       bool      `json:"is_banned,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `bun:"table:users"`
}

type FindUserOptions struct {
	SSHKeyFingerprint string
}

type UserRepository interface {
	Create(context.Context, *User) error
	Find(context.Context, *FindUserOptions) (*User, error)
}
