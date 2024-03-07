// Package conditiongrp Package productgrp maintains the group of handlers for condition access.
package conditiongrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/condition"
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
	condition *condition.Core
	user      *user.Core
}

func new(condition *condition.Core, user *user.Core) *handlers {
	return &handlers{
		condition: condition,
		user:      user,
	}
}

// create adds a new product to the system.
func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewCondition
	if err := web.Decode(r, &app); err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	cn, err := h.condition.Create(ctx, toCoreNewCondition(ctx, app))
	if err != nil {
		return fmt.Errorf("create: app[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppCondition(cn), http.StatusCreated)
}

// update updates a condition in the system.
func (h *handlers) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateCondition
	if err := web.Decode(r, &app); err != nil {
		return v1.NewTrustedError(err, http.StatusBadRequest)
	}

	cn := mid.GetCondition(ctx)

	updPn, err := h.condition.Update(ctx, cn, toCoreUpdateCondition(app))
	if err != nil {
		return fmt.Errorf("update: productID[%s] app[%+v]: %w", updPn.ID, app, err)
	}

	return web.Respond(ctx, w, toAppCondition(updPn), http.StatusOK)
}

// delete removes a product from the system.
func (h *handlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cn := mid.GetCondition(ctx)

	if err := h.condition.Delete(ctx, cn); err != nil {
		return fmt.Errorf("delete: conditionID[%s]: %w", cn.ID, err)
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

	cns, err := h.condition.Query(ctx, filter, orderBy, pg.Number, pg.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.condition.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, v1.NewPageDocument(toAppConditions(cns), total, pg.Number, pg.RowsPerPage), http.StatusOK)
}

// queryByID returns a product by its ID.
func (h *handlers) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.Respond(ctx, w, toAppCondition(mid.GetCondition(ctx)), http.StatusOK)
}
