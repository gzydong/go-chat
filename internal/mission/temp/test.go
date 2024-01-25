package temp

import (
	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/repository/repo"
)

type TestCommand struct {
	UserRepo *repo.Users
}

func (t *TestCommand) Run(ctx *cli.Context, conf *config.Config) error {
	return nil
}
