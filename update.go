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

// UpdateBuilder is the Buildable interface wrapping of Update
type UpdateBuilder interface {
	Set(col Column, val interface{}) UpdateBuilder
	Where(cond Condition) UpdateBuilder
	Limit(limit int) UpdateBuilder
	Offset(offset int) UpdateBuilder
	OrderBy(desc bool, columns ...Column) UpdateBuilder
	ToSql() (query string, args []interface{}, err error)

	privateUpdate()
}

// cUpdate represents a UPDATE statement.
type cUpdate struct {
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
func Update(tbl Table) UpdateBuilder {
	return update(tbl, dialect())
}

func update(tbl Table, d Dialect) *cUpdate {
	if d == nil {
		d = dialect()
	}
	if tbl == nil {
		return &cUpdate{
			err: newError("table is nil."),
		}
	}
	return &cUpdate{
		table:   tbl,
		set:     make([]serializable, 0),
		dialect: d,
	}
}

func (c *cUpdate) privateUpdate() {
	// nop
}

// Set sets SETS clause like col=val.  Call many time for update multi columns.
func (c *cUpdate) Set(col Column, val interface{}) UpdateBuilder {
	if c.err != nil {
		return c
	}
	if !c.table.hasColumn(col) {
		c.err = newError("column not found in FROM.")
		return c
	}
	c.set = append(c.set, newUpdateValue(col, val))
	return c
}

// Where sets WHERE clause.  The cond is filter condition.
func (c *cUpdate) Where(cond Condition) UpdateBuilder {
	if c.err != nil {
		return c
	}
	c.where = cond
	return c
}

// Limit sets LIMIT clause
func (c *cUpdate) Limit(limit int) UpdateBuilder {
	if c.err != nil {
		return c
	}
	c.limit = limit
	return c
}

// Offset sets OFFSET clause
func (c *cUpdate) Offset(offset int) UpdateBuilder {
	if c.err != nil {
		return c
	}
	c.offset = offset
	return c
}

// OrderBy sets "ORDER BY" clause. Use descending order if the desc is true, by the columns
func (c *cUpdate) OrderBy(desc bool, columns ...Column) UpdateBuilder {
	if c.err != nil {
		return c
	}
	if c.orderBy == nil {
		c.orderBy = make([]serializable, 0)
	}

	for _, column := range columns {
		c.orderBy = append(c.orderBy, newOrderBy(desc, column))
	}
	return c
}

// ToSql generates query string, placeholder arguments, and returns err on errors
func (c *cUpdate) ToSql() (query string, args []interface{}, err error) {
	b := newBuilder(c.dialect)
	defer func() {
		query, args, err = b.Query(), b.Args(), b.Err()
	}()
	if c.err != nil {
		b.SetError(c.err)
		return
	}

	// UPDATE TABLE SET (COLUMN=VALUE)
	b.Append("UPDATE ")
	b.AppendItem(c.table)

	b.Append(" SET ")
	if len(c.set) != 0 {
		b.AppendItems(c.set, ", ")
	} else {
		b.SetError(newError("length of sets is 0."))
	}

	// WHERE
	if c.where != nil {
		b.Append(" WHERE ")
		b.AppendItem(c.where)
	}

	// ORDER BY
	if c.orderBy != nil {
		b.Append(" ORDER BY ")
		b.AppendItems(c.orderBy, ", ")
	}

	// LIMIT
	if c.limit != 0 {
		b.Append(" LIMIT ")
		b.AppendValue(c.limit)
	}

	// Offset
	if c.offset != 0 {
		b.Append(" OFFSET ")
		b.AppendValue(c.offset)
	}
	return
}
