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

var _dialect Dialect = nil

// Dialect encapsulates behaviors that differ across SQL database.
type Dialect interface {
	Name() string
	QuerySuffix() string
	BindVar(i int) string
	QuoteField(field interface{}) string
	ColumnTypeToString(ColumnConfig) (string, error)
	ColumnOptionToString(*ColumnOption) (string, error)
	TableOptionToString(*TableOption) (string, error)
}

// SetDialect sets dialect for SQL server.
// Must set dialect at first.
func SetDialect(opt Dialect) {
	_dialect = opt
}

func dialect() Dialect {
	if _dialect == nil {
		panic(newError("default dialect is not set. Please call SetDialect() first."))
	}
	return _dialect
}
