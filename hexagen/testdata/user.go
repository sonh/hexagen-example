// Code generated by "hexagen User ../internal/usermgmt/modules/user/core/entity"; DO NOT EDIT.

package testdata

import (
	"backend/internal/usermgmt/modules/user/core/entity"
	"backend/internal/usermgmt/pkg/field"
)

type User struct {
}

// This statement will fail to compile if *User ever stops matching the interface.
var _ entity.User = (*User)(nil)

func (User User) UserID() field.String {
	return field.NewNullString()
}
func (User User) Email() field.String {
	return field.NewNullString()
}
