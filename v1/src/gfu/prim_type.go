package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type PrimType struct {
  BasicType
}

func (t *PrimType) Init(id *Sym) *PrimType {
  t.BasicType.Init(id)
  return t
}

func (t *PrimType) Call(g *G, pos Pos, val Val, args VecForm, env *Env) (Val, Error) {
  p := val.AsPrim()
  g.prim = p
  return p.imp(g, pos, args, env)
}

func (t *PrimType) Dump(val Val, out *strings.Builder) {
  fmt.Fprintf(out, "(prim %v)", val.AsPrim().id)
}

func (v Val) AsPrim() *Prim {
  return v.imp.(*Prim)
}
