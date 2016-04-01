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

package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template/parse"

	"github.com/seehuhn/gotmpljs/lib"
)

var OutputFileName string
var TmplNamespace string

const JsTmplEscapes = "seehuhn.gotmpl"

type Parsival struct {
	functions    map[string]bool
	hasVariables bool
	loopIndex    int
	prefix       string
	lines        []string

	hasEscapes bool
	completed  []string
}

func (p *Parsival) write(format string, args ...interface{}) {
	p.lines = append(p.lines, fmt.Sprintf(p.prefix+format, args...))
}

func (p *Parsival) String() string {
	return strings.Join(p.lines, "")
}

func (p *Parsival) indent() {
	p.prefix += "  "
}

func (p *Parsival) outdent() {
	p.prefix = p.prefix[:len(p.prefix)-2]
}

func (p *Parsival) evalField(dot string, fieldName string, args []parse.Node,
	final, receiver string) string {
	hasArgs := len(args) > 1 || final != ""
	res := receiver + "['" + fieldName + "']"
	if hasArgs {
		return p.evalCall(dot, res, args, final)
	}
	return res
}

func (p *Parsival) evalFieldChain(dot, receiver string, node parse.Node,
	ident []string, args []parse.Node, final string) string {
	n := len(ident)
	for i := 0; i < n-1; i++ {
		receiver = p.evalField(dot, ident[i], nil, "", receiver)
	}
	return p.evalField(dot, ident[n-1], args, final, receiver)
}

func (p *Parsival) evalFieldNode(dot string, field *parse.FieldNode,
	args []parse.Node, final string) string {
	return p.evalFieldChain(dot, dot, field, field.Ident, args, final)
}

func (p *Parsival) evalArg(dot string, n parse.Node) string {
	switch arg := n.(type) {
	case *parse.DotNode:
		return dot
	// case *parse.NilNode:
	//	if canBeNil(typ) {
	//		return reflect.Zero(typ)
	//	}
	//	s.errorf("cannot assign nil to %s", typ)
	case *parse.FieldNode:
		return p.evalFieldNode(dot, arg, []parse.Node{n}, "")
	// case *parse.VariableNode:
	//	return s.validateType(s.evalVariableNode(dot, arg, nil, zero), typ)
	// case *parse.PipeNode:
	//	return s.validateType(s.evalPipeline(dot, arg), typ)
	// case *parse.IdentifierNode:
	//	return s.evalFunction(dot, arg, arg, nil, zero)
	case *parse.StringNode:
		return lib.JsQuote(arg.Text)
	default:
		fmt.Printf("not implemented %#v\n", arg)
		return arg.String()
	}
}

func (p *Parsival) evalCall(dot string, name string, args []parse.Node,
	final string) string {
	args = args[1:]
	var argv []string
	for i := 0; i < len(args); i++ {
		arg := p.evalArg(dot, args[i])
		argv = append(argv, arg)
	}
	if final != "" {
		argv = append(argv, final)
	}
	var jsFunc string
	escapePrefix := "html_template_"
	if strings.HasPrefix(name, escapePrefix) {
		p.hasEscapes = true
		jsFunc = JsTmplEscapes + "." + name[len(escapePrefix):]
	} else {
		p.functions[name] = true
		jsFunc = "functions." + name
	}
	return jsFunc + "(" + strings.Join(argv, ", ") + ")"
}

func (p *Parsival) evalCommand(dot string, cmd *parse.CommandNode, final string) string {
	firstWord := cmd.Args[0]
	switch word := firstWord.(type) {
	case *parse.FieldNode:
		return p.evalFieldNode(dot, word, cmd.Args, final)
	case *parse.ChainNode:
		// return s.evalChainNode(dot, word, cmd.Args, final)
		p.write("#cmd %#v\n", cmd)
		return "???"
	case *parse.IdentifierNode:
		name := word.Ident
		return p.evalCall(dot, name, cmd.Args, final)
	case *parse.PipeNode:
		return p.evalPipeline(dot, word)
	case *parse.VariableNode:
		// return s.evalVariableNode(dot, word, cmd.Args, final)
		p.write("cmd %#v\n", cmd)
		return "???"
	case *parse.BoolNode:
		// return reflect.ValueOf(word.True)
		p.write("cmd %#v\n", cmd)
		return "???"
	case *parse.DotNode:
		return dot
	case *parse.NilNode:
		// s.errorf("nil is not a command")
		p.write("cmd %#v\n", cmd)
		return "???"
	case *parse.NumberNode:
		return word.Text
	case *parse.StringNode:
		// return reflect.ValueOf(word.Text)
		p.write("cmd %#v\n", cmd)
		return "???"
	}
	panic("not reached!")
}

func (p *Parsival) evalPipeline(dot string, pipe *parse.PipeNode) string {
	var res string
	for _, cmd := range pipe.Cmds {
		res = p.evalCommand(dot, cmd, res)
	}
	for _, variable := range pipe.Decl {
		p.hasVariables = true
		p.write("vars['%s'] = %s;\n", variable.Ident[0][1:], res)
	}
	return res
}

func (p *Parsival) walkRange(dot string, r *parse.RangeNode) {
	val := p.evalPipeline(dot, r.Pipe)
	p.loopIndex++
	loopVar := "i"
	if p.loopIndex > 1 {
		loopVar += strconv.Itoa(p.loopIndex)
	}
	p.write("for (var %s = 0; %s < %s.length; %s++) {\n",
		loopVar, loopVar, val, loopVar)
	p.indent()
	elem := val + "[" + loopVar + "]"
	if len(r.Pipe.Decl) > 0 {
		p.write("vars.??? = %s;\n", elem)
	}
	if len(r.Pipe.Decl) > 1 {
		p.write("vars.??? = %s;\n", loopVar)
	}
	p.walk(elem, r.List)
	p.outdent()
	p.write("}\n")
	p.loopIndex--
}

