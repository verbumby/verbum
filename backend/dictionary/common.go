package dictionary

type Common struct {
	id                      string
	indexID                 string
	aliases                 []string
	title                   string
	abbrevs                 *Abbrevs
	prependContentWithTitle bool
	slugifier               string
}

func (d Common) ID() string {
	return d.id
}

func (d Common) IndexID() string {
	if d.indexID == "" {
		return d.id
	}
	return d.indexID
}

func (d Common) Aliases() []string {
	return d.aliases
}

func (d Common) Title() string {
	return d.title
}

func (d Common) Abbrevs() *Abbrevs {
	return d.abbrevs
}

func (d Common) PrependContentWithTitle() bool {
	return d.prependContentWithTitle
}

func (d Common) Slugifier() string {
	return d.slugifier
}