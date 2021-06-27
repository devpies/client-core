package users

import (
	"time"
)

// User represent user data from the database.
type User struct {
	ID            string    `db:"user_id" json:"id"`
	Auth0ID       string    `db:"auth0_id" json:"auth0Id"`
	Email         string    `db:"email" json:"email"`
	EmailVerified bool      `db:"email_verified" json:"emailVerified"`
	FirstName     *string   `db:"first_name" json:"firstName"`
	LastName      *string   `db:"last_name" json:"lastName"`
	Picture       *string   `db:"picture" json:"picture"`
	Locale        *string   `db:"locale" json:"locale"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
}

type NewUser struct {
	Auth0ID       string  `json:"auth0Id" validate:"required"`
	Email         string  `json:"email" validate:"required"`
	FirstName     *string `json:"firstName" validate:"required"`
	EmailVerified bool    `json:"emailVerified"`
	LastName      *string `json:"lastName"`
	Picture       *string `json:"picture"`
	Locale        *string `json:"locale"`
}

type UpdateUser struct {
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
	Picture   *string   `json:"picture"`
	Locale    *string   `json:"locale"`
	UpdatedAt time.Time `json:"updatedAt"`
}
