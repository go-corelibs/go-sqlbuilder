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
	"testing"
)

func TestCreate(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey:    true,
			AutoIncrement: true,
		}),
		IntColumn("test1", &ColumnOption{
			Unique: true,
		}),
		StringColumn("test2", &ColumnOption{
			Size: 255,
		}),
	)
	table2 := NewTable(
		"TABLE_B",
		&TableOption{},
		StringColumn("id", &ColumnOption{
			PrimaryKey:    true,
			AutoIncrement: true,
			SqlType:       "VARCHAR(255)",
		}),
		AnyColumn("test1", &ColumnOption{
			Unique:  true,
			SqlType: "INTEGER",
		}),
	)
	table3 := NewTable(
		"TABLE_C",
		&TableOption{
			Unique: [][]string{{"test1", "test2"}},
		},
		IntColumn("id", &ColumnOption{
			PrimaryKey:    true,
			AutoIncrement: true,
		}),
		IntColumn("test1", &ColumnOption{
			Unique: true,
		}),
		StringColumn("test2", &ColumnOption{
			Size: 255,
		}),
	)
	tableJoined := table1.InnerJoin(table2, table1.C("test1").Eq(table2.C("id")))
	tableZeroColumns := &cTable{
		name:    "ZERO_TABLE",
		columns: make([]Column, 0),
	}

	var cases = []statementTestCase{{
		stmt:   CreateTable(table1).IfNotExists(),
		query:  `CREATE TABLE IF NOT EXISTS "TABLE_A" ( "id" INTEGER PRIMARY KEY AUTOINCREMENT, "test1" INTEGER UNIQUE, "test2" TEXT );`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   CreateTable(table2).IfNotExists(),
		query:  `CREATE TABLE IF NOT EXISTS "TABLE_B" ( "id" VARCHAR(255) PRIMARY KEY AUTOINCREMENT, "test1" INTEGER UNIQUE );`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   CreateTable(table3).IfNotExists(),
		query:  `CREATE TABLE IF NOT EXISTS "TABLE_C" ( "id" INTEGER PRIMARY KEY AUTOINCREMENT, "test1" INTEGER UNIQUE, "test2" TEXT, UNIQUE("test1", "test2") );`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   CreateIndex(table1).Name("I_TABLE_A").IfNotExists().Columns(table1.C("test1"), table1.C("test2")),
		query:  `CREATE INDEX IF NOT EXISTS "I_TABLE_A" ON "TABLE_A" ( "test1", "test2" );`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   CreateTable(tableZeroColumns),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: CreateTable needs one or more columns.",
	}, {
		stmt:   CreateTable(nil),
		query:  ``,
		args:   []interface{}(nil),
		errmsg: "sqlbuilder: table is nil.",
	}, {
		stmt:   CreateTable(tableJoined),
		query:  ``,
		args:   []interface{}(nil),
		errmsg: "sqlbuilder: CreateTable can use only natural table.",
	}, {
		stmt:   CreateIndex(table1).Columns(table1.C("test1"), table1.C("test2")),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: name was not set.",
	}, {
		stmt:   CreateIndex(table1).Name("I_TABLE_A"),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: columns was not set.",
	}}

	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
