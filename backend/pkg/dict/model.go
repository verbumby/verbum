package dict

//go:generate reform

//Dict represents a dictionary
//
//reform:dicts
type Dict struct {
	ID    int32  `reform:"id,pk"`
	Title string `reform:"title"`
}
