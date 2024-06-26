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
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
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
	acol_id := table1.C("id").As("tbl1id")

	var cases = []statementTestCase{{
		stmt: Select(table1).
			Columns(table1.C("test1"), table1.C("test2")).
			Where(And(
				table1.C("id").Eq(1),
				table1.C("test1").Eq(2),
			)).
			Distinct().
			OrderBy(false, table1.C("id")).
			GroupBy(table1.C("id")).
			Having(table1.C("id").Eq(1)).
			Limit(10).
			Offset(20),
		query: `SELECT DISTINCT "TABLE_A"."test1", "TABLE_A"."test2" ` +
			`FROM "TABLE_A" WHERE "TABLE_A"."id"=? AND "TABLE_A"."test1"=? ` +
			`GROUP BY "TABLE_A"."id" HAVING "TABLE_A"."id"=? ORDER BY "TABLE_A"."id" ASC ` +
			`LIMIT ? OFFSET ?;`,
		args:   []interface{}{int64(1), int64(2), int64(1), 10, 20},
		errmsg: "",
	}, {
		stmt: Select(table1).
			Columns(table1.C("test1"), table1.C("test2")),
		query:  `SELECT "TABLE_A"."test1", "TABLE_A"."test2" FROM "TABLE_A";`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt: Select(table1).Columns(acol_id).
			Where(acol_id.Eq(1)),
		query:  `SELECT "TABLE_A"."id" AS "tbl1id" FROM "TABLE_A" WHERE "tbl1id"=?;`,
		args:   []interface{}{int64(1)},
		errmsg: "",
	}, {
		stmt: Select(table1).Columns(acol_id).
			Where(acol_id.Eq(1)).
			OrderBy(false, table1.C("test1")).
			OrderBy(true, table1.C("test2")),
		query:  `SELECT "TABLE_A"."id" AS "tbl1id" FROM "TABLE_A" WHERE "tbl1id"=? ORDER BY "TABLE_A"."test1" ASC, "TABLE_A"."test2" DESC;`,
		args:   []interface{}{int64(1)},
		errmsg: "",
	}, {
		stmt:   Select(table1).Columns(Star),
		query:  `SELECT * FROM "TABLE_A";`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt:   Select(table1),
		query:  `SELECT * FROM "TABLE_A";`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt: Select(tableJoined).
			Columns(Star),
		query:  `SELECT * FROM "TABLE_A" INNER JOIN "TABLE_B" ON "TABLE_A"."test1"="TABLE_B"."id";`,
		args:   []interface{}{},
		errmsg: "",
	}, {
		stmt: Select(nil).
			Columns(table1.C("test1"), table1.C("test2")),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: table is nil.",
	}, {
		stmt: Select(table1).
			Columns(table1.C("test1"), table1.C("test2")).
			Having(table1.C("id").Eq(1)),
		query:  ``,
		args:   []interface{}{},
		errmsg: "sqlbuilder: GROUP BY clause is not found.",
	}}

	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}

func TestSubquery(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)

	subquery := Select(table1).Columns(table1.C("id")).ToSubquery("SQ1")
	query, attrs, err := Select(subquery).
		Columns(subquery.C("id")).
		Where(subquery.C("id").Eq(1)).ToSql()

	if `SELECT "SQ1"."id" FROM ( SELECT "TABLE_A"."id" FROM "TABLE_A" ) AS SQ1 WHERE "SQ1"."id"=?;` != query {
		t.Errorf("failed \ngot %s", query)
	}
	if !reflect.DeepEqual([]interface{}{int64(1)}, attrs) {
		t.Errorf("failed \ngot %#v", attrs)
	}
	if err != nil {
		t.Errorf("failed \ngot %#v", err)
	}
}

func BenchmarkSelect(b *testing.B) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)

	for i := 0; i < b.N; i++ {
		Select(table1).
			Columns(table1.C("test1"), table1.C("test2")).
			Where(And(
				table1.C("id").Eq(1),
				table1.C("test1").Eq(2),
			)).
			ToSql()
	}
}
