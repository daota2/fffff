package gfu

import (
  //"log"
  "strings"
)

type Mac struct {
  BasicVal
  
  env      *Env
  arg_list ArgList
  body     Vec
}

func NewMac(g *G, env *Env, args []Arg) *Mac {
  return new(Mac).Init(g, env, args)
}

func (m *Mac) Init(g *G, env *Env, args []Arg) *Mac {
  m.BasicVal.Init(&g.MacType, m)

  m.env = env
  m.arg_list.Init(g, args)
  return m
}

func (m *Mac) ExpandCall(g *G, task *Task, env *Env, args Vec) (Val, E) {
  avs := make(Vec, len(args))
  var e E

  for i, a := range args {
    if avs[i], e = a.Quote(g, task, env); e != nil {
      return nil, e
    }
  }

  if e = m.arg_list.Check(g, args); e != nil {
    return nil, e
  }

  var be Env
  m.arg_list.LetVars(g, &be, args)

  if e = m.body.Extenv(g, m.env, &be, false); e != nil {
    return nil, e
  }

  var v Val

  if v, e = m.body.EvalExpr(g, task, &be); e != nil {
    return nil, e
  }

  return v, nil
}

func (m *Mac) Call(g *G, task *Task, env *Env, args Vec) (v Val, e E) {  
  if v, e = m.ExpandCall(g, task, env, args); e != nil {
    return nil, e
  }

  var be Env

  if e = v.Extenv(g, m.env, &be, false); e != nil {
    return nil, e
  }

  if e = v.Extenv(g, env, &be, false); e != nil {
    return nil, e
  }

  return v.Eval(g, task, &be)
}

func (m *Mac) Dump(out *strings.Builder) {
  out.WriteString("(mac (")

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
