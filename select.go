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

// SelectBuilder is the Buildable interface wrapping of Select
type SelectBuilder interface {
	// Distinct adds a DISTINCT clause to the Columns selection
	Distinct() SelectBuilder
	// Columns specifies one or more Columns to return
	Columns(columns ...Column) SelectBuilder

	Where(cond Condition) SelectBuilder
	Having(cond Condition) SelectBuilder

	GroupBy(columns ...Column) SelectBuilder
	OrderBy(desc bool, columns ...Column) SelectBuilder
	Limit(limit int) SelectBuilder
	Offset(offset int) SelectBuilder

	ToSql() (query string, args []interface{}, err error)
	ToSubquery(alias string) Table

	// Describe returns a description of the configured Table instance, useful
	// when unit testing and needing to confirm that the correct values were
	// given during Table construction
	Describe() (output string)

	serialize(b *builder)
	privateSelect()
}

// cSelect represents a SELECT statement.
type cSelect struct {
	columns  cSelectColumnList
	from     Table
	where    Condition
	distinct bool
	groupBy  []serializable
	orderBy  []serializable
	limit    int
	offset   int
	having   Condition

	err error

	dialect Dialect
}

// Select returns new SELECT statement with from as FROM clause.
func Select(from Table) SelectBuilder {
	return selectFn(from, dialect())
}

func selectFn(from Table, d Dialect) *cSelect {
	if d == nil {
		d = dialect()
	}
	if from == nil {
		return &cSelect{
			err: newError("table is nil."),
		}
	}
	return &cSelect{
		from:    from,
		dialect: d,
	}
}

func (s *cSelect) privateSelect() {
	// nop
}

// Columns set columns for select.
// Get all columns (use *) if it is not set.
func (s *cSelect) Columns(columns ...Column) SelectBuilder {
	if s.err != nil {
		return s
	}
	for _, col := range columns {
		if !s.from.hasColumn(col) {
			s.err = newError("column not found in FROM: %q", col.column_name())
			return s
		}
	}

	s.columns = columns
	return s
}

// Where sets WHERE clause.  The cond is filter condition.
func (s *cSelect) Where(cond Condition) SelectBuilder {
	if s.err != nil {
		return s
	}
	for _, col := range cond.columns() {
		if !s.from.hasColumn(col) {
			s.err = newError("column not found in FROM: %q", col.column_name())
			return s
		}
	}

	s.where = cond
	return s
}

// Distinct sets DISTINCT clause.
func (s *cSelect) Distinct() SelectBuilder {
	if s.err != nil {
		return s
	}
	s.distinct = true
	return s
}

// GroupBy sets "GROUP BY" clause by the columns.
func (s *cSelect) GroupBy(columns ...Column) SelectBuilder {
	if s.err != nil {
		return s
	}
	ex_column := make([]serializable, len(columns))
	for i := range columns {
		ex_column[i] = columns[i]
	}
	s.groupBy = ex_column
	return s
}

// Having sets "HAVING" clause with the cond.
func (s *cSelect) Having(cond Condition) SelectBuilder {
	if s.err != nil {
		return s
	}
	s.having = cond
	return s
}

// OrderBy sets "ORDER BY" clause. Use descending order if the desc is true, by the columns.
func (s *cSelect) OrderBy(desc bool, columns ...Column) SelectBuilder {
	if s.err != nil {
		return s
	}
	if s.orderBy == nil {
		s.orderBy = make([]serializable, 0)
	}

	for _, c := range columns {
		s.orderBy = append(s.orderBy, newOrderBy(desc, c))
	}
	return s
}

// Limit sets LIMIT clause.
func (s *cSelect) Limit(limit int) SelectBuilder {
	if s.err != nil {
		return s
	}
	s.limit = limit
	return s
}

// Offset sets OFFSET clause.
func (s *cSelect) Offset(offset int) SelectBuilder {
	if s.err != nil {
		return s
	}
	s.offset = offset
	return s
}

func (s *cSelect) serialize(b *builder) {
	if s.err != nil {
		b.SetError(s.err)
		return
	}

	// SELECT COLUMN
	b.Append("SELECT ")
	if s.distinct {
		b.Append("DISTINCT ")
	}
	b.AppendItem(s.columns)

	// FROM
	b.Append(" FROM ")
	b.AppendItem(s.from)

	// WHERE
	if s.where != nil {
		b.Append(" WHERE ")
		b.AppendItem(s.where)
	}

	// GROUP BY
	if s.groupBy != nil {
		b.Append(" GROUP BY ")
		b.AppendItems(s.groupBy, ",")
	}

	// HAVING
	if s.having != nil {
		if s.groupBy == nil {
			b.SetError(newError("GROUP BY clause is not found."))
		}
		b.Append(" HAVING ")
		b.AppendItem(s.having)
	}

	// ORDER BY
	if s.orderBy != nil {
		b.Append(" ORDER BY ")
		b.AppendItems(s.orderBy, ", ")
	}

	// LIMIT
	if s.limit != 0 {
		b.Append(" LIMIT ")
		b.AppendValue(s.limit)
	}

	// Offset
	if s.offset != 0 {
		b.Append(" OFFSET ")
		b.AppendValue(s.offset)
	}
	return
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (s *cSelect) ToSql() (query string, args []interface{}, err error) {
	b := newBuilder(s.dialect)
	b.AppendItem(s)
	return b.Query(), b.Args(), b.Err()
}

func (s *cSelect) ToSubquery(alias string) Table {
	return newSubQuery(s, alias)
}

func (s *cSelect) Describe() (output string) {
	output = s.from.Describe()
	return
}
