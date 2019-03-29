package htmlui

import (
	"fmt"
	"net/url"
	"strconv"
)

// URLQuery url query
type URLQuery []URLQueryParam

// Clone clones url query params
func (ps URLQuery) Clone() URLQuery {
	result := make([]URLQueryParam, len(ps))
	for i, p := range ps {
		result[i] = p.Clone()
	}
	return result
}

// From loads url query param value from the specified url.Values
func (ps URLQuery) From(values url.Values) URLQuery {
	result := ps.Clone()

	for _, p := range result {
		p.Decode(values.Get(p.Name()))
	}

	return result
}

// With sets name param to the specified value
func (ps URLQuery) With(name, value string) URLQuery {
	result := ps.Clone()

	p := result.Get(name)
	if p == nil {
		panic(fmt.Errorf("failed to find %s url query param", name))
	}
	p.Decode(value)

	return result
}

// Get returns a query param by it's name
func (ps URLQuery) Get(name string) URLQueryParam {
	for _, p := range ps {
		if p.Name() == name {
			return p
		}
	}
	return nil
}

// Encode encodes url query
func (ps URLQuery) Encode() string {
	v := url.Values{}
	for _, p := range ps {
		v.Set(p.Name(), p.Encode())
	}
	return v.Encode()
}

// URLQueryParam url query param
type URLQueryParam interface {
	Name() string
	Decode(value string)
	Encode() string
	Clone() URLQueryParam
}

// StringURLQueryParam string url query param
type StringURLQueryParam struct {
	name  string
	value string
	def   string
}

// NewStringURLQueryParam creates new StringURLQueryParam
func NewStringURLQueryParam(name, defaultValue string) *StringURLQueryParam {
	return &StringURLQueryParam{
		name:  name,
		value: defaultValue,
		def:   defaultValue,
	}
}

// Name implements interface URLQueryParam
func (p *StringURLQueryParam) Name() string {
	return p.name
}

// Value returns param value
func (p *StringURLQueryParam) Value() string {
	return p.value
}

// Decode implements interface URLQueryParam
func (p *StringURLQueryParam) Decode(value string) {
	p.value = value
}

// Encode implements interface URLQueryParam
func (p *StringURLQueryParam) Encode() string {
	return p.value
}

// Clone implements interface URLQueryParam
func (p *StringURLQueryParam) Clone() URLQueryParam {
	result := *p
	return &result
}

// IntegerURLQueryParam integer url query param
type IntegerURLQueryParam struct {
	name  string
	value int
	def   int
}

// NewIntegerURLQueryParam creates new IntegerURLQueryParam
func NewIntegerURLQueryParam(name string, defaultValue int) *IntegerURLQueryParam {
	return &IntegerURLQueryParam{
		name:  name,
		value: defaultValue,
		def:   defaultValue,
	}
}

// Name implements interface URLQueryParam
func (p *IntegerURLQueryParam) Name() string {
	return p.name
}

// Value returns param value
func (p *IntegerURLQueryParam) Value() int {
	return p.value
}

// Decode implements interface URLQueryParam
func (p *IntegerURLQueryParam) Decode(value string) {
	v, err := strconv.ParseInt(value, 10, 32)
	if err == nil {
		p.value = int(v)
	} else {
		p.value = p.def
	}
}

// Encode implements interface URLQueryParam
func (p *IntegerURLQueryParam) Encode() string {
	return strconv.FormatInt(int64(p.value), 10)
}

// Clone implements interface URLQueryParam
func (p *IntegerURLQueryParam) Clone() URLQueryParam {
	result := *p
	return &result
}
