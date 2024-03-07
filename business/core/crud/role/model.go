package role

import (
	"time"

	"github.com/google/uuid"
)

// Role represents an individual role.
type Role struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	DateCreated time.Time
	DateUpdated time.Time
}

// NewRole is what we require from clients when adding a Role.
type NewRole struct {
	UserID uuid.UUID
	Name   string
}

// UpdateRole defines what informaton may be provided to modify an existing
// Role. All fields are optional so clients can send only the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicity blank. Normally
// we do not want to use pointers to basic types but we make exepction around
// marshalling/unmarshalling.
type UpdateRole struct {
	Name *string
}
