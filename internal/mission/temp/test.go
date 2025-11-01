package temp

import (
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/urfave/cli/v2"
)

type TestCommand struct {
	UserRepo *repo.Users
}

func (t *TestCommand) Do(ctx *cli.Context) error {
	return nil
}
