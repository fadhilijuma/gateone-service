package conditiongrp

import (
	"github.com/fadhilijuma/gateone-service/business/core/crud/condition"
	"github.com/fadhilijuma/gateone-service/foundation/validate"
	"github.com/google/uuid"
	"net/http"
)

func parseFilter(r *http.Request) (condition.QueryFilter, error) {
	const (
		filterByConditionID = "condition_id"
		filterByName        = "name"
	)

	values := r.URL.Query()

	var filter condition.QueryFilter

	if conditionID := values.Get(filterByConditionID); conditionID != "" {
		id, err := uuid.Parse(conditionID)
		if err != nil {
			return condition.QueryFilter{}, validate.NewFieldsError(conditionID, err)
		}
		filter.WithConditionID(id)
	}

	if name := values.Get(filterByName); name != "" {
		filter.WithConditionName(name)
	}

	return filter, nil
}
