package queue

import "github.com/google/wire"

type Consumers struct {
	UserLoginConsumer *UserLoginConsumer
}

var ProviderSet = wire.NewSet(
	wire.Struct(new(Consumers), "*"),
	wire.Struct(new(UserLoginConsumer), "*"),
)
