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
	"database/sql/driver"
	"reflect"
)

type cColumnImpl struct {
	*cColumnImplConfig
	table Table
}

func (c *cColumnImpl) table_name() string {
	return c.table.Name()
}

func (c *cColumnImpl) column_name() string {
	return c.name
}

func (c *cColumnImpl) config() ColumnConfig {
	return c.cColumnImplConfig
}

func (c *cColumnImpl) hasColumn(t *cTable) (present bool) {
	if present = c == Star; present {
		return
	}
	for _, col := range t.Columns() {
		if present = SameColumn(col, c); present {
			return
		}
	}
	return
}

func (c *cColumnImpl) acceptType(val interface{}) bool {
	lit, ok := val.(literal)
	if !ok || lit == nil {
		return false
	}
	if lit.Raw() == nil {
		return !c.opt.NotNull
	}
	if c.Type() == ColumnTypeAny {
		return true
	}
	if _, ok := lit.Raw().(driver.Valuer); ok {
		return true
	}

	valt := reflect.TypeOf(lit.Raw())
	for _, t := range c.typ.CapableTypes() {
		if t == valt {
			return true
		}
	}
	return false
}

func (c *cColumnImpl) serialize(bldr *builder) {
	if c == Star {
		bldr.Append("*")
	} else {
		bldr.Append(bldr.dialect.QuoteField(c.table.Name()) + "." + bldr.dialect.QuoteField(c.name))
	}
	return
}

func (c *cColumnImpl) As(alias string) Column {
	return &cColumnAlias{
		column: c,
		alias:  alias,
	}
}

func (c *cColumnImpl) Eq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "=")
}

func (c *cColumnImpl) NotEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<>")
}

func (c *cColumnImpl) Gt(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, ">")
}

func (c *cColumnImpl) GtEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, ">=")
}

func (c *cColumnImpl) Lt(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<")
}

func (c *cColumnImpl) LtEq(right interface{}) Condition {
	return newBinaryOperationCondition(c, right, "<=")
}

func (c *cColumnImpl) Like(right string) Condition {
	return newBinaryOperationCondition(c, right, " LIKE ")
}

func (c *cColumnImpl) NotLike(right string) Condition {
	return newBinaryOperationCondition(c, right, " NOT LIKE ")
}

func (c *cColumnImpl) Between(lower, higher interface{}) Condition {
	return newBetweenCondition(c, lower, higher)
}

func (c *cColumnImpl) In(val ...interface{}) Condition {
	return newInCondition(false, c, val...)
}

func (c *cColumnImpl) NotIn(val ...interface{}) Condition {
	return newInCondition(true, c, val...)
}

func (c *cColumnImpl) Describe() (output string) {
	// not implemented yet
	//output = c.opt.Describe()
	output = c.name
	return
}
