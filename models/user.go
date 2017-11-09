package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
)

type Timestamps struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Email  string `json:"email" db:"email"`
	Active bool   `json:"active" db:"active"`
	Timestamps
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

type UsersJsonResponse struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Email      string `json:"email" db:"email"`
	Active     bool   `json:"active" db:"active"`
	PostsCount int    `json:"posts_count" db:"posts_count"`
}

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "Title"},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
