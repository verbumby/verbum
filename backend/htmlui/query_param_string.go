package htmlui

type QueryParam interface {
	Name() string
	Decode(value string)
	Encode() string
	Clone() QueryParam
	Reset()
	ValueEqualsDefault() bool
}

type StringQueryParam struct {
	name  string
	value string
	def   string
}

func NewStringQueryParam(name, defaultValue string) *StringQueryParam {
	return &StringQueryParam{
		name:  name,
		value: defaultValue,
		def:   defaultValue,
	}
}

func (p *StringQueryParam) Name() string {
	return p.name
}

func (p *StringQueryParam) Value() string {
	return p.value
}

func (p *StringQueryParam) SetValue(v string) {
	p.value = v
}

func (p *StringQueryParam) Decode(value string) {
	p.value = value
}

func (p *StringQueryParam) Encode() string {
	return p.value
}

func (p *StringQueryParam) Clone() QueryParam {
	result := *p
	return &result
}

func (p *StringQueryParam) Reset() {
	p.value = p.def
}

func (p *StringQueryParam) ValueEqualsDefault() bool {
	return p.def == p.value
}
