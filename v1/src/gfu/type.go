package gfu

import (
  "strings"
)

type Type struct {
  BasicVal
  id *Sym
}

func (t *Type) Init(g *G, id *Sym) *Type {
  t.BasicVal.Init(&g.MetaType, t)
  t.id = id
  return t
}

func (t *Type) Bool(g *G) bool {
  return true
}

func (t *Type) Dump(out *strings.Builder) {
  out.WriteString(t.id.name)
}

func (t *Type) Eq(g *G, rhs Val) bool {
  return t == rhs
}

func (t *Type) Eval(g *G, task *Task, env *Env) (Val, E) {
  return t, nil
}

func (t *Type) Id() *Sym {
  return t.id
}

func (t *Type) Is(g *G, rhs Val) bool {
  return t == rhs
}

func (t *Type) Quote(g *G, task *Task, env *Env) (Val, E) {
  return t, nil
}

func (t *Type) Splat(g *G, out Vec) Vec {
  return append(out, t)
}

func (t *Type) String() string {
  return t.id.name
}

func (e *Env) AddType(g *G, t *Type, id string) {
  t.Init(g, g.Sym(id))
  e.Let(t.Id(), t)
}
