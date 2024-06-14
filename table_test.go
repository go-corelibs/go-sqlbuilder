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

	. "github.com/smartystreets/goconvey/convey"
)

func TestTable(t *testing.T) {
	var table1 Table
	var fnPanic = func(fn func()) (ok bool) {
		defer func() {
			if r := recover(); r != nil {
				ok = true
			} else {
				ok = false
			}
		}()
		fn()
		return
	}

	Convey("panic testing", t, func() {
		// panic
		So(fnPanic(func() {
			table1 = NewTable(
				"TABLE_NAME",
				&TableOption{},
			)
		}), ShouldBeTrue)
		So(table1, ShouldBeNil)

		// not panic
		So(fnPanic(func() {
			table1 = NewTable(
				"TABLE_NAME",
				&TableOption{},
				IntColumn("id", nil),
			)
		}), ShouldBeFalse)
		So(table1, ShouldNotBeNil)

		// not panic
		So(fnPanic(func() {
			table1 = NewTable(
				"TABLE_NAME",
				nil,
				IntColumn("id", nil),
			)
		}), ShouldBeFalse)
		So(table1, ShouldNotBeNil)

		// not panic
		So(fnPanic(func() {
			table1 = NewTable(
				"TABLE_NAME",
				nil,
				IntColumn("id", nil),
			)
		}), ShouldBeFalse)
		So(table1, ShouldNotBeNil)

	})

}

func TestJoinTable(t *testing.T) {
	l_table := NewTable(
		"LEFT_TABLE",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("right_id", nil),
	)
	r_table := NewTable(
		"RIGHT_TABLE",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("value", nil),
	)
	rr_table := NewTable(
		"RIGHTRIGHT_TABLE",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
	)

	Convey("joins", t, func() {
		// inner join
		b := newBuilder(TestDialect{})
		joinedTable := l_table.InnerJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
		joinedTable.serialize(b)
		So(b.query.String(), ShouldEqual, `"LEFT_TABLE" INNER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`)
		So(b.err, ShouldBeNil)
		So(len(b.args), ShouldEqual, 0)

		// left outer join
		b = newBuilder(TestDialect{})
		joinedTable = l_table.LeftOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
		joinedTable.serialize(b)
		So(b.query.String(), ShouldEqual, `"LEFT_TABLE" LEFT OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`)
		So(b.err, ShouldBeNil)
		So(len(b.args), ShouldEqual, 0)

		// right outer join
		b = newBuilder(TestDialect{})
		joinedTable = l_table.RightOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
		joinedTable.serialize(b)
		So(b.query.String(), ShouldEqual, `"LEFT_TABLE" RIGHT OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`)
		So(b.err, ShouldBeNil)
		So(len(b.args), ShouldEqual, 0)

		// full outer join
		b = newBuilder(TestDialect{})
		joinedTable = l_table.FullOuterJoin(r_table, l_table.C("right_id").Eq(r_table.C("id")))
		joinedTable.serialize(b)
		So(b.query.String(), ShouldEqual, `"LEFT_TABLE" FULL OUTER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id"`)
		So(b.err, ShouldBeNil)
		So(len(b.args), ShouldEqual, 0)

		// joined table column
		if !reflect.DeepEqual(l_table.C("right_id"), joinedTable.C("right_id")) {
			t.Error("failed")
		}
		if !reflect.DeepEqual(r_table.C("value"), joinedTable.C("value")) {
			t.Error("failed")
		}
		if _, ok := joinedTable.C("not_exist_column").(*errorColumn); !ok {
			t.Error("failed")
		}
		if _, ok := joinedTable.C("id").(*errorColumn); !ok {
			t.Error("failed")
		}

		// combination
		b = newBuilder(TestDialect{})
		joinedTable = l_table.InnerJoin(r_table, l_table.C("right_id").Eq(r_table.C("id"))).InnerJoin(rr_table, l_table.C("right_id").Eq(rr_table.C("id")))
		joinedTable.serialize(b)
		So(b.query.String(), ShouldEqual, `"LEFT_TABLE" INNER JOIN "RIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHT_TABLE"."id" INNER JOIN "RIGHTRIGHT_TABLE" ON "LEFT_TABLE"."right_id"="RIGHTRIGHT_TABLE"."id"`)
		So(b.err, ShouldBeNil)
		So(len(b.args), ShouldEqual, 0)

	})

}

