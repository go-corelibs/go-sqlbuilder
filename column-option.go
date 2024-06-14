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
	"fmt"
	"strconv"
	"strings"
)

// ColumnOption represents option for a column. ex: primary key.
type ColumnOption struct {
	PrimaryKey    bool
	NotNull       bool
	Unique        bool
	AutoIncrement bool
	Size          int
	SqlType       string
	Default       interface{}
}

func (c ColumnOption) Describe() (output string) {
	var parts []string

	if c.PrimaryKey {
		parts = append(parts, "PrimaryKey")
	}

	if c.NotNull {
		parts = append(parts, "NotNull")
	}

	if c.Unique {
		parts = append(parts, "Unique")
	}

	if c.AutoIncrement {
		parts = append(parts, "AutoIncrement")
	}

	if c.Size > 0 {
		parts = append(parts, "Size("+strconv.Itoa(c.Size)+")")
	}

	if c.SqlType != "" {
		parts = append(parts, "SqlType("+strconv.Quote(c.SqlType)+")")
	}

	if c.Default != nil {
		parts = append(parts, fmt.Sprintf("Default(%q)", c.Default))
	}

	return strings.Join(parts, ", ")
}
