// Convert Go templates to JavaScript.
// Copyright (C) 2014  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lib

import (
	"testing"
)

func TestJsQuote(t *testing.T) {
	data := []struct{ in, expected string }{
		{"this is a test", "'this is a test'"},
		{"", "''"},
		{"\000", "'\\x00'"},
		{"natürlich", "'natürlich'"},
		{"she said \"hello\"", "'she said \"hello\"'"},
		{"she said 'hello'", "'she said \\'hello\\''"},
		{"\\hello", "'\\\\hello'"},
		{"運", "'運'"},
		{"\a\b\f\n\r\t\v", "'\\x07\\b\\f\\n\\r\\t\\v'"},
		{"\U00002028", "'\\u2028'"},
		{"\U00101234", "'\\ufffd'"},
		{string([]byte{0x81}), "'\\x81'"},
	}

	for _, run := range data {
		out := JsQuote(run.in)
		if out != run.expected {
			t.Errorf("expected %s, got %s", run.expected, out)
		}
	}
}
