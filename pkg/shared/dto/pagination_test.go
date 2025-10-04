package dto_test

import (
	"testing"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/dto"
	"github.com/stretchr/testify/require"
)

func TestPagination_GetPage(t *testing.T) {
	testCases := []struct {
		name     string
		input    dto.Pagination
		expected int
	}{
		{
			name:     "positive page",
			input:    dto.Pagination{Page: 5},
			expected: 5,
		},
		{
			name:     "zero page",
			input:    dto.Pagination{Page: 0},
			expected: dto.DefaultPage,
		},
		{
			name:     "negative page",
			input:    dto.Pagination{Page: -1},
			expected: dto.DefaultPage,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.input.GetPage())
		})
	}
}

func TestPagination_GetLimit(t *testing.T) {
	testCases := []struct {
		name     string
		input    dto.Pagination
		expected int
	}{
		{
			name:     "positive limit within range",
			input:    dto.Pagination{Limit: 50},
			expected: 50,
		},
		{
			name:     "zero limit",
			input:    dto.Pagination{Limit: 0},
			expected: dto.DefaultLimit,
		},
		{
			name:     "negative limit",
			input:    dto.Pagination{Limit: -10},
			expected: dto.DefaultLimit,
		},
		{
			name:     "limit greater than max",
			input:    dto.Pagination{Limit: 200},
			expected: dto.MaxLimit,
		},
		{
			name:     "limit equal to max",
			input:    dto.Pagination{Limit: 100},
			expected: dto.MaxLimit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.input.GetLimit())
		})
	}
}

func TestPagination_GetOffset(t *testing.T) {
	testCases := []struct {
		name     string
		input    dto.Pagination
		expected int
	}{
		{
			name:     "first page",
			input:    dto.Pagination{Page: 1, Limit: 20},
			expected: 0,
		},
		{
			name:     "second page",
			input:    dto.Pagination{Page: 2, Limit: 20},
			expected: 20,
		},
		{
			name:     "tenth page with default limit",
			input:    dto.Pagination{Page: 10},
			expected: 90,
		},
		{
			name:     "invalid page uses default",
			input:    dto.Pagination{Page: 0, Limit: 50},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.input.GetOffset())
		})
	}
}
