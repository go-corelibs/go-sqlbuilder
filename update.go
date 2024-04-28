package sqlbuilder

// IUpdateStatement is the Buildable interface wrapping of Update
type IUpdateStatement interface {
	Set(col Column, val interface{}) IUpdateStatement
	Where(cond Condition) IUpdateStatement
	Limit(limit int) IUpdateStatement
	Offset(offset int) IUpdateStatement
	OrderBy(desc bool, columns ...Column) IUpdateStatement
	ToSql() (query string, args []interface{}, err error)

	privateUpdate()
}

// UpdateStatement represents a UPDATE statement.
type UpdateStatement struct {
	table   Table
	set     []serializable
	where   Condition
	orderBy []serializable
	limit   int
	offset  int

	err error

	dialect Dialect
}

// Update returns new UPDATE statement. The table is Table object to update.
func Update(tbl Table) IUpdateStatement {
	return update(tbl, dialect())
}

func update(tbl Table, d Dialect) *UpdateStatement {
	if d == nil {
		d = dialect()
	}
	if tbl == nil {
		return &UpdateStatement{
			err: newError("table is nil."),
		}
	}
	return &UpdateStatement{
		table:   tbl,
		set:     make([]serializable, 0),
		dialect: d,
	}
}

func (b *UpdateStatement) privateUpdate() {
	// nop
}

// Set sets SETS clause like col=val.  Call many time for update multi columns.
func (b *UpdateStatement) Set(col Column, val interface{}) IUpdateStatement {
	if b.err != nil {
		return b
	}
	if !b.table.hasColumn(col) {
		b.err = newError("column not found in FROM.")
		return b
	}
	b.set = append(b.set, newUpdateValue(col, val))
	return b
}

// Where sets WHERE clause.  The cond is filter condition.
func (b *UpdateStatement) Where(cond Condition) IUpdateStatement {
	if b.err != nil {
		return b
	}
	b.where = cond
	return b
}

// Limit sets LIMIT clause.
func (b *UpdateStatement) Limit(limit int) IUpdateStatement {
	if b.err != nil {
		return b
	}
	b.limit = limit
	return b
}

// Limit sets OFFSET clause.
func (b *UpdateStatement) Offset(offset int) IUpdateStatement {
	if b.err != nil {
		return b
	}
	b.offset = offset
	return b
}

// OrderBy sets "ORDER BY" clause. Use descending order if the desc is true, by the columns.
func (b *UpdateStatement) OrderBy(desc bool, columns ...Column) IUpdateStatement {
	if b.err != nil {
		return b
	}
	if b.orderBy == nil {
		b.orderBy = make([]serializable, 0)
	}

	for _, c := range columns {
		b.orderBy = append(b.orderBy, newOrderBy(desc, c))
	}
	return b
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *UpdateStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder(b.dialect)
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()
	if b.err != nil {
		bldr.SetError(b.err)
		return
	}

	// UPDATE TABLE SET (COLUMN=VALUE)
	bldr.Append("UPDATE ")
	bldr.AppendItem(b.table)

	bldr.Append(" SET ")
	if len(b.set) != 0 {
		bldr.AppendItems(b.set, ", ")
	} else {
		bldr.SetError(newError("length of sets is 0."))
	}

	// WHERE
	if b.where != nil {
		bldr.Append(" WHERE ")
		bldr.AppendItem(b.where)
	}

	// ORDER BY
	if b.orderBy != nil {
		bldr.Append(" ORDER BY ")
		bldr.AppendItems(b.orderBy, ", ")
	}

	// LIMIT
	if b.limit != 0 {
		bldr.Append(" LIMIT ")
		bldr.AppendValue(b.limit)
	}

	// Offset
	if b.offset != 0 {
		bldr.Append(" OFFSET ")
		bldr.AppendValue(b.offset)
	}
	return
}

type updateValue struct {
	col Column
	val literal
}

func newUpdateValue(col Column, val interface{}) updateValue {
	return updateValue{
		col: col,
		val: toLiteral(val),
	}
}

func (m updateValue) serialize(bldr *builder) {
	if !m.col.acceptType(m.val) {
		bldr.SetError(newError("%s column not accept %T.",
			m.col.config().Type().String(),
			m.val.Raw(),
		))
		return
	}

	bldr.Append(bldr.dialect.QuoteField(m.col.column_name()))
	bldr.Append("=")
	bldr.AppendItem(m.val)
}
