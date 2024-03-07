package patient_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"github.com/fadhilijuma/gateone-service/business/core/crud/user"
	"github.com/fadhilijuma/gateone-service/business/data/dbtest"
	"github.com/fadhilijuma/gateone-service/business/data/sqldb"
	"github.com/fadhilijuma/gateone-service/business/data/transaction"
	"github.com/fadhilijuma/gateone-service/foundation/docker"
	"net/mail"
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

func Test_Patient(t *testing.T) {
	t.Run("crud", crud)
	t.Run("paging", paging)
	t.Run("transaction", tran)
}

func crud(t *testing.T) {
	seed := func(ctx context.Context, usrCore *user.Core, pnCore *patient.Core) ([]patient.Patient, error) {
		var filter user.QueryFilter
		filter.WithName("Admin Gate One")

		usrs, err := usrCore.Query(ctx, filter, user.DefaultOrderBy, 1, 1)
		if err != nil {
			return nil, fmt.Errorf("seeding users : %w", err)
		}

		prds, err := patient.TestGenerateSeedPatients(1, pnCore, usrs[0].ID)
		if err != nil {
			return nil, fmt.Errorf("seeding patients : %w", err)
		}

		return prds, nil
	}

	// -------------------------------------------------------------------------

	test := dbtest.NewTest(t, c, "Test_Patient/crud")
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

	prds, err := seed(ctx, api.User, api.Patient)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	saved, err := api.Patient.QueryByID(ctx, prds[0].ID)
	if err != nil {
		t.Fatalf("Should be able to retrieve patient by ID: %s", err)
	}

	if prds[0].DateCreated.UnixMilli() != saved.DateCreated.UnixMilli() {
		t.Logf("got: %v", saved.DateCreated)
		t.Logf("exp: %v", prds[0].DateCreated)
		t.Logf("dif: %v", saved.DateCreated.Sub(prds[0].DateCreated))
		t.Errorf("Should get back the same date created")
	}

	if prds[0].DateUpdated.UnixMilli() != saved.DateUpdated.UnixMilli() {
		t.Logf("got: %v", saved.DateUpdated)
		t.Logf("exp: %v", prds[0].DateUpdated)
		t.Logf("dif: %v", saved.DateUpdated.Sub(prds[0].DateUpdated))
		t.Errorf("Should get back the same date updated")
	}

	prds[0].DateCreated = time.Time{}
	prds[0].DateUpdated = time.Time{}
	saved.DateCreated = time.Time{}
	saved.DateUpdated = time.Time{}

	if diff := cmp.Diff(prds[0], saved); diff != "" {
		t.Errorf("Should get back the same patient, dif:\n%s", diff)
	}

	// -------------------------------------------------------------------------

	upd := patient.UpdatePatient{
		Name:       dbtest.StringPointer("James"),
		Age:        dbtest.IntPointer(50),
		VideoLinks: []string{"https://www.youtube.com/watch?v=1234"},
		Condition:  dbtest.StringPointer("Deaf"),
		Healed:     dbtest.BoolPointer(false),
	}

	if _, err := api.Patient.Update(ctx, saved, upd); err != nil {
		t.Errorf("Should be able to update patient : %s", err)
	}

	saved, err = api.Patient.QueryByID(ctx, prds[0].ID)
	if err != nil {
		t.Fatalf("Should be able to retrieve updated patient : %s", err)
	}

	diff := prds[0].DateUpdated.Sub(saved.DateUpdated)
	if diff > 0 {
		t.Fatalf("Should have a larger DateUpdated : sav %v, prd %v, dif %v", saved.DateUpdated, saved.DateUpdated, diff)
	}

	patients, err := api.Patient.Query(ctx, patient.QueryFilter{}, patient.DefaultOrderBy, 1, 3)
	if err != nil {
		t.Fatalf("Should be able to retrieve updated patient : %s", err)
	}

	// Check specified fields were updated. Make a copy of the original patient
	// and change just the fields we expect then diff it with what was saved.

	var idx int
	for i, p := range patients {
		if p.ID == saved.ID {
			idx = i
		}
	}

	patients[idx].DateCreated = time.Time{}
	patients[idx].DateUpdated = time.Time{}
	saved.DateCreated = time.Time{}
	saved.DateUpdated = time.Time{}

	if diff := cmp.Diff(saved, patients[idx]); diff != "" {
		t.Fatalf("Should get back the same patient, dif:\n%s", diff)
	}

	// -------------------------------------------------------------------------

	upd = patient.UpdatePatient{
		Name: dbtest.StringPointer("Graphic Novels"),
	}

	if _, err := api.Patient.Update(ctx, saved, upd); err != nil {
		t.Fatalf("Should be able to update just some fields of patient : %s", err)
	}

	saved, err = api.Patient.QueryByID(ctx, prds[0].ID)
	if err != nil {
		t.Fatalf("Should be able to retrieve updated patient : %s", err)
	}

	diff = prds[0].DateUpdated.Sub(saved.DateUpdated)
	if diff > 0 {
		t.Fatalf("Should have a larger DateUpdated : sav %v, pn %v, dif %v", saved.DateUpdated, prds[0].DateUpdated, diff)
	}

	if saved.Name != *upd.Name {
		t.Fatalf("Should be able to see updated Name field : got %q want %q", saved.Name, *upd.Name)
	}

	if err := api.Patient.Delete(ctx, saved); err != nil {
		t.Fatalf("Should be able to delete patient : %s", err)
	}

	_, err = api.Patient.QueryByID(ctx, prds[0].ID)
	if !errors.Is(err, patient.ErrNotFound) {
		t.Fatalf("Should NOT be able to retrieve deleted patient : %s", err)
	}
}

