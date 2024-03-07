// Package conditiongrp Package productgrp maintains the group of handlers for condition access.
package rolegrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/role"
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
	role *role.Core
	user *user.Core
}

func new(role *role.Core, user *user.Core) *handlers {
	return &handlers{
		role: role,
		user: user,
	}
}

// create adds a new role to the system.
func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewRole
	if err := web.Decode(r, &app); err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	rl, err := h.role.Create(ctx, toCoreNewPatient(ctx, app))
	if err != nil {
		return fmt.Errorf("create: app[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppRole(rl), http.StatusCreated)
}

// update updates a condition in the system.
func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateRole
	if err := web.Decode(r, &app); err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	rl := mid.GetRole(ctx)

	upRole, err := h.role.Update(ctx, rl, toCoreUpdateRole(app))
	if err != nil {
		return fmt.Errorf("update: roleID[%s] app[%+v]: %w", upRole.ID, app, err)
	}

	return web.Respond(ctx, w, toAppRole(upRole), http.StatusOK)
}

// delete removes a product from the system.
func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rl := mid.GetRole(ctx)

	if err := h.role.Delete(ctx, rl); err != nil {
		return fmt.Errorf("delete: roleID[%s]: %w", rl.ID, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// query returns a list of products with paging.
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

	roles, err := h.role.Query(ctx, filter, orderBy, pg.Number, pg.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.role.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, v1.NewPageDocument(toAppRoles(roles), total, pg.Number, pg.RowsPerPage), http.StatusOK)
}

// queryByID returns a product by its ID.
func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.Respond(ctx, w, toAppRole(mid.GetRole(ctx)), http.StatusOK)
}
