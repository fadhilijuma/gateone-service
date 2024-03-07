package region_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/region"
	"github.com/fadhilijuma/gateone-service/business/core/crud/user"
	"github.com/fadhilijuma/gateone-service/business/data/dbtest"
	"github.com/fadhilijuma/gateone-service/foundation/docker"
	"os"
	"runtime/debug"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var c *docker.Container

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	var err error

	c, err = dbtest.StartDB()
	if err != nil {
		return 1, err
	}
	defer dbtest.StopDB(c)

	return m.Run(), nil
}

func Test_Region(t *testing.T) {
	t.Run("crud", crud)
}

func crud(t *testing.T) {
	seed := func(ctx context.Context, usrCore *user.Core, regionCore *region.Core) ([]region.Region, error) {
		var filter user.QueryFilter
		filter.WithName("Admin Gopher")

		usrs, err := usrCore.Query(ctx, filter, user.DefaultOrderBy, 1, 1)
		if err != nil {
			return nil, fmt.Errorf("seeding users : %w", err)
		}

		Regions, err := region.TestGenerateSeedRegions(1, regionCore, usrs[0].ID)
		if err != nil {
			return nil, fmt.Errorf("seeding Regions : %w", err)
		}

		return Regions, nil
	}

	// ---------------------------------------------------------------------------

	test := dbtest.NewTest(t, c, "Test_Region/crud")

	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		test.Teardown()
	}()

	api := test.CoreAPIs

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Log("Go seeding ...")

	cns, err := seed(ctx, api.User, api.Region)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// ---------------------------------------------------------------------------

	saved, err := api.Region.QueryByID(ctx, cns[0].ID)
	if err != nil {
		t.Fatalf("Should be able to retrieve Region by ID: %s", err)
	}

	if cns[0].DateCreated.UnixMilli() != saved.DateCreated.UnixMilli() {
		t.Logf("got: %v", saved.DateCreated)
		t.Logf("exp: %v", cns[0].DateCreated)
		t.Logf("dif: %v", saved.DateCreated.Sub(cns[0].DateCreated))
		t.Errorf("Should get back the same date created")
	}

	if cns[0].DateUpdated.UnixMilli() != saved.DateUpdated.UnixMilli() {
		t.Logf("got: %v", saved.DateUpdated)
		t.Logf("exp: %v", cns[0].DateUpdated)
		t.Logf("dif: %v", saved.DateUpdated.Sub(cns[0].DateUpdated))
		t.Errorf("Should get back the same date updated")
	}

	cns[0].DateCreated = time.Time{}
	cns[0].DateUpdated = time.Time{}
	saved.DateCreated = time.Time{}
	saved.DateUpdated = time.Time{}

	if diff := cmp.Diff(cns[0], saved); diff != "" {
		t.Errorf("Should get back the same Region, dif:\n%s", diff)
	}

	// ---------------------------------------------------------------------------

	upd := region.UpdateRegion{
		Name: dbtest.StringPointer("Kisumu"),
	}

	if _, err := api.Region.Update(ctx, saved, upd); err != nil {
		t.Errorf("Should be able to update Region : %s", err)
	}

	saved, err = api.Region.QueryByID(ctx, cns[0].ID)
	if err != nil {
		t.Fatalf("Should be able to retrieve updated Region : %s", err)
	}

	diff := cns[0].DateUpdated.Sub(saved.DateUpdated)
	if diff > 0 {
		t.Fatalf("Should have a larger DateUpdated : sav %v, Region %v, dif %v", saved.DateUpdated, saved.DateUpdated, diff)
	}

	Regions, err := api.Region.Query(ctx, region.QueryFilter{}, user.DefaultOrderBy, 1, 3)
	if err != nil {
		t.Fatalf("Should be able to retrieve updated Region : %s", err)
	}

	// Check specified fields were updated. Make a copy of the original Region
	// and change just the fields we expect then diff it with what was saved.

	var idx int
	for i, h := range Regions {
		if h.ID == saved.ID {
			idx = i
		}
	}

	Regions[idx].DateCreated = time.Time{}
	Regions[idx].DateUpdated = time.Time{}
	saved.DateCreated = time.Time{}
	saved.DateUpdated = time.Time{}

	if diff := cmp.Diff(saved, Regions[idx]); diff != "" {
		t.Fatalf("Should get back the same Region, dif:\n%s", diff)
	}

	// -------------------------------------------------------------------------

	upd = region.UpdateRegion{
		Name: dbtest.StringPointer("Deaf"),
	}

	if _, err := api.Region.Update(ctx, saved, upd); err != nil {
		t.Fatalf("Should be able to update just some fields of Region : %s", err)
	}

	saved, err = api.Region.QueryByID(ctx, cns[0].ID)
	if err != nil {
		t.Fatalf("Should be able to retrieve updated Region : %s", err)
	}

	diff = cns[0].DateUpdated.Sub(saved.DateUpdated)
	if diff > 0 {
		t.Fatalf("Should have a larger DateUpdated : sav %v, Region %v, dif %v", saved.DateUpdated, cns[0].DateUpdated, diff)
	}

	if saved.Name != *upd.Name {
		t.Fatalf("Should be able to see updated name field : got %q want %q", saved.Name, *upd.Name)
	}

	if err := api.Region.Delete(ctx, saved); err != nil {
		t.Fatalf("Should be able to delete Region : %s", err)
	}

	_, err = api.Region.QueryByID(ctx, cns[0].ID)
	if !errors.Is(err, region.ErrNotFound) {
		t.Fatalf("Should NOT be able to retrieve deleted Region : %s", err)
	}
}
