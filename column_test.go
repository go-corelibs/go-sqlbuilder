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

func TestColumnImplements(t *testing.T) {
	fnImplColumn := func(i interface{}) bool {
		return reflect.TypeOf(i).Implements(reflect.TypeOf(new(Column)).Elem())
	}
	if !fnImplColumn(&cColumnImpl{}) {
		t.Errorf("fail")
	}
	if !fnImplColumn(&cErrorColumn{}) {
		t.Errorf("fail")
	}
	if !fnImplColumn(&cColumnAlias{}) {
		t.Errorf("fail")
	}
}

func TestColumnOptionImpl(t *testing.T) {
	if !reflect.DeepEqual(&cColumnImplConfig{
		name: "name",
		typ:  ColumnTypeBytes,
		opt: &ColumnOption{
			Unique: true,
		}}, newColumnImplConfig("name", ColumnTypeBytes, &ColumnOption{Unique: true})) {
		t.Errorf("fail")
	}
	if !reflect.DeepEqual(&cColumnImplConfig{
		name: "name",
		typ:  ColumnTypeBytes,
		opt:  &ColumnOption{},
	}, newColumnImplConfig("name", ColumnTypeBytes, nil)) {
		t.Errorf("fail")
	}
}
