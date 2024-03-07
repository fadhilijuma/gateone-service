package mid

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	v1 "github.com/fadhilijuma/gateone-service/business/web/v1"
	"github.com/fadhilijuma/gateone-service/business/web/v1/auth"
	"github.com/fadhilijuma/gateone-service/foundation/web"
	"net/http"

	"github.com/google/uuid"
)

type ctxPatientKey int

const patientKey ctxPatientKey = 1

// GetPatient returns the patient from the context.
func GetPatient(ctx context.Context) patient.Patient {
	v, ok := ctx.Value(patientKey).(patient.Patient)
	if !ok {
		return patient.Patient{}
	}
	return v
}

func setPatient(ctx context.Context, pn patient.Patient) context.Context {
	return context.WithValue(ctx, patientKey, pn)
}

// AuthorizePatient executes the specified role and extracts the specified
// patient from the DB if a patient id is specified in the call. Depending on
// the rule specified, the userid from the claims may be compared with the
// specified user id from the patient.
func AuthorizePatient(a *auth.Auth, rule string, prdCore *patient.Core) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID uuid.UUID

			if id := web.Param(r, "patient_id"); id != "" {
				var err error
				patientID, err := uuid.Parse(id)
				if err != nil {
					return v1.NewTrustedError(ErrInvalidID, http.StatusBadRequest)
				}

				prd, err := prdCore.QueryByID(ctx, patientID)
				if err != nil {
					switch {
					case errors.Is(err, patient.ErrNotFound):
						return v1.NewTrustedError(err, http.StatusNoContent)
					default:
						return fmt.Errorf("querybyid: patientID[%s]: %w", patientID, err)
					}
				}

				userID = prd.UserID
				ctx = setPatient(ctx, prd)
			}

			claims := getClaims(ctx)

			if err := a.Authorize(ctx, claims, userID, rule); err != nil {
				return auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
