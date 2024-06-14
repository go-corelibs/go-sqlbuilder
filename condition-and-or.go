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
	"strings"
)

type cConditionAndOr struct {
	connector  string
	conditions []Condition
}

// And creates a combined condition with "AND" operator.
func And(conditions ...Condition) Condition {
	return &cConditionAndOr{
		connector:  "AND",
		conditions: conditions,
	}
}

// Or creates a combined condition with "OR" operator.
func Or(conditions ...Condition) Condition {
	return &cConditionAndOr{
		connector:  "OR",
		conditions: conditions,
	}
}

func (c *cConditionAndOr) serialize(b *builder) {
	first := true
	for _, cond := range c.conditions {
		if first {
			first = false
		} else {
			b.Append(" " + c.connector + " ")
		}
		if _, ok := cond.(*cConditionAndOr); ok {
			// if condition is "AND" or "OR"
			b.Append("( ")
			b.AppendItem(cond)
			b.Append(" )")
		} else {
			b.AppendItem(cond)
		}
	}
	return
}

func (c *cConditionAndOr) columns() []Column {
	list := make([]Column, 0)
	for _, cond := range c.conditions {
		list = append(list, cond.columns()...)
	}
	return list
}

func (c *cConditionAndOr) Describe() (output string) {
	var parts []string
	for _, condition := range c.conditions {
		parts = append(parts, condition.Describe())
	}
	return strings.Join(parts, " "+c.connector+" ")
}
