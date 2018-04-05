package article

import (
	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/db"
)

//go:generate reform

// Article represents an article
//
// reform:articles
type Article struct {
	ID      int32  `reform:"id,pk"`
	Title   string `reform:"title"`
	Content string `reform:"content"`
	DictID  int32  `reform:"dict_id"`
	Tasks   []Task
}

// Task article task
type Task struct {
	Task struct {
		ID    int32
		Title string
	}
	Status string
}

// LoadRelationships implements app.Model interface
func (a *Article) LoadRelationships() error {
	query := "SELECT t.id, t.title, tar.status FROM tasks_articles_rel tar " +
		"INNER JOIN tasks t ON t.id = tar.task_id " +
		"WHERE tar.article_id = " + db.DB.Placeholder(1)
	rows, err := db.DB.Query(query, a.ID)
	if err != nil {
		return errors.Wrap(err, "select article tasks")
	}
	defer rows.Close()

	a.Tasks = []Task{}
	for rows.Next() {
		at := Task{}
		if err := rows.Scan(&at.Task.ID, &at.Task.Title, &at.Status); err != nil {
			return errors.Wrap(err, "scan article related task")
		}
		a.Tasks = append(a.Tasks, at)
	}
	return nil
}

// UpdateRelationships implements app.Model interface
func (a *Article) UpdateRelationships() error {
	query := "UPDATE tasks_articles_rel SET status = ? WHERE task_id = ? AND article_id = ?"
	for _, t := range a.Tasks {
		if _, err := db.DB.Exec(query, t.Status, t.Task.ID, a.ID); err != nil {
			return errors.Wrap(err, "update article task status")
		}
	}
	return nil
}
