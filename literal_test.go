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
	"time"
)

func TestLiteralConvert(t *testing.T) {
	str := "makise-kurisu"
	var cases = []struct {
		lit    literal
		out    interface{}
		errmes string
	}{
		{
			lit:    toLiteral(int(10)),
			out:    int64(10),
			errmes: "",
		}, {
			lit:    toLiteral(int64(10)),
			out:    int64(10),
			errmes: "",
		}, {
			lit:    toLiteral(uint(10)),
			out:    int64(10),
			errmes: "",
		}, {
			lit:    toLiteral(uint64(10)),
			out:    int64(10),
			errmes: "",
		}, {
			lit:    toLiteral(float32(10)),
			out:    float64(10),
			errmes: "",
		}, {
			lit:    toLiteral(float64(10)),
			out:    float64(10),
			errmes: "",
		}, {
			lit:    toLiteral(bool(true)),
			out:    bool(true),
			errmes: "",
		}, {
			lit:    toLiteral([]byte{0x11}),
			out:    []byte{0x11},
			errmes: "",
		}, {
			lit:    toLiteral(string("makise-kurisu")),
			out:    string("makise-kurisu"),
			errmes: "",
		}, {
			lit:    toLiteral(&str),
			out:    str,
			errmes: "",
		}, {
			lit:    toLiteral((*string)(nil)),
			out:    nil,
			errmes: "",
		}, {
			lit:    toLiteral(time.Unix(0, 0)),
			out:    time.Unix(0, 0),
			errmes: "",
		}, {
			lit:    toLiteral(nil),
			out:    nil,
			errmes: "",
		}, {
			lit:    toLiteral(complex(0, 0)),
			out:    nil,
			errmes: "sqlbuilder: got complex128 type, but literal is not supporting this.",
		}}

	for num, c := range cases {
		val, err := c.lit.(*cLiteralImpl).converted()
		if !reflect.DeepEqual(c.out, val) {
			t.Errorf("failed on %d", num)
		}
		if len(c.errmes) != 0 {
			if err == nil {
				t.Errorf("failed on %d", num)
			}
			if err.Error() != c.errmes {
				t.Errorf("failed on %d", num)
				panic(err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("failed on %d", num)
			}
		}
	}
}

func TestLiteralString(t *testing.T) {
	var cases = []struct {
		lit    literal
		out    string
		errmes string
	}{
		{
			lit:    toLiteral(int(10)),
			out:    "10",
			errmes: "",
		}, {
			lit:    toLiteral(int64(10)),
			out:    "10",
			errmes: "",
		}, {
			lit:    toLiteral(uint(10)),
			out:    "10",
			errmes: "",
		}, {
			lit:    toLiteral(uint64(10)),
			out:    "10",
			errmes: "",
		}, {
			lit:    toLiteral(float32(10)),
			out:    "10.0000000000",
			errmes: "",
		}, {
			lit:    toLiteral(float64(10)),
			out:    "10.0000000000",
			errmes: "",
		}, {
			lit:    toLiteral(bool(true)),
			out:    "true",
			errmes: "",
		}, {
			lit:    toLiteral([]byte{0x11}),
			out:    string([]byte{0x11}),
			errmes: "",
		}, {
			lit:    toLiteral(string("shibuya-rin")),
			out:    "shibuya-rin",
			errmes: "",
		}, {
			lit:    toLiteral(time.Unix(0, 0).UTC()),
			out:    "1970-01-01 00:00:00",
			errmes: "",
		}, {
			lit:    toLiteral(nil),
			out:    "NULL",
			errmes: "",
		}, {
			lit:    toLiteral(complex(0, 0)),
			out:    "",
			errmes: "aaa",
		}}

	for num, c := range cases {
		val := c.lit.(*cLiteralImpl).string()
		if c.out != val {
			t.Errorf("failed on %d", num)
		}
	}
}

func TestLiteralIsNil(t *testing.T) {
	var cases = []struct {
		in  literal
		out bool
	}{
		{toLiteral(int(10)), false},
		{toLiteral([]byte{}), false},
		{toLiteral(nil), true},
		{toLiteral([]byte(nil)), true},
	}

	for num, c := range cases {
		isnil := c.in.IsNil()
		if c.out != isnil {
			t.Errorf("failed on %d", num)
		}
	}
}
