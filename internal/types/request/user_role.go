package request

type UserRoles struct {
	UserID  string   `json:"userId" form:"userId" binding:"required"`
	RoleIDs []string `json:"roleIDs" form:"roleIDs" binding:"required"`
}
