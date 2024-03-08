package regiondb

import (
	"bytes"
	"github.com/fadhilijuma/gateone-service/business/core/crud/region"
	"strings"
)

func (s *Store) applyFilter(filter region.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["Region_id"] = *filter.ID
		wc = append(wc, "Region_id = :Region_id")
	}

	if filter.UserID != nil {
		data["user_id"] = *filter.UserID
		wc = append(wc, "user_id = :user_id")
	}

	if filter.Name != nil {
		data["name"] = filter.Name
		wc = append(wc, "name = :name")
	}

	if filter.StartCreatedDate != nil {
		data["start_date_created"] = *filter.StartCreatedDate
		wc = append(wc, "date_created >= :start_date_created")
	}

	if filter.EndCreatedDate != nil {
		data["end_date_created"] = *filter.EndCreatedDate
		wc = append(wc, "date_created <= :end_date_created")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
