package patientgrp

import (
	"context"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"github.com/fadhilijuma/gateone-service/business/web/v1/mid"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"time"
)

// AppPatient represents information about an individual patient.
type AppPatient struct {
	ID          string   `json:"id"`
	UserID      string   `json:"userID"`
	Name        string   `json:"name"`
	Age         int      `json:"age"`
	VideoLinks  []string `json:"video_links"`
	Condition   string   `json:"condition"`
	Healed      bool     `json:"healed"`
	DateCreated string   `json:"dateCreated"`
	DateUpdated string   `json:"dateUpdated"`
}

func toAppPatient(pn patient.Patient) AppPatient {
	return AppPatient{
		ID:          pn.ID.String(),
		UserID:      pn.UserID.String(),
		Name:        pn.Name,
		Age:         pn.Age,
		VideoLinks:  pn.VideoLinks,
		Condition:   pn.Condition,
		Healed:      pn.Healed,
		DateCreated: pn.DateCreated.Format(time.RFC3339),
		DateUpdated: pn.DateUpdated.Format(time.RFC3339),
	}
}

func toAppPatients(pns []patient.Patient) []AppPatient {
	items := make([]AppPatient, len(pns))
	for i, pn := range pns {
		items[i] = toAppPatient(pn)
	}

	return items
}

// AppNewPatient defines the data needed to add a new patient.
type AppNewPatient struct {
	Name       string   `json:"name" validate:"required"`
	Age        int      `json:"age" validate:"required"`
	VideoLinks []string `json:"video_links" validate:"required"`
	Condition  string   `json:"condition" validate:"required"`
	Healed     bool     `json:"healed" validate:"required"`
}

func toCoreNewPatient(ctx context.Context, app AppNewPatient) patient.NewPatient {
	pn := patient.NewPatient{
		UserID:     mid.GetUserID(ctx),
		Name:       app.Name,
		Age:        app.Age,
		VideoLinks: app.VideoLinks,
		Condition:  app.Condition,
		Healed:     app.Healed,
	}

	return pn
}

// Validate checks the data in the model is considered clean.
func (app AppNewPatient) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

// AppUpdatePatient defines the data needed to update a patient.
type AppUpdatePatient struct {
	Name       *string  `json:"name"`
	Age        *int     `json:"age"`
	VideoLinks []string `json:"video_links"`
	Condition  *string  `json:"condition"`
	Healed     *bool    `json:"healed"`
}

func toCoreUpdatePatient(app AppUpdatePatient) patient.UpdatePatient {
	core := patient.UpdatePatient{
		Name:       app.Name,
		Age:        app.Age,
		VideoLinks: app.VideoLinks,
		Condition:  app.Condition,
		Healed:     app.Healed,
	}

	return core
}

// Validate checks the data in the model is considered clean.
func (app AppUpdatePatient) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}
