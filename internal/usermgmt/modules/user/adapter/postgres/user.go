package postgres

import "backend/internal/usermgmt/modules/user/core/entity"

const (
	UserTable = "user"
)

const (
	UserTableUserIDColumn = "user_id"
	UserTableEmailColumn  = "email"
)

type UserRepo struct {
}

type User struct {
	entity.NullUser
}

func (user *User) TableName() string {
	return UserTable
}
