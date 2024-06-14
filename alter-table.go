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

// AlterTableBuilder is the Buildable interface wrapping of AlterTable
type AlterTableBuilder interface {
	RenameTo(name string) AlterTableBuilder
	AddColumn(col ColumnConfig) AlterTableBuilder
	AddColumnAfter(col ColumnConfig, after Column) AlterTableBuilder
	AddColumnFirst(col ColumnConfig) AlterTableBuilder
	DropColumn(col Column) AlterTableBuilder
	ChangeColumn(old_column Column, new_column ColumnConfig) AlterTableBuilder
	ChangeColumnAfter(old_column Column, new_column ColumnConfig, after Column) AlterTableBuilder
	ChangeColumnFirst(old_column Column, new_column ColumnConfig) AlterTableBuilder
	ToSql() (query string, args []interface{}, err error)
	ApplyToTable() error

	privateAlterTable()
}

func AlterTable(tbl Table) AlterTableBuilder {
	return alterTable(tbl, dialect())
}

func alterTable(tbl Table, d Dialect) *cAlterTable {
	if d == nil {
		d = dialect()
	}
	if tbl == nil {
		return &cAlterTable{
			err: newError("table is nil."),
		}
	}

	t, ok := tbl.(*cTable)
	if !ok {
		return &cAlterTable{
			err: newError("AlterTable can use only natural table."),
		}
	}
	return &cAlterTable{
		table:          t,
		add_columns:    make([]*cAlterTableAddColumn, 0),
		change_columns: make([]*cAlterTableChangeColumn, 0),
		dialect:        d,
	}
}

type cAlterTable struct {
	table          *cTable
	rename_to      string
	add_columns    []*cAlterTableAddColumn
	drop_columns   []Column
	change_columns []*cAlterTableChangeColumn

	err     error
	dialect Dialect
}

func (b *cAlterTable) privateAlterTable() {
	// nop
}

func (b *cAlterTable) RenameTo(name string) AlterTableBuilder {
	if b.err != nil {
		return b
	}

	b.rename_to = name
	return b
}

func (b *cAlterTable) AddColumn(col ColumnConfig) AlterTableBuilder {
	if b.err != nil {
		return b
	}

	b.add_columns = append(b.add_columns, &cAlterTableAddColumn{
		table:   b.table,
		column:  col,
		first:   false,
		after:   nil,
		dialect: b.dialect,
	})
	return b
}

func (b *cAlterTable) AddColumnAfter(col ColumnConfig, after Column) AlterTableBuilder {
	if b.err != nil {
		return b
	}

	b.add_columns = append(b.add_columns, &cAlterTableAddColumn{
		table:   b.table,
		column:  col,
		first:   false,
		after:   after,
		dialect: b.dialect,
	})
	return b
}

func (b *cAlterTable) AddColumnFirst(col ColumnConfig) AlterTableBuilder {
	if b.err != nil {
		return b
	}

	b.add_columns = append(b.add_columns, &cAlterTableAddColumn{
		table:   b.table,
		column:  col,
		first:   true,
		after:   nil,
		dialect: b.dialect,
	})
	return b
}

func (b *cAlterTable) DropColumn(col Column) AlterTableBuilder {
	if b.err != nil {
		return b
	}

	b.drop_columns = append(b.drop_columns, col)
	return b
}

func (b *cAlterTable) ChangeColumn(old_column Column, new_column ColumnConfig) AlterTableBuilder {
	if b.err != nil {
		return b
	}

	b.change_columns = append(b.change_columns, &cAlterTableChangeColumn{
		table:      b.table,
		old_column: old_column,
		new_column: new_column,
		first:      false,
		after:      nil,
		dialect:    b.dialect,
	})
	return b
}

func (b *cAlterTable) ChangeColumnAfter(old_column Column, new_column ColumnConfig, after Column) AlterTableBuilder {
	if b.err != nil {
		return b
	}

	b.change_columns = append(b.change_columns, &cAlterTableChangeColumn{
		table:      b.table,
		old_column: old_column,
		new_column: new_column,
		first:      false,
		after:      after,
		dialect:    b.dialect,
	})
	return b
}

func (b *cAlterTable) ChangeColumnFirst(old_column Column, new_column ColumnConfig) AlterTableBuilder {
	if b.err != nil {
		return b
	}

	b.change_columns = append(b.change_columns, &cAlterTableChangeColumn{
		table:      b.table,
		old_column: old_column,
		new_column: new_column,
		first:      true,
		after:      nil,
		dialect:    b.dialect,
	})
	return b
}

func (b *cAlterTable) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder(b.dialect)
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()
	if b.err != nil {
		bldr.SetError(b.err)
		return
	}

	bldr.Append("ALTER TABLE ")
	bldr.AppendItem(b.table)
	bldr.Append(" ")

	first := true
	for _, add_column := range b.add_columns {
		if !first {
			bldr.Append(", ")
		}
		first = false
		bldr.AppendItem(add_column)
	}
	for _, change_column := range b.change_columns {
		if !first {
			bldr.Append(", ")
		}
		first = false
		bldr.AppendItem(change_column)
	}
	for _, drop_column := range b.drop_columns {
		if !first {
			bldr.Append(", ")
		}
		first = false
		bldr.Append("DROP COLUMN ")
		if colname := drop_column.column_name(); len(colname) != 0 {
			bldr.Append(b.dialect.QuoteField(colname))
		} else {
			bldr.AppendItem(drop_column)
		}
	}
	if len(b.rename_to) != 0 {
		if !first {
			bldr.Append(", ")
		}
		bldr.Append("RENAME TO ")
		bldr.Append(b.dialect.QuoteField(b.rename_to))
	}

	return "", nil, nil
}

func (b *cAlterTable) ApplyToTable() error {
	for _, add_column := range b.add_columns {
		err := add_column.applyToTable()
		if err != nil {
			return err
		}
	}
	for _, change_column := range b.change_columns {
		err := change_column.applyToTable()
		if err != nil {
			return err
		}
	}
	for _, drop_column := range b.drop_columns {
		err := b.table.DropColumn(drop_column)
		if err != nil {
			return err
		}
	}
	if len(b.rename_to) != 0 {
		b.table.SetName(b.rename_to)
	}
	return nil
}
