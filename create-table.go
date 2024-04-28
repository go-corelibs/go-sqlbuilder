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
		dialect: dialect(),
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
