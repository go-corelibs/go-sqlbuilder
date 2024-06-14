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

func IsColumnError(c Column) bool {
	_, ok := c.(*cErrorColumn)
	return ok
}

func GetColumnError(c Column) (err error) {
	if v, ok := c.(*cErrorColumn); ok {
		return v.err
	}
	return
}

type cErrorColumn struct {
	err error
}

func newErrorColumn(err error) Column {
	return &cErrorColumn{
		err: err,
	}
}

func (c *cErrorColumn) table_name() string {
	return ""
}

func (c *cErrorColumn) column_name() string {
	return ""
}

func (c *cErrorColumn) config() ColumnConfig {
	return nil
}

func (c *cErrorColumn) acceptType(interface{}) bool {
	return false
}

func (c *cErrorColumn) serialize(b *builder) {
	b.SetError(c.err)
	return
}

func (c *cErrorColumn) As(string) Column {
	return c
}

func (c *cErrorColumn) Eq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "=")
}

func (c *cErrorColumn) NotEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<>")
}

func (c *cErrorColumn) Gt(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, ">")
}

func (c *cErrorColumn) GtEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, ">=")
}

func (c *cErrorColumn) Lt(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<")
}

func (c *cErrorColumn) LtEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<=")
}

func (c *cErrorColumn) Like(right string) Condition {
	return newBinaryOperationCondition(c, right, " LIKE ")
}

func (c *cErrorColumn) NotLike(right string) Condition {
	return newBinaryOperationCondition(c, right, " NOT LIKE ")
}

func (c *cErrorColumn) Between(lower, higher interface{}) Condition {
	return newBetweenCondition(c, lower, higher)
}

func (c *cErrorColumn) In(val ...interface{}) Condition {
	return newInCondition(false, c, val...)
}

func (c *cErrorColumn) NotIn(val ...interface{}) Condition {
	return newInCondition(true, c, val...)
}

func (c *cErrorColumn) Describe() (output string) {
	// not implemented yet
	return
}
