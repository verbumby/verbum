package typeahead

//go:generate reform

//Typeahead represents a dictionary
//
//reform:typeaheads
type Typeahead struct {
	ID        int32  `reform:"id,pk"`
	Typeahead string `reform:"typeahead"`
	ArticleID int32  `reform:"article_id"`
}
