package htmlui

import (
	"net/url"
)

// Query url query
type Query []QueryParam

// Clone clones url query params
func (ps Query) Clone() Query {
	result := make([]QueryParam, len(ps))
	for i, p := range ps {
		result[i] = p.Clone()
	}
	return result
}

// From loads url query param value from the specified url.Values
func (ps Query) From(values url.Values) {
	for _, p := range ps {
		p.Decode(values.Get(p.Name()))
	}
}

// Get returns a query param by it's name
func (ps Query) Get(name string) QueryParam {
	for _, p := range ps {
		if p.Name() == name {
			return p
		}
	}
	return nil
}

// Encode encodes url query
func (ps Query) Encode() string {
	v := url.Values{}
	for _, p := range ps {
		if !p.ValueEqualsDefault() {
			v.Set(p.Name(), p.Encode())
		}
	}
	return v.Encode()
}
