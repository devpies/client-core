package invites

import "time"

type Invite struct {
	ID         string    `db:"invite_id" json:"id"`
	UserID     string    `db:"user_id" json:"userId"`
	TeamID     string    `db:"team_id" json:"teamId"`
	Read       bool      `db:"read" json:"read"`
	Accepted   bool      `db:"accepted" json:"accepted"`
	Expiration time.Time `db:"expiration" json:"expiration"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}

type InviteEnhanced struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	TeamID     string    `json:"teamId"`
	TeamName   string    `json:"teamName"`
	Read       bool      `json:"read"`
	Accepted   bool      `json:"accepted"`
	Expiration time.Time `json:"expiration"`
	UpdatedAt  time.Time `json:"updatedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

type NewInvite struct {
	UserID string `json:"userId" validate:"required"`
	TeamID string `json:"teamId" validate:"required"`
}

type NewList struct {
	Emails []string `json:"emailList" validate:"required"`
}

type UpdateInvite struct {
	Accepted bool `json:"accepted" validate:"required"`
}
