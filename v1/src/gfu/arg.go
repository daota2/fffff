package gfu

import (
  //"log"
  "strings"
)

type ArgType int

const (
  ARG_PLAIN ArgType = 0
  ARG_OPT   ArgType = 1
  ARG_SPLAT ArgType = 2
)

type Arg struct {
  arg_type ArgType
  str_id   string
  id       *Sym
  opt_val  Val
}

func (a *Arg) Init(id *Sym) *Arg {
  a.id = id
  return a
}

func A(id string) (a Arg) {
  a.str_id = id
  return a
}

func OptA(id string, val Val) (a Arg) {
  a.str_id = id
  a.opt_val = val
  a.arg_type = ARG_OPT
  return a
}

func SplatA(id string) (a Arg) {
  a.str_id = id
  a.arg_type = ARG_SPLAT
  return a
}

func (a Arg) String() string {
  var out strings.Builder

  switch a.arg_type {
  case ARG_OPT:
    out.WriteRune('(')
    out.WriteString(a.id.name)
    a.opt_val.Dump(&out)
    out.WriteRune(')')
  case ARG_SPLAT:
    out.WriteString(a.id.name)
    out.WriteString("..")
  default:
    out.WriteString(a.id.name)
  }

  return out.String()
}

type ArgList struct {
  items    []Arg
  min, max int
}

func (l *ArgList) Init(g *G, args []Arg) *ArgList {
  nargs := len(args)

  if nargs == 0 {
    return l
  }

  l.items = args
  l.min, l.max = nargs, nargs

  for i, a := range l.items {
    if a.arg_type == ARG_OPT {
      l.min--
    }

    if a.id == nil {
      l.items[i].id = g.Sym(a.str_id)
    }
  }

  a := l.items[nargs-1]

  if a.arg_type == ARG_SPLAT {
    l.min--
    l.max = -1
  }

  return l
}

func (l *ArgList) Check(g *G, args Vec) E {
  nargs := len(args)

  if (l.min != -1 && nargs < l.min) || (l.max != -1 && nargs > l.max) {
    return g.E("Arg mismatch")
  }

  return nil
}

func (l *ArgList) Fill(g *G, args Vec) Vec {
  for i := len(args); i < len(l.items); i++ {
    a := l.items[i]

    if a.arg_type != ARG_OPT {
      break
    }

    args = append(args, a.OptVal(g))
  }

  return args
}

func (a Arg) OptVal(g *G) Val {
  v := a.opt_val

  if v == nil {
    v = &g.NIL
  }

  return v
}

func (l *ArgList) LetVars(g *G, env *Env, args Vec) {
  nargs := len(args)

  for i, a := range l.items {
    if a.arg_type == ARG_SPLAT {
      var v Vec

      if i < nargs {
        v = make(Vec, nargs-i)
        copy(v, args[i:])
      }

      env.Let(a.id, v)
      break
    }

    if i < nargs {
      env.Let(a.id, args[i])
    } else {
      env.Let(a.id, a.OptVal(g))
    }
  }
}

func ParseArgs(g *G, task *Task, env *Env, in Vec) ([]Arg, E) {
  var e E
  var out []Arg

  for _, v := range in {
    var a Arg

    if id, ok := v.(*Sym); ok {
      a.id = id
    } else if vv, ok := v.(Vec); ok {
      a.arg_type = ARG_OPT
      a.id = vv[0].(*Sym)

      if len(vv) > 1 {
        if a.opt_val, e = vv[1].Eval(g, task, env); e != nil {
          return nil, e
        }
      }
    } else if sv, ok := v.(Splat); ok {
      a.arg_type = ARG_SPLAT
      a.id = sv.val.(*Sym)
    } else {
      return nil, g.E("Invalid arg: %v", v)
    }

    out = append(out, a)
  }

  return out, nil
}
