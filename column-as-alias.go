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

type cColumnAlias struct {
	column Column
	alias  string
}

func (c *cColumnAlias) table_name() string {
	return c.column.table_name()
}

func (c *cColumnAlias) column_name() string {
	return c.alias
}

func (c *cColumnAlias) config() ColumnConfig {
	return c.column.config()
}

func (c *cColumnAlias) acceptType(val interface{}) bool {
	return c.column.acceptType(val)
}

func (c *cColumnAlias) hasColumn(t *cTable) (present bool) {
	for _, col := range t.columns {
		if fncol, ok := c.column.(*cSqlFunc); ok {
			if present = fncol.hasColumn(t); present {
				return
			}
		} else if present = SameColumn(col, c.column); present {
			return
		}
	}
	return
}

func (c *cColumnAlias) As(alias string) Column {
	return &cColumnAlias{
		column: c,
		alias:  alias,
	}
}

func (c *cColumnAlias) serialize(b *builder) {
	b.Append(b.dialect.QuoteField(c.alias))
	return
}

func (c *cColumnAlias) column_alias() string {
	return c.alias
}

func (c *cColumnAlias) source() Column {
	return c.column
}

func (c *cColumnAlias) Eq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "=")
}

func (c *cColumnAlias) NotEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<>")
}

func (c *cColumnAlias) Gt(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, ">")
}

func (c *cColumnAlias) GtEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, ">=")
}

func (c *cColumnAlias) Lt(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<")
}

func (c *cColumnAlias) LtEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<=")
}

func (c *cColumnAlias) Like(right string) Condition {
	return newBinaryOperationCondition(c, right, " LIKE ")
}

func (c *cColumnAlias) NotLike(right string) Condition {
	return newBinaryOperationCondition(c, right, " NOT LIKE ")
}

func (c *cColumnAlias) Between(lower, higher interface{}) Condition {
	return newBetweenCondition(c, lower, higher)
}

func (c *cColumnAlias) In(val ...interface{}) Condition {
	return newInCondition(false, c, val...)
}

func (c *cColumnAlias) NotIn(val ...interface{}) Condition {
	return newInCondition(true, c, val...)
}

func (c *cColumnAlias) Describe() (output string) {
	// not implemented yet
	return
}
