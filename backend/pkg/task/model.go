package task

//go:generate reform

// Task represents an article
//
// reform:tasks
type Task struct {
	ID    int32  `reform:"id,pk"`
	Title string `reform:"title"`
}
