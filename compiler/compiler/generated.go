// Code generated by a lame .js file, DO NOT EDIT.

package zoe



func (n Node) Repr() string {
  switch n.Kind() {

  case NODE_EMPTY: return grey("~")

  case NODE_FILE: return "file"

  case NODE_BLOCK: return "block"

  case NODE_TUPLE: return "tuple"

  case NODE_FN: return bblue("fn")

  case NODE_METHOD: return bblue("method")

  case NODE_TYPE: return bblue("type")

  case NODE_NAMESPACE: return bblue("namespace")

  case NODE_VAR: return bblue("var")

  case NODE_SIGNATURE: return "signature"

  case NODE_RETURN: return "return"

  case NODE_STRUCT: return bblue("struct")

  case NODE_UNION: return "union"

  case NODE_STRING: return "str"

  case NODE_ARRAY_LITERAL: return "array"

  case NODE_IF: return "if"

  case NODE_FOR: return "for"

  case NODE_WHILE: return "while"

  case NODE_IMPORT: return bblue("import")

  case NODE_UNA_ELLIPSIS: return "..."

  case NODE_UNA_PLUS: return "+"

  case NODE_UNA_PLUSPLUS: return "++"

  case NODE_UNA_MIN: return "-"

  case NODE_UNA_MINMIN: return "--"

  case NODE_UNA_NOT: return "!"

  case NODE_UNA_POINTER: return "ptr"

  case NODE_UNA_REF: return "ref"

  case NODE_UNA_BITNOT: return "~"

  case NODE_BIN_ASSIGN: return "="

  case NODE_BIN_PLUS: return "+"

  case NODE_BIN_MIN: return "-"

  case NODE_BIN_DIV: return "/"

  case NODE_BIN_MUL: return "*"

  case NODE_BIN_MOD: return "%"

  case NODE_BIN_EQ: return "=="

  case NODE_BIN_NEQ: return "!="

  case NODE_BIN_GTEQ: return ">="

  case NODE_BIN_GT: return ">"

  case NODE_BIN_LTEQ: return "<="

  case NODE_BIN_LT: return "<"

  case NODE_BIN_LSHIFT: return "<<"

  case NODE_BIN_RSHIFT: return ">>"

  case NODE_BIN_BITANDEQ: return "&="

  case NODE_BIN_BITAND: return "&"

  case NODE_BIN_BITOR: return "|"

  case NODE_BIN_BITXOR: return "^"

  case NODE_BIN_OR: return "||"

  case NODE_BIN_AND: return "&&"

  case NODE_BIN_IS: return "is"

  case NODE_BIN_CAST: return "cast"

  case NODE_BIN_CALL: return "call"

  case NODE_BIN_INDEX: return "index"

  case NODE_BIN_DOT: return "."

  case NODE_LIT_NULL: return mag("null")

  case NODE_LIT_VOID: return mag("void")

  case NODE_LIT_FALSE: return mag("false")

  case NODE_LIT_TRUE: return mag("true")

  case NODE_LIT_CHAR: return green(n.GetText())

  case NODE_LIT_RAWSTR: return green("'",n.GetText(),"'")

  case NODE_LIT_NUMBER: return mag(n.GetText())

  case NODE_ID: return cyan(n.InternedString())

  }
  return "<!!!>"
}

func (tk Tk) createFile(scope Scope, contents Node) Node {
  return tk.createNode(scope, NODE_FILE, contents)
}

func (tk Tk) createBlock(scope Scope, contents Node) Node {
  return tk.createNode(scope, NODE_BLOCK, contents)
}

func (tk Tk) createTuple(scope Scope, contents Node) Node {
  return tk.createNode(scope, NODE_TUPLE, contents)
}

func (tk Tk) createFn(scope Scope, name Node, signature Node, definition Node) Node {
  return tk.createNode(scope, NODE_FN, name, signature, definition)
}

func (tk Tk) createMethod(scope Scope, name Node, signature Node, definition Node) Node {
  return tk.createNode(scope, NODE_METHOD, name, signature, definition)
}

func (tk Tk) createType(scope Scope, name Node, template Node, typeexp Node) Node {
  return tk.createNode(scope, NODE_TYPE, name, template, typeexp)
}

func (tk Tk) createNamespace(scope Scope, name Node, block Node) Node {
  return tk.createNode(scope, NODE_NAMESPACE, name, block)
}

func (tk Tk) createVar(scope Scope, name Node, typeexp Node, assign Node) Node {
  return tk.createNode(scope, NODE_VAR, name, typeexp, assign)
}

func (tk Tk) createSignature(scope Scope, template Node, args Node, rettype Node) Node {
  return tk.createNode(scope, NODE_SIGNATURE, template, args, rettype)
}

func (tk Tk) createReturn(scope Scope, exp Node) Node {
  return tk.createNode(scope, NODE_RETURN, exp)
}

func (tk Tk) createStruct(scope Scope, varlist Node) Node {
  return tk.createNode(scope, NODE_STRUCT, varlist)
}

func (tk Tk) createUnion(scope Scope, members Node) Node {
  return tk.createNode(scope, NODE_UNION, members)
}

func (tk Tk) createString(scope Scope, contents Node) Node {
  return tk.createNode(scope, NODE_STRING, contents)
}

func (tk Tk) createArrayLiteral(scope Scope, contents Node) Node {
  return tk.createNode(scope, NODE_ARRAY_LITERAL, contents)
}

func (tk Tk) createIf(scope Scope, cond Node, thenarm Node, elsearm Node) Node {
  return tk.createNode(scope, NODE_IF, cond, thenarm, elsearm)
}

func (tk Tk) createFor(scope Scope, vardecl Node, rng Node, block Node) Node {
  return tk.createNode(scope, NODE_FOR, vardecl, rng, block)
}

func (tk Tk) createWhile(scope Scope, cond Node, block Node) Node {
  return tk.createNode(scope, NODE_WHILE, cond, block)
}

func (tk Tk) createImport(scope Scope, module Node, id Node, exp Node) Node {
  return tk.createNode(scope, NODE_IMPORT, module, id, exp)
}

func (tk Tk) createUnaPointer(scope Scope, pointed Node) Node {
  return tk.createNode(scope, NODE_UNA_POINTER, pointed)
}

func (tk Tk) createUnaRef(scope Scope, variable Node) Node {
  return tk.createNode(scope, NODE_UNA_REF, variable)
}


