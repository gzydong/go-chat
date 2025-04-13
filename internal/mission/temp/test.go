package temp

import (
	"log"

	"github.com/urfave/cli/v2"
	"go-chat/internal/repository/repo"
)

type TestCommand struct {
	UserRepo *repo.Users
}

func (t *TestCommand) Do(ctx *cli.Context) error {
	log.Println("test command")
	return nil
}
