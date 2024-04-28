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
		dialect: dialect(),
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
