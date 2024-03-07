// Package patient Package product provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package patient

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
	ErrNotFound     = errors.New("patient not found")
	ErrUserDisabled = errors.New("user disabled")
	ErrInvalidCost  = errors.New("cost not valid")
)

// Storer interface declares the behavior this package needs to persists and
// retrieve data.
type Storer interface {
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
	Create(ctx context.Context, pn Patient) error
	Update(ctx context.Context, pn Patient) error
	Delete(ctx context.Context, pn Patient) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Patient, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, pnID uuid.UUID) (Patient, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Patient, error)
}

// Core manages the set of APIs for product access.
type Core struct {
	log      *logger.Logger
	usrCore  *user.Core
	delegate *delegate.Delegate
	storer   Storer
}

// NewCore constructs a product core API for use.
func NewCore(log *logger.Logger, usrCore *user.Core, delegate *delegate.Delegate, storer Storer) *Core {
	c := Core{
		log:      log,
		usrCore:  usrCore,
		delegate: delegate,
		storer:   storer,
	}

	c.registerDelegateFunctions()

	return &c
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

// Create adds a new product to the system.
func (c *Core) Create(ctx context.Context, np NewPatient) (Patient, error) {
	usr, err := c.usrCore.QueryByID(ctx, np.UserID)
	if err != nil {
		return Patient{}, fmt.Errorf("user.querybyid: %s: %w", np.UserID, err)
	}

	if !usr.Enabled {
		return Patient{}, ErrUserDisabled
	}

	now := time.Now()

	prd := Patient{
		ID:          uuid.New(),
		Name:        np.Name,
		Age:         np.Age,
		Condition:   np.Condition,
		VideoLinks:  np.VideoLinks,
		Healed:      np.Healed,
		UserID:      np.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, prd); err != nil {
		return Patient{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}

// Update modifies information about a product.
func (c *Core) Update(ctx context.Context, pn Patient, up UpdatePatient) (Patient, error) {
	if up.Name != nil {
		pn.Name = *up.Name
	}

	if up.Age != nil {
		pn.Age = *up.Age
	}

	if up.VideoLinks != nil {
		pn.VideoLinks = up.VideoLinks
	}
	if up.Condition != nil {
		pn.Condition = *up.Condition
	}
	if up.Healed != nil {
		pn.Healed = *up.Healed
	}

	pn.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, pn); err != nil {
		return Patient{}, fmt.Errorf("update: %w", err)
	}

	return pn, nil
}

// Delete removes the specified patient.
func (c *Core) Delete(ctx context.Context, prd Patient) error {
	if err := c.storer.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing patients.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Patient, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	prds, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}

// Count returns the total number of patients.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}

	return c.storer.Count(ctx, filter)
}

// QueryByID finds the patient by the specified ID.
func (c *Core) QueryByID(ctx context.Context, patientID uuid.UUID) (Patient, error) {
	prd, err := c.storer.QueryByID(ctx, patientID)
	if err != nil {
		return Patient{}, fmt.Errorf("query: patientID[%s]: %w", patientID, err)
	}

	return prd, nil
}

// QueryByUserID finds the patients by a specified User ID.
func (c *Core) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Patient, error) {
	prds, err := c.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}
