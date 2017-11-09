package grifts

import (
	"recipes/models"

	. "github.com/markbates/grift/grift"
	"github.com/markbates/pop"
)

var _ = Namespace("seed", func() {
	Desc("posts", "Seed posts table")
	Add("posts", func(c *Context) error {
		return models.DB.Transaction(func(db *pop.Connection) error {
			db.RawQuery("DELETE FROM posts WHERE 1=1")

			var posts []*models.Post

			post1 := &models.Post{
				Title:  "3 ways to organize an Event",
				Body:   "Have you ever organized en event in IT?\nDo you remember what rush is that?",
				UserID: 2,
			}

			post2 := &models.Post{
				Title:  "10 ways to not to shoot in your leg",
				Body:   "Different programming languages offer a bunch of very good staff.",
				UserID: 2,
			}

			posts = append(posts, []*models.Post{post1, post2}...)

			for _, post := range posts {
				if err := db.Create(post); err != nil {
					return err
				}
			}

			return nil
		})
	})
})
