package gfu

import (
  //"log"
  "strings"
)

type MetaType struct {
  BasicType
}

func (t *MetaType) Init(id *Sym) *MetaType {
  t.BasicType.Init(id)
  return t
}

func (t *MetaType) AsBool(g *G, val Val) bool {
  return false
}

func (t *MetaType) Call(g *G, val Val, args ListForm, env *Env, pos Pos) (Val, Error) {
  return val.AsMeta().New(g, val, args, env, pos)
}
  
func (t *MetaType) Dump(val Val, out *strings.Builder) {
  out.WriteString(val.AsMeta().Id().name)
}

func (v Val) AsMeta() Type {
  return v.imp.(Type)
}
