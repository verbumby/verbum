package dictionary

type Common struct {
	id            string
	indexID       string
	boost         float32
	aliases       []string
	title         string
	abbrevs       *Abbrevs
	slugifier     string
	unlisted      bool
	indexSettings IndexSettings
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

func (d Common) Boost() float32 {
	return d.boost
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

func (d Common) Slugifier() string {
	return d.slugifier
}

func (d Common) Unlisted() bool {
	return d.unlisted
}

func (d Common) IndexSettings() IndexSettings {
	return d.indexSettings
}
