#!/usr/bin/env node
/**
  The objective of this generator is to automatically write boilerplate code
  that helps in
    - Generating constructors when using the diamond pattern
    - Generating disambiguated methods with the "diamond" pattern
*/
declare var console: any
declare var require: <T>(m: string) => T
declare var __dirname: string
const process = require<any>('process')
const fs = require<{
  readFileSync(...a: any[]): string}>('fs')
const path = require<{
  join(...a: string[]): string
}>('path')

const input = fs.readFileSync(path.join(__dirname, './nodes.go'), 'utf-8')

const names = []

const output = (msg: string) => console.log(msg)
const debug = (msg: any) => {
  process.stderr.write(typeof msg === 'string' ? msg + '\n' : JSON.stringify(msg))
}

output(`// Code generated by a lame .js file, DO NOT EDIT.

package zoe

`)

class Var {
  name: string
  type: string
  constructor(public exp: string) {
    [this.name, this.type] = exp.split(/\s+/g)
    this.type = this.type ?? ''
  }
}

class Type {
  methods: Map<string, Func> = new Map()
  members: Map<string, Var> = new Map()
  supers: Map<string, Type> = new Map()
  is_super = false

  constructor(
    public name: string,
    public kind: string,
    public body: string[],
  ) { }

  get all_members(): Var[] {
    var res: Var[] = []
    for (let s of this.supers.values()) {
      if (s.name === 'Located') continue
      res = [...res, ...s.all_members]
    }
    for (let m of this.members.values()) {
      res.push(m)
    }
    return res
  }

  get all_members_str(): string {
    return this.all_members.map(m => `${m.name} ${m.type}`).join(', ')
  }

  get all_types(): Set<Type> {
    var res = new Set<Type>()
    function process(t: Type) {
      res.add(t)
      for (var s of t.supers.values()) process(s)
    }
    process(this)
    return res
  }

  get creators(): Func[] {
    var res: Func[] = []
    for (var s of this.all_types) {
      var create = s.methods.get('create')
      if (create) {
        res.push(create)
      }
    }
    return res
  }

  get conflict_methods(): Func[] {
    var already_seen = new Map<string, Func>()
    var overriden = new Map<string, Func>()
    var self = this

    function process(t: Type) {
      for (let f of t.methods.values()) {
        if (f.name === 'create') continue
        if (!already_seen.has(f.name)) {
          already_seen.set(f.name, f)
        } else {
          if (!self.methods.has(f.name)) {
            overriden.set(f.name, f)
          }
        }
      }
      for (let s of t.supers.values()) {
        process(s)
      }
    }
    process(this)

    return [...overriden.values()]
  }

  hasSuper(supname: string) {
    if (this.supers.has(supname)) return true
    for (let s of this.supers.values()) {
      if (s.hasSuper(supname)) return true
    }
    return false
  }
}

class Func {
  self: string | null = null

  constructor(
    public name: string,
    public args: Var[],
    public ret: string
  ) { }
}

// output(input)
const re_type = /^type\s+(\w+)\s+(struct|interface)\s*\{([^\}]*)\}/gm
const re_func = /^func(?: \(\w+ \*?(\w+)\))? (\w+)\(([^\)]*)\)(?: ((?:\[\])?\*?\w+))?/gm

var match: RegExpExecArray

const types = new Map<string, Type>()
while (match = re_type.exec(input)) {
  const type = new Type(
    match[1],
    match[2],
    match[3].split(/\s*\n\s*/g).map(l => l.trim()).filter(t => !!t)
  )

  types.set(type.name, type)
}

for (let t of types.values()) {
  if (t.kind === 'interface') continue
  for (let l of t.body) {
    l = l.replace(/\/\/.*$/, '')
    let splitted = l.split(/\s+/g).filter(f => !!f)
    if (splitted.length === 1) {
      var sup = types.get(splitted[0])
      sup.is_super = true
      t.supers.set(splitted[0], sup)
    } else if (splitted.length > 1) {
      t.members.set(splitted[0], new Var(l))
    }
  }
}

