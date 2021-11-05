package gfu

import (
  "io/ioutil"
  //"log"
  "strings"
)

type Syms map[string]*Sym

type G struct {
  Debug bool
  RootEnv Env
  
  sym_tag Tag
  syms Syms

  recall_args []Val
  
  MetaType,
  BoolType, FormType, FunType, IntType, MacroType, NilType, PrimType, SplatType,
  SymType, VecType Type
  
  NIL, T, F Val
}

func NewG() (*G, E) {
  return new(G).Init()
}

func (g *G) Init() (*G, E) {
  g.syms = make(Syms)
  return g, nil
}

func (g *G) EvalString(pos Pos, s string, env *Env) (Val, E) {
  in := strings.NewReader(s)
  var out Forms
  
  for {
    fs, e := g.Read(&pos, in, out, 0)
    
    if e != nil {
      return g.NIL, e
    }
    
    if fs == nil {
      break
    }

    out = fs
  }

  return out.Eval(g, env)  
}

func (g *G) Load(pos Pos, fname string, env *Env) (Val, E) {
  s, e := ioutil.ReadFile(fname)
  
  if e != nil {
    return g.NIL, g.E(pos, "Failed loading file: %v\n%v", fname, e)
  }

  var fpos Pos
  fpos.Init(fname)
  return g.EvalString(fpos, string(s), env)
}
