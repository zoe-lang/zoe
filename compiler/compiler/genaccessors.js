#!/usr/bin/env node
const fs = require('fs')
const path = require('path')

const input = fs.readFileSync(path.join(__dirname, './nodes.go'), 'utf-8')

const re_type = /type (\w+) struct \{\s*NodeBase\s*([^\}]+)\}/g
const re_prop = /(\w+)\s+((?:\[\])?\*?\w+)/g

const names = []

console.log(`// Code generated by a lame .js file, DO NOT EDIT.

package zoe

import "io"
`)

// 1 trim left
// 2 exp
// 3 trim right
// 4 trim left
// 5 fnbody
// 6 trim right
const re_tag = /\{\{(\-)?((?:(?!\}\})[^])+)?(\-)?\}\}|%%(\-)?((?:(?!%%)[^])*)(\-)?%%/g

/**
 * A simple templating function stolen from https://krasimirtsonev.com/blog/article/Javascript-template-engine-in-just-20-line
 */
function Template(tpl) {
  var fncode = ['var __res = [];'], start = 0, match;

  function add(m) {
    fncode.push('__res.push(`' + tpl.substr(start, m.index - start) + '`)')
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
  return (data) => fn(data);
}

var tpl_create = Template(`
func (p *Position) Create{{v.type}}() *{{v.type}} {
  res := &{{v.type}}{}
  res.ExtendPosition(p)
  return res
}

func (tk *Token) Create{{v.type}}() *{{v.type}} {
  return tk.Position.Create{{v.type}}()
}

func (r *{{v.type}}) EnsureTuple() *Tuple {
%% if (v.type === 'Tuple') { %%
  return r
%% } else { %%
  res := &Tuple{}
  res.AddChildren(r)
  return res
%% } %%
}

%% if (v.fields.length === 0) { %%
func (r *{{v.type}}) Dump(w io.Writer) {
  w.Write([]byte(cyan(r.GetText())))
}
%% } else { %%

func (r *{{v.type}}) Dump(w io.Writer) {
%% if (v.lower === 'operation') { %%
  w.Write([]byte("("))
%% } else if (v.lower === 'block') { %%
  w.Write([]byte("{"))
%% } else if (v.lower === 'tuple' || v.lower === 'vartuple') { %%
  w.Write([]byte("["))
%% } else { %%
  w.Write([]byte("({{v.lower}} "))
%% } %%

%% var __i = 0 %%
%% for (var f of v.fields) { %%
%%   __i++ %%
%%   if (f.is_list) { %%
      for i, n := range r.{{f.name}} {
        n.Dump(w)
        if i < len(r.{{f.name}}) - 1 {
          w.Write([]byte(" "))
        }
      }
%%   } else { %%
      if r.{{f.name}} != nil {
        r.{{f.name}}.Dump(w)
      } else {
        w.Write([]byte(mag("<nil>")))
      }
%%      if (__i < v.fields.length) { %%
      w.Write([]byte(" "))
%%      } %%
%%   } %%
%% } %%

%% if (v.lower === 'block') { %%
  w.Write([]byte("}"))
%% } else if (v.lower === 'tuple' || v.lower === 'vartuple') { %%
  w.Write([]byte("]"))
%% } else { %%
  w.Write([]byte(")"))
%% } %%
}

%% } %%

%% for (var f of v.fields) { -%%

%% if (f.is_list) { %%
func (r *{{v.type}}) Add{{f.name}}(other ...{{f.simple_type}}) *{{v.type}} {
  for _, c := range other {
    if c != nil {
    %% if (f.simple_type === 'Node') { %%
      switch v := c.(type) {
      case *Fragment:
        r.Add{{f.name}}(v.Children...)
      default:
        r.{{f.name}} = append(r.{{f.name}}, c)
        r.ExtendPosition(c)
      }
    %% } else { %%
      r.{{f.name}} = append(r.{{f.name}}, c)
      r.ExtendPosition(c)
    %% } %%
    }
  }
  return r
}
%% } else { %%

%% if (f.type !== 'Node' && f.type !== '*Token') { %%
func (r *{{v.type}}) Ensure{{f.name}}(fn func ({{f.name[0].toLowerCase()}} {{f.type}})) *{{v.type}} {
  if r.{{f.name}} == nil {
    r.{{f.name}} = &{{f.nonptr_type}}{}
  }
  fn(r.{{f.name}})
  r.ExtendPosition(r.{{f.name}})
  return r
}
%% } %%

func (r *{{v.type}}) Set{{f.name}}(other {{f.type}}) *{{v.type}} {
  r.{{f.name}} = other
  if other != nil {
    r.ExtendPosition(other)
  }
  return r
}
%% } %%

%% } %%
`)

var types = []
var match, pmatch
while (match = re_type.exec(input)) {
  const [_, type, src] = match
  if (type === 'NodeBase') continue
  const finalsrc = src.replace(/\/\/[^\n]*\n/g, '')
  const lower = type.toLowerCase()

  const fields = []

  while (pmatch = re_prop.exec(finalsrc)) {
    const [_, name, proptype] = pmatch
    if (!proptype.includes('*') && !proptype.includes('Node')) continue

    fields.push({
      name,
      type: proptype,
      // lower,
      simple_type: proptype.replace(/\[\]/g, ''),
      nonptr_type: proptype.replace(/\*/g, ''),
      is_list: proptype.includes('[]')
    })

    // console.log(type, name, proptype, is_list)
  }
  types.push({
    type,
    lower,
    fields,
  })
}

for (var t of types) {
  console.log(tpl_create(t))
}
// console.log(types)