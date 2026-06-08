package poryscriptz

import (
	"fmt"

	gts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/taproot"
	"m31labs.dev/poryscriptz/target"
	"m31labs.dev/poryscriptz/vocab"
)

// Diag is a resolver diagnostic (error or warning) with source position.
type Diag struct {
	Line, Col int
	Msg       string
}

// Program is the resolved IR produced by Resolve.
type Program struct{ Scripts []*Script }

// Script is a resolved script declaration.
type Script struct {
	Name  string
	Stmts []Stmt
}

// Stmt is a resolved statement (sealed interface).
type Stmt interface{ isStmt() }

// Call is a resolved command-call statement.
type Call struct {
	Cmd  *vocab.Command
	Args []string // raw arg source text; formatted at emit time
}

func (Call) isStmt() {}

// Resolve walks each script_declaration in the CST, type-checks command calls
// against the target vocabulary, and returns the Program plus collected
// diagnostics (collect-and-continue; does not stop at the first error).
func Resolve(w *taproot.Walker, root *gts.Node, tgt target.GameTarget) (*Program, []Diag) {
	prog := &Program{}
	var diags []Diag

	var walk func(n *gts.Node)
	walk = func(n *gts.Node) {
		if n == nil {
			return
		}
		if w.Type(n) == "script_declaration" {
			s, ds := resolveScript(w, n, tgt)
			prog.Scripts = append(prog.Scripts, s)
			diags = append(diags, ds...)
			return // don't recurse into the script — resolveScript handles it
		}
		for i := 0; i < n.ChildCount(); i++ {
			walk(n.Child(i))
		}
	}
	walk(root)
	return prog, diags
}

// resolveScript resolves a single script_declaration node.
func resolveScript(w *taproot.Walker, n *gts.Node, tgt target.GameTarget) (*Script, []Diag) {
	nameNode := w.Field(n, "name")
	name := w.Text(nameNode)
	s := &Script{Name: name}
	var diags []Diag

	// The block is the 3rd child (index 2): script <name> <block>
	block := w.ChildByType(n, "block")
	if block == nil {
		return s, diags
	}

	// Inside the block there is a statement_list child.
	stmtList := w.ChildByType(block, "statement_list")
	if stmtList == nil {
		return s, diags
	}

	// Iterate over expression_statement children of statement_list.
	for i := 0; i < stmtList.ChildCount(); i++ {
		child := stmtList.Child(i)
		if w.Type(child) != "expression_statement" {
			continue
		}
		// The call_expression is the only child of expression_statement.
		callExpr := w.ChildByType(child, "call_expression")
		if callExpr == nil {
			continue
		}
		stmt, ds := resolveCall(w, callExpr, tgt)
		diags = append(diags, ds...)
		if stmt != nil {
			s.Stmts = append(s.Stmts, stmt)
		}
	}
	return s, diags
}

// resolveCall type-checks a single call_expression node against the vocabulary.
func resolveCall(w *taproot.Walker, callExpr *gts.Node, tgt target.GameTarget) (Stmt, []Diag) {
	var diags []Diag

	funcNode := w.Field(callExpr, "function")
	name := w.Text(funcNode)
	line, col := w.Pos(funcNode)

	cmd, ok := tgt.Vocabulary().ByName(name)
	if !ok {
		diags = append(diags, Diag{
			Line: line,
			Col:  col,
			Msg:  fmt.Sprintf("command %q not in HGSS vocabulary", name),
		})
		return nil, diags
	}

	argList := w.Field(callExpr, "arguments")
	gotArgCount := 0
	var rawArgs []string
	if argList != nil {
		gotArgCount = argList.NamedChildCount()
		for i := 0; i < gotArgCount; i++ {
			argNode := argList.NamedChild(i)
			rawArgs = append(rawArgs, w.Text(argNode))
		}
	}

	if gotArgCount != len(cmd.Args) {
		diags = append(diags, Diag{
			Line: line,
			Col:  col,
			Msg:  fmt.Sprintf("%s expects %d args, got %d", name, len(cmd.Args), gotArgCount),
		})
		return nil, diags
	}

	return Call{Cmd: cmd, Args: rawArgs}, diags
}
