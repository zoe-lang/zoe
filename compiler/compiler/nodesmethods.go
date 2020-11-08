// Code generated by a lame .js file, DO NOT EDIT.

package zoe

import "io"


func (t *Token) CreateNamespace() *Namespace {
  res := &Namespace{}
  res.ExtendPosition(t)
  return res
}

func (r *Namespace) AddChildren(other ...Node) *Namespace {
  for _, c := range other {
    if fragment, ok := c.(*Fragment); ok {
      for _, c2 := range fragment.Children {
        r.Children = append(r.Children, c2)
      }
    } else {
      r.Children = append(r.Children, c)
    }
    r.ExtendPosition(c)
  }
  return r
}

func (r *Namespace) Dump(w io.Writer) {
  w.Write([]byte("(namespace"))
  for _, c := range r.Children {
    w.Write([]byte(" "))
    c.Dump(w)
  }
  w.Write([]byte(")"))
}

func (t *Token) CreateFragment() *Fragment {
  res := &Fragment{}
  res.ExtendPosition(t)
  return res
}

func (r *Fragment) AddChildren(other ...Node) *Fragment {
  for _, c := range other {
    if fragment, ok := c.(*Fragment); ok {
      for _, c2 := range fragment.Children {
        r.Children = append(r.Children, c2)
      }
    } else {
      r.Children = append(r.Children, c)
    }
    r.ExtendPosition(c)
  }
  return r
}

func (r *Fragment) Dump(w io.Writer) {
  w.Write([]byte("(fragment"))
  for _, c := range r.Children {
    w.Write([]byte(" "))
    c.Dump(w)
  }
  w.Write([]byte(")"))
}

func (t *Token) CreateTypeDecl() *TypeDecl {
  res := &TypeDecl{}
  res.ExtendPosition(t)
  return res
}

func (r *TypeDecl) SetTemplate(other *Template) *TypeDecl {
  r.Template = other
  r.ExtendPosition(other)
  return r
}

func (r *TypeDecl) SetIdent(other *Ident) *TypeDecl {
  r.Ident = other
  r.ExtendPosition(other)
  return r
}

func (r *TypeDecl) SetDef(other Node) *TypeDecl {
  r.Def = other
  r.ExtendPosition(other)
  return r
}

func (r *TypeDecl) Dump(w io.Writer) {
  w.Write([]byte("(type-decl"))
  w.Write([]byte(" "))
  r.Template.Dump(w)
  w.Write([]byte(" "))
  r.Ident.Dump(w)
  w.Write([]byte(" "))
  r.Def.Dump(w)
  w.Write([]byte(")"))
}

func (t *Token) CreateUnion() *Union {
  res := &Union{}
  res.ExtendPosition(t)
  return res
}

func (r *Union) AddTypeExprs(other ...Node) *Union {
  for _, c := range other {
    if fragment, ok := c.(*Fragment); ok {
      for _, c2 := range fragment.Children {
        r.TypeExprs = append(r.TypeExprs, c2)
      }
    } else {
      r.TypeExprs = append(r.TypeExprs, c)
    }
    r.ExtendPosition(c)
  }
  return r
}

func (r *Union) Dump(w io.Writer) {
  w.Write([]byte("(union"))
  for _, c := range r.TypeExprs {
    w.Write([]byte(" "))
    c.Dump(w)
  }
  w.Write([]byte(")"))
}

func (t *Token) CreateVar() *Var {
  res := &Var{}
  res.ExtendPosition(t)
  return res
}

func (r *Var) SetIdent(other *Ident) *Var {
  r.Ident = other
  r.ExtendPosition(other)
  return r
}

func (r *Var) SetTypeExp(other Node) *Var {
  r.TypeExp = other
  r.ExtendPosition(other)
  return r
}

func (r *Var) SetExp(other Node) *Var {
  r.Exp = other
  r.ExtendPosition(other)
  return r
}

func (r *Var) Dump(w io.Writer) {
  w.Write([]byte("(var"))
  w.Write([]byte(" "))
  r.Ident.Dump(w)
  w.Write([]byte(" "))
  r.TypeExp.Dump(w)
  w.Write([]byte(" "))
  r.Exp.Dump(w)
  w.Write([]byte(")"))
}

func (t *Token) CreateOperation() *Operation {
  res := &Operation{}
  res.ExtendPosition(t)
  return res
}

func (r *Operation) SetToken(other *Token) *Operation {
  r.Token = other
  r.ExtendPosition(other)
  return r
}

func (r *Operation) AddOperands(other ...Node) *Operation {
  for _, c := range other {
    if fragment, ok := c.(*Fragment); ok {
      for _, c2 := range fragment.Children {
        r.Operands = append(r.Operands, c2)
      }
    } else {
      r.Operands = append(r.Operands, c)
    }
    r.ExtendPosition(c)
  }
  return r
}

