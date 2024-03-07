package condition

import (
	"time"

	"github.com/google/uuid"
)

// Condition represents an individual condition.
type Condition struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	DateCreated time.Time
	DateUpdated time.Time
}

// NewCondition is what we require from clients when adding a condition.
type NewCondition struct {
	UserID uuid.UUID
	Name   string
}

// UpdateCondition defines what informaton may be provided to modify an existing
// condition. All fields are optional so clients can send only the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicity blank. Normally
// we do not want to use pointers to basic types but we make exepction around
// marshalling/unmarshalling.
type UpdateCondition struct {
	Name *string
}
