// Package Regiondb contains Region related CRUD functionality.
package regiondb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/region"
	"github.com/fadhilijuma/gateone-service/business/data/sqldb"
	"github.com/fadhilijuma/gateone-service/business/data/transaction"
	"github.com/fadhilijuma/gateone-service/business/web/v1/order"
	"github.com/fadhilijuma/gateone-service/foundation/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for Region database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// ExecuteUnderTransaction constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (region.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// Create inserts a new Region into the database.
func (s *Store) Create(ctx context.Context, rn region.Region) error {
	const q = `
    INSERT INTO regions
        (region_id, user_id, name, date_created, date_updated)
    VALUES
        (:region_id, :user_id, :name, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBRegion(rn)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a Region from the database.
func (s *Store) Delete(ctx context.Context, rn region.Region) error {
	data := struct {
		ID string `db:"region_id"`
	}{
		ID: rn.ID.String(),
	}

	const q = `
    DELETE FROM
	    regions
	WHERE
	  	region_id = :region_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a Region document in the database.
func (s *Store) Update(ctx context.Context, rn region.Region) error {
	const q = `
    UPDATE
        regions
    SET
        "name"          = :name,
        "date_updated"  = :date_updated
    WHERE
        region_id = :Region_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBRegion(rn)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing Regions from the database.
func (s *Store) Query(ctx context.Context, filter region.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]region.Region, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
    SELECT
	    region_id, user_id, name, date_created, date_updated
	FROM
	  	regions`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbRegions []dbRegion
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbRegions); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	Regions, err := toCoreRegionsSlice(dbRegions)
	if err != nil {
		return nil, err
	}

	return Regions, nil
}

// Count returns the total number of Regions in the DB.
func (s *Store) Count(ctx context.Context, filter region.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
    SELECT
        count(1)
    FROM
        regions`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified Region from the database.
func (s *Store) QueryByID(ctx context.Context, RegionID uuid.UUID) (region.Region, error) {
	data := struct {
		ID string `db:"region_id"`
	}{
		ID: RegionID.String(),
	}

	const q = `
    SELECT
	  	region_id, user_id, name, date_created, date_updated
    FROM
        regions
    WHERE
        region_id = :region_id`

	var dbRl dbRegion
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbRl); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return region.Region{}, fmt.Errorf("namedquerystruct: %w", region.ErrNotFound)
		}
		return region.Region{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreRegion(dbRl)
}

// QueryByUserID gets the specified Region from the database by user id.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]region.Region, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
	    region_id, user_id, name, date_created, date_updated
	FROM
		regions
	WHERE
		user_id = :user_id`

	var dbRegions []dbRegion
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbRegions); err != nil {
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreRegionsSlice(dbRegions)
}