func paging(t *testing.T) {
	seed := func(ctx context.Context, usrCore *user.Core, prdCore *patient.Core) ([]patient.Patient, error) {
		var filter user.QueryFilter
		filter.WithName("Admin Gopher")

		usrs, err := usrCore.Query(ctx, filter, user.DefaultOrderBy, 1, 1)
		if err != nil {
			return nil, fmt.Errorf("seeding patients : %w", err)
		}

		prds, err := patient.TestGenerateSeedPatients(2, prdCore, usrs[0].ID)
		if err != nil {
			return nil, fmt.Errorf("seeding patients : %w", err)
		}

		return prds, nil
	}

	// -------------------------------------------------------------------------

	test := dbtest.NewTest(t, c, "Test_Patient/paging")
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

	prds, err := seed(ctx, api.User, api.Patient)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	name := prds[0].Name
	prd1, err := api.Patient.Query(ctx, patient.QueryFilter{Name: &name}, patient.DefaultOrderBy, 1, 1)
	if err != nil {
		t.Fatalf("Should be able to retrieve patients %q : %s", name, err)
	}

	n, err := api.Patient.Count(ctx, patient.QueryFilter{Name: &name})
	if err != nil {
		t.Fatalf("Should be able to retrieve patients count %q : %s", name, err)
	}

	if len(prd1) != n && prd1[0].Name == name {
		t.Log("got:", len(prd1))
		t.Log("exp:", n)
		t.Fatalf("Should have a single patients for %q", name)
	}

	name = prds[1].Name
	prd2, err := api.Patient.Query(ctx, patient.QueryFilter{Name: &name}, patient.DefaultOrderBy, 1, 1)
	if err != nil {
		t.Fatalf("Should be able to retrieve patients %q : %s", name, err)
	}

	n, err = api.Patient.Count(ctx, patient.QueryFilter{Name: &name})
	if err != nil {
		t.Fatalf("Should be able to retrieve patient count %q : %s", name, err)
	}

	if len(prd2) != n && prd2[0].Name == name {
		t.Log("got:", len(prd2))
		t.Log("exp:", n)
		t.Fatalf("Should have a single patient for %q", name)
	}

	prd3, err := api.Patient.Query(ctx, patient.QueryFilter{}, patient.DefaultOrderBy, 1, 2)
	if err != nil {
		t.Fatalf("Should be able to retrieve 2 patients for page 1 : %s", err)
	}

	n, err = api.Patient.Count(ctx, patient.QueryFilter{})
	if err != nil {
		t.Fatalf("Should be able to retrieve patient count %q : %s", name, err)
	}

	if len(prd3) != n {
		t.Logf("got: %v", len(prd3))
		t.Logf("exp: %v", n)
		t.Fatalf("Should have 2 patients for page ")
	}

	if prd3[0].ID == prd3[1].ID {
		t.Logf("patient1: %v", prd3[0].ID)
		t.Logf("patient2: %v", prd3[1].ID)
		t.Fatalf("Should have different patient")
	}
}

