package util

type Page struct {
	Page     int
	PageSize int
}

func Normalize(p, ps int) Page {
	if p <= 0 { p = 1 }
	if ps <= 0 { ps = 20 }
	if ps > 200 { ps = 200 }
	return Page{Page: p, PageSize: ps}
}

func (p Page) Offset() int { return (p.Page - 1) * p.PageSize }
