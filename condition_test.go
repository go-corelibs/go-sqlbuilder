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

func TestBinaryCondition(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)
	var cases = []conditionTestCase{
		{
			cond:   table1.C("id").Eq(table1.C("test1")),
			query:  `"TABLE_A"."id"="TABLE_A"."test1"`,
			args:   []interface{}{},
			errmsg: "",
		}, {
			cond:   table1.C("id").Eq(1),
			query:  `"TABLE_A"."id"=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").NotEq(1),
			query:  `"TABLE_A"."id"<>?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").Gt(1),
			query:  `"TABLE_A"."id">?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").GtEq(1),
			query:  `"TABLE_A"."id">=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").Lt(1),
			query:  `"TABLE_A"."id"<?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").LtEq(1),
			query:  `"TABLE_A"."id"<=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   table1.C("id").Like("hoge"),
			query:  `"TABLE_A"."id" LIKE ?`,
			args:   []interface{}{"hoge"},
			errmsg: "",
		}, {
			cond:   table1.C("id").Between(1, 2),
			query:  `"TABLE_A"."id" BETWEEN ? AND ?`,
			args:   []interface{}{int64(1), int64(2)},
			errmsg: "",
		}, {
			cond:   table1.C("id").In(1, 2),
			query:  `"TABLE_A"."id" IN ( ?, ? )`,
			args:   []interface{}{int64(1), int64(2)},
			errmsg: "",
		}, {
			cond:   table1.C("id").Eq(nil),
			query:  `"TABLE_A"."id" IS NULL`,
			args:   []interface{}{},
			errmsg: "",
		}, {
			cond:   table1.C("id").NotEq([]byte(nil)),
			query:  `"TABLE_A"."id" IS NOT NULL`,
			args:   []interface{}{},
			errmsg: "",
		}, {
			cond:   table1.C("id").Gt([]byte(nil)),
			query:  `"TABLE_A"."id"`,
			args:   []interface{}{},
			errmsg: "sqlbuilder: NULL can not be used with > operator.",
		}, {
			// case for fail
			cond:   table1.C("id").In(NewTable("DUMMY TABLE", &TableOption{}, StringColumn("id", nil))),
			query:  `"TABLE_A"."id" IN ( `,
			args:   []interface{}{},
			errmsg: "sqlbuilder: got sqlbuilder.cTable type, but literal is not supporting this.",
		},
	}

	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}

func TestBinaryConditionForSqlFunctions(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)
	var cases = []conditionTestCase{
		{
			cond:   Func("count", table1.C("id")).Eq(table1.C("test1")),
			query:  `count("TABLE_A"."id")="TABLE_A"."test1"`,
			args:   []interface{}{},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).Eq(1),
			query:  `count("TABLE_A"."id")=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).NotEq(1),
			query:  `count("TABLE_A"."id")<>?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).Gt(1),
			query:  `count("TABLE_A"."id")>?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).GtEq(1),
			query:  `count("TABLE_A"."id")>=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).Lt(1),
			query:  `count("TABLE_A"."id")<?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).LtEq(1),
			query:  `count("TABLE_A"."id")<=?`,
			args:   []interface{}{int64(1)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).Like("hoge"),
			query:  `count("TABLE_A"."id") LIKE ?`,
			args:   []interface{}{"hoge"},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).Between(1, 2),
			query:  `count("TABLE_A"."id") BETWEEN ? AND ?`,
			args:   []interface{}{int64(1), int64(2)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).In(1, 2),
			query:  `count("TABLE_A"."id") IN ( ?, ? )`,
			args:   []interface{}{int64(1), int64(2)},
			errmsg: "",
		}, {
			cond:   Func("count", table1.C("id")).In(NewTable("DUMMY TABLE", &TableOption{}, StringColumn("id", nil))),
			query:  `count("TABLE_A"."id") IN ( `,
			args:   []interface{}{},
			errmsg: "sqlbuilder: got sqlbuilder.cTable type, but literal is not supporting this.",
		},
	}

	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}

}

func TestConnectCondition(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)
	cases := []conditionTestCase{{
		cond: And(
			table1.C("id").Eq(table1.C("test1")),
			table1.C("id").Eq(1),
			table1.C("id").Eq(2),
		),
		query:  `"TABLE_A"."id"="TABLE_A"."test1" AND "TABLE_A"."id"=? AND "TABLE_A"."id"=?`,
		args:   []interface{}{int64(1), int64(2)},
		errmsg: "",
	}, {
		cond: Or(
			table1.C("id").Eq(table1.C("test1")),
			table1.C("id").Eq(1),
		),
		query:  `"TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=?`,
		args:   []interface{}{int64(1)},
		errmsg: "",
	}, {
		cond: And(
			Or(
				table1.C("id").Eq(table1.C("test1")),
				table1.C("id").Eq(1),
			),
			Or(
				table1.C("id").Eq(table1.C("test1")),
				table1.C("id").Eq(1),
			),
		),
		query:  `( "TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=? ) AND ( "TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=? )`,
		args:   []interface{}{int64(1), int64(1)},
		errmsg: "",
	}, {
		cond: And(
			Or(
				table1.C("id").Eq(table1.C("test1")),
				table1.C("id").Eq(1),
			),
			table1.C("id").Eq(table1.C("test1")),
		),
		query:  `( "TABLE_A"."id"="TABLE_A"."test1" OR "TABLE_A"."id"=? ) AND "TABLE_A"."id"="TABLE_A"."test1"`,
		args:   []interface{}{int64(1)},
		errmsg: "",
	}}
	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
