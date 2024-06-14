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
	"errors"
	"fmt"
	"strconv"
	"time"

	sb "github.com/go-corelibs/go-sqlbuilder"
)

var _ sb.Dialect = Sqlite{}

type Sqlite struct{}

func (m Sqlite) Name() string {
	return "sqlite3"
}

func (m Sqlite) QuerySuffix() string {
	return ";"
}

func (m Sqlite) BindVar(i int) string {
	return "?"
}

func (m Sqlite) QuoteField(field interface{}) string {
	switch t := field.(type) {
	case string:
		return strconv.Quote(t)
	case []byte:
		return strconv.Quote(string(t))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return strconv.Quote(fmt.Sprint(t))
	case float32, float64:
		return strconv.Quote(fmt.Sprint(t))
	case time.Time:
		return strconv.Quote(t.Format("2006-01-02 15:04:05"))
	case bool:
		if t {
			return "TRUE"
		}
		return "FALSE"
	case nil:
		return "NULL"
	}
	return ""
}

func (m Sqlite) ColumnTypeToString(cc sb.ColumnConfig) (string, error) {
	if cc.Option().SqlType != "" {
		return cc.Option().SqlType, nil
	}

	switch cc.Type() {
	case sb.ColumnTypeInt:
		return "INTEGER", nil
	case sb.ColumnTypeString:
		return "TEXT", nil
	case sb.ColumnTypeDate:
		return "DATETIME", nil
	case sb.ColumnTypeFloat:
		return "REAL", nil
	case sb.ColumnTypeBool:
		return "BOOLEAN", nil
	case sb.ColumnTypeBytes:
		return "BLOB", nil
	default:
		return "", errors.New("dialects: unknown column type")
	}
}

func (m Sqlite) ColumnOptionToString(co *sb.ColumnOption) (string, error) {
	opt := ""
	if co.PrimaryKey {
		opt = str_append(opt, "PRIMARY KEY")
	}
	if co.AutoIncrement {
		opt = str_append(opt, "AUTOINCREMENT")
	}
	if co.NotNull {
		opt = str_append(opt, "NOT NULL")
	}
	if co.Unique {
		opt = str_append(opt, "UNIQUE")
	}
	if co.Default == nil {
		if !co.PrimaryKey {
			opt = str_append(opt, "DEFAULT NULL")
		}
	} else {
		opt = str_append(opt, "DEFAULT "+m.QuoteField(co.Default))
	}

	return opt, nil
}

func (m Sqlite) TableOptionToString(to *sb.TableOption) (string, error) {
	opt := ""
	if to.Unique != nil {
		opt = str_append(opt, m.tableOptionUnique(to.Unique))
	}

	return opt, nil
}

func (m Sqlite) tableOptionUnique(op [][]string) (opt string) {
	for idx, unique := range op {
		if idx > 0 {
			opt += ", "
		}
		opt += "UNIQUE("
		for jdx, col := range unique {
			if jdx > 0 {
				opt += ", "
			}
			opt += m.QuoteField(col)
		}
		opt += ")"
	}
	return opt
}
