package gfu

import (
  "fmt"
  "strings"
)

type Type interface {
  fmt.Stringer
  Val

  Init(*G, *Sym) E
  Bool(*G, Val) (bool, E)
  Call(*G, *Task, *Env, Val, Vec) (Val, E)
  Clone(*G, Val) (Val, E)
  Drop(*G, Val, Int) (Val, E)
  Dump(*G, Val, *strings.Builder) E
  Dup(*G, Val) (Val, E)
  Eq(*G, Val, Val) (bool, E)
  Eval(*G, *Task, *Env, Val) (Val, E)
  Expand(*G, *Task, *Env, Val, Int) (Val, E)
  Extenv(*G, *Env, *Env, Val, bool) E
  Id() *Sym
  Is(*G, Val, Val) bool
  Iter(*G, Val) (Val, E)
  Len(*G, Val) (Int, E)
  Pop(*G, Val) (Val, Val, E)
  Print(*G, Val, *strings.Builder)
  Push(*G, Val, ...Val) (Val, E)
  Quote(*G, *Task, *Env, Val) (Val, E)
  Splat(*G, Val, Vec) (Vec, E)
}

type BasicType struct {
  id *Sym
}

type MetaType struct {
  BasicType
}

func (t *BasicType) Init(g *G, id *Sym) E {
  t.id = id
  return nil
}

func (_ *BasicType) Bool(g *G, val Val) (bool, E) {
  return true, nil
}

func (t *BasicType) Call(g *G, task *Task, env *Env, val Val, args Vec) (Val, E) {
  return nil, g.E("Call not supported: %v", t.id)
}

func (_ *BasicType) Clone(g *G, val Val) (Val, E) {
  return g.Dup(val)
}

func (_ *BasicType) Drop(g *G, val Val, n Int) (out Val, e E) {
  for i := Int(0); i < n; i++ {
    if _, out, e = g.Pop(val); e != nil {
      return nil, e
    }
  }

  return out, nil
}

func (t *BasicType) Dump(g *G, val Val, out *strings.Builder) E {
  out.WriteString(t.id.name)
  return nil
}

func (_ *BasicType) Dup(g *G, val Val) (Val, E) {
  return val, nil
}

func (_ *BasicType) Eq(g *G, lhs, rhs Val) (bool, E) {
  return g.Is(lhs, rhs), nil
}

func (_ *BasicType) Eval(g *G, task *Task, env *Env, val Val) (Val, E) {
  return val, nil
}

func (_ *BasicType) Expand(g *G, task *Task, env *Env, val Val, depth Int) (Val, E) {
  return val, nil
}

func (_ *BasicType) Extenv(g *G, src, dst *Env, val Val, clone bool) E {
  return nil
}

func (t *BasicType) Id() *Sym {
  return t.id
}

func (_ *BasicType) Is(g *G, lhs, rhs Val) bool {
  return lhs == rhs
}

func (t *BasicType) Iter(g *G, val Val) (Val, E) {
  return nil, g.E("Iter not supported: %v", t.id)
}

func (t *BasicType) Len(g *G, val Val) (Int, E) {
  return -1, g.E("Len not supported: %v", t.id)
}

func (t *BasicType) Pop(g *G, val Val) (Val, Val, E) {
  return nil, nil, g.E("Pop not supported: %v", t.id)
}

func (_ *BasicType) Print(g *G, val Val, out *strings.Builder) {
  g.Dump(val, out)
}

func (t *BasicType) Push(g *G, val Val, its ...Val) (Val, E) {
  return nil, g.E("Push not supported: %v", t.id)
}

func (_ *BasicType) Quote(g *G, task *Task, env *Env, val Val) (Val, E) {
  return val, nil
}

func (_ *BasicType) Splat(g *G, val Val, out Vec) (Vec, E) {
  return append(out, val), nil
}

func (t *BasicType) String() string {
  return t.id.name
}

func (_ *BasicType) Type(g *G) Type {
  return &g.MetaType
}

func (e *Env) AddType(g *G, t Type, id string) E {
  t.Init(g, g.Sym(id))
  return e.Let(g, t.Id(), t)
}
