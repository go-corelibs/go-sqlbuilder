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
)

type statementTestCase struct {
	stmt   Statement
	query  string
	args   []interface{}
	errmsg string
}

func (testCase statementTestCase) Run() (message string, args []interface{}, ok bool) {
	query, args, err := testCase.stmt.ToSql()
	if len(testCase.errmsg) != 0 {
		if err == nil {
			return "error: expect returns error but got nil.", []interface{}{}, false
		}
		if err.Error() != testCase.errmsg {
			return "error: expect error message is '%s' but got '%s'.", []interface{}{err.Error(), testCase.errmsg}, false
		}
	} else {
		if err != nil {
			return "error: expect returns no error got %s.", []interface{}{err.Error()}, false
		}
	}
	if testCase.query != query {
		return "expect returns query \n%s \nbut got\n%s.", []interface{}{testCase.query, query}, false
	}
	if !reflect.DeepEqual(testCase.args, args) {
		return "expect returns arguments \n%s \nbut got\n%s.", []interface{}{testCase.args, args}, false
	}
	return "", nil, true
}

type conditionTestCase struct {
	cond   Condition
	query  string
	args   []interface{}
	errmsg string
}

func (testCase conditionTestCase) Run() (message string, args []interface{}, ok bool) {
	bldr := newBuilder(TestDialect{})
	testCase.cond.serialize(bldr)
	if len(testCase.errmsg) != 0 {
		if bldr.err == nil {
			return "error: expect returns error but got nil.", []interface{}{}, false
		}
		if bldr.err.Error() != testCase.errmsg {
			return "error: expect error message is '%s' but got '%s'.", []interface{}{bldr.err.Error(), testCase.errmsg}, false
		}
	} else {
		if bldr.err != nil {
			return "error: expect returns no error got %s.", []interface{}{bldr.err.Error()}, false
		}
	}
	if bldr.query.String() != testCase.query {
		return "expect returns query \n%s \nbut got\n%s.", []interface{}{testCase.query, bldr.query.String()}, false
	}
	if !reflect.DeepEqual(bldr.args, testCase.args) {
		return "expect returns arguments \n%s \nbut got\n%s.", []interface{}{testCase.args, args}, false
	}
	return "", nil, true
}
