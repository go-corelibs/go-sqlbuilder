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
	"time"
)

func TestInsert(t *testing.T) {
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		StringColumn("str", &ColumnOption{
			Size: 255,
		}),
		BoolColumn("bool", nil),
		FloatColumn("float", nil),
		DateColumn("date", nil),
		BytesColumn("bytes", nil),
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
		stmt: Insert(table1).
			Columns(table1.C("str"), table1.C("bool"), table1.C("float"), table1.C("date"), table1.C("bytes")).
			Values("hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}),
		query:  `INSERT INTO "TABLE_A" ( "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ? );`,
		args:   []interface{}{"hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		errmsg: "",
	}, {
		stmt: Insert(table1).
			Set(table1.C("str"), "hoge").
			Set(table1.C("bool"), true).
			Set(table1.C("float"), 0.1).
			Set(table1.C("date"), time.Unix(0, 0).UTC()).
			Set(table1.C("bytes"), []byte{0x01}),
		query:  `INSERT INTO "TABLE_A" ( "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ? );`,
		args:   []interface{}{"hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		errmsg: "",
	}, {
		stmt:   Insert(table1).Values(1, "hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}),
		query:  `INSERT INTO "TABLE_A" ( "id", "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ?, ? );`,
		args:   []interface{}{int64(1), "hoge", true, 0.1, time.Unix(0, 0).UTC(), []byte{0x01}},
		errmsg: "",
	}, {
		stmt:   Insert(table1).Columns(table1.C("id")).Values(1, 2, 3),
		query:  "",
		args:   []interface{}{},
		errmsg: "sqlbuilder: 1 values needed, but got 3.",
	}, {
		stmt:   Insert(nil).Columns(table1.C("id")).Values(1),
		query:  "",
		args:   []interface{}{},
		errmsg: "sqlbuilder: table is nil.",
	}, {
		stmt:   Insert(table1).Columns(table1.C("str")).Values(1),
		query:  "",
		args:   []interface{}{},
		errmsg: "sqlbuilder: string column not accept int.",
	}, {
		stmt:   Insert(tableJoined).Columns(table1.C("str")).Values(1),
		query:  "",
		args:   []interface{}{},
		errmsg: "sqlbuilder: table is not natural table.",
	}}

	for num, c := range cases {
		mes, args, ok := c.Run()
		if !ok {
			t.Errorf(mes+" (case no.%d)", append(args, num)...)
		}
	}
}
