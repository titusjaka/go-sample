package snippets

import (
	"fmt"
	"math"
	"strings"

	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

const (
	maxLimit uint = 100
)

// NewPagination creates new service.Pagination with default maximum limit
func NewPagination(limit uint, offset uint, total uint) service.Pagination {
	pagination := service.Pagination{
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}

	if pagination.Limit == 0 || pagination.Limit > maxLimit {
		pagination.Limit = maxLimit
	}

	pagination.TotalPages = uint(math.Ceil(float64(pagination.Total) / float64(pagination.Limit)))
	pagination.CurrentPage = (pagination.Offset / pagination.Limit) + 1

	return pagination
}

// ConvertPaginationToSQLExpression converts service.Pagination to SQL expression
func ConvertPaginationToSQLExpression(pagination service.Pagination) string {
	conditions := make([]string, 0, 2)

	if pagination.Limit != 0 {
		conditions = append(conditions, fmt.Sprintf("LIMIT %d", pagination.Limit))
	}

	if pagination.Offset != 0 {
		conditions = append(conditions, fmt.Sprintf("OFFSET %d", pagination.Offset))
	}

	return strings.Join(conditions, " ")
}
