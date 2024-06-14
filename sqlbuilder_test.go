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
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestMain(m *testing.M) {
	SetDialect(TestingDialect{})
	os.Exit(m.Run())
}

func TestError(t *testing.T) {
	err := newError("hogehogestring")
	if "sqlbuilder: hogehogestring" != err.Error() {
		t.Errorf("failed\ngot %#v", err)
	}
}

func Example() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Set dialect first
	// dialects are in github.com/go-corelibs/go-sqlbuilder/dialects
	SetDialect(TestingDialect{})

	// Define a table
	tbl_person := NewTable(
		"PERSON",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		StringColumn("name", &ColumnOption{
			Unique:  true,
			Size:    255,
			Default: "no_name",
		}),
		DateColumn("birth", nil),
	)

	// Create Table
	query, args, err := CreateTable(tbl_person).ToSql()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = db.Exec(query, args...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Insert data
	// (Table).C function returns a column object.
	query, args, err = Insert(tbl_person).
		Set(tbl_person.C("name"), "Kurisu Makise").
		Set(tbl_person.C("birth"), time.Date(1992, time.July, 25, 0, 0, 0, 0, time.UTC)).
		ToSql()
	_, err = db.Exec(query, args...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Query
	var birth time.Time
	query, args, err = Select(tbl_person).Columns(
		tbl_person.C("birth"),
	).Where(
		tbl_person.C("name").Eq("Kurisu Makise"),
	).ToSql()
	err = db.QueryRow(query, args...).Scan(&birth)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Kurisu's birthday is %s,%d %d", birth.Month().String(), birth.Day(), birth.Year())

	// Output:
	// Kurisu's birthday is July,25 1992
}
