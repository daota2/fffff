package gfu

import (
  //"log"
  "strings"
)

type Str string

func (s Str) Bool(g *G) bool {
  return len(s) > 0
}

func (s Str) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return s, nil
}

func (s Str) Clone() Val {
  return s
}

func (s Str) Dump(out *strings.Builder) {
  out.WriteRune('"')
  out.WriteString(string(s))
  out.WriteRune('"')
}

func (s Str) Eq(g *G, rhs Val) bool {
  rs, ok := rhs.(Str)
  return ok && rs == s
}

func (s Str) Eval(g *G, task *Task, env *Env) (Val, E) {
  return s, nil
}

func (s Str) Is(g *G, rhs Val) bool {
  return s.Eq(g, rhs)
}

func (s Str) Len(g *G) (Int, E) {
  return Int(len(s)), nil
}

func (_ Str) Pop(g *G) (Val, Val, E) {
  return nil, nil, g.E("Pop not supported: Str")
}

func (_ Str) Push(g *G, its...Val) (Val, E) {
  return nil, g.E("Push not supported: Str")
}

func (s Str) Quote(g *G, task *Task, env *Env) (Val, E) {
  return s, nil
}

func (s Str) Splat(g *G, out Vec) Vec {
  return append(out, s)
}

func (s Str) String() string {
  return DumpString(s)
}

func (_ Str) Type(g *G) *Type {
  return &g.StrType
}
