package gfu

import (
)

type FunImp func(*G, Pos, []Val, *Env) (Val, E)

type Fun struct {
  env *Env
  arg_list ArgList
  body []Form
  imp FunImp
}

func NewFun(g *G, env *Env, args []*Sym) *Fun {
  return new(Fun).Init(g, env, args)
}

func (f *Fun) Init(g *G, env *Env, args []*Sym) *Fun {
  f.env = env
  f.arg_list.Init(g, args)
  return f
}

func (e *Env) AddFun(g *G, id string, imp FunImp, args...string) {
  as := make([]*Sym, len(args))

  for i, a := range args {
    as[i] = g.S(a)
  }
  
  f := NewFun(g, e, as)
  f.imp = imp
  
  var v Val
  v.Init(g.Fun, f)
  e.Put(g.S(id), v)
}
