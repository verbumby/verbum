package htmlui

import (
	"net/url"
	"strconv"
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
		v.Set(p.Name(), p.Encode())
	}
	return v.Encode()
}

// QueryParam url query param
type QueryParam interface {
	Name() string
	Decode(value string)
	Encode() string
	Clone() QueryParam
}

// StringQueryParam string url query param
type StringQueryParam struct {
	name  string
	value string
	def   string
}

// NewStringQueryParam creates new StringURLQueryParam
func NewStringQueryParam(name, defaultValue string) *StringQueryParam {
	return &StringQueryParam{
		name:  name,
		value: defaultValue,
		def:   defaultValue,
	}
}

// Name implements interface URLQueryParam
func (p *StringQueryParam) Name() string {
	return p.name
}

// Value returns param value
func (p *StringQueryParam) Value() string {
	return p.value
}

// SetValue set's value
func (p *StringQueryParam) SetValue(v string) {
	p.value = v
}

// Decode implements interface URLQueryParam
func (p *StringQueryParam) Decode(value string) {
	p.value = value
}

// Encode implements interface URLQueryParam
func (p *StringQueryParam) Encode() string {
	return p.value
}

// Clone implements interface URLQueryParam
func (p *StringQueryParam) Clone() QueryParam {
	result := *p
	return &result
}

// IntegerQueryParam integer url query param
type IntegerQueryParam struct {
	name  string
	value int
	def   int
}

// NewIntegerQueryParam creates new IntegerURLQueryParam
func NewIntegerQueryParam(name string, defaultValue int) *IntegerQueryParam {
	return &IntegerQueryParam{
		name:  name,
		value: defaultValue,
		def:   defaultValue,
	}
}

// Name implements interface URLQueryParam
func (p *IntegerQueryParam) Name() string {
	return p.name
}

// Value returns param value
func (p *IntegerQueryParam) Value() int {
	return p.value
}

// SetValue set's value
func (p *IntegerQueryParam) SetValue(v int) {
	p.value = v
}

// Decode implements interface URLQueryParam
func (p *IntegerQueryParam) Decode(value string) {
	v, err := strconv.ParseInt(value, 10, 32)
	if err == nil {
		p.value = int(v)
	} else {
		p.value = p.def
	}
}

// Encode implements interface URLQueryParam
func (p *IntegerQueryParam) Encode() string {
	return strconv.FormatInt(int64(p.value), 10)
}

// Clone implements interface URLQueryParam
func (p *IntegerQueryParam) Clone() QueryParam {
	result := *p
	return &result
}
