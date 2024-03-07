package mid

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/condition"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	v1 "github.com/fadhilijuma/gateone-service/business/web/v1"
	"github.com/fadhilijuma/gateone-service/business/web/v1/auth"
	"github.com/fadhilijuma/gateone-service/foundation/web"
	"net/http"

	"github.com/google/uuid"
)

type ctxConditionKey int

const conditionKey ctxConditionKey = 1

// GetCondition returns the condition from the context.
func GetCondition(ctx context.Context) condition.Condition {
	v, ok := ctx.Value(conditionKey).(condition.Condition)
	if !ok {
		return condition.Condition{}
	}
	return v
}

func setCondition(ctx context.Context, cn condition.Condition) context.Context {
	return context.WithValue(ctx, conditionKey, cn)
}

// AuthorizeCondition executes the specified role and extracts the specified
// condition from the DB if a condition id is specified in the call. Depending on
// the rule specified, the userid from the claims may be compared with the
// specified user id from the condition.
func AuthorizeCondition(a *auth.Auth, rule string, cnCore *condition.Core) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID uuid.UUID

			if id := web.Param(r, "patient_id"); id != "" {
				var err error
				conditionID, err := uuid.Parse(id)
				if err != nil {
					return v1.NewTrustedError(ErrInvalidID, http.StatusBadRequest)
				}

				cn, err := cnCore.QueryByID(ctx, conditionID)
				if err != nil {
					switch {
					case errors.Is(err, patient.ErrNotFound):
						return v1.NewTrustedError(err, http.StatusNoContent)
					default:
						return fmt.Errorf("querybyid: conditionID[%s]: %w", conditionID, err)
					}
				}

				userID = cn.UserID
				ctx = setCondition(ctx, cn)
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
