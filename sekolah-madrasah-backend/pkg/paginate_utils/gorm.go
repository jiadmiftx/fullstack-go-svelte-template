package paginate_utils

import "gorm.io/gorm"

func Paginate(p *PaginateData) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p == nil {
			return db
		}
		offset := (p.Page - 1) * p.Limit
		return db.Offset(offset).Limit(p.Limit)
	}
}
