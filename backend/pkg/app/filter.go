package app

// Filter filter interface
type Filter interface {
	Name() string
	ToSQL(param string) (from string, fromArgs []interface{}, where string, whereArgs []interface{}, err error)
}
