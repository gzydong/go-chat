package temp

import (
	"github.com/urfave/cli/v2"
	"go-chat/internal/repository/repo"
)

type TestCommand struct {
	UserRepo *repo.Users
}

func (t *TestCommand) Do(ctx *cli.Context) error {
	return nil
}
