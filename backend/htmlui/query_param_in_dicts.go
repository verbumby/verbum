package htmlui

import (
	"slices"
	"strings"

	"github.com/verbumby/verbum/backend/dictionary"
)

type InDictsQueryParam struct {
	name  string
	value []string
	def   []string
}

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

func (p *InDictsQueryParam) Name() string {
	return p.name
}

func (p *InDictsQueryParam) Value() []string {
	return p.value
}

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

func (p *InDictsQueryParam) Decode(value string) {
	if value == "" {
		p.Reset()
		return
	}

	vs := strings.Split(value, ",")
	p.SetValue(vs)
}

func (p *InDictsQueryParam) Encode() string {
	return strings.Join(p.value, ",")
}

func (p *InDictsQueryParam) Clone() QueryParam {
	result := InDictsQueryParam{}
	result.name = p.name
	result.value = make([]string, len(p.value))
	copy(result.value, p.value)
	result.def = make([]string, len(p.def))
	copy(result.def, p.def)
	return &result
}

func (p *InDictsQueryParam) Reset() {
	p.value = make([]string, len(p.def))
	copy(p.value, p.def)
}

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
