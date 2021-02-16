package snippets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/titusjaka/go-sample/internal/business/snippets"
	"github.com/titusjaka/go-sample/internal/infrastructure/service"
)

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name               string
		limit              uint
		offset             uint
		total              uint
		expectedPagination service.Pagination
	}{
		{
			name:               "Empty fields",
			limit:              0,
			offset:             0,
			total:              0,
			expectedPagination: service.Pagination{Limit: 100, CurrentPage: 1},
		},
		{
			name:               "Fields in range",
			limit:              10,
			offset:             100,
			total:              0,
			expectedPagination: service.Pagination{Limit: 10, Offset: 100, CurrentPage: 11},
		},
		{
			name:               "Limit out of range",
			limit:              1000,
			offset:             0,
			total:              0,
			expectedPagination: service.Pagination{Limit: 100, CurrentPage: 1},
		},
		{
			name:               "Set total and offset",
			limit:              10,
			offset:             50,
			total:              100,
			expectedPagination: service.Pagination{Limit: 10, Offset: 50, Total: 100, TotalPages: 10, CurrentPage: 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPagination := snippets.NewPagination(tt.limit, tt.offset, tt.total)
			assert.Equal(t, tt.expectedPagination, gotPagination)
		})
	}
}

func TestConvertPaginationToSQLExpression(t *testing.T) {
	tests := []struct {
		name               string
		limit              uint
		offset             uint
		expectedExpression string
	}{
		{
			name:               "Empty fields",
			limit:              0,
			offset:             0,
			expectedExpression: "",
		},
		{
			name:               "LIMIT 10",
			limit:              10,
			offset:             0,
			expectedExpression: "LIMIT 10",
		},
		{
			name:               "OFFSET 100",
			limit:              0,
			offset:             100,
			expectedExpression: "OFFSET 100",
		},
		{
			name:               "LIMIT 10 OFFSET 100",
			limit:              10,
			offset:             100,
			expectedExpression: "LIMIT 10 OFFSET 100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExpression := snippets.ConvertPaginationToSQLExpression(service.Pagination{Limit: tt.limit, Offset: tt.offset})
			assert.Equal(t, tt.expectedExpression, gotExpression)
		})
	}
}
