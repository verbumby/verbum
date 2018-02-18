package headword

//go:generate reform

//Headword represents a dictionary
//
//reform:headwords
type Headword struct {
	ID        int32  `reform:"id,pk"`
	Headword  string `reform:"headword"`
	ArticleID int32  `reform:"article_id"`
}
