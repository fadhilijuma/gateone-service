package patient

import (
	"fmt"
	"github.com/fadhilijuma/gateone-service/foundation/validate"

	"github.com/google/uuid"
)

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ID         *uuid.UUID
	UserID     *uuid.UUID
	Name       *string `validate:"omitempty,min=3"`
	Age        *int
	Condition  *string
	Healed     *bool
	VideoLinks []string
}

// Validate can perform a check of the data against the validate tags.
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// WithPatientID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithPatientID(patientID uuid.UUID) {
	qf.ID = &patientID
}

// WithUserID sets the User ID field of the QueryFilter value.
func (qf *QueryFilter) WithUserID(userID uuid.UUID) {
	qf.ID = &userID
}

// WithName sets the Name field of the QueryFilter value.
func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}

// WithAge sets the Age field of the QueryFilter value.
func (qf *QueryFilter) WithAge(age int) {
	qf.Age = &age
}

// WithCondition sets the Condition field of the QueryFilter value.
func (qf *QueryFilter) WithCondition(condition string) {
	qf.Condition = &condition
}

// WithVideoLink sets the VideoLink field of the QueryFilter value.
func (qf *QueryFilter) WithVideoLink(videoLinks []string) {
	qf.VideoLinks = videoLinks
}

// WithHealed sets the Healed field of the QueryFilter value.
func (qf *QueryFilter) WithHealed(healed bool) {
	qf.Healed = &healed
}
