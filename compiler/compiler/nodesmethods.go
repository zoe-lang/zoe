// Code generated by a lame .js file, DO NOT EDIT.

package zoe

import (
  "io"
  "bytes"
)


func (p *Position) CreateNamespace() *Namespace {
  res := &Namespace{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateNamespace() *Namespace {
  return tk.Position.CreateNamespace()
}

func (r *Namespace) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Namespace) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Namespace) Dump(w io.Writer) {

  w.Write([]byte("(namespace "))






      if r.Identifier != nil {
        r.Identifier.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Block != nil {
        r.Block.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}









func (r *Namespace) SetIdentifier(other Node) *Namespace {
  r.Identifier = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}







func (r *Namespace) EnsureBlock(fn func (b *Block)) *Namespace {
  if r.Block == nil {
    r.Block = &Block{}
  }
  fn(r.Block)
  r.ExtendPosition(r.Block)
  return r
}


func (r *Namespace) SetBlock(other *Block) *Namespace {
  r.Block = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateFragment() *Fragment {
  res := &Fragment{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateFragment() *Fragment {
  return tk.Position.CreateFragment()
}

func (r *Fragment) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Fragment) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Fragment) Dump(w io.Writer) {

  w.Write([]byte("(fragment "))






      for i, n := range r.Children {
        n.Dump(w)
        if i < len(r.Children) - 1 {
          w.Write([]byte(" "))
        }
      }




  w.Write([]byte(")"))

}






func (r *Fragment) AddChildren(other ...Node) *Fragment {
  for _, c := range other {
    if c != nil {
    
      switch v := c.(type) {
      case *Fragment:
        r.AddChildren(v.Children...)
      default:
        r.Children = append(r.Children, c)
        r.ExtendPosition(c)
      }
    
    }
  }
  return r
}





func (p *Position) CreateTypeDecl() *TypeDecl {
  res := &TypeDecl{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateTypeDecl() *TypeDecl {
  return tk.Position.CreateTypeDecl()
}

func (r *TypeDecl) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *TypeDecl) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *TypeDecl) Dump(w io.Writer) {

  w.Write([]byte("(typedecl "))






      if r.Ident != nil {
        r.Ident.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Template != nil {
        r.Template.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Def != nil {
        r.Def.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}








func (r *TypeDecl) EnsureIdent(fn func (i *BaseIdent)) *TypeDecl {
  if r.Ident == nil {
    r.Ident = &BaseIdent{}
  }
  fn(r.Ident)
  r.ExtendPosition(r.Ident)
  return r
}


func (r *TypeDecl) SetIdent(other *BaseIdent) *TypeDecl {
  r.Ident = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}







func (r *TypeDecl) EnsureTemplate(fn func (t *Template)) *TypeDecl {
  if r.Template == nil {
    r.Template = &Template{}
  }
  fn(r.Template)
  r.ExtendPosition(r.Template)
  return r
}


func (r *TypeDecl) SetTemplate(other *Template) *TypeDecl {
  r.Template = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}








func (r *TypeDecl) SetDef(other Node) *TypeDecl {
  r.Def = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateUnion() *Union {
  res := &Union{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateUnion() *Union {
  return tk.Position.CreateUnion()
}

func (r *Union) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Union) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Union) Dump(w io.Writer) {

  w.Write([]byte("(union "))






      for i, n := range r.TypeExprs {
        n.Dump(w)
        if i < len(r.TypeExprs) - 1 {
          w.Write([]byte(" "))
        }
      }




  w.Write([]byte(")"))

}






func (r *Union) AddTypeExprs(other ...Node) *Union {
  for _, c := range other {
    if c != nil {
    
      switch v := c.(type) {
      case *Fragment:
        r.AddTypeExprs(v.Children...)
      default:
        r.TypeExprs = append(r.TypeExprs, c)
        r.ExtendPosition(c)
      }
    
    }
  }
  return r
}





func (p *Position) CreateImportAs() *ImportAs {
  res := &ImportAs{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateImportAs() *ImportAs {
  return tk.Position.CreateImportAs()
}

func (r *ImportAs) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *ImportAs) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *ImportAs) Dump(w io.Writer) {

  w.Write([]byte("(importas "))






      if r.Path != nil {
        r.Path.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.As != nil {
        r.As.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}









func (r *ImportAs) SetPath(other Node) *ImportAs {
  r.Path = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}







func (r *ImportAs) EnsureAs(fn func (a *BaseIdent)) *ImportAs {
  if r.As == nil {
    r.As = &BaseIdent{}
  }
  fn(r.As)
  r.ExtendPosition(r.As)
  return r
}


func (r *ImportAs) SetAs(other *BaseIdent) *ImportAs {
  r.As = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateImportList() *ImportList {
  res := &ImportList{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateImportList() *ImportList {
  return tk.Position.CreateImportList()
}

func (r *ImportList) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *ImportList) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *ImportList) Dump(w io.Writer) {

  w.Write([]byte("(importlist "))






      if r.Path != nil {
        r.Path.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Names != nil {
        r.Names.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}









func (r *ImportList) SetPath(other Node) *ImportList {
  r.Path = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}







func (r *ImportList) EnsureNames(fn func (n *Tuple)) *ImportList {
  if r.Names == nil {
    r.Names = &Tuple{}
  }
  fn(r.Names)
  r.ExtendPosition(r.Names)
  return r
}


func (r *ImportList) SetNames(other *Tuple) *ImportList {
  r.Names = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateVar() *Var {
  res := &Var{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateVar() *Var {
  return tk.Position.CreateVar()
}

func (r *Var) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Var) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Var) Dump(w io.Writer) {

  w.Write([]byte("(var "))






      if r.Ident != nil {
        r.Ident.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.TypeExp != nil {
        r.TypeExp.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Exp != nil {
        r.Exp.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}








func (r *Var) EnsureIdent(fn func (i *BaseIdent)) *Var {
  if r.Ident == nil {
    r.Ident = &BaseIdent{}
  }
  fn(r.Ident)
  r.ExtendPosition(r.Ident)
  return r
}


func (r *Var) SetIdent(other *BaseIdent) *Var {
  r.Ident = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}








func (r *Var) SetTypeExp(other Node) *Var {
  r.TypeExp = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}








func (r *Var) SetExp(other Node) *Var {
  r.Exp = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateOperation() *Operation {
  res := &Operation{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateOperation() *Operation {
  return tk.Position.CreateOperation()
}

func (r *Operation) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Operation) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Operation) Dump(w io.Writer) {

  w.Write([]byte("("))






      if r.Token != nil {
        r.Token.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      for i, n := range r.Operands {
        n.Dump(w)
        if i < len(r.Operands) - 1 {
          w.Write([]byte(" "))
        }
      }




  w.Write([]byte(")"))

}









func (r *Operation) SetToken(other *Token) *Operation {
  r.Token = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (r *Operation) AddOperands(other ...Node) *Operation {
  for _, c := range other {
    if c != nil {
    
      switch v := c.(type) {
      case *Fragment:
        r.AddOperands(v.Children...)
      default:
        r.Operands = append(r.Operands, c)
        r.ExtendPosition(c)
      }
    
    }
  }
  return r
}





func (p *Position) CreateTemplate() *Template {
  res := &Template{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateTemplate() *Template {
  return tk.Position.CreateTemplate()
}

func (r *Template) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Template) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Template) Dump(w io.Writer) {

  w.Write([]byte("(template "))






      if r.Args != nil {
        r.Args.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}








func (r *Template) EnsureArgs(fn func (a *VarTuple)) *Template {
  if r.Args == nil {
    r.Args = &VarTuple{}
  }
  fn(r.Args)
  r.ExtendPosition(r.Args)
  return r
}


func (r *Template) SetArgs(other *VarTuple) *Template {
  r.Args = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateFnDef() *FnDef {
  res := &FnDef{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateFnDef() *FnDef {
  return tk.Position.CreateFnDef()
}

func (r *FnDef) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *FnDef) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *FnDef) Dump(w io.Writer) {

  w.Write([]byte("(fndef "))






      if r.Template != nil {
        r.Template.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Signature != nil {
        r.Signature.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Definition != nil {
        r.Definition.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}








func (r *FnDef) EnsureTemplate(fn func (t *Template)) *FnDef {
  if r.Template == nil {
    r.Template = &Template{}
  }
  fn(r.Template)
  r.ExtendPosition(r.Template)
  return r
}


func (r *FnDef) SetTemplate(other *Template) *FnDef {
  r.Template = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}







func (r *FnDef) EnsureSignature(fn func (s *Signature)) *FnDef {
  if r.Signature == nil {
    r.Signature = &Signature{}
  }
  fn(r.Signature)
  r.ExtendPosition(r.Signature)
  return r
}


func (r *FnDef) SetSignature(other *Signature) *FnDef {
  r.Signature = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}







func (r *FnDef) EnsureDefinition(fn func (d *Block)) *FnDef {
  if r.Definition == nil {
    r.Definition = &Block{}
  }
  fn(r.Definition)
  r.ExtendPosition(r.Definition)
  return r
}


func (r *FnDef) SetDefinition(other *Block) *FnDef {
  r.Definition = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateSignature() *Signature {
  res := &Signature{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateSignature() *Signature {
  return tk.Position.CreateSignature()
}

func (r *Signature) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Signature) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Signature) Dump(w io.Writer) {

  w.Write([]byte("(signature "))






      if r.Args != nil {
        r.Args.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.ReturnTypeExp != nil {
        r.ReturnTypeExp.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}








func (r *Signature) EnsureArgs(fn func (a *VarTuple)) *Signature {
  if r.Args == nil {
    r.Args = &VarTuple{}
  }
  fn(r.Args)
  r.ExtendPosition(r.Args)
  return r
}


func (r *Signature) SetArgs(other *VarTuple) *Signature {
  r.Args = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}








func (r *Signature) SetReturnTypeExp(other Node) *Signature {
  r.ReturnTypeExp = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateFnCall() *FnCall {
  res := &FnCall{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateFnCall() *FnCall {
  return tk.Position.CreateFnCall()
}

func (r *FnCall) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *FnCall) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *FnCall) Dump(w io.Writer) {

  w.Write([]byte("(fncall "))






      if r.Left != nil {
        r.Left.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Args != nil {
        r.Args.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}









func (r *FnCall) SetLeft(other Node) *FnCall {
  r.Left = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}







func (r *FnCall) EnsureArgs(fn func (a *Tuple)) *FnCall {
  if r.Args == nil {
    r.Args = &Tuple{}
  }
  fn(r.Args)
  r.ExtendPosition(r.Args)
  return r
}


func (r *FnCall) SetArgs(other *Tuple) *FnCall {
  r.Args = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateGetIndex() *GetIndex {
  res := &GetIndex{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateGetIndex() *GetIndex {
  return tk.Position.CreateGetIndex()
}

func (r *GetIndex) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *GetIndex) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *GetIndex) Dump(w io.Writer) {

  w.Write([]byte("(getindex "))






      if r.Left != nil {
        r.Left.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Index != nil {
        r.Index.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}









func (r *GetIndex) SetLeft(other Node) *GetIndex {
  r.Left = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}








func (r *GetIndex) SetIndex(other Node) *GetIndex {
  r.Index = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateSetIndex() *SetIndex {
  res := &SetIndex{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateSetIndex() *SetIndex {
  return tk.Position.CreateSetIndex()
}

func (r *SetIndex) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *SetIndex) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *SetIndex) Dump(w io.Writer) {

  w.Write([]byte("(setindex "))






      if r.Left != nil {
        r.Left.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Index != nil {
        r.Index.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Value != nil {
        r.Value.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}









func (r *SetIndex) SetLeft(other Node) *SetIndex {
  r.Left = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}








func (r *SetIndex) SetIndex(other Node) *SetIndex {
  r.Index = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}








func (r *SetIndex) SetValue(other Node) *SetIndex {
  r.Value = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateIf() *If {
  res := &If{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateIf() *If {
  return tk.Position.CreateIf()
}

func (r *If) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *If) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *If) Dump(w io.Writer) {

  w.Write([]byte("(if "))






      if r.Cond != nil {
        r.Cond.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Then != nil {
        r.Then.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.Else != nil {
        r.Else.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}









func (r *If) SetCond(other Node) *If {
  r.Cond = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}








func (r *If) SetThen(other Node) *If {
  r.Then = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}








func (r *If) SetElse(other Node) *If {
  r.Else = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateFnDecl() *FnDecl {
  res := &FnDecl{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateFnDecl() *FnDecl {
  return tk.Position.CreateFnDecl()
}

func (r *FnDecl) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *FnDecl) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *FnDecl) Dump(w io.Writer) {

  w.Write([]byte("(fndecl "))






      if r.Ident != nil {
        r.Ident.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }

      w.Write([]byte(" "))





      if r.FnDef != nil {
        r.FnDef.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}








func (r *FnDecl) EnsureIdent(fn func (i *BaseIdent)) *FnDecl {
  if r.Ident == nil {
    r.Ident = &BaseIdent{}
  }
  fn(r.Ident)
  r.ExtendPosition(r.Ident)
  return r
}


func (r *FnDecl) SetIdent(other *BaseIdent) *FnDecl {
  r.Ident = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}







func (r *FnDecl) EnsureFnDef(fn func (f *FnDef)) *FnDecl {
  if r.FnDef == nil {
    r.FnDef = &FnDef{}
  }
  fn(r.FnDef)
  r.ExtendPosition(r.FnDef)
  return r
}


func (r *FnDecl) SetFnDef(other *FnDef) *FnDecl {
  r.FnDef = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateTuple() *Tuple {
  res := &Tuple{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateTuple() *Tuple {
  return tk.Position.CreateTuple()
}

func (r *Tuple) EnsureTuple() *Tuple {

  return r

}

func (r *Tuple) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Tuple) Dump(w io.Writer) {

  w.Write([]byte("["))






      for i, n := range r.Children {
        n.Dump(w)
        if i < len(r.Children) - 1 {
          w.Write([]byte(" "))
        }
      }




  w.Write([]byte("]"))

}






func (r *Tuple) AddChildren(other ...Node) *Tuple {
  for _, c := range other {
    if c != nil {
    
      switch v := c.(type) {
      case *Fragment:
        r.AddChildren(v.Children...)
      default:
        r.Children = append(r.Children, c)
        r.ExtendPosition(c)
      }
    
    }
  }
  return r
}





func (p *Position) CreateVarTuple() *VarTuple {
  res := &VarTuple{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateVarTuple() *VarTuple {
  return tk.Position.CreateVarTuple()
}

func (r *VarTuple) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *VarTuple) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *VarTuple) Dump(w io.Writer) {

  w.Write([]byte("["))






      for i, n := range r.Vars {
        n.Dump(w)
        if i < len(r.Vars) - 1 {
          w.Write([]byte(" "))
        }
      }




  w.Write([]byte("]"))

}






func (r *VarTuple) AddVars(other ...*Var) *VarTuple {
  for _, c := range other {
    if c != nil {
    
      r.Vars = append(r.Vars, c)
      r.ExtendPosition(c)
    
    }
  }
  return r
}





func (p *Position) CreateBlock() *Block {
  res := &Block{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateBlock() *Block {
  return tk.Position.CreateBlock()
}

func (r *Block) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Block) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Block) Dump(w io.Writer) {

  w.Write([]byte("{"))






      for i, n := range r.Children {
        n.Dump(w)
        if i < len(r.Children) - 1 {
          w.Write([]byte(" "))
        }
      }




  w.Write([]byte("}"))

}






func (r *Block) AddChildren(other ...Node) *Block {
  for _, c := range other {
    if c != nil {
    
      switch v := c.(type) {
      case *Fragment:
        r.AddChildren(v.Children...)
      default:
        r.Children = append(r.Children, c)
        r.ExtendPosition(c)
      }
    
    }
  }
  return r
}





func (p *Position) CreateReturn() *Return {
  res := &Return{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateReturn() *Return {
  return tk.Position.CreateReturn()
}

func (r *Return) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Return) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Return) Dump(w io.Writer) {

  w.Write([]byte("(return "))






      if r.Expr != nil {
        r.Expr.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }





  w.Write([]byte(")"))

}









func (r *Return) SetExpr(other Node) *Return {
  r.Expr = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}





func (p *Position) CreateIdent() *Ident {
  res := &Ident{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateIdent() *Ident {
  return tk.Position.CreateIdent()
}

func (r *Ident) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Ident) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Ident) Dump(w io.Writer) {

  w.Write([]byte("(ident "))






      for i, n := range r.Path {
        n.Dump(w)
        if i < len(r.Path) - 1 {
          w.Write([]byte(" "))
        }
      }




  w.Write([]byte(")"))

}






func (r *Ident) AddPath(other ...*BaseIdent) *Ident {
  for _, c := range other {
    if c != nil {
    
      r.Path = append(r.Path, c)
      r.ExtendPosition(c)
    
    }
  }
  return r
}





func (p *Position) CreateBaseIdent() *BaseIdent {
  res := &BaseIdent{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateBaseIdent() *BaseIdent {
  return tk.Position.CreateBaseIdent()
}

func (r *BaseIdent) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *BaseIdent) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}


func (r *BaseIdent) Dump(w io.Writer) {
  w.Write([]byte(cyan(r.GetText())))
}





func (p *Position) CreateNull() *Null {
  res := &Null{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateNull() *Null {
  return tk.Position.CreateNull()
}

func (r *Null) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Null) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}


func (r *Null) Dump(w io.Writer) {
  w.Write([]byte(cyan(r.GetText())))
}





func (p *Position) CreateFalse() *False {
  res := &False{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateFalse() *False {
  return tk.Position.CreateFalse()
}

func (r *False) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *False) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}


func (r *False) Dump(w io.Writer) {
  w.Write([]byte(cyan(r.GetText())))
}





func (p *Position) CreateTrue() *True {
  res := &True{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateTrue() *True {
  return tk.Position.CreateTrue()
}

func (r *True) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *True) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}


func (r *True) Dump(w io.Writer) {
  w.Write([]byte(cyan(r.GetText())))
}





func (p *Position) CreateChar() *Char {
  res := &Char{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateChar() *Char {
  return tk.Position.CreateChar()
}

func (r *Char) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Char) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}


func (r *Char) Dump(w io.Writer) {
  w.Write([]byte(cyan(r.GetText())))
}





func (p *Position) CreateStr() *Str {
  res := &Str{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateStr() *Str {
  return tk.Position.CreateStr()
}

func (r *Str) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Str) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}



func (r *Str) Dump(w io.Writer) {

  w.Write([]byte("(str "))






      for i, n := range r.Children {
        n.Dump(w)
        if i < len(r.Children) - 1 {
          w.Write([]byte(" "))
        }
      }




  w.Write([]byte(")"))

}






func (r *Str) AddChildren(other ...Node) *Str {
  for _, c := range other {
    if c != nil {
    
      switch v := c.(type) {
      case *Fragment:
        r.AddChildren(v.Children...)
      default:
        r.Children = append(r.Children, c)
        r.ExtendPosition(c)
      }
    
    }
  }
  return r
}





func (p *Position) CreateString() *String {
  res := &String{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateString() *String {
  return tk.Position.CreateString()
}

func (r *String) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *String) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}


func (r *String) Dump(w io.Writer) {
  w.Write([]byte(cyan(r.GetText())))
}





func (p *Position) CreateInteger() *Integer {
  res := &Integer{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateInteger() *Integer {
  return tk.Position.CreateInteger()
}

func (r *Integer) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Integer) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}


func (r *Integer) Dump(w io.Writer) {
  w.Write([]byte(cyan(r.GetText())))
}





func (p *Position) CreateFloat() *Float {
  res := &Float{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateFloat() *Float {
  return tk.Position.CreateFloat()
}

func (r *Float) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Float) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}


func (r *Float) Dump(w io.Writer) {
  w.Write([]byte(cyan(r.GetText())))
}





func (p *Position) CreateEof() *Eof {
  res := &Eof{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) CreateEof() *Eof {
  return tk.Position.CreateEof()
}

func (r *Eof) EnsureTuple() *Tuple {

  res := &Tuple{}
  res.AddChildren(r)
  return res

}

func (r *Eof) DumpString() string {
  var res bytes.Buffer
  r.Dump(&res)
  return res.String()
}


func (r *Eof) Dump(w io.Writer) {
  w.Write([]byte(cyan(r.GetText())))
}




