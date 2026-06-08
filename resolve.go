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

// If is a resolved if-statement with a flag or var-equality condition.
// Kind is "flag" or "vareq"; the corresponding fields (Flag / Var+Value) are set.
type If struct {
	Kind  string // "flag" or "vareq"
	Flag  string // Kind=="flag": identifier text, e.g. "FLAG_X"
	Var   string // Kind=="vareq": var identifier text, e.g. "VAR_X"
	Value string // Kind=="vareq": literal text, e.g. "3"
	Body  []Stmt
}

func (If) isStmt() {}

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

	stmts, ds := resolveStmtList(w, stmtList, tgt)
	diags = append(diags, ds...)
	s.Stmts = append(s.Stmts, stmts...)
	return s, diags
}

// resolveStmtList resolves all statements in a statement_list node.
func resolveStmtList(w *taproot.Walker, stmtList *gts.Node, tgt target.GameTarget) ([]Stmt, []Diag) {
	var stmts []Stmt
	var diags []Diag

	for i := 0; i < stmtList.ChildCount(); i++ {
		child := stmtList.Child(i)
		switch w.Type(child) {
		case "expression_statement":
			// The call_expression is the only child of expression_statement.
			callExpr := w.ChildByType(child, "call_expression")
			if callExpr == nil {
				continue
			}
			stmt, ds := resolveCall(w, callExpr, tgt)
			diags = append(diags, ds...)
			if stmt != nil {
				stmts = append(stmts, stmt)
			}
		case "if_statement":
			stmt, ds := resolveIf(w, child, tgt)
			diags = append(diags, ds...)
			if stmt != nil {
				stmts = append(stmts, stmt)
			}
		}
	}
	return stmts, diags
}

// resolveIf resolves an if_statement node into an If statement.
//
// Recognised condition shapes:
//   - flag(<IDENT>)                   → If{Kind:"flag", Flag:<ident>}
//   - var(<IDENT>) == <int_literal>   → If{Kind:"vareq", Var:<ident>, Value:<lit>}
//
// Any other shape emits an "unsupported if-condition" diagnostic and returns nil.
func resolveIf(w *taproot.Walker, ifNode *gts.Node, tgt target.GameTarget) (Stmt, []Diag) {
	var diags []Diag

	condNode := w.Field(ifNode, "condition")
	conseqNode := w.Field(ifNode, "consequence")

	line, col := w.Pos(condNode)

	var result If

	switch w.Type(condNode) {
	case "call_expression":
		// flag(<IDENT>) condition
		funcNode := w.Field(condNode, "function")
		if w.Text(funcNode) != "flag" {
			diags = append(diags, Diag{Line: line, Col: col, Msg: "unsupported if-condition: only flag() and var()==N supported"})
			return nil, diags
		}
		argList := w.Field(condNode, "arguments")
		if argList == nil || argList.NamedChildCount() != 1 {
			diags = append(diags, Diag{Line: line, Col: col, Msg: "unsupported if-condition: flag() requires exactly one argument"})
			return nil, diags
		}
		flagIdent := argList.NamedChild(0)
		result = If{Kind: "flag", Flag: w.Text(flagIdent)}

	case "binary_expression":
		// var(<IDENT>) == <int_literal> condition
		leftNode := w.Field(condNode, "left")
		rightNode := w.Field(condNode, "right")
		opNode := w.Field(condNode, "operator")

		if w.Text(opNode) != "==" {
			diags = append(diags, Diag{Line: line, Col: col, Msg: "unsupported if-condition: only == operator supported"})
			return nil, diags
		}
		if w.Type(leftNode) != "call_expression" {
			diags = append(diags, Diag{Line: line, Col: col, Msg: "unsupported if-condition: left side must be var()"})
			return nil, diags
		}
		funcNode := w.Field(leftNode, "function")
		if w.Text(funcNode) != "var" {
			diags = append(diags, Diag{Line: line, Col: col, Msg: "unsupported if-condition: left call must be var()"})
			return nil, diags
		}
		argList := w.Field(leftNode, "arguments")
		if argList == nil || argList.NamedChildCount() != 1 {
			diags = append(diags, Diag{Line: line, Col: col, Msg: "unsupported if-condition: var() requires exactly one argument"})
			return nil, diags
		}
		varIdent := argList.NamedChild(0)
		if w.Type(rightNode) != "int_literal" {
			diags = append(diags, Diag{Line: line, Col: col, Msg: "unsupported if-condition: right side must be int literal"})
			return nil, diags
		}
		result = If{Kind: "vareq", Var: w.Text(varIdent), Value: w.Text(rightNode)}

	default:
		diags = append(diags, Diag{Line: line, Col: col, Msg: fmt.Sprintf("unsupported if-condition: %s", w.Type(condNode))})
		return nil, diags
	}

	// Resolve the consequence block's statement list.
	if conseqNode != nil {
		bodyStmtList := w.ChildByType(conseqNode, "statement_list")
		if bodyStmtList != nil {
			bodyStmts, ds := resolveStmtList(w, bodyStmtList, tgt)
			diags = append(diags, ds...)
			result.Body = bodyStmts
		}
	}

	return result, diags
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
