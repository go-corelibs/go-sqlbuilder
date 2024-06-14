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
	"github.com/go-corelibs/go-sqlbuilder"
)

// Parse examines the name given to determine which dialect to use, accepts
// the following cases:
//
//	Dialect  | Names, aliases...
//	-------------------------------------
//	MySQL    | mysql, mariadb
//	Postgres | postgres, postgresql, pg
//	Sqlite   | sqlite, sqlite3
//	Testing  | testing, test
func Parse(name string) (d sqlbuilder.Dialect, ok bool) {
	if ok = name == "mysql" || name == "mariadb"; ok {
		d = MySql{}
	} else if ok = name == "pg" || name == "postgres" || name == "postgresql"; ok {
		d = Postgresql{}
	} else if ok = name == "sqlite" || name == "sqlite3"; ok {
		d = Sqlite{}
	} else if ok = name == "testing" || name == "test"; ok {
		d = sqlbuilder.TestingDialect{}
	}
	return
}
