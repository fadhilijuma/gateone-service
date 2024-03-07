// Package condition provides a business access to condition data in the system.
package condition

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
	ErrNotFound     = errors.New("condition not found")
	ErrUserDisabled = errors.New("user disabled")
)

// Storer interface declares the behaviour this package needs to persist and
// retrieve data.
type Storer interface {
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
	Create(ctx context.Context, cn Condition) error
	Update(ctx context.Context, cn Condition) error
	Delete(ctx context.Context, cn Condition) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Condition, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, conditionID uuid.UUID) (Condition, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Condition, error)
}

// Core manages the set of APIs for condition api access.
type Core struct {
	log      *logger.Logger
	usrCore  *user.Core
	delegate *delegate.Delegate
	storer   Storer
}

// NewCore constructs a condition core API for use.
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

// Create adds a new condition to the system.
func (c *Core) Create(ctx context.Context, nr NewCondition) (Condition, error) {
	usr, err := c.usrCore.QueryByID(ctx, nr.UserID)
	if err != nil {
		return Condition{}, fmt.Errorf("user.querybyid: %s: %w", nr.UserID, err)
	}

	if !usr.Enabled {
		return Condition{}, ErrUserDisabled
	}

	now := time.Now()

	hme := Condition{
		ID:          uuid.New(),
		Name:        nr.Name,
		UserID:      nr.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, hme); err != nil {
		return Condition{}, fmt.Errorf("create: %w", err)
	}

	return hme, nil
}

// Update modifies information about a condition.
func (c *Core) Update(ctx context.Context, condition Condition, ur UpdateCondition) (Condition, error) {
	if ur.Name != nil {
		condition.Name = *ur.Name
	}

	condition.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, condition); err != nil {
		return Condition{}, fmt.Errorf("update: %w", err)
	}

	return condition, nil
}

// Delete removes the specified condition.
func (c *Core) Delete(ctx context.Context, cn Condition) error {
	if err := c.storer.Delete(ctx, cn); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing conditions.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Condition, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	hmes, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return hmes, nil
}

// Count returns the total number of conditions.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}

	return c.storer.Count(ctx, filter)
}

// QueryByID finds the condition by the specified ID.
func (c *Core) QueryByID(ctx context.Context, conditionID uuid.UUID) (Condition, error) {
	hme, err := c.storer.QueryByID(ctx, conditionID)
	if err != nil {
		return Condition{}, fmt.Errorf("query: conditionID[%s]: %w", conditionID, err)
	}

	return hme, nil
}

// QueryByUserID finds the conditions by a specified User ID.
func (c *Core) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Condition, error) {
	conditions, err := c.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return conditions, nil
}
