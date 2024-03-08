package regiongrp

import (
	"github.com/fadhilijuma/gateone-service/business/core/crud/region"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"github.com/google/uuid"
	"net/http"
)

func parseFilter(r *http.Request) (region.QueryFilter, error) {
	const (
		filterByRoleID = "role_id"
		filterByName   = "name"
	)

	values := r.URL.Query()

	var filter region.QueryFilter

	if conditionID := values.Get(filterByRoleID); conditionID != "" {
		id, err := uuid.Parse(conditionID)
		if err != nil {
			return region.QueryFilter{}, validate.NewFieldsError(conditionID, err)
		}
		filter.WithRegionID(id)
	}

	if name := values.Get(filterByName); name != "" {
		filter.WithRegionName(name)
	}

	return filter, nil
}
