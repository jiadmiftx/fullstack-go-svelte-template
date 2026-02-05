package common

import (
	"sekolah-madrasah/pkg/paginate_utils"
	"gorm.io/gorm"
)

// ApplyPagination applies pagination to GORM query
func ApplyPagination(query *gorm.DB, paginate *paginate_utils.PaginateData) *gorm.DB {
	if paginate == nil {
		return query
	}
	return query.Scopes(paginate_utils.Paginate(paginate))
}

// ApplyOrderBy applies default ordering to query
func ApplyOrderBy(query *gorm.DB, orderBy string) *gorm.DB {
	if orderBy == "" {
		orderBy = "created_at DESC"
	}
	return query.Order(orderBy)
}

// CountTotal gets total count for pagination
func CountTotal(query *gorm.DB, paginate *paginate_utils.PaginateData) error {
	if paginate == nil {
		return nil
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return err
	}

	paginate.TotalData = total
	return nil
}
