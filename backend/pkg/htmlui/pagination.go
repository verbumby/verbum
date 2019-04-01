package htmlui

import "strconv"

// Pagination pagination
type Pagination struct {
	Current   int
	Total     int
	PageToURL func(n int) string
}

// PaginationLink pagination link
type PaginationLink struct {
	URL      string
	Text     string
	Active   bool
	Disabled bool
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

	const d = 2
	if max(p.Current-d, 1) > 1 {
		result = append(result, PaginationLink{
			URL:  p.PageToURL(1),
			Text: strconv.FormatInt(1, 10),
		})
	}

	if max(p.Current-d, 1) > 2 {
		result = append(result, PaginationLink{
			URL:      "",
			Disabled: true,
			Text:     "...",
		})
	}

	for i := max(p.Current-d, 1); i < p.Current; i++ {
		result = append(result, PaginationLink{
			URL:  p.PageToURL(i),
			Text: strconv.FormatInt(int64(i), 10),
		})
	}

	result = append(result, PaginationLink{
		URL:    p.PageToURL(p.Current),
		Text:   strconv.FormatInt(int64(p.Current), 10),
		Active: true,
	})

	for i := p.Current + 1; i <= min(p.Current+d, p.Total); i++ {
		result = append(result, PaginationLink{
			URL:  p.PageToURL(i),
			Text: strconv.FormatInt(int64(i), 10),
		})
	}

	if min(p.Current+d, p.Total) < p.Total-1 {
		result = append(result, PaginationLink{
			URL:      "",
			Disabled: true,
			Text:     "...",
		})
	}

	if min(p.Current+d, p.Total) < p.Total {
		result = append(result, PaginationLink{
			URL:  p.PageToURL(p.Total),
			Text: strconv.FormatInt(int64(p.Total), 10),
		})
	}

	return result
}
