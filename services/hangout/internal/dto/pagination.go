package dto

import (
	"fmt"
	"strings"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/google/uuid"
)

type CursorPagination struct {
	Limit   int        `json:"limit"`
	AfterID *uuid.UUID `json:"after_id"`
	SortBy  string     `json:"sort_by"`
	SortDir string     `json:"sort_dir"`
}

func (p *CursorPagination) GetLimit() int {
	if p.Limit <= 0 {
		return constants.DefaultLimit
	}
	if p.Limit > constants.MaxLimit {
		return constants.MaxLimit
	}
	return p.Limit
}

func (p *CursorPagination) GetSortBy() string {
	switch p.SortBy {
	case constants.SortByDate, constants.SortByCreatedAt:
		return p.SortBy
	default:
		return constants.SortByCreatedAt
	}
}

func (p *CursorPagination) GetSortDir() string {
	switch strings.ToLower(p.SortDir) {
	case constants.SortDirectionAsc:
		return constants.SortDirectionAsc
	case constants.SortDirectionDesc:
		return constants.SortDirectionDesc
	default:
		return constants.SortDirectionDesc
	}
}

func (p *CursorPagination) GetOrderByClause() string {
	return fmt.Sprintf("%s %s, id %s", p.GetSortBy(), p.GetSortDir(), p.GetSortDir())
}

type PaginatedHangouts struct {
	Data       []*HangoutListItemResponse `json:"data"`
	NextCursor *uuid.UUID                 `json:"next_cursor"`
	HasMore    bool                       `json:"has_more"`
}