func TestTableColumnOperation(t *testing.T) {
	var fnEqualColumnName = func(cols []Column, expect []string) bool {
		if len(cols) != len(expect) {
			return false
		}
		for i, col := range cols {
			if col.column_name() != expect[i] {
				return false
			}
		}
		return true
	}

	table1 := NewTable(
		"TABLE_NAME",
		nil,
		IntColumn("id", nil),
	).(*table)

	Convey("column operations", t, func() {

		// initial check
		So(fnEqualColumnName(table1.Columns(), []string{"id"}), ShouldBeTrue)

		// AddColumnLast
		err := table1.AddColumnLast(IntColumn("test1", nil))
		So(err, ShouldBeNil)
		So(fnEqualColumnName(table1.Columns(), []string{"id", "test1"}), ShouldBeTrue)

		// AddColumnFirst
		err = table1.AddColumnFirst(IntColumn("first", nil))
		So(err, ShouldBeNil)
		So(fnEqualColumnName(table1.Columns(), []string{"first", "id", "test1"}), ShouldBeTrue)

		// AddColumnAfter
		err = table1.AddColumnAfter(IntColumn("second", nil), table1.C("first"))
		So(err, ShouldBeNil)
		err = table1.AddColumnAfter(IntColumn("aaa", nil), table1.C("invalid"))
		So(err, ShouldNotBeNil)
		So(fnEqualColumnName(table1.Columns(), []string{"first", "second", "id", "test1"}), ShouldBeTrue)

		// ChangeColumn
		err = table1.ChangeColumn(table1.C("id"), IntColumn("third", nil))
		So(err, ShouldBeNil)
		err = table1.ChangeColumn(table1.C("invalid"), IntColumn("third", nil))
		So(err, ShouldNotBeNil)
		So(fnEqualColumnName(table1.Columns(), []string{"first", "second", "third", "test1"}), ShouldBeTrue)

		// ChangeColumnFirst
		err = table1.ChangeColumnFirst(table1.C("test1"), IntColumn("new_first", nil))
		So(err, ShouldBeNil)
		err = table1.ChangeColumnFirst(table1.C("invalid"), IntColumn("new_first", nil))
		So(err, ShouldNotBeNil)
		So(fnEqualColumnName(table1.Columns(), []string{"new_first", "first", "second", "third"}), ShouldBeTrue)

		// ChangeColumnAfter
		err = table1.ChangeColumnAfter(table1.C("new_first"), IntColumn("fourth", nil), table1.C("third"))
		So(err, ShouldBeNil)
		err = table1.ChangeColumnAfter(table1.C("invalid"), IntColumn("fourth", nil), table1.C("third"))
		So(err, ShouldNotBeNil)
		err = table1.ChangeColumnAfter(table1.C("second"), IntColumn("fourth", nil), table1.C("invalid"))
		So(err, ShouldNotBeNil)
		So(fnEqualColumnName(table1.Columns(), []string{"first", "second", "third", "fourth"}), ShouldBeTrue)

		// ChangeColumnAfter
		err = table1.DropColumn(table1.C("fourth"))
		So(err, ShouldBeNil)
		err = table1.DropColumn(table1.C("invalid"))
		So(err, ShouldNotBeNil)
		So(fnEqualColumnName(table1.Columns(), []string{"first", "second", "third"}), ShouldBeTrue)

		// SetName
		table1.SetName("TABLE_MODIFIED")
		So(table1.Name(), ShouldEqual, "TABLE_MODIFIED")

	})

}
