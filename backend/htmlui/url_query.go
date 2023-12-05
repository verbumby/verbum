package htmlui

import (
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/verbumby/verbum/backend/dictionary"
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

// QueryParam url query param
type QueryParam interface {
	Name() string
	Decode(value string)
	Encode() string
	Clone() QueryParam
	Reset()
	ValueEqualsDefault() bool
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

// Reset implements interface URLQueryParam
func (p *StringQueryParam) Reset() {
	p.value = p.def
}

// ValueEqualsDefault implements interface URLQueryParam
func (p *StringQueryParam) ValueEqualsDefault() bool {
	return p.def == p.value
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

// Reset implements interface URLQueryParam
func (p *IntegerQueryParam) Reset() {
	p.value = p.def
}

// ValueEqualsDefault implements interface URLQueryParam
func (p *IntegerQueryParam) ValueEqualsDefault() bool {
	return p.def == p.value
}

// StringQueryParam string url query param
type InDictsQueryParam struct {
	name  string
	value []string
	def   []string
}

// NewStringQueryParam creates new StringURLQueryParam
func NewInDictsQueryParam(name string) *InDictsQueryParam {
	dicts := dictionary.GetAllListed()
	def := make([]string, len(dicts))
	for i, d := range dicts {
		def[i] = d.ID()
	}
	value := make([]string, len(def))
	copy(value, def)

	return &InDictsQueryParam{
		name:  name,
		value: value,
		def:   def,
	}
}

// Name implements interface URLQueryParam
func (p *InDictsQueryParam) Name() string {
	return p.name
}

// Value returns param value
func (p *InDictsQueryParam) Value() []string {
	return p.value
}

// SetValue set's value
func (p *InDictsQueryParam) SetValue(vs []string) {
	value := []string{}
	for _, d := range dictionary.GetAll() {
		contains := false
		for _, v := range vs {
			if d.ID() == v || slices.Contains(d.Aliases(), v) {
				contains = true
				break
			}
		}
		if contains {
			value = append(value, d.ID())
		}
	}
	p.value = value
}

// Decode implements interface URLQueryParam
func (p *InDictsQueryParam) Decode(value string) {
	if value == "" {
		p.Reset()
		return
	}

	vs := strings.Split(value, ",")
	p.SetValue(vs)
}

// Encode implements interface URLQueryParam
func (p *InDictsQueryParam) Encode() string {
	return strings.Join(p.value, ",")
}

// Clone implements interface URLQueryParam
func (p *InDictsQueryParam) Clone() QueryParam {
	result := InDictsQueryParam{}
	result.name = p.name
	result.value = make([]string, len(p.value))
	copy(result.value, p.value)
	result.def = make([]string, len(p.def))
	copy(result.def, p.def)
	return &result
}

// Reset implements interface URLQueryParam
func (p *InDictsQueryParam) Reset() {
	p.value = make([]string, len(p.def))
	copy(p.value, p.def)
}

// ValueEqualsDefault implements interface URLQueryParam
func (p *InDictsQueryParam) ValueEqualsDefault() bool {
	if len(p.value) != len(p.def) {
		return false
	}

	for i := range p.value {
		if p.value[i] != p.def[i] {
			return false
		}
	}

	return true
}
