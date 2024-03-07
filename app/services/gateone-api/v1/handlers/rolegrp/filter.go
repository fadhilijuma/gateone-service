package rolegrp

import (
	"github.com/fadhilijuma/gateone-service/business/core/crud/role"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"github.com/google/uuid"
	"net/http"
)

func parseFilter(r *http.Request) (role.QueryFilter, error) {
	const (
		filterByRoleID = "role_id"
		filterByName   = "name"
	)

	values := r.URL.Query()

	var filter role.QueryFilter

	if conditionID := values.Get(filterByRoleID); conditionID != "" {
		id, err := uuid.Parse(conditionID)
		if err != nil {
			return role.QueryFilter{}, validate.NewFieldsError(conditionID, err)
		}
		filter.WithRoleID(id)
	}

	if name := values.Get(filterByName); name != "" {
		filter.WithRoleName(name)
	}

	return filter, nil
}
