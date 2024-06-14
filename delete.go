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

// DeleteBuilder is the Buildable interface wrapping of Delete
type DeleteBuilder interface {
	Where(cond Condition) DeleteBuilder
	ToSql() (query string, args []interface{}, err error)

	privateDelete()
}

// DeleteStatement represents a DELETE statement.
type DeleteStatement struct {
	from  Table
	where Condition

	err error

	dialect Dialect
}

// Delete returns new DELETE statement. The table is Table object to delete from.
func Delete(from Table) DeleteBuilder {
	return deleteFn(from, dialect())
}

func deleteFn(from Table, d Dialect) *DeleteStatement {
	if d == nil {
		d = dialect()
	}
	if from == nil {
		return &DeleteStatement{
			err: newError("from is nil."),
		}
	}
	if _, ok := from.(*table); !ok {
		return &DeleteStatement{
			err: newError("CreateTable can use only natural table."),
		}
	}
	return &DeleteStatement{
		from:    from,
		dialect: d,
	}
}

func (b *DeleteStatement) privateDelete() {
	// nop
}

// Where sets WHERE clause. cond is filter condition.
func (b *DeleteStatement) Where(cond Condition) DeleteBuilder {
	if b.err != nil {
		return b
	}
	for _, col := range cond.columns() {
		if !b.from.hasColumn(col) {
			b.err = newError("column not found in FROM")
			return b
		}
	}
	b.where = cond
	return b
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *DeleteStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder(b.dialect)
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()
	if b.err != nil {
		bldr.SetError(b.err)
		return
	}

	bldr.Append("DELETE FROM ")
	bldr.AppendItem(b.from)

	if b.where != nil {
		bldr.Append(" WHERE ")
		bldr.AppendItem(b.where)
	}
	return
}
