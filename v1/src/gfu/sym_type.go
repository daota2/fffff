package gfu

import (
  "strings"
)

type SymType struct {
  BasicType
}

func (t *SymType) Dump(val Val, out *strings.Builder) {
  out.WriteString(val.AsSym().name)
}

func (t *SymType) Eval(g *G, val Val, env *Env) (Val, E) {
  s := val.AsSym()
  _, found := env.Find(s)

  if found == nil {
    return g.NIL, g.E("Unknown: %v", s)
  }

  return found.Val, nil
}

func (v Val) AsSym() *Sym {
  return v.imp.(*Sym)
}
