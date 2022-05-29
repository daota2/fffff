package gfu

import (
  "fmt"
  "strings"
)

type Chan chan Val

func NewChan(buf Int) Chan {
  return make(Chan, buf)
}

func (c Chan) Bool(g *G) bool {
  return len(c) != 0
}

func (c Chan) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return c, nil
}

func (c Chan) Clone() Val {
  return c
}

func (c Chan) Dump(out *strings.Builder) {
  fmt.Fprintf(out, "(Chan %v)", (chan Val)(c))
}

func (c Chan) Eq(g *G, rhs Val) bool {
  return c.Is(g, rhs)
}

func (c Chan) Eval(g *G, task *Task, env *Env) (Val, E) {
  return c, nil
}

func (c Chan) Is(g *G, rhs Val) bool {
  return c == rhs
}

func (c Chan ) Pop(g *G) (Val, Val, E) {
  v := <- c

  if v == nil {
    v = &g.NIL
  }

  return v, c, nil
}

func (c Chan) Push(g *G, its...Val) (Val, E) {
  for _, v := range its {
    c <- v
  }

  return c, nil
}

func (c Chan) Quote(g *G, task *Task, env *Env) (Val, E) {
  return c, nil
}

func (c Chan) Splat(g *G, out Vec) Vec {
  return append(out, c)
}

func (c Chan) String() string {
  return DumpString(c)
}

func (c Chan) Type(g *G) *Type {
  return &g.ChanType
}
