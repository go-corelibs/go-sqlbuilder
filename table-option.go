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
	"strconv"
)

// TODO: TableOption.ForeignKey map[string]Column
// TODO: convert TableOption to a buildable design pattern

// TableOption represents constraint of a table.
type TableOption struct {
	Unique [][]string
}

// Describe returns a string representation of the TableOption
func (t TableOption) Describe() (output string) {
	if len(t.Unique) == 0 {
		return
	}
	for idx, list := range t.Unique {
		output += ".Unique[" + strconv.Itoa(idx) + "]("
		for jdx, key := range list {
			if jdx > 0 {
				output += ", "
			}
			output += strconv.Quote(key)
		}
		output += ")"
	}
	return
}
