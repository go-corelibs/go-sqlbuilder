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
