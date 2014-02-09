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
	"html/template"
	"io/ioutil"
	"regexp"
)

// ParseAndEscapeTemplate parses a Go template and executes the
// template once to make sure that the escaping functions from the
// html/template module are installed.
func ParseAndEscapeTemplate(name, body string) (*template.Template, error) {
	// Install the template we are interested in as an associated
	// template for a trivial, outer template.  This way we are sure that
	// we can execute the template without causing a panic.

	// First create the outer template.
	outerTmplText := `{{if .}}{{template "` + name + `" .}}{{end}}`
	outerTmpl, err := template.New("OUTER").Parse(outerTmplText)
	if err != nil {
		panic(err)
	}

	// Attach the real template, so that it is escaped together with
	// the outerTmpl.  In case the template relies on template
	// functions, install dummy functions with the correct names as
	// needed.  Is there a less ugly way to do this?
	pat, _ := regexp.Compile("^.* function \"([^ ]+)\" not defined$")
retry:
	innerTmpl, err := outerTmpl.New(name).Parse(body)
	if err != nil {
		match := pat.FindStringSubmatch(err.Error())
		if len(match) == 2 {
			funcName := match[1]
			funcMap := template.FuncMap{}
			funcMap[funcName] = func() string { return "" }
			outerTmpl = outerTmpl.Funcs(funcMap)
			goto retry
		}
		return nil, err
	}

	// Execute outerTmpl.  This triggers the code in html/template
	// which modifies both outerTmpl and the associated innerTmpl, by
	// installing the required escaping functions.
	outerTmpl.Execute(ioutil.Discard, nil)

	// Finally, return the sub-template, as modified by the .Execute()
	// call above.
	return innerTmpl, nil
}
