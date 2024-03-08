package regiongrp

import (
	"context"
	"github.com/fadhilijuma/gateone-service/business/core/crud/region"
	"github.com/fadhilijuma/gateone-service/business/web/v1/mid"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"time"
)

// AppRegion represents information about an individual region.
type AppRegion struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	Name        string `json:"name"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

func toAppRegion(rn region.Region) AppRegion {
	return AppRegion{
		ID:          rn.ID.String(),
		UserID:      rn.UserID.String(),
		Name:        rn.Name,
		DateCreated: rn.DateCreated.Format(time.RFC3339),
		DateUpdated: rn.DateUpdated.Format(time.RFC3339),
	}
}

func toAppRegions(rns []region.Region) []AppRegion {
	items := make([]AppRegion, len(rns))
	for i, rn := range rns {
		items[i] = toAppRegion(rn)
	}

	return items
}

// AppNewRegion defines the data needed to add a new region.
type AppNewRegion struct {
	Name string `json:"name" validate:"required"`
}

func toCoreNewRegion(ctx context.Context, app AppNewRegion) region.NewRegion {
	rn := region.NewRegion{
		UserID: mid.GetUserID(ctx),
		Name:   app.Name,
	}

	return rn
}

// Validate checks the data in the model is considered clean.
func (app AppNewRegion) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

// AppUpdateRegion defines the data needed to update a role.
type AppUpdateRegion struct {
	Name *string `json:"name"`
}

func toCoreUpdateRegion(app AppUpdateRegion) region.UpdateRegion {
	core := region.UpdateRegion{
		Name: app.Name,
	}

	return core
}

// Validate checks the data in the model is considered clean.
func (app AppUpdateRegion) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}
