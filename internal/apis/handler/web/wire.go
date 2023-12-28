package web

import (
	"github.com/google/wire"
	v1 "go-chat/internal/apis/handler/web/v1"
	"go-chat/internal/apis/handler/web/v1/article"
	"go-chat/internal/apis/handler/web/v1/contact"
	"go-chat/internal/apis/handler/web/v1/group"
	"go-chat/internal/apis/handler/web/v1/talk"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(v1.Auth), "*"),
	wire.Struct(new(v1.Common), "*"),
	wire.Struct(new(v1.User), "*"),
	wire.Struct(new(v1.Organize), "*"),
	wire.Struct(new(v1.Upload), "*"),
	wire.Struct(new(v1.Emoticon), "*"),

	wire.Struct(new(contact.Contact), "*"),
	wire.Struct(new(contact.Apply), "*"),
	wire.Struct(new(contact.Group), "*"),

	wire.Struct(new(group.Group), "*"),
	wire.Struct(new(group.Apply), "*"),
	wire.Struct(new(group.Notice), "*"),

	wire.Struct(new(talk.Session), "*"),
	wire.Struct(new(talk.Message), "*"),
	wire.Struct(new(talk.Records), "*"),
	wire.Struct(new(talk.Publish), "*"),

	wire.Struct(new(article.Article), "*"),
	wire.Struct(new(article.Annex), "*"),
	wire.Struct(new(article.Class), "*"),
	wire.Struct(new(article.Tag), "*"),

	wire.Struct(new(V1), "*"),
)
