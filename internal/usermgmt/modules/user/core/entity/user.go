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
	OrganizationID() field.String
}

func ValidateUser(user User) error {
	if user.Email().String() == "" {
		return errors.New("email can not be empty")
	}

	return nil
}
