
# Specifications of the AST outputed by Zoe

Zoe parses a file into a simplified AST that can be represented into an unambiguous tree.
The type in the source code representing the AST is `Node`. The parser will try to emit
a valid tree, but the typechecker does further check that everything makes sense.

The goal of this document is to detail what are all the forms that the ast can take, using
its representation.

This mostly "in-memory" representation is a middle step before a last transformation pass
that converts all operations to an SSA form more apt to be later compiled by LLVM, once
typechecking validated that everything made sense.

Any node can potentially be in an error state, in which case it must not be processed by the type
checker.

Once validated, the produced trees can be interpreted and run at compile time, provided there
is no syscall done (except through compiler intrinsics), no external function called and no
unsafe pointer shenanigans.

# Types

All different types are `node`s in the AST.

- Identifiers
- Numbers
- The booleans `true` and `false`
- `null`
- strings enclosed in single quotes
- `(identifier nodes...)` an operation
- `{ node node ... }` a block of instructions
- `[ ]` a list of nodes

# "Grammar"

```bash

top-level = {...exp}

declaration = fn | type | var | namespace | import

import = (import )
namespace = (namespace block)
var = (var identifier type:exp? value:exp?) -- id, type, default value
type = (type identifier template? value:exp) -- exp must resolve to a valid type

# A fn with a nil block is just a signature/function pointer declaration
#   - it is only valid in a type expression
# A fn without identifier is a lambda expression and is only valid as a function
# value, or as a type.
# Function pointer declarations may have a name (?)
fn = (fn identifier? signature block?)

# template variables must be prefixed by '$'
# args names must always be lower-case
# return type is optional only in lambda or inline functions as it can be infered
signature = (signature template:[...var] args:[...var] return-type:exp?)


template = [...var] # vars

block = {...exp}

typedef = ( type identifier struct-def | union-def | exp ) -- exp must refer to a valid type or a valid constant
struct = ( struct [...var] )
union-def =  ( union [...exp] ) -- exp must resolve to a valid type or a valid constant

exp =
  | identifier
  | declaration
  | number
  | raw-string
  | character
  | 'false' | 'true'
  | 'null'
  | '*'
  | (
      | =
      | + | - | * | / | %
      | @ | . | ::
      | == | != | ! | < | > | <= | >= | && | '||'
      | & | '|'
      exp
      exp
    )
  | ( +|- exp )
  | ( get-index exp exp ) -- exp1 will be indexed, exp2 is the index
  | ( set-index exp exp exp ) -- exp1 will be indexed, exp2 is the index, exp3 the new value
  | ( call exp [...exp] )
```

# Fragment

The compiler may emit a node of type `fragment`, which will only stay in the AST in the case
of an incorrect tree ; its only objective is to be merged in its enclosing block.

They are typically emitted by a `var` with identifiers separated by commas.

# Expressions

Expressions can be any of these nodes. Syntactic sugar nodes such as nullish-coalescing operators,
`+=`, `++` and the likes are expanded into their final form in the AST.

- `(<op> exp exp)`, where op is `+`, `-`, `*`, `/`, `=`, `.`, ...
- `(call exp [exp...])` the function call
- `(get-index exp exp)` which is the array/map index operator
- `(set-index exp exp exp)` index set
- `(infer)` a token to tell the compiler to guess the value

# Function calls

For instance :
- `(call (. obj add) [43, '22'])` is `obj.add(43, '22')`

# Type declaration

These are the result of `type ... is` statements. They can be enclosed `template` blocks.

- `(decl:type <identifier> <typedef>)`

# Type Definition

A type can be any of those :

- `(union [type+])`
- `(enum [member+])`
  - member: `(= identifier exp) | identifier`
- `(struct [field field...])`
  - field: `(: identifier type exp?)`
- `null`

## Infer

This node is special. It tells the compiler that it will need to guess what was the type based
on for instance the block definition of a function, or the expected callback type at the call
site where a function is called.

- `(infer)`

# Variable declaration

The variable declaration can take several forms.

- `(decl:var identifier type exp?)`

# Function declaration

- `(decl:fn <identifier> <fndef>)`

# Function definition

- `(fndef <signature> <body?>)`
- `(signature [<arg...>] <type>)`
- arg: `(: <identifier> <type> <default?>)`

# Templates

Templates graft themselves onto type and function *definitions*, not their declaration and enclose them :

If templates finds a `(decl:type)` or `(decl:fn)` node, it will insert itself into their definition

- `(template [<args...>] <definition>)`

