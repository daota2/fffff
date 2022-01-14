package gfu

import (
  //"log"
)

type Vec []Val

func (v Vec) EvalExpr(g *G, env *Env) (Val, E) {
  out := g.NIL
  
  for _, it := range v {
    var e E
    
    if out, e = it.Eval(g, env); e != nil {
      return g.NIL, e
    }

    if g.recall {
      break
    }
  }

  return out, nil
}

func (v Vec) EvalVec(g *G, env *Env) (Vec, E) {
  var out Vec
  
  for _, it := range v {
    it, e := it.Eval(g, env)

    if e != nil {
      return nil, g.E("Arg eval failed: %v", e)
    }

    if g.recall {
      break
    }
    
    if it.val_type == g.SplatType {
      out = it.Splat(g, out)
    } else {
      if it.val_type == g.VecType {
        it.imp = it.Splat(g, nil)
      }
      
      out = append(out, it)
    }
  }

  return out, nil
}

func (v Vec) Push(its...Val) Vec {
  return append(v, its...)
}

func (v Vec) Peek(g *G) Val {
  n := len(v)
  
  if n == 0 {
    return g.NIL
  }

  return v[n-1]
}

func (v Vec) Pop(g *G) (Val, Vec) {
  n := len(v)

  if n == 0 {
    return g.NIL, v
  }

  return v[n-1], v[:n-1]
}
