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

func TestUpdate(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)
	table2 := NewTable(
		"TABLE_B",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
	)
	tableJoined := table1.InnerJoin(table2, table1.C("test1").Eq(table2.C("id")))

	var cases = []statementTestCase{{
		stmt: Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20).
			OrderBy(true, table1.C("test1")).
			Limit(1).
			Offset(2),
		query:  `UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=? ORDER BY "TABLE_A"."test1" DESC LIMIT ? OFFSET ?;`,
		args:   []interface{}{int64(10), int64(20), int64(1), 1, 2},
		errmsg: "",
	}, {
		stmt: Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20),
		query:  `UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=?;`,
		args:   []interface{}{int64(10), int64(20), int64(1)},
		errmsg: "",
	}, {
		stmt: Update(nil).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), 10).
			Set(table1.C("test2"), 20),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: table is nil.",
	}, {
		stmt:   Update(table1).Where(table1.C("id").Eq(1)),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: length of sets is 0.",
	}, {
		stmt: Update(table1).Where(table1.C("id").Eq(1)).
			Set(table1.C("test1"), "foo"),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: int column not accept string.",
	}, {
		stmt:   Update(tableJoined),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: length of sets is 0.",
	}}
	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
