package gfu

type BoolType struct {
  BasicType
}

func (t *BoolType) Bool(g *G, val Val) bool {
  return val.imp.(bool)
}

func (v Val) Bool(g *G) bool {
  vt := v.val_type
  
  if vt == g.Bool {
    return v.imp.(bool)
  }

  return vt.Bool(g, v)
}
