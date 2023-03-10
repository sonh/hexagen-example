package entity

import "backend/internal/usermgmt/pkg/field"

type user struct {
	userID field.String
	email  field.String
	HasOrganizationID
}

func (user *user) UserID() field.String {
	return user.userID
}
func (user *user) Email() field.String {
	return user.email
}
func (user *user) OrganizationID() field.String {
	return user.email
}

type EntOpt func(*user)

func WithOrganizationID(organizationID HasOrganizationID) EntOpt {
	return func(u *user) {
		u.HasOrganizationID = organizationID
	}
}

func DelegateUser(userToBeDelegate User, userFieldsToDelegate ...EntOpt) User {
	user := &user{}

	if userToBeDelegate != nil {
		user.userID = userToBeDelegate.UserID()
		user.email = userToBeDelegate.Email()
		user.HasOrganizationID = userToBeDelegate
	}

	for _, userFieldToDelegate := range userFieldsToDelegate {
		userFieldToDelegate(user)
	}

	return user
}