const funcs = new Map<string, Func>()
while (match = re_func.exec(input)) {
  const fn = new Func(
    match[2],
    match[3].split(/\s*,\s*/g).map(f => f.trim()).map(v => new Var(v)), // args
    (match[4] ?? '').trim(), // return type
  )
  if (match[1]) {
    fn.self = match[1]
  }
  if (fn.self && types.has(fn.self)) {
    // output(fn.self)
    types.get(fn.self).methods.set(fn.name, fn)
  } else {
    funcs.set(fn.name, fn)
  }
  // output(fn)
}

const tpl_creates = Template(`
func (parser *Parser) create{{v.name}}(scope *Scope) {{ !v.is_super ? '*' : '' }}{{v.name}} {
  var res = {{ !v.is_super ? '&' : '' }}{{ v.name }}{}
  %%- for (let c of v.creators) { %%
  res.{{ c.self !== v.name ? (c.self + '.') : '' }}create(parser, scope)
  %%- } %%
  return res
}

`)

const tpl_override = Template(`
func (n *{{v.type.name}}) {{v.fn.name}}({{v.fn.args.map(a => a.exp).join(', ')}}) {{v.fn.ret}} {
  {{ v.fn.ret ? 'return ' : '' }}n.{{ v.final.name }}.{{v.fn.name}}({{ v.fn.args.map(a => a.name).join(', ') }})
}
`)

// const tpl_registers = Template(`
// func (n *{{v.name}}) Register(other Node) {
//   n.Located.Register(n)
//   other.SetParent(n)
// }
// `)


for (let t of types.values()) {
  if (t.kind === 'interface') continue

  if (t.is_super) continue

  if (!funcs.has(`Create${t.name}`)) {
    output(tpl_creates(t))
  }

  // debug(JSON.stringify(t.conflict_methods))
  for (var f of t.conflict_methods) {
    debug(`overriding ${t.name}.${f.name}() with method from ${f.self}`)
    output(tpl_override({
      type: t,
      final: types.get(f.self),
      fn: f,
    }))
  }
}



/**
 * A simple templating function stolen from https://krasimirtsonev.com/blog/article/Javascript-template-engine-in-just-20-line
 */
function Template(tpl: string) {

  // 1 trim left
  // 2 exp
  // 3 trim right
  // 4 trim left
  // 5 fnbody
  // 6 trim right
  const re_tag = /\{\{(\-)?((?:(?!\}\})[^])+)?(\-)?\}\}|%%(\-)?((?:(?!%%)[^])*)(\-)?%%/g

  var fncode = ['var __res = [];'], start = 0, match;
  var trim_left = false
  var trim_right = false

  function add(m) {
    var str = tpl.substr(start, m.index - start)
    trim_left = !!(m[1] || m[4])
    if (trim_left) {
      str = (str as any).trimEnd()
    }
    if (trim_right) {
      str = (str as any).trimStart()
    }

    trim_right = !!(m[3] || m[6])

    if (m[3] || m[6]) str = (str as any).trimEnd()
    fncode.push('__res.push(`' + str + '`)')
    if (m[2]) {
      fncode.push(`__res.push(${m[2]})`)
    } else if (m[5]) {
      fncode.push(m[5])
    }
    // console.log(cursor, match)
  }

  while(match = re_tag.exec(tpl)) {
      add(match);
      start = match.index + match[0].length;
  }

  fncode.push('__res.push(`' + tpl.substr(start, tpl.length - start) + '`)');
  fncode.push(`return __res.join('');`);
  // code = code.replace(/[\r\t\n]+/g, ' ')
  var rescode = fncode.join('\n')
  // console.log(rescode)
  var fn = new Function('v', rescode)
  return (data: any) => fn(data);
}