func (r *Operation) Dump(w io.Writer) {
  w.Write([]byte("(operation"))
  w.Write([]byte(" "))
  r.Token.Dump(w)
  for _, c := range r.Operands {
    w.Write([]byte(" "))
    c.Dump(w)
  }
  w.Write([]byte(")"))
}

func (t *Token) CreateTemplate() *Template {
  res := &Template{}
  res.ExtendPosition(t)
  return res
}

func (r *Template) AddArgs(other ...*Var) *Template {
  for _, c := range other {
    r.Args = append(r.Args, c)
    r.ExtendPosition(c)
  }
  return r
}

func (r *Template) Dump(w io.Writer) {
  w.Write([]byte("(template"))
  for _, c := range r.Args {
    w.Write([]byte(" "))
    c.Dump(w)
  }
  w.Write([]byte(")"))
}

func (t *Token) CreateFnDef() *FnDef {
  res := &FnDef{}
  res.ExtendPosition(t)
  return res
}

func (r *FnDef) SetTemplate(other *Template) *FnDef {
  r.Template = other
  r.ExtendPosition(other)
  return r
}

func (r *FnDef) SetSignature(other *FnSignature) *FnDef {
  r.Signature = other
  r.ExtendPosition(other)
  return r
}

func (r *FnDef) SetDefinition(other *Block) *FnDef {
  r.Definition = other
  r.ExtendPosition(other)
  return r
}

func (r *FnDef) Dump(w io.Writer) {
  w.Write([]byte("(fn-def"))
  w.Write([]byte(" "))
  r.Template.Dump(w)
  w.Write([]byte(" "))
  r.Signature.Dump(w)
  w.Write([]byte(" "))
  r.Definition.Dump(w)
  w.Write([]byte(")"))
}

func (t *Token) CreateFnSignature() *FnSignature {
  res := &FnSignature{}
  res.ExtendPosition(t)
  return res
}

func (r *FnSignature) AddArgs(other ...*Var) *FnSignature {
  for _, c := range other {
    r.Args = append(r.Args, c)
    r.ExtendPosition(c)
  }
  return r
}

func (r *FnSignature) SetReturnTypeExp(other Node) *FnSignature {
  r.ReturnTypeExp = other
  r.ExtendPosition(other)
  return r
}

func (r *FnSignature) Dump(w io.Writer) {
  w.Write([]byte("(fn-signature"))
  for _, c := range r.Args {
    w.Write([]byte(" "))
    c.Dump(w)
  }
  w.Write([]byte(" "))
  r.ReturnTypeExp.Dump(w)
  w.Write([]byte(")"))
}

func (t *Token) CreateFnCall() *FnCall {
  res := &FnCall{}
  res.ExtendPosition(t)
  return res
}

func (r *FnCall) SetLeft(other Node) *FnCall {
  r.Left = other
  r.ExtendPosition(other)
  return r
}

func (r *FnCall) SetArgs(other *Tuple) *FnCall {
  r.Args = other
  r.ExtendPosition(other)
  return r
}

func (r *FnCall) Dump(w io.Writer) {
  w.Write([]byte("(fn-call"))
  w.Write([]byte(" "))
  r.Left.Dump(w)
  w.Write([]byte(" "))
  r.Args.Dump(w)
  w.Write([]byte(")"))
}

func (t *Token) CreateGetIndex() *GetIndex {
  res := &GetIndex{}
  res.ExtendPosition(t)
  return res
}

func (r *GetIndex) SetLeft(other Node) *GetIndex {
  r.Left = other
  r.ExtendPosition(other)
  return r
}

func (r *GetIndex) SetIndex(other Node) *GetIndex {
  r.Index = other
  r.ExtendPosition(other)
  return r
}

func (r *GetIndex) Dump(w io.Writer) {
  w.Write([]byte("(get-index"))
  w.Write([]byte(" "))
  r.Left.Dump(w)
  w.Write([]byte(" "))
  r.Index.Dump(w)
  w.Write([]byte(")"))
}

func (t *Token) CreateSetIndex() *SetIndex {
  res := &SetIndex{}
  res.ExtendPosition(t)
  return res
}

func (r *SetIndex) SetLeft(other Node) *SetIndex {
  r.Left = other
  r.ExtendPosition(other)
  return r
}

func (r *SetIndex) SetIndex(other Node) *SetIndex {
  r.Index = other
  r.ExtendPosition(other)
  return r
}

func (r *SetIndex) SetValue(other Node) *SetIndex {
  r.Value = other
  r.ExtendPosition(other)
  return r
}

func (r *SetIndex) Dump(w io.Writer) {
  w.Write([]byte("(set-index"))
  w.Write([]byte(" "))
  r.Left.Dump(w)
  w.Write([]byte(" "))
  r.Index.Dump(w)
  w.Write([]byte(" "))
  r.Value.Dump(w)
  w.Write([]byte(")"))
}

func (t *Token) CreateIf() *If {
  res := &If{}
  res.ExtendPosition(t)
  return res
}

func (r *If) SetCond(other Node) *If {
  r.Cond = other
  r.ExtendPosition(other)
  return r
}

