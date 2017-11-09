package grifts

import (
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {

		grift.Run("seed:users", c)
		grift.Run("seed:posts", c)

		return nil
	})

})
