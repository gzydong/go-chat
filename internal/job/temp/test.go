package temp

import (
	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/repository/repo"
)

type TestCommand struct {
	UserRepo        *repo.Users
	TalkRecordsRepo *repo.TalkRecords
}

func (t *TestCommand) Run(ctx *cli.Context, conf *config.Config) error {

	// var i int32
	//
	// for {
	// 	info, err := t.TalkRecordsRepo.FindAll(ctx.Context, func(db *gorm.DB) {
	// 		db.Where("id > ?", i)
	// 		db.Where("msg_type in (?)", []int{1, 1000})
	// 		db.Limit(1000)
	// 	})
	//
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	for _, v := range info {
	// 		t.TalkRecordsRepo.UpdateWhere(ctx.Context, map[string]any{
	// 			"extra": jsonutil.Encode(model.TalkRecordExtraText{
	// 				Content: v.Content,
	// 			}),
	// 		}, "id = ?", v.Id)
	// 	}
	//
	// 	if len(info) < 1000 {
	// 		break
	// 	}
	//
	// 	i = int32(info[len(info)-1].Id)
	// }

	return nil
}
