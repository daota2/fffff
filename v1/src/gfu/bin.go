package gfu

import (
  "bytes"
  "fmt"
  //"log"
  "strings"
)

type Bin []byte

type BinType struct {
  BasicType
}

func NewBin(len int) Bin {
  return make(Bin, len)
}

func (b Bin) Type(g *G) Type {
  return &g.BinType
}

func (_ *BinType) Bool(g *G, val Val) (bool, E) {
  return len(val.(Bin)) > 0, nil
}

func (_ *BinType) Dump(g *G, val Val, out *strings.Builder) E {
  out.WriteString("(0x")

  for _, v := range val.(Bin) {
    fmt.Fprintf(out, " %02x", v)
  }

  out.WriteRune(')')
  return nil
}

func (_ *BinType) Dup(g *G, val Val) (Val, E) {
  var dst Bin
  src := val.(Bin)
  
  if len(src) > 0 {
    dst = NewBin(len(src))
    copy(dst, src)
  }

  return dst, nil
}

func (_ *BinType) Eq(g *G, lhs, rhs Val) (bool, E) {
  return bytes.Compare(lhs.(Bin), rhs.(Bin)) == 0, nil
}

func (_ *BinType) Index(g *G, val Val, key Vec) (Val, E) {
  if len(key) > 1 {
    return nil, g.E("Invalid index: %v", key.Type(g))
  }

  b := val.(Bin)
  i, ok := key[0].(Int)

  if !ok {
    return nil, g.E("Invalid index: %v", key[0].Type(g))
  }
  
  if i := int(i); i < 0 || i > len(b) {
    return nil, g.E("Index out of bounds: %v", i)
  }

  return Byte(b[i]), nil
}

func (_ *BinType) Len(g *G, val Val) (Int, E) {
  return Int(len(val.(Bin))), nil
}

func (_ *BinType) Print(g *G, val Val, out *strings.Builder) {
  out.WriteString(string(val.(Bin)))
}

func (_ *BinType) SetIndex(g *G, val Val, key Vec, set Setter) (Val, E) {
  if len(key) > 1 {
    return nil, g.E("Invalid index: %v", key.Type(g))
  }

  b := val.(Bin)
  i, ok := key[0].(Int)

  if !ok {
    return nil, g.E("Invalid index: %v", key[0].Type(g))
  }

  if i := int(i); i < 0 || i > len(b) {
    return nil, g.E("Index out of bounds: %v", i)
  }

  v, e := set(Byte(b[i]))

  if e != nil {
    return nil, e
  }

  bv, ok := v.(Byte)
  
  if !ok {
    return nil, g.E("Expected Byte: %v", v.Type(g))
  }

  b[i] = byte(bv)
  return bv, nil
}
