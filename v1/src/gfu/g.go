package gfu

import (
  "io/ioutil"
  //"log"
  "strings"
)

type Syms map[string]*Sym

type G struct {
  syms Syms

  Debug bool  
  MainTask Task
  RootEnv Env

  MetaType,
  FalseType, FunType, IntType, MacroType, NilType, OptType, PrimType,
  QuoteType, SplatType, SpliceType, SymType, TrueType, VecType Type
  
  NIL Nil
  T True
  F False
}

func NewG() (*G, E) {
  return new(G).Init()
}

func (g *G) Init() (*G, E) {
  g.syms = make(Syms)
  g.MainTask.Init(NewChan(0), nil)
  return g, nil
}

func (g *G) NewEnv() *Env {
  var env Env 
  g.RootEnv.Clone(&env)
  return &env
}

func (g *G) EvalString(task *Task, env *Env, pos Pos, s string) (Val, E) {
  in := strings.NewReader(s)
  var out Vec

  for {
    vs, e := g.Read(&pos, in, Vec(out), 0)
    
    if e != nil {
      return nil, e
    }
    
    if vs == nil {
      break
    }

    out = vs
  }

  return out.EvalExpr(g, &g.MainTask, env)  
}

func (g *G) Load(task *Task, env *Env, fname string) (Val, E) {
  s, e := ioutil.ReadFile(fname)
  
  if e != nil {
    return nil, g.E("Failed loading file: %v\n%v", fname, e)
  }

  var pos Pos
  pos.Init(fname)
  return g.EvalString(task, env, pos, string(s))
}
