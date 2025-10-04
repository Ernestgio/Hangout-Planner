package dto_test

import (
	"fmt"
	"testing"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/stretchr/testify/require"
)

func TestCursorPagination_GetLimit(t *testing.T) {
	testCases := []struct {
		name          string
		pagination    dto.CursorPagination
		expectedLimit int
	}{
		{
			name:          "limit is zero, should default",
			pagination:    dto.CursorPagination{Limit: 0},
			expectedLimit: constants.DefaultLimit,
		},
		{
			name:          "limit is negative, should default",
			pagination:    dto.CursorPagination{Limit: -10},
			expectedLimit: constants.DefaultLimit,
		},
		{
			name:          "limit is within range",
			pagination:    dto.CursorPagination{Limit: 50},
			expectedLimit: 50,
		},
		{
			name:          "limit exceeds max, should cap at max",
			pagination:    dto.CursorPagination{Limit: 200},
			expectedLimit: constants.MaxLimit,
		},
		{
			name:          "limit is exactly max",
			pagination:    dto.CursorPagination{Limit: 100},
			expectedLimit: constants.MaxLimit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedLimit, tc.pagination.GetLimit())
		})
	}
}

func TestCursorPagination_GetSortBy(t *testing.T) {
	testCases := []struct {
		name           string
		pagination     dto.CursorPagination
		expectedSortBy string
	}{
		{
			name:           "sort by is date",
			pagination:     dto.CursorPagination{SortBy: constants.SortByDate},
			expectedSortBy: constants.SortByDate,
		},
		{
			name:           "sort by is created_at",
			pagination:     dto.CursorPagination{SortBy: constants.SortByCreatedAt},
			expectedSortBy: constants.SortByCreatedAt,
		},
		{
			name:           "sort by is empty, should default",
			pagination:     dto.CursorPagination{SortBy: ""},
			expectedSortBy: constants.SortByCreatedAt,
		},
		{
			name:           "sort by is invalid, should default",
			pagination:     dto.CursorPagination{SortBy: "invalid_column"},
			expectedSortBy: constants.SortByCreatedAt,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedSortBy, tc.pagination.GetSortBy())
		})
	}
}

func TestCursorPagination_GetSortDir(t *testing.T) {
	testCases := []struct {
		name            string
		pagination      dto.CursorPagination
		expectedSortDir string
	}{
		{
			name:            "sort dir is asc",
			pagination:      dto.CursorPagination{SortDir: "asc"},
			expectedSortDir: constants.SortDirectionAsc,
		},
		{
			name:            "sort dir is desc",
			pagination:      dto.CursorPagination{SortDir: "desc"},
			expectedSortDir: constants.SortDirectionDesc,
		},
		{
			name:            "sort dir is uppercase ASC",
			pagination:      dto.CursorPagination{SortDir: "ASC"},
			expectedSortDir: constants.SortDirectionAsc,
		},
		{
			name:            "sort dir is empty, should default",
			pagination:      dto.CursorPagination{SortDir: ""},
			expectedSortDir: constants.SortDirectionDesc,
		},
		{
			name:            "sort dir is invalid, should default",
			pagination:      dto.CursorPagination{SortDir: "up"},
			expectedSortDir: constants.SortDirectionDesc,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedSortDir, tc.pagination.GetSortDir())
		})
	}
}

func TestCursorPagination_GetOrderByClause(t *testing.T) {
	testCases := []struct {
		name            string
		pagination      dto.CursorPagination
		expectedOrderBy string
	}{
		{
			name:            "default sort",
			pagination:      dto.CursorPagination{},
			expectedOrderBy: fmt.Sprintf("%s %s, id %s", constants.SortByCreatedAt, constants.SortDirectionDesc, constants.SortDirectionDesc),
		},
		{
			name:            "sort by date ascending",
			pagination:      dto.CursorPagination{SortBy: "date", SortDir: "asc"},
			expectedOrderBy: fmt.Sprintf("%s %s, id %s", constants.SortByDate, constants.SortDirectionAsc, constants.SortDirectionAsc),
		},
		{
			name:            "invalid sort falls back to default",
			pagination:      dto.CursorPagination{SortBy: "name", SortDir: "down"},
			expectedOrderBy: fmt.Sprintf("%s %s, id %s", constants.SortByCreatedAt, constants.SortDirectionDesc, constants.SortDirectionDesc),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedOrderBy, tc.pagination.GetOrderByClause())
		})
	}
}
