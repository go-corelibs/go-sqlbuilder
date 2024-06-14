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

// CreateTableStatement represents a "CREATE TABLE" statement.
type CreateTableStatement struct {
	table       Table
	ifNotExists bool

	err error

	dialect Dialect
}

// CreateTable returns new "CREATE TABLE" statement. The table is Table object to create.
func CreateTable(tbl Table) CreateTableBuilder {
	return createTable(tbl, dialect())
}

func createTable(tbl Table, d Dialect) *CreateTableStatement {
	if d == nil {
		d = dialect()
	}
	if tbl == nil {
		return &CreateTableStatement{
			err: newError("table is nil."),
		}
	}
	if _, ok := tbl.(*table); !ok {
		return &CreateTableStatement{
			err: newError("CreateTable can use only natural table."),
		}
	}

	return &CreateTableStatement{
		table:   tbl,
		dialect: d,
	}
}

func (b *CreateTableStatement) privateCreateTable() {
	// nop
}

// IfNotExists sets "IF NOT EXISTS" clause.
func (b *CreateTableStatement) IfNotExists() CreateTableBuilder {
	if b.err != nil {
		return b
	}
	b.ifNotExists = true
	return b
}

// ToSql generates query string, placeholder arguments, and error.
func (b *CreateTableStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder(b.dialect)
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()
	if b.err != nil {
		bldr.SetError(b.err)
		return
	}

	bldr.Append("CREATE TABLE ")
	if b.ifNotExists {
		bldr.Append("IF NOT EXISTS ")
	}
	bldr.AppendItem(b.table)

	if len(b.table.Columns()) != 0 {
		bldr.Append(" ( ")
		bldr.AppendItem(createTableColumnList(b.table.Columns()))
		bldr.Append(" )")
	} else {
		bldr.SetError(newError("CreateTableStatement needs one or more columns."))
		return
	}

	// table option
	if tabopt, err := b.dialect.TableOptionToString(b.table.Option()); err == nil {
		if len(tabopt) != 0 {
			bldr.Append(" " + tabopt)
		}
	} else {
		bldr.SetError(err)
	}

	return
}

type createTableColumnList []Column

func (m createTableColumnList) serialize(bldr *builder) {
	first := true
	for _, column := range m {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		cc := column.config()

		// Column name
		bldr.AppendItem(cc)
		bldr.Append(" ")

		// SQL data name
		str, err := bldr.dialect.ColumnTypeToString(cc)
		if err != nil {
			bldr.SetError(err)
		}
		bldr.Append(str)

		str, err = bldr.dialect.ColumnOptionToString(cc.Option())
		if err != nil {
			bldr.SetError(err)
		}
		if len(str) != 0 {
			bldr.Append(" " + str)
		}
	}
}
