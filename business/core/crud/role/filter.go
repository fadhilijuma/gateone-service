package role

import (
	"fmt"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"time"

	"github.com/google/uuid"
)

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ID               *uuid.UUID
	UserID           *uuid.UUID
	Name             *string
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}

// Validate can perform a check of the data against the validate tags.
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// WithRoleID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithRoleID(roleID uuid.UUID) {
	qf.ID = &roleID
}

// WithUserID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithUserID(userID uuid.UUID) {
	qf.UserID = &userID
}

// WithRoleName sets the name field of the QueryFilter value.
func (qf *QueryFilter) WithRoleName(name string) {
	qf.Name = &name
}

// WithStartDateCreated sets the StartCreatedDate field of the QueryFilter value.
func (qf *QueryFilter) WithStartDateCreated(startDate time.Time) {
	d := startDate.UTC()
	qf.StartCreatedDate = &d
}

// WithEndCreatedDate sets the EndCreatedDate field of the QueryFilter value.
func (qf *QueryFilter) WithEndCreatedDate(endDate time.Time) {
	d := endDate.UTC()
	qf.EndCreatedDate = &d
}
