package queue

import (
	"context"
	"encoding/json"

	"go-chat/api/pb/queue/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/consumer"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"go-chat/internal/service/message"
)

var _ consumer.IConsumerHandle = (*UserLoginConsumer)(nil)

type UserLoginConsumer struct {
	RobotRepo          *repo.Robot
	IpAddressService   service.IIpAddressService
	TalkSessionService service.ITalkSessionService
	Message            message.IService
}

func (u *UserLoginConsumer) Touch() bool {
	return false
}

func (u *UserLoginConsumer) Topic() string {
	return entity.LoginTopic
}

func (u *UserLoginConsumer) Channel() string {
	return "default"
}

func (u *UserLoginConsumer) Do(ctx context.Context, msg []byte, attempts uint16) error {
	var in queue.UserLoginRequest
	if err := json.Unmarshal(msg, &in); err != nil {
		return err
	}

	root, err := u.RobotRepo.GetLoginRobot(ctx)
	if err != nil {
		return nil
	}

	if root == nil {
		return nil
	}

	address, err := u.IpAddressService.FindAddress(in.IpAddr)
	if err != nil {
		return nil
	}

	_, _ = u.TalkSessionService.Create(ctx, &service.TalkSessionCreateOpt{
		UserId:     int(in.UserId),
		TalkType:   entity.ChatPrivateMode,
		ReceiverId: root.UserId,
		IsBoot:     true,
	})

	// 推送登录消息
	return u.Message.CreateLoginMessage(ctx, message.CreateLoginMessageOption{
		UserId:   int(in.UserId),
		Ip:       in.IpAddr,
		Address:  address,
		Platform: in.Platform,
		Agent:    in.Agent,
		Reason:   "常用设备登录",
		LoginAt:  in.LoginAt,
	})
}
