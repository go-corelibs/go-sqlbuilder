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
	"errors"
	"fmt"
	"time"
)

var _ Dialect = TestingDialect{}

type TestingDialect struct{}

func (td TestingDialect) Name() string {
	return "testing"
}

func (td TestingDialect) QuerySuffix() string {
	return ";"
}

func (td TestingDialect) BindVar(i int) string {
	return "?"
}

func (td TestingDialect) QuoteField(field interface{}) string {
	str := ""
	bracket := true
	switch t := field.(type) {
	case string:
		str = t
	case []byte:
		str = string(t)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		str = fmt.Sprint(field)
	case float32, float64:
		str = fmt.Sprint(field)
	case time.Time:
		str = t.Format("2006-01-02 15:04:05")
	case bool:
		if t {
			str = "TRUE"
		} else {
			str = "FALSE"
		}
		bracket = false
	case nil:
		return "NULL"
	}
	if bracket {
		str = "\"" + str + "\""
	}
	return str
}

func (td TestingDialect) ColumnTypeToString(cc ColumnConfig) (string, error) {
	if cc.Option().SqlType != "" {
		return cc.Option().SqlType, nil
	}

	typ := ""
	switch cc.Type() {
	case ColumnTypeInt:
		typ = "INTEGER"
	case ColumnTypeString:
		typ = "TEXT"
	case ColumnTypeDate:
		typ = "DATE"
	case ColumnTypeFloat:
		typ = "REAL"
	case ColumnTypeBool:
		typ = "BOOLEAN"
	case ColumnTypeBytes:
		typ = "BLOB"
	case ColumnTypeAny:
	}
	if typ == "" {
		return "", errors.New("dialects: unknown column type")
	} else {
		return typ, nil
	}
}

func (td TestingDialect) ColumnOptionToString(co *ColumnOption) (string, error) {
	apnd := func(str, opt string) string {
		if len(str) != 0 {
			str += " "
		}
		str += opt
		return str
	}

	opt := ""
	if co.PrimaryKey {
		opt = apnd(opt, "PRIMARY KEY")
	}
	if co.AutoIncrement {
		opt = apnd(opt, "AUTOINCREMENT")
	}
	if co.NotNull {
		opt = apnd(opt, "NOT NULL")
	}
	if co.Unique {
		opt = apnd(opt, "UNIQUE")
	}

	// TestingDialect omitted handling DEFAULT keyword

	return opt, nil
}

func (td TestingDialect) TableOptionToString(to *TableOption) (string, error) {
	opt := ""
	apnd := func(str, opt string) string {
		if len(str) != 0 {
			str += " "
		}
		str += opt
		return str
	}

	if to.Unique != nil {
		opt = apnd(opt, td.tableOptionUnique(to.Unique))
	}
	return opt, nil
}

func (td TestingDialect) tableOptionUnique(op [][]string) string {
	opt := ""
	for idx, unique := range op {
		if idx > 0 {
			opt += " "
		}

		opt += "UNIQUE("
		first := true
		for _, col := range unique {
			if first {
				first = false
			} else {
				opt += ", "
			}
			opt += td.QuoteField(col)
		}
		opt += ")"
	}
	return opt
}
