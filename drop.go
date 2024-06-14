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

// IDropTableStatement is the Buildable interface wrapping of DeleteTable
type IDropTableStatement interface {
	ToSql() (query string, args []interface{}, err error)

	privateDropTable()
}

// DropTableStatement represents a "DROP TABLE" statement.
type DropTableStatement struct {
	table Table

	err error

	dialect Dialect
}

// DropTable returns new "DROP TABLE" statement. The table is Table object to drop.
func DropTable(tbl Table) IDropTableStatement {
	return dropTable(tbl, dialect())
}

func dropTable(tbl Table, d Dialect) *DropTableStatement {
	if d == nil {
		d = dialect()
	}
	if tbl == nil {
		return &DropTableStatement{
			err: newError("table is nil."),
		}
	}
	if _, ok := tbl.(*table); !ok {
		return &DropTableStatement{
			err: newError("table is not natural table."),
		}
	}
	return &DropTableStatement{
		table:   tbl,
		dialect: d,
	}
}

func (b *DropTableStatement) privateDropTable() {
	// nop
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *DropTableStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder(b.dialect)
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()
	if b.err != nil {
		bldr.SetError(b.err)
		return
	}

	bldr.Append("DROP TABLE ")
	bldr.AppendItem(b.table)
	return
}
