package admin

import (
	"go-chat/internal/apis/handler/admin/system"
	"go-chat/internal/repository/repo"
)

type Handler struct {
	Auth      *Auth
	Totp      *Totp
	Admin     *system.Admin
	Role      *system.Role
	Resource  *system.Resource
	Menu      *system.Menu
	AdminRepo *repo.Admin
}
