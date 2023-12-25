package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Int int64

type IntType struct {
  BasicType
}

type IntIter struct {
  pos, max Int
}

type IntIterType struct {
  BasicIterType
}

func (i Int) Abs() Int {
  if i < 0 {
    return -i
  }

  return i
}

func (_ Int) Type(g *G) Type {
  return &g.IntType
}

func (_ *IntType) Bool(g *G, val Val) (bool, E) {
  return val.(Int) != 0, nil
}

func (_ *IntType) Dump(g *G, val Val, out *strings.Builder) E {
  fmt.Fprintf(out, "%v", int64(val.(Int)))
  return nil
}

func (_ *IntType) Iter(g *G, val Val) (Val, E) {
  return new(IntIter).Init(g, val.(Int)), nil
}

func (i *IntIter) Init(g *G, max Int) *IntIter {
  i.max = max
  return i
}

func (_ *IntIter) Type(g *G) Type {
  return &g.IntIterType
}

func (_ *IntIterType) Bool(g *G, val Val) (bool, E) {
  i := val.(*IntIter)
  return i.pos < i.max, nil
}

func (_ *IntIterType) Drop(g *G, val Val, n Int) (Val, E) {
  i := val.(*IntIter)
  
  if i.max-i.pos < n {
    return nil, g.E("Nothing to drop")
  }

  i.pos += n
  return i, nil
}

func (_ *IntIterType) Dup(g *G, val Val) (Val, E) {
  out := *val.(*IntIter)
  return &out, nil
}

func (_ *IntIterType) Eq(g *G, lhs, rhs Val) (bool, E) {
  li := lhs.(*IntIter)
  ri, ok := rhs.(*IntIter)
  return ok && ri.max == li.max && ri.pos == li.pos, nil
}

func (_ *IntIterType) Pop(g *G, val Val) (Val, Val, E) {
  i := val.(*IntIter)
  
  if i.pos >= i.max {
    return &g.NIL, i, nil
  }

  v := i.pos
  i.pos++
  return v, i, nil
}

func (_ *IntIterType) Splat(g *G, val Val, out Vec) (Vec, E) {
  i := val.(*IntIter)

  for {
    v, _, e := g.Pop(i)

    if e != nil {
      return nil, e
    }

    if v == nil {
      break
    }

    out = append(out, v)
  }

  return out, nil
}
