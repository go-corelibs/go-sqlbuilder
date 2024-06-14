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

package dialects

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-corelibs/go-sqlbuilder"
)

func TestParse(t *testing.T) {
	Convey("general", t, func() {
		tests := []struct {
			input  string
			parsed sqlbuilder.Dialect
			ok     bool
		}{
			{"mysql", MySql{}, true},
			{"mariadb", MySql{}, true},
			{"mysql5", nil, false},
			{"postgres", Postgresql{}, true},
			{"postgresql", Postgresql{}, true},
			{"pg", Postgresql{}, true},
			{"pgsql", nil, false},
			{"sqlite", Sqlite{}, true},
			{"sqlite3", Sqlite{}, true},
			{"sqlite2", nil, false},
			{"nope", nil, false},
		}

		for _, test := range tests {
			parsed, ok := Parse(test.input)
			So(ok, ShouldEqual, test.ok)
			So(parsed, ShouldEqual, test.parsed)
		}
	})
}

var (
	gTestInputs = []string{
		`ALTER TABLE "TABLE_A" ADD COLUMN "test0" INTEGER AFTER "id";`,
		`ALTER TABLE "TABLE_A" ADD COLUMN "test0" INTEGER FIRST;`,
		`ALTER TABLE "TABLE_A" ADD COLUMN "test3" INTEGER UNIQUE;`,
		`ALTER TABLE "TABLE_A" ADD COLUMN "test3" INTEGER, ADD COLUMN "test4" INTEGER, CHANGE COLUMN "test1" "test1a" INTEGER, DROP COLUMN "test1", RENAME TO "TABLE_AAA";`,
		`ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER AFTER "test2";`,
		`ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER FIRST;`,
		`ALTER TABLE "TABLE_A" CHANGE COLUMN "test1" "test1a" INTEGER UNIQUE;`,
		`ALTER TABLE "TABLE_A" DROP COLUMN "test1";`,
		`ALTER TABLE "TABLE_A" RENAME TO "TABLE_AAA";`,
		`CREATE INDEX IF NOT EXISTS "I_TABLE_A" ON "TABLE_A" ( "test1", "test2" );`,
		`CREATE TABLE IF NOT EXISTS "TABLE_A" ( "id" INTEGER PRIMARY KEY AUTOINCREMENT, "test1" INTEGER UNIQUE, "test2" TEXT );`,
		`CREATE TABLE IF NOT EXISTS "TABLE_B" ( "id" VARCHAR(255) PRIMARY KEY AUTOINCREMENT, "test1" INTEGER UNIQUE );`,
		`CREATE TABLE IF NOT EXISTS "TABLE_C" ( "id" INTEGER PRIMARY KEY AUTOINCREMENT, "test1" INTEGER UNIQUE, "test2" TEXT ) UNIQUE("test1", "test2");`,
		`DELETE FROM "TABLE_A" WHERE "TABLE_A"."id"=?;`,
		`DROP TABLE "TABLE_A";`,
		`INSERT INTO "TABLE_A" ( "id", "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ?, ? );`,
		`INSERT INTO "TABLE_A" ( "str", "bool", "float", "date", "bytes" ) VALUES ( ?, ?, ?, ?, ? );`,
		`SELECT "TABLE_A"."id" AS "tbl1id" FROM "TABLE_A" WHERE "tbl1id"=? ORDER BY "TABLE_A"."test1" ASC, "TABLE_A"."test2" DESC;`,
		`SELECT "TABLE_A"."id" AS "tbl1id" FROM "TABLE_A" WHERE "tbl1id"=?;`,
		`SELECT "TABLE_A"."test1", "TABLE_A"."test2" FROM "TABLE_A";`,
		`SELECT * FROM "TABLE_A" INNER JOIN "TABLE_B" ON "TABLE_A"."test1"="TABLE_B"."id";`,
		`SELECT * FROM "TABLE_A";`,
		`UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=? ORDER BY "TABLE_A"."test1" DESC LIMIT ? OFFSET ?;`,
		`UPDATE "TABLE_A" SET "test1"=?, "test2"=? WHERE "TABLE_A"."id"=?;`,
	}
)
