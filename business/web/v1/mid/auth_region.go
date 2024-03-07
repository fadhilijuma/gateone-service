package mid

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"github.com/fadhilijuma/gateone-service/business/core/crud/region"
	v1 "github.com/fadhilijuma/gateone-service/business/web/v1"
	"github.com/fadhilijuma/gateone-service/business/web/v1/auth"
	"github.com/fadhilijuma/gateone-service/foundation/web"
	"net/http"

	"github.com/google/uuid"
)

type ctxRegionKey int

const regionKey ctxRegionKey = 1

// GetRegion returns the region from the context.
func GetRegion(ctx context.Context) region.Region {
	v, ok := ctx.Value(regionKey).(region.Region)
	if !ok {
		return region.Region{}
	}
	return v
}

func setRegion(ctx context.Context, rn region.Region) context.Context {
	return context.WithValue(ctx, regionKey, rn)
}

// AuthorizeRegion executes the specified role and extracts the specified
// region from the DB if a role id is specified in the call. Depending on
// the rule specified, the userid from the claims may be compared with the
// specified user id from the region.
func AuthorizeRegion(a *auth.Auth, rule string, rCore *region.Core) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID uuid.UUID

			if id := web.Param(r, "patient_id"); id != "" {
				var err error
				roleID, err := uuid.Parse(id)
				if err != nil {
					return v1.NewTrustedError(ErrInvalidID, http.StatusBadRequest)
				}

				reg, err := rCore.QueryByID(ctx, roleID)
				if err != nil {
					switch {
					case errors.Is(err, patient.ErrNotFound):
						return v1.NewTrustedError(err, http.StatusNoContent)
					default:
						return fmt.Errorf("querybyid: roleID[%s]: %w", roleID, err)
					}
				}

				userID = reg.UserID
				ctx = setRegion(ctx, reg)
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
