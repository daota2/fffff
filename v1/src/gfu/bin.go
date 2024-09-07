package gfu

import (
  "bytes"
  "fmt"
  //"log"
  "strings"
)

type Bin struct {
  data []byte
  len  int
}

type BinType struct {
  BasicType
}

func NewBin(len int) *Bin {
  b := new(Bin)
  b.data = make([]byte, len)
  b.len = len
  return b
}

func (b Bin) Type(g *G) Type {
  return &g.BinType
}

func (_ *BinType) Bool(g *G, val Val) (bool, E) {
  return val.(*Bin).len > 0, nil
}

func (_ *BinType) Dump(g *G, val Val, out *strings.Builder) E {
  out.WriteString("(0x")

  for _, v := range val.(*Bin).data {
    fmt.Fprintf(out, " %02x", v)
  }

  out.WriteRune(')')
  return nil
}

func (_ *BinType) Dup(g *G, val Val) (Val, E) {
  b := val.(*Bin)
  src := b.data
  dst := new(Bin)

  if len(src) > 0 {
    dst.data = make([]byte, len(src))
    copy(dst.data, src)
  }

  return dst, nil
}

func (_ *BinType) Eq(g *G, lhs, rhs Val) (bool, E) {
  lb, rb := lhs.(*Bin), rhs.(*Bin)
  return bytes.Compare(lb.data[:lb.len], rb.data[:rb.len]) == 0, nil
}

func (_ *BinType) Index(g *G, val Val, key Vec) (Val, E) {
  if len(key) > 1 {
    return nil, g.E("Invalid index: %v", key.Type(g))
  }

  b := val.(*Bin)
  i, ok := key[0].(Int)

  if !ok {
    return nil, g.E("Invalid index: %v", key[0].Type(g))
  }
  
  if i := int(i); i < 0 || i > b.len {
    return nil, g.E("Index out of bounds: %v", i)
  }

  return Byte(b.data[i]), nil
}

func (_ *BinType) Len(g *G, val Val) (Int, E) {
  return Int(val.(*Bin).len), nil
}

func (_ *BinType) Print(g *G, val Val, out *strings.Builder) {
  out.WriteString(string(val.(*Bin).data))
}

func (_ *BinType) SetIndex(g *G, val Val, key Vec, set Setter) (Val, E) {
  if len(key) > 1 {
    return nil, g.E("Invalid index: %v", key.Type(g))
  }

  b := val.(*Bin)
  i, ok := key[0].(Int)

  if !ok {
    return nil, g.E("Invalid index: %v", key[0].Type(g))
  }

  if i := int(i); i < 0 || i > len(b.data) {
    return nil, g.E("Index out of bounds: %v", i)
  }

  v, e := set(Byte(b.data[i]))

  if e != nil {
    return nil, e
  }

  bv, ok := v.(Byte)
  
  if !ok {
    return nil, g.E("Expected Byte: %v", v.Type(g))
  }

  b.data[i] = byte(bv)
  return bv, nil
}
