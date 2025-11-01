package admin

import (
	"github.com/gzydong/go-chat/internal/apis/handler/admin/system"
	"github.com/gzydong/go-chat/internal/apis/handler/admin/user"
	"github.com/gzydong/go-chat/internal/repository/repo"
)

type Handler struct {
	Auth      *Auth
	Totp      *Totp
	Admin     *system.Admin
	Role      *system.Role
	Resource  *system.Resource
	Menu      *system.Menu
	AdminRepo *repo.Admin
	User      *user.User
}
