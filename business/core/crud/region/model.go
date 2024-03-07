package region

import (
	"time"

	"github.com/google/uuid"
)

// Region represents an individual Region.
type Region struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	DateCreated time.Time
	DateUpdated time.Time
}

// NewRegion is what we require from clients when adding a Region.
type NewRegion struct {
	UserID uuid.UUID
	Name   string
}

// UpdateRegion defines what informaton may be provided to modify an existing
// Region. All fields are optional so clients can send only the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicity blank. Normally
// we do not want to use pointers to basic types but we make exepction around
// marshalling/unmarshalling.
type UpdateRegion struct {
	Name *string
}
