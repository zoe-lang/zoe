// Code generated by a lame .js file, DO NOT EDIT.

package zoe


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

func (t *Token) CreateIfThen() *IfThen {
  res := &IfThen{}
  res.ExtendPosition(t)
  return res
}

func (r *IfThen) SetCond(other Node) *IfThen {
  r.Cond = other
  r.ExtendPosition(other)
  return r
}

func (r *IfThen) SetThen(other Node) *IfThen {
  r.Then = other
  r.ExtendPosition(other)
  return r
}

func (r *IfThen) SetElse(other Node) *IfThen {
  r.Else = other
  r.ExtendPosition(other)
  return r
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

func (t *Token) CreateIdent() *Ident {
  res := &Ident{}
  res.ExtendPosition(t)
  return res
}

func (t *Token) CreateNull() *Null {
  res := &Null{}
  res.ExtendPosition(t)
  return res
}

func (t *Token) CreateFalse() *False {
  res := &False{}
  res.ExtendPosition(t)
  return res
}

func (t *Token) CreateTrue() *True {
  res := &True{}
  res.ExtendPosition(t)
  return res
}

func (t *Token) CreateString() *String {
  res := &String{}
  res.ExtendPosition(t)
  return res
}

func (t *Token) CreateInteger() *Integer {
  res := &Integer{}
  res.ExtendPosition(t)
  return res
}

func (t *Token) CreateFloat() *Float {
  res := &Float{}
  res.ExtendPosition(t)
  return res
}