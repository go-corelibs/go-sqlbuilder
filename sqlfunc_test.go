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
	"testing"
)

func TestSqlFuncImplements(t *testing.T) {
	fnImplColumn := func(i interface{}) bool {
		return reflect.TypeOf(i).Implements(reflect.TypeOf(new(Column)).Elem())
	}
	fnImplColumn(&columnImpl{})
}

func TestSqlFunc(t *testing.T) {
	b := newBuilder(TestDialect{})
	table1 := NewTable(
		"TABLE_A",
		&TableOption{},
		IntColumn("id", &ColumnOption{
			PrimaryKey: true,
		}),
		IntColumn("test1", nil),
		IntColumn("test2", nil),
	)

	Func("funcname", table1.C("id")).serialize(b)
	if `funcname("TABLE_A"."id")` != b.query.String() {
		t.Errorf("failed")
	}
	if len(b.Args()) != 0 {
		t.Errorf("failed")
	}
	if b.Err() != nil {
		t.Errorf("failed")
	}
}
