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

type cConditionBetween struct {
	left   serializable
	lower  serializable
	higher serializable
}

func newBetweenCondition(left Column, low, high interface{}) Condition {
	low_literal := toLiteral(low)
	high_literal := toLiteral(high)

	return &cConditionBetween{
		left:   left,
		lower:  low_literal,
		higher: high_literal,
	}
}

func (c *cConditionBetween) serialize(b *builder) {
	b.AppendItem(c.left)
	b.Append(" BETWEEN ")
	b.AppendItem(c.lower)
	b.Append(" AND ")
	b.AppendItem(c.higher)
	return
}

func (c *cConditionBetween) columns() []Column {
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

func (c *cConditionBetween) Describe() (output string) {
	output += fmt.Sprintf("%s BETWEEN (%q, %q)", c.left.Describe(), c.lower, c.higher)
	return
}
