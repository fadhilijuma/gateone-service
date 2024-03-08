// Package regiongrp maintains the group of handlers for region access.
package regiongrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/region"
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
	region *region.Core
	user   *user.Core
}

func new(region *region.Core, user *user.Core) *handlers {
	return &handlers{
		region: region,
		user:   user,
	}
}

// create adds a new role to the system.
func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewRegion
	if err := web.Decode(r, &app); err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	rl, err := h.region.Create(ctx, toCoreNewRegion(ctx, app))
	if err != nil {
		return fmt.Errorf("create: app[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppRegion(rl), http.StatusCreated)
}

// update updates a condition in the system.
func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateRegion
	if err := web.Decode(r, &app); err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	rn := mid.GetRegion(ctx)

	upRegion, err := h.region.Update(ctx, rn, toCoreUpdateRegion(app))
	if err != nil {
		return fmt.Errorf("update: regionID[%s] app[%+v]: %w", upRegion.ID, app, err)
	}

	return web.Respond(ctx, w, toAppRegion(upRegion), http.StatusOK)
}

// delete removes a region from the system.
func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rn := mid.GetRegion(ctx)

	if err := h.region.Delete(ctx, rn); err != nil {
		return fmt.Errorf("delete: regionID[%s]: %w", rn.ID, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// query returns a list of patients with paging.
func (h *handlers) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	pg, err := page.Parse(r)
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

	regions, err := h.region.Query(ctx, filter, orderBy, pg.Number, pg.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.region.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, v1.NewPageDocument(toAppRegions(regions), total, pg.Number, pg.RowsPerPage), http.StatusOK)
}

// queryByID returns a region by its ID.
func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.Respond(ctx, w, toAppRegion(mid.GetRegion(ctx)), http.StatusOK)
}
