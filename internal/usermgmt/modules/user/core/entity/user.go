package entity

import (
	"github.com/pkg/errors"

	"backend/internal/usermgmt/pkg/field"
)

//go:generate hexagen ent-mock --type=User -o ../../../../../../mock/usermgmt
//go:generate hexagen ent-impl --type=User -o .
//go:generate hexagen ent-repo --type=User -o ../../adapter/postgres

type User interface {
	UserID() field.String
	Email() field.String
}

func ValidateUser(user User) error {
	if user.Email().String() == "" {
		return errors.New("email can not be empty")
	}

	DelegateUser(NullUser{}, WithUserEmail(field.NewNullString()))
	return nil
}

// //go:generate hexagen ent-proto --type=User --outpkg=upb -o ../../../../../../pkg/manabuf/usermgmt/

type user struct {
	userID field.String
	email  field.String
}

func (user *user) UserID() field.String {
	return user.userID
}
func (user *user) Email() field.String {
	return user.email
}

type UserOpt func(*user)

func WithUserEmail(email field.String) UserOpt {
	return func(u *user) {
		u.email = email
	}
}

func DelegateUser(userToBeDelegate User, userFieldsToDelegate ...UserOpt) User {
	user := &user{}

	if userToBeDelegate != nil {
		user.userID = userToBeDelegate.UserID()
		user.email = userToBeDelegate.Email()
	}

	for _, userFieldToDelegate := range userFieldsToDelegate {
		userFieldToDelegate(user)
	}

	return user
}
