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

// Condition represents a condition for WHERE clause and other.
type Condition interface {
	serializable

	columns() []Column
}

type connectCondition struct {
	connector string
	conds     []Condition
}

func (c *connectCondition) serialize(bldr *builder) {
	first := true
	for _, cond := range c.conds {
		if first {
			first = false
		} else {
			bldr.Append(" " + c.connector + " ")
		}
		if _, ok := cond.(*connectCondition); ok {
			// if condition is AND or OR
			bldr.Append("( ")
			bldr.AppendItem(cond)
			bldr.Append(" )")
		} else {
			bldr.AppendItem(cond)
		}
	}
	return
}

func (c *connectCondition) columns() []Column {
	list := make([]Column, 0)
	for _, cond := range c.conds {
		list = append(list, cond.columns()...)
	}
	return list
}

// And creates a combined condition with "AND" operator.
func And(conds ...Condition) Condition {
	return &connectCondition{
		connector: "AND",
		conds:     conds,
	}
}

// And creates a combined condition with "OR" operator.
func Or(conds ...Condition) Condition {
	return &connectCondition{
		connector: "OR",
		conds:     conds,
	}
}

type binaryOperationCondition struct {
	left     serializable
	right    serializable
	operator string
	err      error
}

func newBinaryOperationCondition(left, right interface{}, operator string) *binaryOperationCondition {
	cond := &binaryOperationCondition{
		operator: operator,
	}
	column_exist := false
	switch t := left.(type) {
	case Column:
		column_exist = true
		cond.left = t
	case nil:
		cond.err = newError("left-hand side of binary operator is null.")
	default:
		cond.left = toLiteral(t)
	}
	switch t := right.(type) {
	case Column:
		column_exist = true
		cond.right = t
	default:
		cond.right = toLiteral(t)
	}
	if !column_exist {
		cond.err = newError("binary operation is need column.")
	}

	return cond
}

func newBetweenCondition(left Column, low, high interface{}) Condition {
	low_literal := toLiteral(low)
	high_literal := toLiteral(high)

	return &betweenCondition{
		left:   left,
		lower:  low_literal,
		higher: high_literal,
	}
}

func (c *binaryOperationCondition) serialize(bldr *builder) {
	bldr.AppendItem(c.left)

	switch t := c.right.(type) {
	case literal:
		if t.IsNil() {
			switch c.operator {
			case "=":
				bldr.Append(" IS ")
			case "<>":
				bldr.Append(" IS NOT ")
			default:
				bldr.SetError(newError("NULL can not be used with %s operator.", c.operator))
			}
			bldr.Append("NULL")
		} else {
			bldr.Append(c.operator)
			bldr.AppendItem(c.right)
		}
	default:
		bldr.Append(c.operator)
		bldr.AppendItem(c.right)
	}
	return
}

func (c *binaryOperationCondition) columns() []Column {
	list := make([]Column, 0)
	if col, ok := c.left.(Column); ok {
		list = append(list, col)
	}
	if col, ok := c.right.(Column); ok {
		list = append(list, col)
	}
	return list
}

type betweenCondition struct {
	left   serializable
	lower  serializable
	higher serializable
}

func (c *betweenCondition) serialize(bldr *builder) {
	bldr.AppendItem(c.left)
	bldr.Append(" BETWEEN ")
	bldr.AppendItem(c.lower)
	bldr.Append(" AND ")
	bldr.AppendItem(c.higher)
	return
}

func (c *betweenCondition) columns() []Column {
	list := make([]Column, 0)
	if col, ok := c.left.(Column); ok {
		list = append(list, col)
	}
	if col, ok := c.lower.(Column); ok {
		list = append(list, col)
	}
	if col, ok := c.higher.(Column); ok {
		list = append(list, col)
	}
	return list
}

type inCondition struct {
	left serializable
	in   []serializable
}

func newInCondition(left Column, list ...interface{}) Condition {
	m := &inCondition{
		left: left,
		in:   make([]serializable, 0, len(list)),
	}
	for _, item := range list {
		if c, ok := item.(Column); ok {
			m.in = append(m.in, c)
		} else {
			m.in = append(m.in, toLiteral(item))
		}
	}
	return m
}

func (c *inCondition) serialize(bldr *builder) {
	bldr.AppendItem(c.left)
	bldr.Append(" IN ( ")
	bldr.AppendItems(c.in, ", ")
	bldr.Append(" )")
}

func (c *inCondition) columns() []Column {
	list := make([]Column, 0)
	if col, ok := c.left.(Column); ok {
		list = append(list, col)
	}
	for _, in := range c.in {
		if col, ok := in.(Column); ok {
			list = append(list, col)
		}
	}
	return list
}
