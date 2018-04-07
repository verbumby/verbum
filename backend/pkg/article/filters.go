package article

import (
	"strconv"

	"github.com/pkg/errors"
)

// FilterDictID filter by dict it
type FilterDictID struct{}

// Name filter name
func (f *FilterDictID) Name() string {
	return "DictID"
}

// ToSQL generates sql
func (f *FilterDictID) ToSQL(param string) (from string, fromArgs []interface{}, where string, whereArgs []interface{}, err error) {
	var paramInt64 int64
	paramInt64, err = strconv.ParseInt(param, 10, 32)
	if err != nil {
		err = errors.Wrapf(err, "parse %s to int", param)
	}
	where = "dict_id = ?"
	whereArgs = []interface{}{paramInt64}
	return
}

// FilterTitlePrefix filter by title prefix
type FilterTitlePrefix struct{}

// Name filter name
func (f *FilterTitlePrefix) Name() string {
	return "TitlePrefix"
}

// ToSQL generates sql
func (f *FilterTitlePrefix) ToSQL(param string) (from string, fromArgs []interface{}, where string, whereArgs []interface{}, err error) {
	where = "title LIKE ?"
	whereArgs = []interface{}{param + "%"}
	return
}

// FilterTaks filter by task
type FilterTaks struct{}

// Name filter name
func (f *FilterTaks) Name() string {
	return "TaskID"
}

// ToSQL generates sql
func (f *FilterTaks) ToSQL(param string) (from string, fromArgs []interface{}, where string, whereArgs []interface{}, err error) {
	var paramInt64 int64
	paramInt64, err = strconv.ParseInt(param, 10, 32)
	if err != nil {
		err = errors.Wrapf(err, "parse %s to int", param)
	}

	from = "INNER JOIN `tasks_articles_rel` ON `tasks_articles_rel`.`article_id` = `articles`.`id` " +
		"AND `tasks_articles_rel`.`task_id` = ?  " +
		"AND `tasks_articles_rel`.`status` = 'PENDING' "
	fromArgs = []interface{}{paramInt64}
	return
}
