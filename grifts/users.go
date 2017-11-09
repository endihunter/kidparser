package grifts

import (
	"recipes/models"

	. "github.com/markbates/grift/grift"
	"github.com/markbates/pop"
)

var _ = Namespace("seed", func() {
	Desc("users", "Seed users table")
	Add("users", func(c *Context) error {
		return models.DB.Transaction(func(db *pop.Connection) error {
			db.RawQuery("DELETE FROM users WHERE 1=1")

			users := []*models.User{}

			user1 := &models.User{
				ID:     1,
				Name:   "Administrator",
				Email:  "admin@example.com",
				Active: true,
			}

			user2 := &models.User{
				ID:     2,
				Name:   "Manager",
				Email:  "manager@example.com",
				Active: true,
			}

			users = append(users, []*models.User{user1, user2}...)

			for _, u := range users {
				if err := db.Create(u); err != nil {
					return err
				}
			}

			return nil
		})

		return nil
	})
})
