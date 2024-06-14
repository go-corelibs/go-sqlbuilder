[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/go-corelibs/go-sqlbuilder)
[![codecov](https://codecov.io/gh/go-corelibs/go-sqlbuilder/graph/badge.svg?token=RKKUET0wcB)](https://codecov.io/gh/go-corelibs/go-sqlbuilder)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-corelibs/go-sqlbuilder)](https://goreportcard.com/report/github.com/go-corelibs/go-sqlbuilder)

# go-sqlbuilder - multi-dialect fork of github.com/umisama/go-sqlbuilder

This package is a fork of [umisama/go-sqlbuilder] which primarily adds the
ability to use multiple dialects at the same time with the introduction of
a new Builder instance type.

# Installation

``` shell
> go get github.com/go-corelibs/go-sqlbuilder@latest
```

# Examples

## Buildable

This example does not include any error checking for brevity.

``` go
func main() {

    // connect to a sqlite3 instance
    db, _ := sql.Open("sqlite3", ":memory:")
    defer db.Close()

    // construct a new instance, configured for sqlite usage
    b := sqlbuilder.NewBuildable(dialects.Sqlite)

    // define a new table to work with
    table := sqlbuilder.NewTable(
        "table_name",
        &TableOption{},
        IntColumn("id", &ColumnOption{
            PrimaryKey: true,
        }),
        StringColumn("content", nil),
    )

    // create the newly defined table
    sql, args, _ := b.CreateTable(table).ToSQL()
    // execute the create-table statement
     _, _ = db.Exec(query, args...)

     // insert some rows into the table
     query, args, _ = b.Insert(table).
         Set(table.C("content"), "Hello strange new world")
         ToSQL()
    // execute the insert row statement
     _, _ = db.Exec(query, args...)

    // ... and so on, almost identical to the original
    // package usage, with the exception of using the
    // Buildable instance instead of the package-level
    // functions directly. This allows for multiple
    // Buildable instances, with different dialects,
    // to be used at the same time
}
```

# Go-CoreLibs

[Go-CoreLibs] is a repository of shared code between the [Go-Curses] and
[Go-Enjin] projects.

# License

```
Copyright (c) 2014 umisama <Takaaki IBARAKI>
Copyright (c) 2024 The Go-CoreLibs Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```

[umisama/go-sqlbuilder]: https://github.com/umisama/go-sqlbuilder
[Go-CoreLibs]: https://github.com/go-corelibs
[Go-Curses]: https://github.com/go-curses
[Go-Enjin]: https://github.com/go-enjin
