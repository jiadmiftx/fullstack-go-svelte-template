package paginate_utils

import (
	"errors"
	"strconv"
)

func CheckPaginateFromMap(filter map[string]interface{}, p *PaginateData) error {
	if p == nil {
		return errors.New("paginate struct is nil")
	}
	if v, ok := filter["page"]; ok {
		if s, ok := v.(string); ok {
			if i, err := strconv.Atoi(s); err == nil {
				p.Page = i
			}
		}
	}
	if v, ok := filter["limit"]; ok {
		if s, ok := v.(string); ok {
			if i, err := strconv.Atoi(s); err == nil {
				p.Limit = i
			}
		}
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return nil
}