func (r *If) SetThen(other Node) *If {
  r.Then = other
  r.ExtendPosition(other)
  return r
}

func (r *If) SetElse(other Node) *If {
  r.Else = other
  r.ExtendPosition(other)
  return r
}

func (r *If) Dump(w io.Writer) {
  w.Write([]byte("(if"))
  w.Write([]byte(" "))
  r.Cond.Dump(w)
  w.Write([]byte(" "))
  r.Then.Dump(w)
  w.Write([]byte(" "))
  r.Else.Dump(w)
  w.Write([]byte(")"))
}

func (t *Token) CreateFnDecl() *FnDecl {
  res := &FnDecl{}
  res.ExtendPosition(t)
  return res
}

func (r *FnDecl) SetIdent(other *Ident) *FnDecl {
  r.Ident = other
  r.ExtendPosition(other)
  return r
}

func (r *FnDecl) SetFnDef(other *FnDef) *FnDecl {
  r.FnDef = other
  r.ExtendPosition(other)
  return r
}

func (r *FnDecl) Dump(w io.Writer) {
  w.Write([]byte("(fn-decl"))
  w.Write([]byte(" "))
  r.Ident.Dump(w)
  w.Write([]byte(" "))
  r.FnDef.Dump(w)
  w.Write([]byte(")"))
}

func (t *Token) CreateTuple() *Tuple {
  res := &Tuple{}
  res.ExtendPosition(t)
  return res
}

func (r *Tuple) AddChildren(other ...Node) *Tuple {
  for _, c := range other {
    if fragment, ok := c.(*Fragment); ok {
      for _, c2 := range fragment.Children {
        r.Children = append(r.Children, c2)
      }
    } else {
      r.Children = append(r.Children, c)
    }
    r.ExtendPosition(c)
  }
  return r
}

func (r *Tuple) Dump(w io.Writer) {
  w.Write([]byte("(tuple"))
  for _, c := range r.Children {
    w.Write([]byte(" "))
    c.Dump(w)
  }
  w.Write([]byte(")"))
}

func (t *Token) CreateBlock() *Block {
  res := &Block{}
  res.ExtendPosition(t)
  return res
}

func (r *Block) AddChildren(other ...Node) *Block {
  for _, c := range other {
    if fragment, ok := c.(*Fragment); ok {
      for _, c2 := range fragment.Children {
        r.Children = append(r.Children, c2)
      }
    } else {
      r.Children = append(r.Children, c)
    }
    r.ExtendPosition(c)
  }
  return r
}

func (r *Block) Dump(w io.Writer) {
  w.Write([]byte("(block"))
  for _, c := range r.Children {
    w.Write([]byte(" "))
    c.Dump(w)
  }
  w.Write([]byte(")"))
}

func (t *Token) CreateReturn() *Return {
  res := &Return{}
  res.ExtendPosition(t)
  return res
}

func (r *Return) SetExpr(other Node) *Return {
  r.Expr = other
  r.ExtendPosition(other)
  return r
}

func (r *Return) Dump(w io.Writer) {
  w.Write([]byte("(return"))
  w.Write([]byte(" "))
  r.Expr.Dump(w)
  w.Write([]byte(")"))
}

func (t *Token) CreateIdent() *Ident {
  res := &Ident{}
  res.ExtendPosition(t)
  return res
}

func (r *Ident) Dump(w io.Writer) {
  w.Write([]byte(r.GetText()))
}

func (t *Token) CreateNull() *Null {
  res := &Null{}
  res.ExtendPosition(t)
  return res
}

func (r *Null) Dump(w io.Writer) {
  w.Write([]byte(r.GetText()))
}

func (t *Token) CreateFalse() *False {
  res := &False{}
  res.ExtendPosition(t)
  return res
}

func (r *False) Dump(w io.Writer) {
  w.Write([]byte(r.GetText()))
}

func (t *Token) CreateTrue() *True {
  res := &True{}
  res.ExtendPosition(t)
  return res
}

func (r *True) Dump(w io.Writer) {
  w.Write([]byte(r.GetText()))
}

func (t *Token) CreateString() *String {
  res := &String{}
  res.ExtendPosition(t)
  return res
}

func (r *String) Dump(w io.Writer) {
  w.Write([]byte(r.GetText()))
}

func (t *Token) CreateInteger() *Integer {
  res := &Integer{}
  res.ExtendPosition(t)
  return res
}

func (r *Integer) Dump(w io.Writer) {
  w.Write([]byte(r.GetText()))
}

func (t *Token) CreateFloat() *Float {
  res := &Float{}
  res.ExtendPosition(t)
  return res
}

func (r *Float) Dump(w io.Writer) {
  w.Write([]byte(r.GetText()))
}

func (t *Token) CreateEof() *Eof {
  res := &Eof{}
  res.ExtendPosition(t)
  return res
}

func (r *Eof) Dump(w io.Writer) {
  w.Write([]byte(r.GetText()))
}