func tran(t *testing.T) {
	test := dbtest.NewTest(t, c, "Test_Patient/tran")
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

	// -------------------------------------------------------------------------
	// Execute under a transaction with rollback

	f := func(tx transaction.Transaction) error {
		usrCore, err := api.User.ExecuteUnderTransaction(tx)
		if err != nil {
			t.Fatalf("Should be able to create new user core: %s.", err)
		}

		prdCore, err := api.Patient.ExecuteUnderTransaction(tx)
		if err != nil {
			t.Fatalf("Should be able to create new patient core: %s.", err)
		}

		email, err := mail.ParseAddress("test@test.com")
		if err != nil {
			t.Fatalf("Should be able to parse email: %s.", err)
		}

		nu := user.NewUser{
			Name:            "test user",
			Email:           *email,
			Roles:           []user.Role{user.RoleAdmin},
			Department:      "some",
			Password:        "some",
			PasswordConfirm: "some",
		}

		usr, err := usrCore.Create(ctx, nu)
		if err != nil {
			return err
		}

		np := patient.NewPatient{
			UserID:     usr.ID,
			Name:       "test patient",
			Age:        10,
			Condition:  "deaf",
			VideoLinks: []string{"https://www.youtube.com/watch?v=1234"},
			Healed:     false,
		}

		_, err = prdCore.Create(ctx, np)
		if err != nil {
			return err
		}

		return nil
	}

	err := transaction.ExecuteUnderTransaction(ctx, test.Log, sqldb.NewBeginner(test.DB), f)
	if !errors.Is(err, patient.ErrInvalidCost) {
		t.Fatalf("Should NOT be able to add patient : %s.", err)
	}

	// -------------------------------------------------------------------------
	// Validate rollback

	email, err := mail.ParseAddress("test@test.com")
	if err != nil {
		t.Fatalf("Should be able to parse email: %s.", err)
	}

	usr, err := api.User.QueryByEmail(ctx, *email)
	if err == nil {
		t.Fatalf("Should NOT be able to retrieve user but got: %+v.", usr)
	}
	if !errors.Is(err, user.ErrNotFound) {
		t.Fatalf("Should get ErrNotFound but got: %s.", err)
	}

	count, err := api.Patient.Count(ctx, patient.QueryFilter{})
	if err != nil {
		t.Fatalf("Should be able to count patients: %s.", err)
	}

	if count > 0 {
		t.Fatalf("Should have no patients in the DB, but have: %d.", count)
	}

	// -------------------------------------------------------------------------
	// Good transaction

	f = func(tx transaction.Transaction) error {
		usrCore, err := api.User.ExecuteUnderTransaction(tx)
		if err != nil {
			t.Fatalf("Should be able to create new user core: %s.", err)
		}

		prdCore, err := api.Patient.ExecuteUnderTransaction(tx)
		if err != nil {
			t.Fatalf("Should be able to create new patient core: %s.", err)
		}

		email, err := mail.ParseAddress("test@test.com")
		if err != nil {
			t.Fatalf("Should be able to parse email: %s.", err)
		}

		nu := user.NewUser{
			Name:            "test user",
			Email:           *email,
			Roles:           []user.Role{user.RoleAdmin},
			Department:      "some",
			Password:        "some",
			PasswordConfirm: "some",
		}

		usr, err := usrCore.Create(ctx, nu)
		if err != nil {
			return err
		}

		np := patient.NewPatient{
			UserID:     usr.ID,
			Name:       "test patient",
			Age:        10,
			Condition:  "deaf",
			VideoLinks: []string{"https://www.youtube.com/watch?v=1234"},
			Healed:     false,
		}

		_, err = prdCore.Create(ctx, np)
		if err != nil {
			return err
		}

		return nil
	}

	err = transaction.ExecuteUnderTransaction(ctx, test.Log, sqldb.NewBeginner(test.DB), f)
	if errors.Is(err, patient.ErrInvalidCost) {
		t.Fatalf("Should be able to add patient : %s.", err)
	}

	// -------------------------------------------------------------------------
	// Validate

	usr, err = api.User.QueryByEmail(ctx, *email)
	if err != nil {
		t.Fatalf("Should be able to retrieve user but got: %+v.", usr)
	}

	count, err = api.Patient.Count(ctx, patient.QueryFilter{})
	if err != nil {
		t.Fatalf("Should be able to count patients: %s.", err)
	}

	if count == 0 {
		t.Fatal("Should have prodpatientsucts in the DB.")
	}
}
