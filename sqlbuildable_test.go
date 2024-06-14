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

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildable(t *testing.T) {
	tbl := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		StringColumn("thing", nil),
		IntColumn("count", nil),
	)

	d := Dialect(TestingDialect{})
	b := NewBuildable(d)

	Convey("Buildable Methods", t, func() {

		So(b.Dialect(), ShouldEqual, d)

		//tb2 := b.NewTable("TABLE_2", &TableOption{}, IntColumn("id", &ColumnOption{PrimaryKey: true}))
		//So(tb2, ShouldNotBeNil)
		//So(tb2.Name(), ShouldEqual, "TABLE_2")
		//So(b.T("TABLE_2"), ShouldEqual, tb2)
		//So(b.T("NOPE"), ShouldBeNil)

		//tb2a := NewTable("TABLE_2", &TableOption{}, IntColumn("id", &ColumnOption{PrimaryKey: true}), StringColumn("data", &ColumnOption{}))
		//So(tb2a, ShouldNotBeNil)
		//So(tb2a.Name(), ShouldEqual, "TABLE_2")
		//So(b.T("TABLE_2"), ShouldNotEqual, tb2a)
		//So(b.SetTable(tb2a), ShouldBeTrue)
		//So(b.T("TABLE_2"), ShouldEqual, tb2a)

		//So(b.ListTables(), ShouldEqual, []string{"TABLE_2"})

	})

	Convey("AlterTable", t, func() {
		So(b, ShouldNotBeNil)

		sb := b.AlterTable(tbl).AddColumnAfter(BoolColumn("present", &ColumnOption{Default: true}), tbl.C("thing"))
		So(sb, ShouldNotBeNil)
		sql, argv, err := sb.ToSql()
		So(err, ShouldBeNil)
		So(argv, ShouldBeEmpty)
		So(sql, ShouldEqual, `ALTER TABLE "TABLE_A" ADD COLUMN "present" BOOLEAN AFTER "thing";`)
	})

	Convey("CreateTable", t, func() {
		So(b, ShouldNotBeNil)

		sb := b.CreateTable(tbl).IfNotExists()
		So(sb, ShouldNotBeNil)
		sql, argv, err := sb.ToSql()
		So(err, ShouldBeNil)
		So(argv, ShouldBeEmpty)
		So(sql, ShouldEqual, `CREATE TABLE IF NOT EXISTS "TABLE_A" ( "id" INTEGER PRIMARY KEY, "thing" TEXT, "count" INTEGER );`)
	})

	Convey("CreateIndex", t, func() {
		So(b, ShouldNotBeNil)

		sb := b.CreateIndex(tbl).
			Name("index_thing").
			Columns(tbl.C("thing"))
		So(sb, ShouldNotBeNil)
		sql, argv, err := sb.ToSql()
		So(err, ShouldBeNil)
		So(argv, ShouldBeEmpty)
		So(sql, ShouldEqual, `CREATE INDEX "index_thing" ON "TABLE_A" ( "thing" );`)
	})

	Convey("Delete", t, func() {
		So(b, ShouldNotBeNil)

		sb := b.Delete(tbl)
		So(sb, ShouldNotBeNil)
		sql, argv, err := sb.ToSql()
		So(err, ShouldBeNil)
		So(argv, ShouldBeEmpty)
		So(sql, ShouldEqual, `DELETE FROM "TABLE_A";`)
	})

	Convey("Insert", t, func() {
		So(b, ShouldNotBeNil)

		sb := b.Insert(tbl).
			Set(tbl.C("thing"), "value").
			Set(tbl.C("count"), 10)
		So(sb, ShouldNotBeNil)
		sql, argv, err := sb.ToSql()
		So(err, ShouldBeNil)
		So(len(argv), ShouldEqual, 2)
		So(argv[0], ShouldEqual, "value")
		So(argv[1], ShouldEqual, 10)
		So(sql, ShouldEqual, `INSERT INTO "TABLE_A" ( "thing", "count" ) VALUES ( ?, ? );`)
	})

	Convey("Select", t, func() {
		So(b, ShouldNotBeNil)

		sb := b.Select(tbl)
		So(sb, ShouldNotBeNil)
		sql, argv, err := sb.ToSql()
		So(err, ShouldBeNil)
		So(argv, ShouldBeEmpty)
		So(sql, ShouldEqual, `SELECT * FROM "TABLE_A";`)

	})

	Convey("Columns", t, func() {
		table := NewTable(
			"TABLE_A",
			&TableOption{},
			IntColumn("id", &ColumnOption{
				PrimaryKey: true,
			}),
			StringColumn("thing", nil),
			IntColumn("count", nil),
		)
		c := table.C("thing").As("alias")
		So(c, ShouldNotBeNil)
		Select(table).Columns(c).Where(c.Eq("stuff"))
	})

	Convey("Joined Tables", t, func() {
		t0 := NewTable(
			"TABLE_0",
			&TableOption{},
			IntColumn("id", &ColumnOption{
				PrimaryKey: true,
			}),
			StringColumn("thing", nil),
			IntColumn("count", nil),
		)
		t1 := NewTable(
			"TABLE_1",
			&TableOption{},
			IntColumn("id", &ColumnOption{
				PrimaryKey: true,
			}),
			StringColumn("other", nil),
			IntColumn("num", nil),
		)

		// SELECT TABLE_A.id, TABLE_A.thing, TABLE_0.thing, TABLE_1.thing
		// FROM TABLE_A
		// INNER JOIN TABLE_0 ON TABLE_A.id = TABLE_0.num
		// INNER JOIN TABLE_1 ON TABLE_A.id = TABLE_1.num;

		joined := tbl.InnerJoin(t0, tbl.C("id").Eq(t0.C("count")))
		So(joined.Describe(), ShouldEqual, `TABLE_A
	.Column[int]("id": PrimaryKey)
	.Column[string]("thing")
	.Column[int]("count")
	INNER JOIN TABLE_0 ON (id = count)`)

		again := joined.InnerJoin(t1, tbl.C("id").Eq(t1.C("num")))
		So(again.Describe(), ShouldEqual, `TABLE_A
	.Column[int]("id": PrimaryKey)
	.Column[string]("thing")
	.Column[int]("count")
	INNER JOIN TABLE_0 ON (id = count)
	INNER JOIN TABLE_1 ON (id = num)`)

		sql, argv, err := Select(again).
			Columns(tbl.C("id")).
			ToSql()
		So(err, ShouldBeNil)
		So(argv, ShouldEqual, []interface{}{})
		So(sql, ShouldEqual, ``+
			`SELECT "TABLE_A"."id" `+
			`FROM "TABLE_A" `+
			`INNER JOIN "TABLE_0" ON "TABLE_A"."id"="TABLE_0"."count" `+
			`INNER JOIN "TABLE_1" ON "TABLE_A"."id"="TABLE_1"."num";`,
		)

	})
}
