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
	"reflect"
	"time"
)

// ColumnType reprecents a type of column.
// Dialects handle this for know column options.
type ColumnType int

const (
	ColumnTypeAny ColumnType = iota
	ColumnTypeInt
	ColumnTypeString
	ColumnTypeDate
	ColumnTypeFloat
	ColumnTypeBool
	ColumnTypeBytes
)

func (t ColumnType) String() string {
	switch t {
	case ColumnTypeInt:
		return "int"
	case ColumnTypeString:
		return "string"
	case ColumnTypeDate:
		return "date"
	case ColumnTypeFloat:
		return "float"
	case ColumnTypeBool:
		return "bool"
	case ColumnTypeBytes:
		return "bytes"
	case ColumnTypeAny:
		return "any"
	}
	panic(newError("unknown columnType"))
}

func (t ColumnType) CapableTypes() []reflect.Type {
	switch t {
	case ColumnTypeInt:
		return []reflect.Type{
			reflect.TypeOf(int(0)),
			reflect.TypeOf(int8(0)),
			reflect.TypeOf(int16(0)),
			reflect.TypeOf(int32(0)),
			reflect.TypeOf(int64(0)),
			reflect.TypeOf(uint(0)),
			reflect.TypeOf(uint8(0)),
			reflect.TypeOf(uint16(0)),
			reflect.TypeOf(uint32(0)),
			reflect.TypeOf(uint64(0)),
		}
	case ColumnTypeString:
		return []reflect.Type{
			reflect.TypeOf(""),
		}
	case ColumnTypeDate:
		return []reflect.Type{
			reflect.TypeOf(time.Time{}),
		}
	case ColumnTypeFloat:
		return []reflect.Type{
			reflect.TypeOf(float32(0)),
			reflect.TypeOf(float64(0)),
		}
	case ColumnTypeBool:
		return []reflect.Type{
			reflect.TypeOf(bool(true)),
		}
	case ColumnTypeBytes:
		return []reflect.Type{
			reflect.TypeOf([]byte{}),
		}
	case ColumnTypeAny:
		return []reflect.Type{} // but accept all types
	}
	return []reflect.Type{}
}
