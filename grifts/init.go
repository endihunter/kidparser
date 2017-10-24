package grifts

import (
  "github.com/gobuffalo/buffalo"
	"recipes/actions"
)

func init() {
  buffalo.Grifts(actions.App())
}
