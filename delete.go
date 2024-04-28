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
