package htmlui

import "strconv"

type IntegerQueryParam struct {
	name  string
	value int
	def   int
}

func NewIntegerQueryParam(name string, defaultValue int) *IntegerQueryParam {
	return &IntegerQueryParam{
		name:  name,
		value: defaultValue,
		def:   defaultValue,
	}
}

func (p *IntegerQueryParam) Name() string {
	return p.name
}

func (p *IntegerQueryParam) Value() int {
	return p.value
}

func (p *IntegerQueryParam) SetValue(v int) {
	p.value = v
}

func (p *IntegerQueryParam) Decode(value string) {
	v, err := strconv.ParseInt(value, 10, 32)
	if err == nil {
		p.value = int(v)
	} else {
		p.value = p.def
	}
}

func (p *IntegerQueryParam) Encode() string {
	return strconv.FormatInt(int64(p.value), 10)
}

func (p *IntegerQueryParam) Clone() QueryParam {
	result := *p
	return &result
}

func (p *IntegerQueryParam) Reset() {
	p.value = p.def
}

func (p *IntegerQueryParam) ValueEqualsDefault() bool {
	return p.def == p.value
}
