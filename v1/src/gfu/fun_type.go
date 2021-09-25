package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type FunType struct {
  BasicType
}

func (t *FunType) Call(g *G, pos Pos, val Val, args []Form, env *Env) (Val, E) {
  f := val.AsFun()
  avs, e := VecForm(args).Eval(g, env)

  if e != nil {
    return g.NIL, g.E(pos, "Args eval failed: %v", e)
  }

  if e := f.arg_list.CheckVals(g, pos, avs); e != nil {
    return g.NIL, e
  }

  var v Val
recall:
  if f.imp == nil {
    var be Env
    f.env.Clone(&be)
    f.arg_list.PutEnv(g, &be, avs)
    
    if v, e = Forms(f.body).Eval(g, &be); e != nil {
      g.recall_args = nil
      return g.NIL, e
    }

    if g.recall_args != nil {
      avs, g.recall_args = g.recall_args, nil
      goto recall
    }
  } else {
    if v, e = f.imp(g, pos, avs, env); e != nil {
      g.recall_args = nil
      return g.NIL, e
    }
  }
  
  return v, nil
}

func (t *FunType) Dump(val Val, out *strings.Builder) {
  f := val.AsFun()
  out.WriteString("(fun (")

  for i, a := range f.arg_list.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(a.id.name)
  }

  if f.imp == nil {
    fmt.Fprintf(out, ") %v)", f.imp)
  } else {
    out.WriteString(") ")
    
    for i, bf := range f.body {
      if i > 0 {
        out.WriteRune(' ')
      }
      
      out.WriteString(bf.String())   
    }
  
    out.WriteRune(')')
  }
}

func (v Val) AsFun() *Fun {
  return v.imp.(*Fun)
}
