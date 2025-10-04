package dto

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type Pagination struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		return DefaultLimit
	}
	if p.Limit > MaxLimit {
		return MaxLimit
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		return DefaultPage
	}
	return p.Page
}
