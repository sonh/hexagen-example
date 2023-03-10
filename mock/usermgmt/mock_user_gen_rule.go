package entity

import (
	"fmt"

	"backend/internal/usermgmt/pkg/field"
)

func GenerateUserEmail(randomID string) field.String {
	return field.NewString(fmt.Sprintf("email+%s@example.com", randomID))
}
