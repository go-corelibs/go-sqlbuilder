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
	"strings"
)

// Table represents a table.
type Table interface {
	serializable
	hasColumn(column Column) bool

	// C returns table's column by the name.
	C(name string) Column

	// Name returns table' name.
	// returns empty if it is joined table or subquery.
	Name() string

	// Option returns table's option(table constraint).
	// returns nil if it is joined table or subquery.
	Option() *TableOption

	// Columns returns all columns.
	Columns() []Column

	// InnerJoin returns a joined table use with "INNER JOIN" clause.
	// The joined table can be handled in same way as single table.
	InnerJoin(Table, Condition) Table

	// LeftOuterJoin returns a joined table use with "LEFT OUTER JOIN" clause.
	// The joined table can be handled in same way as single table.
	LeftOuterJoin(Table, Condition) Table

	// RightOuterJoin returns a joined table use with "RIGHT OUTER JOIN" clause.
	// The joined table can be handled in same way as single table.
	RightOuterJoin(Table, Condition) Table

	// FullOuterJoin returns a joined table use with "FULL OUTER JOIN" clause.
	// The joined table can be handled in same way as single table.
	FullOuterJoin(Table, Condition) Table

	// Describe returns a string representation of the complete structure
	Describe() (output string)
}

type cTable struct {
	name    string
	option  *TableOption
	columns []Column
}

// NewTable returns a new table named by the name.  Specify table columns by the column_config.
// Panic if column is empty.
func NewTable(name string, option *TableOption, column_configs ...ColumnConfig) Table {
	if len(column_configs) == 0 {
		panic(newError("column is needed."))
	}
	if option == nil {
		option = &TableOption{}
	}

	t := &cTable{
		name:    name,
		option:  option,
		columns: make([]Column, 0, len(column_configs)),
	}

	for _, column_config := range column_configs {
		err := t.AddColumnLast(column_config)
		if err != nil {
			panic(err)
		}
	}

	return t
}

func (m *cTable) serialize(b *builder) {
	b.Append(b.dialect.QuoteField(m.name))
	return
}

func (m *cTable) C(name string) Column {
	for _, column := range m.columns {
		if column.column_name() == name {
			return column
		}
	}

	return newErrorColumn(newError("column %s.%s was not found.", m.name, name))
}

func (m *cTable) Name() string {
	return m.name
}

func (m *cTable) SetName(name string) {
	m.name = name
}

func (m *cTable) Columns() []Column {
	return m.columns
}

func (m *cTable) Option() *TableOption {
	return m.option
}

func (m *cTable) AddColumnLast(cc ColumnConfig) error {
	return m.addColumn(cc, len(m.columns))
}

func (m *cTable) AddColumnFirst(cc ColumnConfig) error {
	return m.addColumn(cc, 0)
}

func (m *cTable) AddColumnAfter(cc ColumnConfig, after Column) error {
	for i := range m.columns {
		if m.columns[i] == after {
			return m.addColumn(cc, i+1)
		}
	}
	return newError("column not found.")
}

func (m *cTable) ChangeColumn(trg Column, cc ColumnConfig) error {
	for i := range m.columns {
		if m.columns[i] == trg {
			err := m.dropColumn(i)
			if err != nil {
				return err
			}
			err = m.addColumn(cc, i)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return newError("column not found.")
}

func (m *cTable) ChangeColumnFirst(trg Column, cc ColumnConfig) error {
	for i := range m.columns {
		if m.columns[i] == trg {
			err := m.dropColumn(i)
			if err != nil {
				return err
			}
			err = m.addColumn(cc, 0)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return newError("column not found.")
}

func (m *cTable) ChangeColumnAfter(trg Column, cc ColumnConfig, after Column) error {
	backup := make([]Column, len(m.columns))
	copy(backup, m.columns)
	found := false
	for i := range m.columns {
		if m.columns[i] == trg {
			err := m.dropColumn(i)
			if err != nil {
				m.columns = backup
				return err
			}
			found = true
			break
		}
	}
	if !found {
		return newError("column not found.")
	}
	for i := range m.columns {
		if m.columns[i] == after {
			err := m.addColumn(cc, i+1)
			if err != nil {
				m.columns = backup
				return err
			}
			return nil
		}
	}
	m.columns = backup
	return newError("column not found.")
}

func (m *cTable) addColumn(cc ColumnConfig, pos int) error {
	if len(m.columns) < pos || pos < 0 {
		return newError("Invalid position.")
	}

	var (
		u = make([]Column, pos)
		p = make([]Column, len(m.columns)-pos)
	)
	copy(u, m.columns[:pos])
	copy(p, m.columns[pos:])
	c := cc.toColumn(m)
	m.columns = append(u, c)
	m.columns = append(m.columns, p...)
	return nil
}

func (m *cTable) DropColumn(col Column) error {
	for i := range m.columns {
		if m.columns[i] == col {
			return m.dropColumn(i)
		}
	}
	return newError("column not found.")
}

func (m *cTable) dropColumn(pos int) error {
	if len(m.columns) < pos || pos < 0 {
		return newError("Invalid position.")
	}
	var (
		u = make([]Column, pos)
		p = make([]Column, len(m.columns)-pos-1)
	)
	copy(u, m.columns[:pos])
	if len(m.columns) > pos+1 {
		copy(p, m.columns[pos+1:])
	}
	m.columns = append(u, p...)
	return nil
}

func (m *cTable) InnerJoin(right Table, on Condition) Table {
	return &cTableJoin{
		left:  m,
		right: right,
		join:  gInnerJoin,
		on:    on,
	}
}

func (m *cTable) LeftOuterJoin(right Table, on Condition) Table {
	return &cTableJoin{
		left:  m,
		right: right,
		join:  gLeftOuterJoin,
		on:    on,
	}
}

func (m *cTable) RightOuterJoin(right Table, on Condition) Table {
	return &cTableJoin{
		left:  m,
		right: right,
		join:  gRightOuterJoin,
		on:    on,
	}
}

func (m *cTable) FullOuterJoin(right Table, on Condition) Table {
	return &cTableJoin{
		left:  m,
		right: right,
		join:  gFullOuterJoin,
		on:    on,
	}
}

func (m *cTable) hasColumn(target Column) bool {
	if cimpl, ok := target.(*cColumnImpl); ok {
		return cimpl.hasColumn(m)
	}
	if acol, ok := target.(*cColumnAlias); ok {
		return acol.hasColumn(m)
	}
	if sqlfn, ok := target.(*cSqlFunc); ok {
		return sqlfn.hasColumn(m)
	}
	return false
}

func (m *cTable) Describe() (output string) {
	output += m.name
	if v := m.option.Describe(); v != "" {
		output += "\n\t" + v
	}
	if len(m.columns) > 0 {
		var cols []string
		for _, c := range m.columns {
			cols = append(cols, c.config().Describe())
		}
		output += "\n\t" + strings.Join(cols, "\n\t")
	}
	return
}
