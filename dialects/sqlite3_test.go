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
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-corelibs/go-sqlbuilder"
)

func TestSqlite3(t *testing.T) {
	d := Sqlite{}

	Convey("QuerySuffix", t, func() {
		So(d.QuerySuffix(), ShouldEqual, `;`)
	})

	Convey("BindVar", t, func() {
		So(d.BindVar(0), ShouldEqual, `?`)
		So(d.BindVar(1), ShouldEqual, `?`)
		So(d.BindVar(2), ShouldEqual, `?`)
	})

	Convey("QuoteField", t, func() {
		now := time.Now()
		for idx, test := range []struct {
			i interface{}
			o string
		}{
			{"ten", `"ten"`},
			{[]byte("yes"), `"yes"`},
			{10, `"10"`},
			{10.10, `"10.1"`},
			{now, `"` + now.Format("2006-01-02 15:04:05") + `"`},
			{true, `TRUE`},
			{false, `FALSE`},
			{nil, `NULL`},
			{time.Minute, ``},
		} {
			Convey(fmt.Sprintf("case #%d", idx), func() {
				So(d.QuoteField(test.i), ShouldEqual, test.o)
			})
		}
	})

	Convey("ColumnTypeToString", t, func() {

		for idx, test := range []struct {
			input  sqlbuilder.ColumnConfig
			output string
			err    Assertion
		}{
			{
				sqlbuilder.AnyColumn("any_column", &sqlbuilder.ColumnOption{Size: 10}),
				``,
				ShouldNotBeNil,
			},
			{
				sqlbuilder.AnyColumn("any_column", &sqlbuilder.ColumnOption{Size: 10, SqlType: "TEST"}),
				`TEST`,
				ShouldBeNil,
			},
			{
				sqlbuilder.StringColumn("string_column", &sqlbuilder.ColumnOption{Size: 1010}),
				`TEXT`,
				ShouldBeNil,
			},
			{
				sqlbuilder.IntColumn("int_column", &sqlbuilder.ColumnOption{}),
				`INTEGER`,
				ShouldBeNil,
			},
			{
				sqlbuilder.FloatColumn("float_column", &sqlbuilder.ColumnOption{}),
				`REAL`,
				ShouldBeNil,
			},
			{
				sqlbuilder.BoolColumn("bool_column", &sqlbuilder.ColumnOption{}),
				`BOOLEAN`,
				ShouldBeNil,
			},
			{
				sqlbuilder.BytesColumn("bytes_column", &sqlbuilder.ColumnOption{}),
				`BLOB`,
				ShouldBeNil,
			},
			{
				sqlbuilder.DateColumn("date_column", &sqlbuilder.ColumnOption{}),
				`DATE`,
				ShouldBeNil,
			},
		} {
			Convey(fmt.Sprintf("case #%d", idx), func() {
				str, err := d.ColumnTypeToString(test.input)
				So(err, test.err)
				So(str, ShouldEqual, test.output)
			})
		}

	})

	Convey("ColumnOptionToString", t, func() {

		for idx, test := range []struct {
			input  *sqlbuilder.ColumnOption
			output string
			err    Assertion
		}{
			{
				&sqlbuilder.ColumnOption{},
				`DEFAULT NULL`,
				ShouldBeNil,
			},
			{
				&sqlbuilder.ColumnOption{NotNull: true},
				`NOT NULL DEFAULT NULL`,
				ShouldBeNil,
			},
			{
				&sqlbuilder.ColumnOption{Default: "thing"},
				`DEFAULT "thing"`,
				ShouldBeNil,
			},
			{
				&sqlbuilder.ColumnOption{Unique: true},
				`UNIQUE DEFAULT NULL`,
				ShouldBeNil,
			},
			{
				&sqlbuilder.ColumnOption{PrimaryKey: true},
				`PRIMARY KEY`,
				ShouldBeNil,
			},
			{
				&sqlbuilder.ColumnOption{PrimaryKey: true, AutoIncrement: true},
				`PRIMARY KEY AUTOINCREMENT`,
				ShouldBeNil,
			},
		} {
			Convey(fmt.Sprintf("case #%d", idx), func() {
				str, err := d.ColumnOptionToString(test.input)
				So(err, test.err)
				So(str, ShouldEqual, test.output)
			})
		}

	})

	Convey("TableOptionToString", t, func() {

		for idx, test := range []struct {
			input  *sqlbuilder.TableOption
			output string
			err    Assertion
		}{
			{
				&sqlbuilder.TableOption{},
				``,
				ShouldBeNil,
			},
			{
				&sqlbuilder.TableOption{Unique: [][]string{{"one"}}},
				`UNIQUE("one")`,
				ShouldBeNil,
			},
			{
				&sqlbuilder.TableOption{Unique: [][]string{{"one", "two"}}},
				`UNIQUE("one", "two")`,
				ShouldBeNil,
			},
			{
				&sqlbuilder.TableOption{Unique: [][]string{{"one"}, {"two"}}},
				`UNIQUE("one") UNIQUE("two")`,
				ShouldBeNil,
			},
		} {
			Convey(fmt.Sprintf("case #%d", idx), func() {
				str, err := d.TableOptionToString(test.input)
				So(err, test.err)
				So(str, ShouldEqual, test.output)
			})
		}

	})
}
