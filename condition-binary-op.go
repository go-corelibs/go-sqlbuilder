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

import (
	"fmt"
)

type cConditionBinaryOp struct {
	left     serializable
	right    serializable
	operator string
	err      error
}

func newBinaryOperationCondition(left, right interface{}, operator string) *cConditionBinaryOp {
	cond := &cConditionBinaryOp{
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

func (c *cConditionBinaryOp) serialize(b *builder) {
	b.AppendItem(c.left)

	switch t := c.right.(type) {
	case literal:
		if t.IsNil() {
			switch c.operator {
			case "=":
				b.Append(" IS ")
			case "<>":
				b.Append(" IS NOT ")
			default:
				b.SetError(newError("NULL can not be used with %s operator.", c.operator))
			}
			b.Append("NULL")
		} else {
			b.Append(c.operator)
			b.AppendItem(c.right)
		}
	default:
		b.Append(c.operator)
		b.AppendItem(c.right)
	}
	return
}

func (c *cConditionBinaryOp) columns() []Column {
	list := make([]Column, 0)
	if col, ok := c.left.(Column); ok {
		list = append(list, col)
	}
	if col, ok := c.right.(Column); ok {
		list = append(list, col)
	}
	return list
}

func (c *cConditionBinaryOp) Describe() (output string) {
	output += fmt.Sprintf("%v %v %v", c.left.Describe(), c.operator, c.right.Describe())
	return
}
