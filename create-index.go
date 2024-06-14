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

// CreateIndexBuilder is the Buildable interface wrapping of CreateTable
type CreateIndexBuilder interface {
	IfNotExists() CreateIndexBuilder
	Columns(columns ...Column) CreateIndexBuilder
	Name(name string) CreateIndexBuilder
	ToSql() (query string, args []interface{}, err error)

	privateCreateIndex()
}

// CreateIndexStatement represents a "CREATE INDEX" statement.
type CreateIndexStatement struct {
	table       Table
	columns     []Column
	name        string
	ifNotExists bool

	err error

	dialect Dialect
}

// CreateIndex returns new "CREATE INDEX" statement. The table is Table object to create index.
func CreateIndex(tbl Table) CreateIndexBuilder {
	return createIndex(tbl, dialect())
}

func createIndex(tbl Table, d Dialect) *CreateIndexStatement {
	if d == nil {
		d = dialect()
	}
	if tbl == nil {
		return &CreateIndexStatement{
			err: newError("table is nil."),
		}
	}
	if _, ok := tbl.(*table); !ok {
		return &CreateIndexStatement{
			err: newError("CreateTable can use only natural table."),
		}
	}
	return &CreateIndexStatement{
		table:   tbl,
		dialect: d,
	}
}

func (b *CreateIndexStatement) privateCreateIndex() {
	// nop
}

// IfNotExists sets "IF NOT EXISTS" clause.
func (b *CreateIndexStatement) IfNotExists() CreateIndexBuilder {
	if b.err != nil {
		return b
	}
	b.ifNotExists = true
	return b
}

// IfNotExists sets "IF NOT EXISTS" clause. If not set this, returns error on ToSql().
func (b *CreateIndexStatement) Columns(columns ...Column) CreateIndexBuilder {
	if b.err != nil {
		return b
	}
	b.columns = columns
	return b
}

// Name sets name for index.
// If not set this, auto generated name will be used.
func (b *CreateIndexStatement) Name(name string) CreateIndexBuilder {
	if b.err != nil {
		return b
	}
	b.name = name
	return b
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *CreateIndexStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder(b.dialect)
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()
	if b.err != nil {
		bldr.SetError(b.err)
		return
	}

	bldr.Append("CREATE INDEX ")
	if b.ifNotExists {
		bldr.Append("IF NOT EXISTS ")
	}

	if len(b.name) != 0 {
		bldr.Append(b.dialect.QuoteField(b.name))
	} else {
		bldr.SetError(newError("name was not setted."))
		return
	}

	bldr.Append(" ON ")
	bldr.AppendItem(b.table)

	if len(b.columns) != 0 {
		bldr.Append(" ( ")
		bldr.AppendItem(createIndexColumnList(b.columns))
		bldr.Append(" )")
	} else {
		bldr.SetError(newError("columns was not setted."))
		return
	}
	return
}

type createIndexColumnList []Column

func (m createIndexColumnList) serialize(bldr *builder) {
	first := true
	for _, column := range m {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		cc := column.config()
		bldr.AppendItem(cc)
	}
}
