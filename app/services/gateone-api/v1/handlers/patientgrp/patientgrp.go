// Package patientgrp maintains the group of handlers for patient access.
package patientgrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"github.com/fadhilijuma/gateone-service/business/core/crud/user"
	v1 "github.com/fadhilijuma/gateone-service/business/web/v1"
	"github.com/fadhilijuma/gateone-service/business/web/v1/mid"
	"github.com/fadhilijuma/gateone-service/business/web/v1/page"
	"github.com/fadhilijuma/gateone-service/foundation/web"
	"net/http"
)

// Set of error variables for handling patient group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

type handlers struct {
	patient *patient.Core
	user    *user.Core
}

func new(patient *patient.Core, user *user.Core) *handlers {
	return &handlers{
		patient: patient,
		user:    user,
	}
}

// create adds a new patient to the system.
func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewPatient
	if err := web.Decode(r, &app); err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	pn, err := h.patient.Create(ctx, toCoreNewPatient(ctx, app))
	if err != nil {
		return fmt.Errorf("create: app[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppPatient(pn), http.StatusCreated)
}

// update updates a patient in the system.
func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdatePatient
	if err := web.Decode(r, &app); err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	pn := mid.GetPatient(ctx)

	updPn, err := h.patient.Update(ctx, pn, toCoreUpdatePatient(app))
	if err != nil {
		return fmt.Errorf("update: patientID[%s] app[%+v]: %w", updPn.ID, app, err)
	}

	return web.Respond(ctx, w, toAppPatient(updPn), http.StatusOK)
}

// delete removes a patient from the system.
func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	prd := mid.GetPatient(ctx)

	if err := h.patient.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: patientID[%s]: %w", prd.ID, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// query returns a list of patients with paging.
func (h *handlers) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := page.Parse(r)
	if err != nil {
		return err
	}

	filter, err := parseFilter(r)
	if err != nil {
		return err
	}

	orderBy, err := parseOrder(r)
	if err != nil {
		return err
	}

	prds, err := h.patient.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.patient.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, v1.NewPageDocument(toAppPatients(prds), total, page.Number, page.RowsPerPage), http.StatusOK)
}

// queryByID returns a patient by its ID.
func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.Respond(ctx, w, toAppPatient(mid.GetPatient(ctx)), http.StatusOK)
}
