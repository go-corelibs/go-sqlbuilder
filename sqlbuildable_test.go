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

	b := NewBuildable(TestDialect{})

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
}
