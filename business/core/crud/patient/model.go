package patient

import (
	"time"

	"github.com/google/uuid"
)

// Patient represents a patient.
type Patient struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Age         int
	VideoLinks  []string
	Condition   string
	Healed      bool
	DateCreated time.Time
	DateUpdated time.Time
}

// NewPatient is what we require from clients when adding a Patient.
type NewPatient struct {
	UserID     uuid.UUID
	Name       string
	Age        int
	VideoLinks []string
	Condition  string
	Healed     bool
}

// UpdatePatient defines what information may be provided to modify an
// existing Patient. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdatePatient struct {
	Name       *string
	Age        *int
	VideoLinks []string
	Condition  *string
	Healed     *bool
}
