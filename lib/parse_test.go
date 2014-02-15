// Unit tests for the file "lib/parse.go"
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
	"fmt"
	"testing"
	"text/template/parse"
)

func TestParseAndEscapeTemplate(t *testing.T) {
	res, err := ParseAndEscapeTemplate("test", `{{.Name`)
	if res != nil || err == nil {
		t.Error("template error not detected")
	}

	res, err = ParseAndEscapeTemplate("test", `<p>hi {{F .Name}}`)
	if err != nil {
		t.Fatal(err)
	}
	root := res.Tree.Root
	actionNode := root.Nodes[1].(*parse.ActionNode) // .Nodes[0] is the "<p>hi"
	if len(actionNode.Pipe.Cmds) != 2 {
		t.Fatal("template not escaped")
	}
	for _, n := range actionNode.Pipe.Cmds {
		fmt.Printf("%#v\n", n.Args[0])
	}
	cmdNode := actionNode.Pipe.Cmds[0] // .Cmds[0] is "F(.Name)"
	idNode := cmdNode.Args[0].(*parse.IdentifierNode)
	if idNode.Ident != "F" {
		t.Error("function F not called")
	}
	cmdNode = actionNode.Pipe.Cmds[1] // .Cmds[1] is HTML escaping
	idNode = cmdNode.Args[0].(*parse.IdentifierNode)
	if idNode.Ident != "html_template_htmlescaper" {
		t.Error("template not escaped")
	}
}
