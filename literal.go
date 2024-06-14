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
	sqldriver "database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type literal interface {
	serializable
	Raw() interface{}
	IsNil() bool
}

type literalImpl struct {
	raw         interface{}
	placeholder bool
}

func toLiteral(v interface{}) literal {
	refv := reflect.ValueOf(v)
	if v != nil &&
		refv.Kind() == reflect.Ptr &&
		!refv.Type().Implements(reflect.TypeOf((*sqldriver.Valuer)(nil)).Elem()) {
		if refv.IsNil() {
			v = nil
		} else {
			v = reflect.Indirect(refv).Interface()
		}
	}
	return &literalImpl{
		raw:         v,
		placeholder: true,
	}
}

func (l *literalImpl) serialize(bldr *builder) {
	val, err := l.converted()
	if err != nil {
		bldr.SetError(err)
		return
	}

	if l.placeholder {
		bldr.AppendValue(val)
	} else {
		bldr.Append(l.string())
	}
	return
}

func (l *literalImpl) IsNil() bool {
	if l.raw == nil {
		return true
	}

	v := reflect.ValueOf(l.raw)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// convert to sqldriver.Value(int64/float64/bool/[]byte/string/time.Time)
func (l *literalImpl) converted() (interface{}, error) {
	switch t := l.raw.(type) {
	case int, int8, int16, int32, int64:
		return int64(reflect.ValueOf(t).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(t).Uint()), nil
	case float32, float64:
		return reflect.ValueOf(l.raw).Float(), nil
	case bool:
		return t, nil
	case []byte:
		return t, nil
	case string:
		return t, nil
	case time.Time:
		return t, nil
	case sqldriver.Valuer:
		return t, nil
	case nil:
		return nil, nil
	default:
		return nil, newError("got %T type, but literal is not supporting this.", t)
	}
}

func (l *literalImpl) string() string {
	val, err := l.converted()
	if err != nil {
		return ""
	}

	switch t := val.(type) {
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', 10, 64)
	case bool:
		return strconv.FormatBool(t)
	case string:
		return t
	case []byte:
		return string(t)
	case time.Time:
		return t.Format("2006-01-02 15:04:05")
	case fmt.Stringer:
		return t.String()
	case nil:
		return "NULL"
	default:
		return ""
	}
}

func (l *literalImpl) Raw() interface{} {
	return l.raw
}
