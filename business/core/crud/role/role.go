// Package role provides a business access to role data in the system.
package role

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
	ErrNotFound     = errors.New("role not found")
	ErrUserDisabled = errors.New("user disabled")
)

// Storer interface declares the behaviour this package needs to persist and
// retrieve data.
type Storer interface {
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
	Create(ctx context.Context, hme Role) error
	Update(ctx context.Context, hme Role) error
	Delete(ctx context.Context, hme Role) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Role, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, roleID uuid.UUID) (Role, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Role, error)
}

// Core manages the set of APIs for role api access.
type Core struct {
	log      *logger.Logger
	usrCore  *user.Core
	delegate *delegate.Delegate
	storer   Storer
}

// NewCore constructs a role core API for use.
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

// Create adds a new role to the system.
func (c *Core) Create(ctx context.Context, nr NewRole) (Role, error) {
	usr, err := c.usrCore.QueryByID(ctx, nr.UserID)
	if err != nil {
		return Role{}, fmt.Errorf("user.querybyid: %s: %w", nr.UserID, err)
	}

	if !usr.Enabled {
		return Role{}, ErrUserDisabled
	}

	now := time.Now()

	hme := Role{
		ID:          uuid.New(),
		Name:        nr.Name,
		UserID:      nr.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, hme); err != nil {
		return Role{}, fmt.Errorf("create: %w", err)
	}

	return hme, nil
}

// Update modifies information about a role.
func (c *Core) Update(ctx context.Context, role Role, ur UpdateRole) (Role, error) {
	if ur.Name != nil {
		role.Name = *ur.Name
	}

	role.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, role); err != nil {
		return Role{}, fmt.Errorf("update: %w", err)
	}

	return role, nil
}

// Delete removes the specified role.
func (c *Core) Delete(ctx context.Context, hme Role) error {
	if err := c.storer.Delete(ctx, hme); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing roles.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Role, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	hmes, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return hmes, nil
}

// Count returns the total number of roles.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}

	return c.storer.Count(ctx, filter)
}

// QueryByID finds the role by the specified ID.
func (c *Core) QueryByID(ctx context.Context, roleID uuid.UUID) (Role, error) {
	hme, err := c.storer.QueryByID(ctx, roleID)
	if err != nil {
		return Role{}, fmt.Errorf("query: roleID[%s]: %w", roleID, err)
	}

	return hme, nil
}

// QueryByUserID finds the roles by a specified User ID.
func (c *Core) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	roles, err := c.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return roles, nil
}
