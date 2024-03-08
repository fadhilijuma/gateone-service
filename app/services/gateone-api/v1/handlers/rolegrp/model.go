package rolegrp

import (
	"context"
	"github.com/fadhilijuma/gateone-service/business/core/crud/role"
	"github.com/fadhilijuma/gateone-service/business/web/v1/mid"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"time"
)

// AppRole represents information about an individual role.
type AppRole struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	Name        string `json:"name"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

func toAppRole(rl role.Role) AppRole {
	return AppRole{
		ID:          rl.ID.String(),
		UserID:      rl.UserID.String(),
		Name:        rl.Name,
		DateCreated: rl.DateCreated.Format(time.RFC3339),
		DateUpdated: rl.DateUpdated.Format(time.RFC3339),
	}
}

func toAppRoles(rls []role.Role) []AppRole {
	items := make([]AppRole, len(rls))
	for i, rl := range rls {
		items[i] = toAppRole(rl)
	}

	return items
}

// AppNewRole defines the data needed to add a new role.
type AppNewRole struct {
	Name string `json:"name" validate:"required"`
}

func toCoreNewPatient(ctx context.Context, app AppNewRole) role.NewRole {
	rl := role.NewRole{
		UserID: mid.GetUserID(ctx),
		Name:   app.Name,
	}

	return rl
}

// Validate checks the data in the model is considered clean.
func (app AppNewRole) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

// AppUpdateRole defines the data needed to update a role.
type AppUpdateRole struct {
	Name *string `json:"name"`
}

func toCoreUpdateRole(app AppUpdateRole) role.UpdateRole {
	core := role.UpdateRole{
		Name: app.Name,
	}

	return core
}

// Validate checks the data in the model is considered clean.
func (app AppUpdateRole) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}