func (p *Parsival) walkIfOrWith(typ parse.NodeType, dot string,
	pipe *parse.PipeNode, list, elseList *parse.ListNode) {
	val := p.evalPipeline(dot, pipe)
	p.write("if (%s) {\n", val)
	p.indent()
	if typ == parse.NodeWith {
		p.walk(val, list)
	} else {
		p.walk(dot, list)
	}
	if elseList != nil {
		p.outdent()
		p.write("} else {\n")
		p.indent()
		p.walk(dot, elseList)
	}
	p.outdent()
	p.write("}\n")
}

func (p *Parsival) walk(dot string, node parse.Node) {
	switch node := node.(type) {
	case *parse.ActionNode:
		res := p.evalPipeline(dot, node.Pipe)
		if len(node.Pipe.Decl) == 0 {
			p.write("res.push(%s);\n", res)
		}
	case *parse.IfNode:
		p.walkIfOrWith(parse.NodeIf, dot, node.Pipe, node.List,
			node.ElseList)
	case *parse.ListNode:
		for _, node := range node.Nodes {
			p.walk(dot, node)
		}
	case *parse.RangeNode:
		p.walkRange(dot, node)
	case *parse.TemplateNode:
		p.write("TemplateNode: %#v\n", node)
	case *parse.TextNode:
		p.write("res.push(%s);\n", lib.JsQuote(string(node.Text)))
	case *parse.WithNode:
		p.walkIfOrWith(parse.NodeWith, dot, node.Pipe, node.List,
			node.ElseList)
	default:
		p.write("error, unknown node: %#v\n", node)
	}
}

func (p *Parsival) processTemplate(t *template.Template) {
	p.functions = map[string]bool{}
	p.hasVariables = false
	p.loopIndex = 0
	p.lines = nil
	dot := "data"

	// Delay writing the JavaScript function header until we know
	// which template functions are used.
	p.indent()
	p.write("var res = new Array();\n")
	p.walk(dot, t.Tree.Root)
	p.write("return res.join('');\n")
	p.outdent()
	p.write("};\n")

	head := []string{
		"\n\n", // each JavaScript function is preceeded by two empty lines
	}
	head = append(head, "/**\n * Execute the \""+t.Name()+"\" template.\n")
	head = append(head, " * @param {*} "+dot+
		" The data to apply the template to.\n")
	args := dot
	if len(p.functions) > 0 {
		var names []string
		for name := range p.functions {
			names = append(names, name)
		}
		sort.Strings(names)
		head = append(head, " * @param {{"+strings.Join(names, ", ")+
			"}} functions Map names to template functions.\n")
		args += ", functions"
	}
	head = append(head, " * @return {string} The template output.\n")
	head = append(head, " */\n")
	head = append(head, TmplNamespace+"."+t.Name()+" = function("+args+") {\n")
	if p.hasVariables {
		head = append(head, "  var vars = {};")
	}
	jsFunc := strings.Join(head, "") + p.String()
	p.completed = append(p.completed, jsFunc)
}

func main() {
	usage := "filename for the JavaScript output"
	flag.StringVar(&OutputFileName, "output", "", usage)
	flag.StringVar(&OutputFileName, "o", "", usage+" (shorthand)")
	usage = "javascript namespace for the template functions"
	flag.StringVar(&TmplNamespace, "namespace", "", usage)
	flag.StringVar(&TmplNamespace, "n", "", usage+" (shorthand)")
	flag.Parse()

	if TmplNamespace == "" {
		log.Fatal("namespace not set (use -n)")
	}

	p := &Parsival{}

	for _, inputName := range flag.Args() {
		tmplBody, err := ioutil.ReadFile(inputName)
		if err != nil {
			log.Fatalf("cannot read %q: %s", inputName, err.Error())
		}

		tmplName := filepath.Base(inputName)
		for i, c := range tmplName {
			if c == '.' && i > 0 {
				tmplName = tmplName[:i]
				break
			}
		}

		t, err := lib.ParseAndEscapeTemplate(tmplName, string(tmplBody))
		if err != nil {
			log.Fatalf("cannot parse template %q: %s", tmplName, err.Error())
		}

		p.processTemplate(t)
	}

	var err error
	var out *os.File
	var header string
	if OutputFileName != "" {
		out, err = os.Create(OutputFileName)
		if err != nil {
			log.Fatalf("cannot open %q: %s", OutputFileName, err.Error())
		}
		defer out.Close()
		header = filepath.Base(OutputFileName) + " - "
	} else {
		out = os.Stdout
	}

	_, err = out.WriteString("// " + header +
		"generated by github.com/seehuhn/gotmpljs, do not edit\n")
	if err != nil {
		log.Fatalf("cannot write %q: %s", out.Name(), err.Error())
	}
	_, err = out.WriteString("\ngoog.provide('" + TmplNamespace + "');\n")
	if err != nil {
		log.Fatalf("cannot write %q: %s", out.Name(), err.Error())
	}
	if p.hasEscapes {
		_, err := out.WriteString("\ngoog.require('" + JsTmplEscapes + "');\n")
		if err != nil {
			log.Fatalf("cannot write %q: %s", out.Name(), err.Error())
		}
	}
	_, err = out.WriteString(strings.Join(p.completed, ""))
	if err != nil {
		log.Fatalf("cannot write %q: %s", out.Name(), err.Error())
	}
}
