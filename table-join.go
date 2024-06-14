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

type cTableJoin struct {
	join  tableJoinType
	left  Table
	right Table
	on    Condition
}

func (c *cTableJoin) C(name string) Column {
	l_col := c.left.C(name)
	r_col := c.right.C(name)

	_, l_err := l_col.(*cErrorColumn)
	_, r_err := r_col.(*cErrorColumn)

	switch {
	case l_err && r_err:
		return newErrorColumn(newError("column %s was not found.", name))
	case l_err && !r_err:
		return r_col
	case !l_err && r_err:
		return l_col
	default:
		return newErrorColumn(newError("column %s was duplicated.", name))
	}
}

func (c *cTableJoin) Name() string {
	return ""
}

func (c *cTableJoin) LeftName() string {
	if c.left != nil {
		if t, ok := c.left.(*cTableJoin); ok {
			return t.LeftName()
		}
		return c.left.Name()
	}
	return ""
}

func (c *cTableJoin) RightName() string {
	if c.right != nil {
		if t, ok := c.right.(*cTableJoin); ok {
			return t.RightName()
		}
		return c.right.Name()
	}
	return ""
}

func (c *cTableJoin) Columns() []Column {
	return append(c.left.Columns(), c.right.Columns()...)
}

func (c *cTableJoin) Option() *TableOption {
	return nil
}

func (c *cTableJoin) InnerJoin(right Table, on Condition) Table {
	return &cTableJoin{
		left:  c,
		right: right,
		join:  gInnerJoin,
		on:    on,
	}
}

func (c *cTableJoin) LeftOuterJoin(right Table, on Condition) Table {
	return &cTableJoin{
		left:  c,
		right: right,
		join:  gLeftOuterJoin,
		on:    on,
	}
}

func (c *cTableJoin) RightOuterJoin(right Table, on Condition) Table {
	return &cTableJoin{
		left:  c,
		right: right,
		join:  gRightOuterJoin,
		on:    on,
	}
}

func (c *cTableJoin) FullOuterJoin(right Table, on Condition) Table {
	return &cTableJoin{
		left:  c,
		right: right,
		join:  gFullOuterJoin,
		on:    on,
	}
}

func (c *cTableJoin) serialize(b *builder) {
	b.AppendItem(c.left)

	switch t := c.right.(type) {
	case *cTable:

		c.writeJoin(b, c.join, t.name, c.on)

	case *cTableJoin:

		c.writeJoin(b, c.join, t.LeftName(), c.on)
	}

	return
}

func (c *cTableJoin) writeJoin(b *builder, join tableJoinType, other string, cond Condition) {
	switch join {
	case gInnerJoin:
		b.Append(" INNER JOIN ")
	case gLeftOuterJoin:
		b.Append(" LEFT OUTER JOIN ")
	case gRightOuterJoin:
		b.Append(" RIGHT OUTER JOIN ")
	case gFullOuterJoin:
		b.Append(" FULL OUTER JOIN ")
	}
	b.Append(b.dialect.QuoteField(other))
	b.Append(" ON ")
	b.AppendItem(cond)
}

func (c *cTableJoin) hasColumn(trg Column) bool {
	if c.left.hasColumn(trg) {
		return true
	}
	if c.right.hasColumn(trg) {
		return true
	}
	return false
}

func (c *cTableJoin) Describe() (output string) {

	output += c.left.Describe()
	output += "\n\t"

	switch t := c.right.(type) {
	case *cTable:
		output += c.join.String() + " " + t.Name()
	case *cTableJoin:
		output += c.join.String() + " " + t.RightName()
	}

	output += " ON (" + c.on.Describe() + ")"
	return
}
