package patientdb

import (
	"bytes"
	"fmt"
	"github.com/fadhilijuma/gateone-service/business/core/crud/patient"
	"strings"
)

func (s *Store) applyFilter(filter patient.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["patient_id"] = *filter.ID
		wc = append(wc, "patient_id = :patient_id")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if filter.Age != nil {
		data["age"] = *filter.Age
		wc = append(wc, "age = :age")
	}

	if filter.Condition != nil {
		data["condition"] = *filter.Condition
		wc = append(wc, "condition = :condition")
	}
	if filter.VideoLink != nil {
		data["video_link"] = *filter.VideoLink
		wc = append(wc, "video_link = :video_link")
	}
	if filter.Healed != nil {
		data["healed"] = *filter.Healed
		wc = append(wc, "healed = :healed")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
