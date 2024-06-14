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

// CreateTableBuilder is the Buildable interface wrapping of CreateIndex
type CreateTableBuilder interface {
	IfNotExists() CreateTableBuilder
	ToSql() (query string, args []interface{}, err error)

	privateCreateTable()
}

// cCreateTable represents a "CREATE TABLE" statement.
type cCreateTable struct {
	table       Table
	ifNotExists bool

	err error

	dialect Dialect
}

// CreateTable returns new "CREATE TABLE" statement. The table is Table object to create.
func CreateTable(tbl Table) CreateTableBuilder {
	return createTable(tbl, dialect())
}

func createTable(tbl Table, d Dialect) *cCreateTable {
	if d == nil {
		d = dialect()
	}
	if tbl == nil {
		return &cCreateTable{
			err: newError("table is nil."),
		}
	}
	if _, ok := tbl.(*cTable); !ok {
		return &cCreateTable{
			err: newError("CreateTable can use only natural table."),
		}
	}

	return &cCreateTable{
		table:   tbl,
		dialect: d,
	}
}

func (c *cCreateTable) privateCreateTable() {
	// nop
}

// IfNotExists sets "IF NOT EXISTS" clause.
func (c *cCreateTable) IfNotExists() CreateTableBuilder {
	if c.err != nil {
		return c
	}
	c.ifNotExists = true
	return c
}

// ToSql generates query string, placeholder arguments, and error.
func (c *cCreateTable) ToSql() (query string, args []interface{}, err error) {
	if c.err != nil {
		err = c.err
		return
	}
	b := newBuilder(c.dialect)
	defer func() {
		query, args, err = b.Query(), b.Args(), b.Err()
	}()
	if c.err != nil {
		b.SetError(c.err)
		return
	}

	if len(c.table.Columns()) == 0 {
		b.SetError(newError("CreateTable needs one or more columns."))
		return
	}

	b.Append("CREATE TABLE ")
	if c.ifNotExists {
		b.Append("IF NOT EXISTS ")
	}
	b.AppendItem(c.table)

	b.Append(" ( ")

	b.AppendItem(cCreateTableColumnList(c.table.Columns()))

	// table option
	if tabopt, err := c.dialect.TableOptionToString(c.table.Option()); err == nil {
		if len(tabopt) != 0 {
			b.Append(", " + tabopt)
		}
	} else {
		b.SetError(err)
	}

	b.Append(" )")
	return
}
