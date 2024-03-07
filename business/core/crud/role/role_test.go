package role_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/role"
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

func Test_Role(t *testing.T) {
	t.Run("crud", crud)
}

func crud(t *testing.T) {
	seed := func(ctx context.Context, usrCore *user.Core, roleCore *role.Core) ([]role.Role, error) {
		var filter user.QueryFilter
		filter.WithName("Admin Gopher")

		usrs, err := usrCore.Query(ctx, filter, user.DefaultOrderBy, 1, 1)
		if err != nil {
			return nil, fmt.Errorf("seeding users : %w", err)
		}

		roles, err := role.TestGenerateSeedRoles(1, roleCore, usrs[0].ID)
		if err != nil {
			return nil, fmt.Errorf("seeding roles : %w", err)
		}

		return roles, nil
	}

	// ---------------------------------------------------------------------------

	test := dbtest.NewTest(t, c, "Test_Role/crud")

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

	hmes, err := seed(ctx, api.User, api.Role)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// ---------------------------------------------------------------------------

	saved, err := api.Role.QueryByID(ctx, hmes[0].ID)
	if err != nil {
		t.Fatalf("Should be able to retrieve role by ID: %s", err)
	}

	if hmes[0].DateCreated.UnixMilli() != saved.DateCreated.UnixMilli() {
		t.Logf("got: %v", saved.DateCreated)
		t.Logf("exp: %v", hmes[0].DateCreated)
		t.Logf("dif: %v", saved.DateCreated.Sub(hmes[0].DateCreated))
		t.Errorf("Should get back the same date created")
	}

	if hmes[0].DateUpdated.UnixMilli() != saved.DateUpdated.UnixMilli() {
		t.Logf("got: %v", saved.DateUpdated)
		t.Logf("exp: %v", hmes[0].DateUpdated)
		t.Logf("dif: %v", saved.DateUpdated.Sub(hmes[0].DateUpdated))
		t.Errorf("Should get back the same date updated")
	}

	hmes[0].DateCreated = time.Time{}
	hmes[0].DateUpdated = time.Time{}
	saved.DateCreated = time.Time{}
	saved.DateUpdated = time.Time{}

	if diff := cmp.Diff(hmes[0], saved); diff != "" {
		t.Errorf("Should get back the same role, dif:\n%s", diff)
	}

	// ---------------------------------------------------------------------------

	upd := role.UpdateRole{
		Name: dbtest.StringPointer("Admin"),
	}

	if _, err := api.Role.Update(ctx, saved, upd); err != nil {
		t.Errorf("Should be able to update role : %s", err)
	}

	saved, err = api.Role.QueryByID(ctx, hmes[0].ID)
	if err != nil {
		t.Fatalf("Should be able to retrieve updated role : %s", err)
	}

	diff := hmes[0].DateUpdated.Sub(saved.DateUpdated)
	if diff > 0 {
		t.Fatalf("Should have a larger DateUpdated : sav %v, role %v, dif %v", saved.DateUpdated, saved.DateUpdated, diff)
	}

	roles, err := api.Role.Query(ctx, role.QueryFilter{}, user.DefaultOrderBy, 1, 3)
	if err != nil {
		t.Fatalf("Should be able to retrieve updated role : %s", err)
	}

	// Check specified fields were updated. Make a copy of the original role
	// and change just the fields we expect then diff it with what was saved.

	var idx int
	for i, h := range roles {
		if h.ID == saved.ID {
			idx = i
		}
	}

	roles[idx].DateCreated = time.Time{}
	roles[idx].DateUpdated = time.Time{}
	saved.DateCreated = time.Time{}
	saved.DateUpdated = time.Time{}

	if diff := cmp.Diff(saved, roles[idx]); diff != "" {
		t.Fatalf("Should get back the same role, dif:\n%s", diff)
	}

	// -------------------------------------------------------------------------

	upd = role.UpdateRole{
		Name: dbtest.StringPointer("User"),
	}

	if _, err := api.Role.Update(ctx, saved, upd); err != nil {
		t.Fatalf("Should be able to update just some fields of role : %s", err)
	}

	saved, err = api.Role.QueryByID(ctx, hmes[0].ID)
	if err != nil {
		t.Fatalf("Should be able to retrieve updated role : %s", err)
	}

	diff = hmes[0].DateUpdated.Sub(saved.DateUpdated)
	if diff > 0 {
		t.Fatalf("Should have a larger DateUpdated : sav %v, role %v, dif %v", saved.DateUpdated, hmes[0].DateUpdated, diff)
	}

	if saved.Name != *upd.Name {
		t.Fatalf("Should be able to see updated name field : got %q want %q", saved.Name, *upd.Name)
	}

	if err := api.Role.Delete(ctx, saved); err != nil {
		t.Fatalf("Should be able to delete role : %s", err)
	}

	_, err = api.Role.QueryByID(ctx, hmes[0].ID)
	if !errors.Is(err, role.ErrNotFound) {
		t.Fatalf("Should NOT be able to retrieve deleted role : %s", err)
	}
}
