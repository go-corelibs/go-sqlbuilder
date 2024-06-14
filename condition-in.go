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
	"strings"
)

type cConditionIn struct {
	not  bool
	left serializable
	in   []serializable
}

func newInCondition(not bool, left Column, list ...interface{}) Condition {
	m := &cConditionIn{
		not:  not,
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

func (c *cConditionIn) serialize(b *builder) {
	b.AppendItem(c.left)
	if c.not {
		b.Append(" NOT ")
	}
	b.Append(" IN ( ")
	b.AppendItems(c.in, ", ")
	b.Append(" )")
}

func (c *cConditionIn) columns() []Column {
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

func (c *cConditionIn) Describe() (output string) {
	output += fmt.Sprintf("%q", c.left)
	if c.not {
		output += " NOT"
	}
	output += " IN ("
	var parts []string
	for _, in := range c.in {
		parts = append(parts, in.Describe())
	}
	output += strings.Join(parts, ", ")
	output += ")"
	return
}
