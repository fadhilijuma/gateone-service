package conditiongrp

import (
	"context"
	"github.com/fadhilijuma/gateone-service/business/core/crud/condition"
	"github.com/fadhilijuma/gateone-service/business/web/v1/mid"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"time"
)

// AppCondition represents information about an individual condition.
type AppCondition struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	Name        string `json:"name"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

func toAppCondition(pn condition.Condition) AppCondition {
	return AppCondition{
		ID:          pn.ID.String(),
		UserID:      pn.UserID.String(),
		Name:        pn.Name,
		DateCreated: pn.DateCreated.Format(time.RFC3339),
		DateUpdated: pn.DateUpdated.Format(time.RFC3339),
	}
}

func toAppConditions(pns []condition.Condition) []AppCondition {
	items := make([]AppCondition, len(pns))
	for i, pn := range pns {
		items[i] = toAppCondition(pn)
	}

	return items
}

// AppNewCondition defines the data needed to add a new condition.
type AppNewCondition struct {
	Name string `json:"name" validate:"required"`
}

func toCoreNewCondition(ctx context.Context, app AppNewCondition) condition.NewCondition {
	pn := condition.NewCondition{
		UserID: mid.GetUserID(ctx),
		Name:   app.Name,
	}

	return pn
}

// Validate checks the data in the model is considered clean.
func (app AppNewCondition) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

// AppUpdateCondition defines the data needed to update a condition.
type AppUpdateCondition struct {
	Name *string `json:"name"`
}

func toCoreUpdateCondition(app AppUpdateCondition) condition.UpdateCondition {
	core := condition.UpdateCondition{
		Name: app.Name,
	}

	return core
}

// Validate checks the data in the model is considered clean.
func (app AppUpdateCondition) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}
