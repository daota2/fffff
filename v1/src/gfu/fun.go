package gfu

import (
  "bufio"
  "fmt"
  //"log"
)

type FunImp func(*G, *Task, *Env, Vec) (Val, E)

type Fun struct {
  id       *Sym
  env      Env
  arg_list ArgList
  body     Vec
  imp      FunImp
  pure     bool
}

type FunType struct {
  BasicType
}

type PunType struct {
  FunType
}

type EImpure struct {
  call Val
}

type EImpureType struct {
  BasicType
}

func NewFun(g *G, env *Env, id *Sym, args...Arg) *Fun {
  return new(Fun).Init(g, env, id, args)
}

func (f *Fun) Init(g *G, env *Env, id *Sym, args []Arg) *Fun {
  f.id = id
  f.arg_list.Init(g, args)
  return f
}

func (f *Fun) InitEnv(g *G, env *Env) E {
  return g.Extenv(env, &f.env, f.body, false)
}

func (f *Fun) CheckArgs(g *G, args Vec) (Vec, E) {
  if e := f.arg_list.Check(g, args); e != nil {
    return nil, e
  }

  if f.imp != nil {
    args = f.arg_list.Fill(g, args)
  }

  return args, nil
}

func (f *Fun) CallArgs(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  args, e := f.CheckArgs(g, args)

  if e != nil {
    return nil, e
  }

  if f.pure {
    task.pure++
    defer func() { task.pure-- }()
  } else {
    if task.pure > 0 {
      return nil, g.EImpure(f)
    }
  }
  
  if f.imp != nil {
    return f.imp(g, task, env, args)
  }

  var be Env
recall:
  f.env.Dup(&be)
  f.arg_list.LetVars(g, &be, args)
  var v Val

  if v, e = f.body.EvalExpr(g, task, &be, args_env); e != nil {
    if r, ok := e.(Recall); ok {
      be.Clear()
      args = r.args
      goto recall
    }

    return nil, e
  }

  return v, e
}

func (f *Fun) Type(g *G) Type {
  if f.pure {
    return &g.PunType
  }
  
  return &g.FunType
}

func (_ *FunType) ArgList(g *G, val Val) (*ArgList, E) {
  return &val.(*Fun).arg_list, nil
}

func (_ *FunType) Call(g *G, task *Task, env *Env, val Val, args Vec, args_env *Env) (Val, E) {
  f := val.(*Fun)
  args, e := args.EvalVec(g, task, args_env, args_env)

  if e != nil {
    return nil, e
  }

  return f.CallArgs(g, task, env, args, args_env)
}

func (_ *FunType) Dump(g *G, val Val, out *bufio.Writer) E {
  f := val.(*Fun)
  
  out.WriteRune('(')

  if f.pure {
    out.WriteString("pun")
  } else {
    out.WriteString("fun")
  }

  if id := f.id; id != nil {
    fmt.Fprintf(out, " %v", f.id)
  }

  nargs := len(f.arg_list.items)

  if nargs > 0 {
    out.WriteString(" (")
  }

  for i, a := range f.arg_list.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    if e := a.Dump(g, out); e != nil {
      return e
    }
  }

  if nargs > 0 {
    out.WriteRune(')')
  }

  out.WriteRune(')')
  return nil
}

func (env *Env) AddFun(g *G, id string, imp FunImp, args ...Arg) (*Fun, E) {
  ids := g.Sym(id)
  f := NewFun(g, env, ids, args...)
  f.imp = imp

  if e := env.Let(g, ids, f); e != nil {
    return nil, e
  }

  return f, nil
}

func (env *Env) AddPun(g *G, id string, imp FunImp, args ...Arg) (*Fun, E) {
  f, e := env.AddFun(g, id, imp, args...)

  if e != nil {
    return nil, e
  }

  f.pure = true
  return f, nil
}

func (g *G) EImpure(call Val) *EImpure {
  e := new(EImpure)
  e.call = call
  return e
}

func (_ EImpure) Type(g *G) Type {
  return &g.EImpureType
}

func (_ *EImpureType) Dump(g *G, val Val, out *bufio.Writer) E {
  fmt.Fprintf(out, "Error: Impure: %v", g.EString(val.(*EImpure).call))
  return nil
}
