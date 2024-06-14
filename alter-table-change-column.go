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

type cAlterTableChangeColumn struct {
	table      *cTable
	old_column Column
	new_column ColumnConfig
	first      bool
	after      Column

	dialect Dialect
}

func (c *cAlterTableChangeColumn) serialize(b *builder) {
	b.Append("CHANGE COLUMN ")
	if name := c.old_column.column_name(); len(name) != 0 {
		b.Append(c.dialect.QuoteField(name))
	} else {
		b.AppendItem(c.old_column)
	}
	b.Append(" ")
	b.AppendItem(c.new_column)

	typ, err := c.dialect.ColumnTypeToString(c.new_column)
	if err != nil {
		b.SetError(err)
	} else if len(typ) == 0 {
		b.SetError(newError("column type is required.(maybe, a bug is in implements of dialect.)"))
	} else {
		b.Append(" ")
		b.Append(typ)
	}

	opt, err := c.dialect.ColumnOptionToString(c.new_column.Option())
	if err != nil {
		b.SetError(err)
	} else if len(opt) != 0 {
		b.Append(" ")
		b.Append(opt)
	}

	if c.first {
		b.Append(" FIRST")
	} else if c.after != nil {
		b.Append(" AFTER ")
		if colname := c.after.column_name(); len(colname) != 0 {
			b.Append(c.dialect.QuoteField(colname))
		} else {
			b.AppendItem(c.after)
		}
	}
}

func (c *cAlterTableChangeColumn) applyToTable() error {
	if c.first {
		return c.table.ChangeColumnFirst(c.old_column, c.new_column)
	}
	if c.after != nil {
		return c.table.ChangeColumnAfter(c.old_column, c.new_column, c.after)
	}
	return c.table.ChangeColumn(c.old_column, c.new_column)
}

func (c *cAlterTableChangeColumn) Describe() (output string) {
	// not implemented yet
	return
}
