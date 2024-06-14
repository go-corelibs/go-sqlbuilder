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

type cUpdateValue struct {
	col Column
	val literal
}

func newUpdateValue(col Column, val interface{}) cUpdateValue {
	return cUpdateValue{
		col: col,
		val: toLiteral(val),
	}
}

func (c cUpdateValue) serialize(b *builder) {
	if !c.col.acceptType(c.val) {
		b.SetError(newError("%s column not accept %T.",
			c.col.config().Type().String(),
			c.val.Raw(),
		))
		return
	}

	b.Append(b.dialect.QuoteField(c.col.column_name()))
	b.Append("=")
	b.AppendItem(c.val)
}

func (c cUpdateValue) Describe() (output string) {
	// not implemented yet
	return
}
