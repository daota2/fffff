package gfu

import (
  "fmt"
  "strings"
)

type Type interface {
  Call(g *G, val Val, args []Val, env *Env, pos Pos) (Val, Error)
  Dump(x Val, out *strings.Builder)
  Eq(x, y Val) bool
}

type BasicType struct {
  id *Sym
}

func (t *BasicType) Init(id *Sym) *BasicType {
  t.id = id
  return t
}

func (t *BasicType) Call(g *G, val Val, args []Val, env *Env, pos Pos) (Val, Error) {
  return val, nil
}

func (t *BasicType) Dump(x Val, out *strings.Builder) {
  fmt.Fprintf(out, "%v", x.imp)
}

func (t *BasicType) Eq(x, y Val) bool {
  return x.imp == y.imp
}
