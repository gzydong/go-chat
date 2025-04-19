package entity

const (
	JwtIssuerWeb   = "web"
	JwtIssuerAdmin = "admin"
)

type WebClaims struct {
	UserId int32 `json:"user_id"`
}

func (w *WebClaims) GetAuthID() int {
	return int(w.UserId)
}

type AdminClaims struct {
	AdminId int32 `json:"admin_id"`
}

func (a *AdminClaims) GetAuthID() int {
	return int(a.AdminId)
}
