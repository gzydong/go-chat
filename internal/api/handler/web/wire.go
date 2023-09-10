package web

import (
	"github.com/google/wire"
	v12 "go-chat/internal/api/handler/web/v1"
	article2 "go-chat/internal/api/handler/web/v1/article"
	contact2 "go-chat/internal/api/handler/web/v1/contact"
	group2 "go-chat/internal/api/handler/web/v1/group"
	talk2 "go-chat/internal/api/handler/web/v1/talk"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(v12.Auth), "*"),
	wire.Struct(new(v12.Common), "*"),
	wire.Struct(new(v12.User), "*"),
	wire.Struct(new(v12.Organize), "*"),
	wire.Struct(new(v12.Upload), "*"),
	wire.Struct(new(v12.Emoticon), "*"),

	wire.Struct(new(contact2.Contact), "*"),
	wire.Struct(new(contact2.Apply), "*"),
	wire.Struct(new(contact2.Group), "*"),

	wire.Struct(new(group2.Group), "*"),
	wire.Struct(new(group2.Apply), "*"),
	wire.Struct(new(group2.Notice), "*"),

	wire.Struct(new(talk2.Session), "*"),
	wire.Struct(new(talk2.Message), "*"),
	wire.Struct(new(talk2.Records), "*"),
	wire.Struct(new(talk2.Publish), "*"),

	wire.Struct(new(article2.Article), "*"),
	wire.Struct(new(article2.Annex), "*"),
	wire.Struct(new(article2.Class), "*"),
	wire.Struct(new(article2.Tag), "*"),

	wire.Struct(new(V1), "*"),
)
