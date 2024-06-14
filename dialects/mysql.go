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
	"time"

	sb "github.com/go-corelibs/go-sqlbuilder"
)

var _ sb.Dialect = MySql{}

type MySql struct{}

func (m MySql) Name() string {
	return "mysql"
}

func (m MySql) QuerySuffix() string {
	return ";"
}

func (m MySql) BindVar(i int) string {
	return "?"
}

func (m MySql) QuoteField(field interface{}) string {
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
		str = "`" + str + "`"
	}
	return str
}

func (m MySql) ColumnTypeToString(cc sb.ColumnConfig) (string, error) {
	if cc.Option().SqlType != "" {
		return cc.Option().SqlType, nil
	}

	typ := ""
	switch cc.Type() {
	case sb.ColumnTypeInt:
		typ = "INTEGER"
	case sb.ColumnTypeString:
		typ = fmt.Sprintf("VARCHAR(%d)", cc.Option().Size)
	case sb.ColumnTypeDate:
		typ = "DATETIME"
	case sb.ColumnTypeFloat:
		typ = "FLOAT"
	case sb.ColumnTypeBool:
		typ = "BOOLEAN"
	case sb.ColumnTypeBytes:
		typ = "BLOB"
	default:
		return "", errors.New("dialects: unknown column type")
	}

	return typ, nil
}

func (m MySql) ColumnOptionToString(co *sb.ColumnOption) (string, error) {
	opt := ""
	if co.PrimaryKey {
		opt = str_append(opt, "PRIMARY KEY")
	}
	if co.AutoIncrement {
		opt = str_append(opt, "AUTO_INCREMENT")
	}
	if co.NotNull {
		opt = str_append(opt, "NOT NULL")
	}
	if co.Unique {
		opt = str_append(opt, "UNIQUE")
	}
	if co.Default != nil {
		opt = str_append(opt, "DEFAULT "+m.QuoteField(co.Default))
	}

	return opt, nil
}

func (m MySql) TableOptionToString(to *sb.TableOption) (string, error) {
	opt := ""
	if to.Unique != nil {
		opt = str_append(opt, m.tableOptionUnique(to.Unique))
	}

	return opt, nil
}

func (m MySql) tableOptionUnique(op [][]string) string {
	opt := ""
	first_op := true
	for _, unique := range op {
		if first_op {
			first_op = false
		} else {
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
			opt += m.QuoteField(col)
		}
		opt += ")"
	}
	return opt
}

func str_append(str, opt string) string {
	if len(str) != 0 {
		str += " "
	}
	str += opt
	return str
}
