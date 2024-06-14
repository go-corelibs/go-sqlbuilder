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

// Package sqlbuilder is a Standard Query Language builder for golang supporting
// Sqlite3, MySQL and Postgres
//
// See https://github.com/go-corelibs/go-sqlbuilder for more information
package sqlbuilder

import (
	"bytes"
	"fmt"
)

var _dialect Dialect = nil

// Star reprecents * column.
var Star Column = &columnImpl{nil, nil}

// Statement reprecents a statement(SELECT/INSERT/UPDATE and other)
type Statement interface {
	ToSql() (query string, attrs []interface{}, err error)
}

type serializable interface {
	serialize(b *builder)
}

// Dialect encapsulates behaviors that differ across SQL database.
type Dialect interface {
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

type builder struct {
	query *bytes.Buffer
	args  []interface{}
	err   error

	dialect Dialect
}

func newBuilder(d Dialect) *builder {
	if d == nil {
		d = dialect()
	}
	return &builder{
		query:   bytes.NewBuffer(make([]byte, 0, 256)),
		args:    make([]interface{}, 0, 8),
		err:     nil,
		dialect: d,
	}
}

func (b *builder) Err() error {
	if b.err != nil {
		return b.err
	}
	return nil
}

func (b *builder) Query() string {
	if b.err != nil {
		return ""
	}
	return b.query.String() + b.dialect.QuerySuffix()
}

func (b *builder) Args() []interface{} {
	if b.err != nil {
		return []interface{}{}
	}
	return b.args
}

func (b *builder) SetError(err error) {
	if b.err != nil {
		return
	}
	b.err = err
	return
}

func (b *builder) Append(query string) {
	if b.err != nil {
		return
	}

	b.query.WriteString(query)
}

func (b *builder) AppendValue(val interface{}) {
	if b.err != nil {
		return
	}

	b.query.WriteString(b.dialect.BindVar(len(b.args) + 1))
	b.args = append(b.args, val)
	return
}

func (b *builder) AppendItems(parts []serializable, sep string) {
	if parts == nil {
		return
	}

	first := true
	for _, part := range parts {
		if first {
			first = false
		} else {
			b.Append(sep)
		}
		part.serialize(b)
	}
	return
}

func (b *builder) AppendItem(part serializable) {
	if part == nil {
		return
	}
	part.serialize(b)
}

type errors struct {
	fmt  string
	args []interface{}
}

func newError(fmt string, args ...interface{}) *errors {
	return &errors{
		fmt:  fmt,
		args: args,
	}
}

func (err *errors) Error() string {
	return fmt.Sprintf("sqlbuilder: "+err.fmt, err.args...)
}
