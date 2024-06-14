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

// InsertBuilder is the Buildable interface wrapping of Insert
type InsertBuilder interface {
	Columns(columns ...Column) InsertBuilder
	Values(values ...interface{}) InsertBuilder
	Set(column Column, value interface{}) InsertBuilder
	ToSql() (query string, args []interface{}, err error)

	privateInsert()
}

// cInsert represents a INSERT statement.
type cInsert struct {
	columns ColumnList
	values  []literal
	into    Table

	err error

	dialect Dialect
}

// Insert returns new INSERT statement. The table is Table object for into.
func Insert(into Table) InsertBuilder {
	return insert(into, dialect())
}

func insert(into Table, d Dialect) *cInsert {
	if d == nil {
		d = dialect()
	}
	if into == nil {
		return &cInsert{
			err: newError("table is nil."),
		}
	}
	if _, ok := into.(*cTable); !ok {
		return &cInsert{
			err: newError("table is not natural table."),
		}
	}
	return &cInsert{
		into:    into,
		columns: make(ColumnList, 0),
		values:  make([]literal, 0),
		dialect: d,
	}
}

func (b *cInsert) privateInsert() {
	// nop
}

// Columns sets columns for insert.  This overwrites old results of Columns() or Set().
// If not set this, get error on ToSql().
func (b *cInsert) Columns(columns ...Column) InsertBuilder {
	if b.err != nil {
		return b
	}
	for _, col := range columns {
		if !b.into.hasColumn(col) {
			b.err = newError("column not found in table.")
			return b
		}
	}
	b.columns = ColumnList(columns)
	return b
}

// Values sets VALUES clause. This overwrites old results of Values() or Set().
func (b *cInsert) Values(values ...interface{}) InsertBuilder {
	if b.err != nil {
		return b
	}
	sl := make([]literal, len(values))
	for i := range values {
		sl[i] = toLiteral(values[i])
	}
	b.values = sl
	return b
}

// Set sets the column and value together.
// Set cannot be called with Columns() or Values() in a statement.
func (b *cInsert) Set(column Column, value interface{}) InsertBuilder {
	if b.err != nil {
		return b
	}
	if !b.into.hasColumn(column) {
		b.err = newError("column not found in FROM.")
		return b
	}
	b.columns = append(b.columns, column)
	b.values = append(b.values, toLiteral(value))
	return b
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *cInsert) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder(b.dialect)
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()
	if b.err != nil {
		bldr.SetError(b.err)
		return
	}

	// INSERT
	bldr.Append("INSERT")

	// INTO Table
	bldr.Append(" INTO ")
	bldr.AppendItem(b.into)

	// (COLUMN)
	if len(b.columns) == 0 {
		b.columns = b.into.Columns()
	}
	bldr.Append(" ( ")
	bldr.AppendItem(b.columns)
	bldr.Append(" )")

	// VALUES
	if len(b.columns) != len(b.values) {
		bldr.SetError(newError("%d values needed, but got %d.", len(b.columns), len(b.values)))
		return
	}
	for i := range b.columns {
		if !b.columns[i].acceptType(b.values[i]) {
			bldr.SetError(newError("%s column not accept %T.",
				b.columns[i].config().Type().String(),
				b.values[i].Raw()))
			return
		}
	}
	bldr.Append(" VALUES ( ")
	values := make([]serializable, len(b.values))
	for i := range values {
		values[i] = b.values[i]
	}
	bldr.AppendItems(values, ", ")
	bldr.Append(" )")

	return
}
