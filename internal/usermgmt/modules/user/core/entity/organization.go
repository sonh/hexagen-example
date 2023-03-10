package entity

import "backend/internal/usermgmt/pkg/field"

type HasOrganizationID interface {
	OrganizationID() field.String
}

type Organization interface {
	OrganizationID() field.String
}
