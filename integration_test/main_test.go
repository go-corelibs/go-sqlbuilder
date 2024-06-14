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

package sqlbuilder_integration

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/ziutek/mymysql/godrv"

	sb "github.com/go-corelibs/go-sqlbuilder"
	"github.com/go-corelibs/go-sqlbuilder/dialects"
)

var db *sql.DB

// Table for testing
var (
	tbl_person, tbl_phone, tbl_email sb.Table
)

// Data for testing
type Person struct {
	Id    int
	Name  string
	Birth time.Time
}

type Phone struct {
	PersonId int
	Number   string
}

type Email struct {
	PersonId int
	Address  string
}

var persons = []Person{{
	Id:    1,
	Name:  "Rintaro Okabe",
	Birth: time.Date(1991, time.December, 14, 0, 0, 0, 0, time.UTC),
}, {
	Id:    2,
	Name:  "Mayuri Shiina",
	Birth: time.Date(1994, time.February, 1, 0, 0, 0, 0, time.UTC),
}, {
	Id:    3,
	Name:  "Itaru Hashida",
	Birth: time.Date(1991, time.May, 19, 0, 0, 0, 0, time.UTC),
}}

var phones = []Phone{{
	PersonId: 1,
	Number:   "000-0000-0000",
}, {
	PersonId: 2,
	Number:   "111-1111-1111",
}, {
	PersonId: 2,
	Number:   "111-1111-2222",
}}

var emails = []Email{{
	PersonId: 1,
	Address:  "sg-epk@jtk93.x29.jp",
}, {
	PersonId: 1,
	Address:  "okarin@example.org",
}, {
	PersonId: 2,
	Address:  "mayusii@example.org",
}, {
	PersonId: 3,
	Address:  "hashida@example.org",
}}

func TestMain(m *testing.M) {
	results := make(map[string]int)
	type testcase struct {
		name    string
		dialect sb.Dialect
		driver  string
		dsn     string
	}

	var cases = []testcase{
		{"sqlite", dialects.Sqlite{}, "sqlite3", ":memory:"},
		{"mysql(ziutek/mymysql)", dialects.MySql{}, "mymysql", "go_sqlbuilder_test1/root/"},
		{"mysql(go-sql-driver/mysql)", dialects.MySql{}, "mysql", "root:@/go_sqlbuilder_test2?parseTime=true"},
		{"postgres", dialects.Postgresql{}, "postgres", "user=postgres dbname=go_sqlbuilder_test sslmode=disable"},
	}

	for _, c := range cases {
		fmt.Println("START unit test for", c.name)

		// tables
		tbl_person = sb.NewTable(
			"PERSON", nil,
			sb.IntColumn("id", &sb.ColumnOption{
				PrimaryKey: true,
			}),
			sb.StringColumn("name", &sb.ColumnOption{
				Unique:  true,
				Size:    255,
				Default: "default_name",
			}),
			sb.DateColumn("birth", nil),
		)
		tbl_phone = sb.NewTable(
			"PHONE",
			&sb.TableOption{
				Unique: [][]string{{"phone_id", "number"}},
			},
			sb.IntColumn("id", &sb.ColumnOption{
				PrimaryKey:    true,
				AutoIncrement: true,
			}),
			sb.IntColumn("person_id", nil),
			sb.StringColumn("number", &sb.ColumnOption{
				Size: 255,
			}),
		)
		tbl_email = sb.NewTable(
			"EMAIL",
			&sb.TableOption{
				Unique: [][]string{{"person_id", "address"}},
			},
			sb.IntColumn("id", &sb.ColumnOption{
				PrimaryKey:    true,
				AutoIncrement: true,
			}),
			sb.IntColumn("person_id", nil),
			sb.StringColumn("address", &sb.ColumnOption{
				Size: 255,
			}),
		)

		var err error
		db, err = sql.Open(c.driver, c.dsn)
		if err != nil {
			fmt.Println(err.Error())
		}
		sb.SetDialect(c.dialect)

		results[c.name] = m.Run()
	}

	for _, v := range results {
		if v != 0 {
			os.Exit(v)
		}
	}
	os.Exit(0)
}
