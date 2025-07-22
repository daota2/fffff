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
}

type FunType struct {
  BasicType
}

func NewFun(g *G, env *Env, id *Sym, args []Arg) (*Fun, E) {
  return new(Fun).Init(g, env, id, args)
}

func (f *Fun) Init(g *G, env *Env, id *Sym, args []Arg) (*Fun, E) {
  if id != nil {
    f.id = id

    if e := env.Let(g, id, f); e != nil {
      return nil, e
    }
  }
  
  f.arg_list.Init(g, args)
  return f, nil
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

func (_ *FunType) Clone(g *G, val Val) (Val, E) {
  f, cf := val.(*Fun), new(Fun)
  *cf = *f
  
  if e := f.env.Clone(g, &cf.env); e != nil {
    return nil, e
  }
  
  return cf, nil
}

func (_ *FunType) Dump(g *G, val Val, out *bufio.Writer) E {
  f := val.(*Fun)

  if id := f.id; id == nil {
    out.WriteString("(fun")
  } else {
    fmt.Fprintf(out, "(fun %v", f.id)
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

func (env *Env) AddFun(g *G, id string, imp FunImp, args ...Arg) E {
  f, e := NewFun(g, env, g.Sym(id), args)
  
  if e != nil {
    return e
  }

  f.imp = imp
  return nil
}

type Recall struct {
  args Vec
}

func NewRecall(args Vec) (r Recall) {
  r.args = args
  return r
}

func (r Recall) Dump(g *G, out *bufio.Writer) {
  out.WriteString("(recall")

  for _, a := range r.args {
    g.Dump(a, out)
  }

  out.WriteRune(')')
}

func (r Recall) String() string {
  return "Recall"
}
