package article

//go:generate reform

// Article represents an article
//
// reform:articles
type Article struct {
	ID      int32  `reform:"id,pk"`
	Content string `reform:"content"`
}
