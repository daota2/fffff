package gfu

import (
  //"log"
  "strings"
)

type Macro struct {
  BasicVal
  
  env      *Env
  arg_list ArgList
  body     Vec
}

func NewMacro(g *G, env *Env, args []Arg) *Macro {
  return new(Macro).Init(g, env, args)
}

func (m *Macro) Init(g *G, env *Env, args []Arg) *Macro {
  m.BasicVal.Init(&g.MacroType, m)

  m.env = env
  m.arg_list.Init(g, args)
  return m
}

func (m *Macro) ExpandCall(g *G, task *Task, env *Env, args Vec) (*Env, Val, E) {
  avs := make(Vec, len(args))
  var e E

  for i, a := range args {
    if avs[i], e = a.Quote(g, task, env); e != nil {
      return nil, nil, e
    }
  }

  if e = m.arg_list.Check(g, args); e != nil {
    return nil, nil, e
  }

  var be Env
  m.arg_list.LetVars(g, &be, args)

  if e = m.body.Extenv(g, m.env, &be, false); e != nil {
    return nil, nil, e
  }

  var v Val

  if v, e = m.body.EvalExpr(g, task, &be); e != nil {
    return nil, nil, e
  }

  if e = v.Extenv(g, env, &be, false); e != nil {
    return nil, nil, e
  }

  return &be, v, nil
}

func (m *Macro) Call(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  var be *Env
  
  if be, v, e = m.ExpandCall(g, task, env, args); e != nil {
    return nil, e
  }

  return v.Eval(g, task, be)
}

func (m *Macro) Dump(out *strings.Builder) {
  out.WriteString("(macro (")

  for i, a := range m.arg_list.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(a.id.name)
  }

  out.WriteString(") ")

  for i, bv := range m.body {
    if i > 0 {
      out.WriteRune(' ')
    }

    bv.Dump(out)
  }

  out.WriteRune(')')
}
