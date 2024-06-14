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

type cCreateTableColumnList []Column

func (c cCreateTableColumnList) serialize(b *builder) {
	first := true
	for _, column := range c {
		if first {
			first = false
		} else {
			b.Append(", ")
		}
		cc := column.config()

		// Column name
		b.AppendItem(cc)
		b.Append(" ")

		// SQL data name
		str, err := b.dialect.ColumnTypeToString(cc)
		if err != nil {
			b.SetError(err)
		}
		b.Append(str)

		str, err = b.dialect.ColumnOptionToString(cc.Option())
		if err != nil {
			b.SetError(err)
		}
		if len(str) != 0 {
			b.Append(" " + str)
		}
	}
}

func (c cCreateTableColumnList) Describe() (output string) {
	// not implemented yet
	return
}
