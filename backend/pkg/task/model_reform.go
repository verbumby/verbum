// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package task

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type taskTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *taskTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("tasks").
func (v *taskTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *taskTableType) Columns() []string {
	return []string{"id", "title"}
}

// NewStruct makes a new struct for that view or table.
func (v *taskTableType) NewStruct() reform.Struct {
	return new(Task)
}

// NewRecord makes a new record for that table.
func (v *taskTableType) NewRecord() reform.Record {
	return new(Task)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *taskTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// TaskTable represents tasks view or table in SQL database.
var TaskTable = &taskTableType{
	s: parse.StructInfo{Type: "Task", SQLSchema: "", SQLName: "tasks", Fields: []parse.FieldInfo{{Name: "ID", Type: "int32", Column: "id"}, {Name: "Title", Type: "string", Column: "title"}}, PKFieldIndex: 0},
	z: new(Task).Values(),
}

// String returns a string representation of this struct or record.
func (s Task) String() string {
	res := make([]string, 2)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Title: " + reform.Inspect(s.Title, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Task) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Title,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Task) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Title,
	}
}

// View returns View object for that struct.
func (s *Task) View() reform.View {
	return TaskTable
}

// Table returns Table object for that record.
func (s *Task) Table() reform.Table {
	return TaskTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Task) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Task) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Task) HasPK() bool {
	return s.ID != TaskTable.z[TaskTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Task) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = int32(i64)
	} else {
		s.ID = pk.(int32)
	}
}

// check interfaces
var (
	_ reform.View   = TaskTable
	_ reform.Struct = (*Task)(nil)
	_ reform.Table  = TaskTable
	_ reform.Record = (*Task)(nil)
	_ fmt.Stringer  = (*Task)(nil)
)

func init() {
	parse.AssertUpToDate(&TaskTable.s, new(Task))
}