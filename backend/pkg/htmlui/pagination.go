package htmlui

// Pagination pagination
type Pagination struct {
	Current   int
	Total     int
	PageToURL func(n int) string
}

// PaginationLink pagination link
type PaginationLink struct {
	URL    string
	PageN  int
	Active bool
}

// Links returns pagination links
func (p Pagination) Links() []PaginationLink {
	result := []PaginationLink{}

	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	for i := max(p.Current-5, 1); i < p.Current; i++ {
		result = append(result, PaginationLink{
			URL:   p.PageToURL(i),
			PageN: i,
		})
	}

	result = append(result, PaginationLink{
		URL:    p.PageToURL(p.Current),
		PageN:  p.Current,
		Active: true,
	})

	for i := p.Current + 1; i <= min(p.Current+5, p.Total); i++ {
		result = append(result, PaginationLink{
			URL:   p.PageToURL(i),
			PageN: i,
		})
	}

	return result
}
