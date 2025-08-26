package htmlui

type BoolQueryParam struct {
	name  string
	value bool
	def   bool
}

func NewBoolQueryParam(name string, defaultValue bool) *BoolQueryParam {
	return &BoolQueryParam{
		name:  name,
		value: defaultValue,
		def:   defaultValue,
	}
}

func (p *BoolQueryParam) Name() string {
	return p.name
}

func (p *BoolQueryParam) Value() bool {
	return p.value
}

func (p *BoolQueryParam) SetValue(v bool) {
	p.value = v
}

func (p *BoolQueryParam) Decode(value string) {
	switch value {
	case "true":
		p.value = true
	case "false":
		p.value = false
	default:
		p.value = p.def
	}
}

func (p *BoolQueryParam) Encode() string {
	if p.value {
		return "true"
	} else {
		return "false"
	}
}

func (p *BoolQueryParam) Clone() QueryParam {
	result := *p
	return &result
}

func (p *BoolQueryParam) Reset() {
	p.value = p.def
}

func (p *BoolQueryParam) ValueEqualsDefault() bool {
	return p.def == p.value
}
