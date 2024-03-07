// Package region provides a business access to Region data in the system.
package region

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/delegate"
	"github.com/fadhilijuma/gateone-service/business/core/crud/user"
	"github.com/fadhilijuma/gateone-service/business/data/transaction"
	"github.com/fadhilijuma/gateone-service/business/web/v1/order"
	"github.com/fadhilijuma/gateone-service/foundation/logger"
	"time"

	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound     = errors.New("region not found")
	ErrUserDisabled = errors.New("user disabled")
)

// Storer interface declares the behaviour this package needs to persist and
// retrieve data.
type Storer interface {
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
	Create(ctx context.Context, rn Region) error
	Update(ctx context.Context, rn Region) error
	Delete(ctx context.Context, rn Region) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Region, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, RegionID uuid.UUID) (Region, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Region, error)
}

// Core manages the set of APIs for Region api access.
type Core struct {
	log      *logger.Logger
	usrCore  *user.Core
	delegate *delegate.Delegate
	storer   Storer
}

// NewCore constructs a Region core API for use.
func NewCore(log *logger.Logger, usrCore *user.Core, delegate *delegate.Delegate, storer Storer) *Core {
	return &Core{
		log:      log,
		usrCore:  usrCore,
		delegate: delegate,
		storer:   storer,
	}
}

// ExecuteUnderTransaction constructs a new Core value that will use the
// specified transaction in any store related calls.
func (c *Core) ExecuteUnderTransaction(tx transaction.Transaction) (*Core, error) {
	storer, err := c.storer.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	usrCore, err := c.usrCore.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	core := Core{
		log:      c.log,
		usrCore:  usrCore,
		delegate: c.delegate,
		storer:   storer,
	}

	return &core, nil
}

// Create adds a new Region to the system.
func (c *Core) Create(ctx context.Context, nr NewRegion) (Region, error) {
	usr, err := c.usrCore.QueryByID(ctx, nr.UserID)
	if err != nil {
		return Region{}, fmt.Errorf("user.querybyid: %s: %w", nr.UserID, err)
	}

	if !usr.Enabled {
		return Region{}, ErrUserDisabled
	}

	now := time.Now()

	hme := Region{
		ID:          uuid.New(),
		Name:        nr.Name,
		UserID:      nr.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, hme); err != nil {
		return Region{}, fmt.Errorf("create: %w", err)
	}

	return hme, nil
}

// Update modifies information about a Region.
func (c *Core) Update(ctx context.Context, rn Region, ur UpdateRegion) (Region, error) {
	if ur.Name != nil {
		rn.Name = *ur.Name
	}

	rn.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, rn); err != nil {
		return Region{}, fmt.Errorf("update: %w", err)
	}

	return rn, nil
}

// Delete removes the specified Region.
func (c *Core) Delete(ctx context.Context, cn Region) error {
	if err := c.storer.Delete(ctx, cn); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing Regions.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Region, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	hmes, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return hmes, nil
}

// Count returns the total number of Regions.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}

	return c.storer.Count(ctx, filter)
}

// QueryByID finds the Region by the specified ID.
func (c *Core) QueryByID(ctx context.Context, RegionID uuid.UUID) (Region, error) {
	hme, err := c.storer.QueryByID(ctx, RegionID)
	if err != nil {
		return Region{}, fmt.Errorf("query: RegionID[%s]: %w", RegionID, err)
	}

	return hme, nil
}

// QueryByUserID finds the Regions by a specified User ID.
func (c *Core) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Region, error) {
	Regions, err := c.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return Regions, nil
}
