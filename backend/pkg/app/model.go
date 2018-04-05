package app

import reform "gopkg.in/reform.v1"

// Model application model interface
type Model interface {
	reform.Record
	LoadRelationships() error
	UpdateRelationships() error
}

// ModelMeta application model metadata inteface
type ModelMeta interface {
	NewModel() Model
}
