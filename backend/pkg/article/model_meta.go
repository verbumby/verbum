package article

import "github.com/verbumby/verbum/backend/pkg/app"

type articleMetaType struct {
}

// ArticleMeta article meta
var ArticleMeta = &articleMetaType{}

func (am *articleMetaType) NewModel() app.Model {
	return ArticleTable.NewRecord().(app.Model)
}
