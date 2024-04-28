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
Copyright 2024 The Go-CoreLibs Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

[umisama/go-sqlbuilder]: https://github.com/umisama/go-sqlbuilder
[Go-CoreLibs]: https://github.com/go-corelibs
[Go-Curses]: https://github.com/go-curses
[Go-Enjin]: https://github.com/go-enjin
