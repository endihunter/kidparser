package actions

import (
	"recipes/models"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	var users []models.UsersJsonResponse

	query := models.DB.Where("users.id IN (?)", 1, 2)
	query.LeftJoin("posts", "posts.user_id=users.id")
	query.GroupBy("users.id")

	sql, args := query.ToSQL(&pop.Model{Value: models.User{}}, "users.id", "name", "email", "active", "COUNT(posts.id) as posts_count")

	err := query.RawQuery(sql, args...).All(&users)

	if err != nil {
		return err
	}

	return c.Render(200, r.JSON(map[string]interface{}{
		"data": users,
	}))
}
