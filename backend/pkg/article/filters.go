package article

import (
	"strconv"

	"github.com/pkg/errors"
)

type FilterDictID struct{}

func (f *FilterDictID) Name() string {
	return "DictID"
}

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

type FilterTitlePrefix struct{}

func (f *FilterTitlePrefix) Name() string {
	return "TitlePrefix"
}

func (f *FilterTitlePrefix) ToSQL(param string) (from string, fromArgs []interface{}, where string, whereArgs []interface{}, err error) {
	where = "title LIKE ?"
	whereArgs = []interface{}{param + "%"}
	return
}
