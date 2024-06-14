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
	"sync"
)

// Buildable is an interface for performing the package-level functions
// which all require a package-level Dialect to be configured using a
// call to SetDialect. Buildable allows for using one or more independent
// Dialects at the same time within the same codebase
//
// Buildable is concurrency-safe
type Buildable interface {
	// Dialect returns the Dialect associated with this Buildable instance
	Dialect() Dialect

	// NewTable wraps the NewTable package-level function and for all named
	// tables, stores the Table reference within the Buildable instance for
	// later retrieval with the T method
	NewTable(name string, option *TableOption, column_configs ...ColumnConfig) (t Table)

	// AlterTable starts a new ALTER TABLE statement builder
	AlterTable(tbl Table) AlterTableBuilder

	// CreateTable starts a new CREATE TABLE statement builder
	CreateTable(tbl Table) CreateTableBuilder

	// CreateIndex starts a new ADD INDEX statement builder
	CreateIndex(tbl Table) CreateIndexBuilder

	// Delete starts a new DELETE statement builder
	Delete(from Table) DeleteBuilder

	// Insert starts a new INSERT statement builder
	Insert(into Table) InsertBuilder

	// Select starts a new SELECT statement builder
	Select(from Table) SelectBuilder
}

type buildable struct {
	dialect Dialect
	//tables  map[string]Table
	//order   []string

	m *sync.RWMutex
}

// NewBuildable constructs a new Buildable instance with the given Dialect
func NewBuildable(d Dialect) Buildable {
	b := &buildable{
		dialect: d,
		//tables:  make(map[string]Table),
		m: &sync.RWMutex{},
	}
	return b
}

func (b *buildable) NewTable(name string, option *TableOption, column_configs ...ColumnConfig) (t Table) {
	t = NewTable(name, option, column_configs...)
	//b.trackTable(t)
	return
}

func (b *buildable) Dialect() Dialect {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.dialect
}

func (b *buildable) AlterTable(tbl Table) AlterTableBuilder {
	return alterTable(tbl, b.Dialect())
}

func (b *buildable) CreateTable(tbl Table) CreateTableBuilder {
	return createTable(tbl, b.Dialect())
}

func (b *buildable) CreateIndex(tbl Table) CreateIndexBuilder {
	return createIndex(tbl, b.Dialect())
}

func (b *buildable) Delete(from Table) DeleteBuilder {
	return deleteFn(from, b.Dialect())
}

func (b *buildable) Insert(into Table) InsertBuilder {
	return insert(into, b.Dialect())
}

func (b *buildable) Select(from Table) SelectBuilder {
	return selectFn(from, b.Dialect())
}
