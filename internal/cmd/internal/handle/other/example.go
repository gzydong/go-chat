package other

import (
	"context"
	"fmt"
	"time"

	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type ExampleHandle struct {
	db       *gorm.DB
	sequence *repo.Sequence
}

func NewExampleHandle(db *gorm.DB, sequence *repo.Sequence) *ExampleHandle {
	return &ExampleHandle{db, sequence}
}

var uids = []int64{1046,
	1058,
	1070,
	1149,
	1887,
	1903,
	2001,
	2016,
	2022,
	2053,
	2055,
	3019,
	3045,
	3054,
	3063,
	3072,
	3084,
	3117,
	4018,
	4031,
	4042,
	4044,
	4045,
	4046,
	4047,
	4048,
	4066,
	4070,
	4071,
	4083,
	4099,
	4100,
	4101,
	4102,
	4103,
	4107,
	4108,
	4109,
	4162,
	4192,
	4214,
	4216,
	4218,
	4225,
	4240,
	4254,
	4261,
	4262,
	4268,
	4279,
	4283,
	4292,
	4362,
	4364,
	4365,
	4370,
	4374,
	4409}

func (e *ExampleHandle) Handle(ctx context.Context) error {

	fmt.Println("Job ExampleHandle Start")

	ex := time.Now().AddDate(0, 0, 10)

	var items []string

	for _, uid := range uids {
		// 生成登录凭证
		token := jwt.GenerateToken("api", "836c3fea9bba4e04d51bd0fbcc5", &jwt.Options{
			ExpiresAt: jwt.NewNumericDate(ex),
			ID:        fmt.Sprintf("%d", uid),
			Issuer:    "im.web",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		})

		items = append(items, token)
	}

	fmt.Println(jsonutil.Encode(items))

	fmt.Println("Job ExampleHandle End")

	return nil
}
