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
