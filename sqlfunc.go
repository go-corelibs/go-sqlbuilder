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

import (
	"fmt"
)

// SqlFunc represents an SQL function (ex:count(*)) which can be use in the same way as a Column
type SqlFunc interface {
	Column

	columns() []Column
}

// Func returns new SQL function.  The name is function name, and the args is arguments of function
func Func(name string, args ...Column) SqlFunc {
	return &cSqlFunc{
		name: name,
		args: args,
	}
}

type cSqlFunc struct {
	name string
	args cSqlFuncColumnList
}

func (c *cSqlFunc) As(alias string) Column {
	return &cColumnAlias{
		column: c,
		alias:  alias,
	}
}

func (c *cSqlFunc) table_name() string {
	return ""
}

func (c *cSqlFunc) column_name() string {
	return c.name
}

func (c *cSqlFunc) not_null() bool {
	return true
}

func (c *cSqlFunc) config() ColumnConfig {
	return nil
}

func (c *cSqlFunc) acceptType(interface{}) bool {
	return false
}

func (c *cSqlFunc) serialize(b *builder) {
	b.Append(c.name)
	b.Append("(")
	b.AppendItem(c.args)
	b.Append(")")
}

func (c *cSqlFunc) hasColumn(t *cTable) (present bool) {
	for _, fncol := range c.columns() {
		if fncolfn, ok := fncol.(*cSqlFunc); ok {
			if present = fncolfn.hasColumn(t); present {
				return
			}
		} else {
			for _, col := range t.columns {
				if present = SameColumn(col, fncol); present {
					return
				}
			}
		}
	}
	return
}

func (c *cSqlFunc) Eq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "=")
}

func (c *cSqlFunc) NotEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<>")
}

func (c *cSqlFunc) Gt(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, ">")
}

func (c *cSqlFunc) GtEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, ">=")
}

func (c *cSqlFunc) Lt(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<")
}

func (c *cSqlFunc) LtEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<=")
}

func (c *cSqlFunc) Like(right string) Condition {
	return newBinaryOperationCondition(c, right, " LIKE ")
}

func (c *cSqlFunc) NotLike(right string) Condition {
	return newBinaryOperationCondition(c, right, " NOT LIKE ")
}

func (c *cSqlFunc) Between(lower, higher interface{}) Condition {
	return newBetweenCondition(c, lower, higher)
}

func (c *cSqlFunc) In(vals ...interface{}) Condition {
	return newInCondition(false, c, vals...)
}

func (c *cSqlFunc) NotIn(vals ...interface{}) Condition {
	return newInCondition(true, c, vals...)
}

func (c *cSqlFunc) columns() []Column {
	return c.args
}

func (c *cSqlFunc) Describe() (output string) {
	output = fmt.Sprintf("%q(%s)", c.name, c.args.Describe())
	return
}
