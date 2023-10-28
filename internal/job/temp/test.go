package temp

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/repository/repo"
)

type TestCommand struct {
	UserRepo *repo.Users
}

func (t *TestCommand) Run(ctx *cli.Context, conf *config.Config) error {

	user, err := t.UserRepo.FindByMobile("18798272054")

	fmt.Println("开始测试了....", user, err)
	return nil
}
