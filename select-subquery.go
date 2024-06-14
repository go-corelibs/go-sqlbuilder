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
	"strconv"
)

type cSubQuery struct {
	stat  *cSelect
	alias string
	err   error
}

func newSubQuery(s *cSelect, alias string) *cSubQuery {
	m := &cSubQuery{
		stat:  s,
		alias: alias,
	}

	if len(alias) == 0 {
		m.err = newError("alias is empty.")
	}
	return m
}

func (c *cSubQuery) serialize(bldr *builder) {
	if c.err != nil {
		bldr.SetError(c.err)
	}

	bldr.Append("( ")
	bldr.AppendItem(c.stat)
	bldr.Append(" ) AS " + c.alias)
	return
}

func (c *cSubQuery) Name() string {
	return c.alias
}

func (c *cSubQuery) C(name string) Column {
	for _, col := range c.stat.columns {
		if ac, ok := col.(iAliasedColumn); ok {
			if ac.column_alias() == name {
				return col.config().toColumn(c)
			}
		}
		if col.column_name() == name {
			return col.config().toColumn(c)
		}
	}
	return newErrorColumn(newError("column %s was not found.", name))
}

func (c *cSubQuery) Columns() []Column {
	l := make([]Column, len(c.stat.columns))
	for _, col := range c.stat.columns {
		if _, ok := col.(iAliasedColumn); ok {
			l = append(l, col.config().toColumn(c))
		}
		l = append(l, col.config().toColumn(c))
	}
	return nil
}

func (c *cSubQuery) Option() *TableOption {
	return nil
}

func (c *cSubQuery) InnerJoin(Table, Condition) Table {
	c.err = newError("subquery can not join.")
	return c
}

func (c *cSubQuery) LeftOuterJoin(Table, Condition) Table {
	c.err = newError("subquery can not join.")
	return c
}

func (c *cSubQuery) RightOuterJoin(Table, Condition) Table {
	c.err = newError("subquery can not join.")
	return c
}

func (c *cSubQuery) FullOuterJoin(Table, Condition) Table {
	c.err = newError("subquery can not join.")
	return c
}

func (c *cSubQuery) hasColumn(trg Column) bool {
	if cimpl, ok := trg.(*cColumnImpl); ok {
		if trg == Star {
			return true
		}
		if cimpl.table != c {
			return false
		}
		for _, col := range c.stat.columns {
			if col.column_name() == trg.column_name() {
				return true
			}
		}
		return false
	}
	if acol, ok := trg.(*cColumnAlias); ok {
		if acol.column.(*cColumnImpl).table != c {
			return false
		}
		for _, col := range c.stat.columns {
			if col.column_name() == trg.column_name() {
				return true
			}
		}
		return false
	}
	if sqlfn, ok := trg.(*cSqlFunc); ok {
		for _, fncol := range sqlfn.columns() {
			find := false
			for _, col := range c.stat.columns {
				if col.column_name() == fncol.column_name() {
					find = true
				}
			}
			if !find {
				return false
			}
		}
		return true
	}
	return false
}

func (c *cSubQuery) Describe() (output string) {
	if c.stat != nil {
		output += c.stat.Describe()
		if c.alias != "" {
			output += " AS " + strconv.Quote(c.alias)
		}
	}
	return
}
