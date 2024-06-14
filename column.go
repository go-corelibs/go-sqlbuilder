// Copyright (c) 2014 umisama <Takaaki IBARAKI>
// Copyright (c)  The Go-CoreLibs Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package sqlbuilder

// Column represents a table column.
type Column interface {
	serializable

	table_name() string
	column_name() string
	config() ColumnConfig
	acceptType(interface{}) bool

	// As creates Column alias.
	As(alias string) Column

	// Eq creates Condition for "column==right".  Type for right is column's one or other Column.
	Eq(right interface{}) Condition

	// NotEq creates Condition for "column<>right".  Type for right is column's one or other Column.
	NotEq(right interface{}) Condition

	// Gt creates Condition for "column>right".  Type for right is column's one or other Column.
	Gt(right interface{}) Condition

	// GtEq creates Condition for "column>=right".  Type for right is column's one or other Column.
	GtEq(right interface{}) Condition

	// Lt creates Condition for "column<right".  Type for right is column's one or other Column.
	Lt(right interface{}) Condition

	// LtEq creates Condition for "column<=right".  Type for right is column's one or other Column.
	LtEq(right interface{}) Condition

	// Like creates Condition for "column LIKE right".  Type for right is column's one or other Column.
	Like(right string) Condition

	// NotLike creates Condition for "column NOT LIKE right".  Type for right is column's one or other Column.
	NotLike(right string) Condition

	// Between creates Condition for "column BETWEEN lower AND higher".  Type for lower/higher is int or time.Time.
	Between(lower, higher interface{}) Condition

	// In creates Condition for "column IN (values[0], values[1] ...)".  Type for values is column's one or other Column.
	In(values ...interface{}) Condition

	// NotIn creates Condition for "column NOT IN (values[0], values[1] ...)".  Type for values is column's one or other Column.
	NotIn(values ...interface{}) Condition
}

func SameColumn(a, b Column) (same bool) {
	return a.table_name() == b.table_name() && a.column_name() == b.column_name()
}

// Star represents * column
var Star Column = &cColumnImpl{nil, nil}

// AnyColumn creates config for any types.
func AnyColumn(name string, opt *ColumnOption) ColumnConfig {
	return newColumnImplConfig(name, ColumnTypeAny, opt)
}

// IntColumn creates config for INTEGER type column.
func IntColumn(name string, opt *ColumnOption) ColumnConfig {
	return newColumnImplConfig(name, ColumnTypeInt, opt)
}

// StringColumn creates config for TEXT or VARCHAR type column.
func StringColumn(name string, opt *ColumnOption) ColumnConfig {
	return newColumnImplConfig(name, ColumnTypeString, opt)
}

// DateColumn creates config for DATETIME type column.
func DateColumn(name string, opt *ColumnOption) ColumnConfig {
	return newColumnImplConfig(name, ColumnTypeDate, opt)
}

// FloatColumn creates config for REAL or FLOAT type column.
func FloatColumn(name string, opt *ColumnOption) ColumnConfig {
	return newColumnImplConfig(name, ColumnTypeFloat, opt)
}

// BoolColumn creates config for BOOLEAN type column.
func BoolColumn(name string, opt *ColumnOption) ColumnConfig {
	return newColumnImplConfig(name, ColumnTypeBool, opt)
}

// BytesColumn creates config for BLOB type column.
func BytesColumn(name string, opt *ColumnOption) ColumnConfig {
	return newColumnImplConfig(name, ColumnTypeBytes, opt)
}

// TimeColumn creates config for DATETIME type column.
func TimeColumn(name string, opt *ColumnOption) ColumnConfig {
	return newColumnImplConfig(name, ColumnTypeBytes, opt)
}
