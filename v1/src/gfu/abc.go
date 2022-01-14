package gfu

import (
  "fmt"
  //"log"
  "os"
  "strings"
  "time"
)

func do_imp(g *G, args Vec, env *Env) (Val, E) {
  return args.EvalExpr(g, env)
}

func fun_imp(g *G, args Vec, env *Env) (Val, E) {
  avs := args[0]
  
  if avs.val_type != g.VecType {
    return g.NIL, g.E("Invalid args: %v", avs)
  }

  as, e := Args(avs.AsVec()).Parse(g)

  if e != nil {
    return g.NIL, e
  }

  f := NewFun(g, env, as)
  f.body = args[1:]
  
  var v Val
  v.Init(g.FunType, f)
  return v, nil
}

func macro_imp(g *G, args Vec, env *Env) (Val, E) {
  avs := args[0]
  
  if avs.val_type != g.VecType {
    return g.NIL, g.E("Invalid args: %v", avs)
  }

  as, e := Args(avs.AsVec()).Parse(g)

  if e != nil {
    return g.NIL, e
  }

  m := NewMacro(g, env, as)
  m.body = args[1:]
  
  var v Val
  v.Init(g.MacroType, m)
  return v, nil
}

func let_imp(g *G, args Vec, env *Env) (Val, E) {
  bsf := args[0]
  var bs Vec
  is_scope := bsf.val_type == g.VecType
  
  if is_scope {
    bs = bsf.AsVec()
  } else {
    bs = args
  }

  var le *Env

  if is_scope {
    le = new(Env)
    env.Clone(le)
  } else {
    le = env
  }
    
  for i := 0; i < len(bs); i += 2 {
    kf, vf := bs[i], bs[i+1]

    if kf.val_type != g.SymType {
      return g.NIL, g.E("Invalid let key: %v", kf)
    }

    k := kf.AsSym()
    v, e := vf.Eval(g, le)

    if e != nil {
      return g.NIL, e
    }

    le.Put(k, v)
  }

  if !is_scope {
    return g.NIL, nil
  }
  
  rv, e := args[1:].EvalExpr(g, le)
  
  if e != nil {
    return g.NIL, e
  }
  
  return rv, nil
}

func if_imp(g *G, args Vec, env *Env) (Val, E) {
  c, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  if c.AsBool(g) {
    return args[1].Eval(g, env)
  }

  if len(args) > 2 {
    return args[2].Eval(g, env)
  }

  return g.NIL, nil
}

func and_imp(g *G, args Vec, env *Env) (Val, E) {
  var e E
  v := g.NIL
  
  for _, in := range args {
    v, e = in.Eval(g, env)

    if e != nil {
      return g.NIL, e
    }
    
    if !v.AsBool(g) {
      return g.F, nil
    }
  }

  return v, nil
}

func or_imp(g *G, args Vec, env *Env) (Val, E) {
  for _, in := range args {
    v, e := in.Eval(g, env)

    if e != nil {
      return g.NIL, e
    }
    
    if v.AsBool(g) {
      return v, nil
    }
  }

  return g.F, nil
}

func inc_imp(g *G, args Vec, env *Env) (Val, E) {
  id := args[0].AsSym()
  _, found := env.Find(id)

  if found == nil {
    return g.NIL, g.E("Unknown var: %v", id)
  }

  d := 1
  
  if len(args) == 2 {
    dv, e := args[1].Eval(g, env)

    if e != nil {
      return g.NIL, e
    }

    d = dv.AsInt()
  }

  v := &found.Val
  v.imp = v.AsInt() + d
  return *v, nil
}

func for_imp(g *G, args Vec, env *Env) (Val, E) {
  nv, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  n := nv.AsInt()
  b := args[1:]
  v := g.NIL
  
  for i := 0; i < n; i++ {
    if v, e = b.EvalExpr(g, env); e != nil {
      return g.NIL, e
    }
  }
  
  return v, nil
}

func test_imp(g *G, args Vec, env *Env) (Val, E) {
  for _, in := range args {
    v, e := in.Eval(g, env)

    if e != nil {
      return g.NIL, e
    }

    if !v.AsBool(g) {
      return g.NIL, g.E("Test failed: %v", in)
    }
  }

  return g.NIL, nil
}

func bench_imp(g *G, args Vec, env *Env) (Val, E) {
  nv, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  n := nv.AsInt()
  b := args[1:]

  for i := 0; i < n; i++ {
    b.EvalExpr(g, env)
  }

  t := time.Now()
  
  for i := 0; i < n; i++ {
    if _, e = b.EvalExpr(g, env); e != nil {
      return g.NIL, e
    }
  }

  var v Val
  v.Init(g.IntType, time.Now().Sub(t).Nanoseconds() / 1000000) 
  return v, nil
}

func dump_imp(g *G, args Vec, env *Env) (Val, E) {
  var out strings.Builder
  
  for _, v := range args {
    v.Dump(&out)
    out.WriteRune('\n')
  }

  os.Stderr.WriteString(out.String())
  return g.NIL, nil
}

func eval_imp(g *G, args Vec, env *Env) (Val, E) {
  return args[0].Eval(g, env)
}

func recall_imp(g *G, args Vec, env *Env) (Val, E) {
  if g.recall {
    return g.NIL, g.E("Recall already in progress")
  }

  g.recall = true
  g.recall_args = args
  return g.NIL, nil
}

func g_sym_imp(g *G, args Vec, env *Env) (Val, E) {
  var n string
  t := g.NextSymTag()
    
  if len(args) > 0 {
    n = fmt.Sprintf("g-%v-%v", args[0], t)
  } else {
    n = fmt.Sprintf("g-%v", t)
  }
  
  var out Val
  out.Init(g.SymType, g.Sym(n))
  return out, nil
}

func bool_imp(g *G, args Vec, env *Env) (Val, E) {
  v := args[0]
  v.Init(g.BoolType, v.AsBool(g))
  return v, nil
}

func not_imp(g *G, args Vec, env *Env) (Val, E) {
  v := args[0]
  v.Init(g.BoolType, !v.AsBool(g))
  return v, nil
}

func eq_imp(g *G, args Vec, env *Env) (Val, E) {
  v := args[0]
  
  for _, iv := range args[1:] {
    if !iv.Eq(g, v) {
      v.Init(g.BoolType, false)
      return v, nil
    }
  }
  
  v.Init(g.BoolType, true)
  return v, nil
}

func is_imp(g *G, args Vec, env *Env) (Val, E) {
  v := args[0]
  
  for _, iv := range args[1:] {
    if !iv.Is(g, v) {
      v.Init(g.BoolType, false)
      return v, nil
    }
  }
  
  v.Init(g.BoolType, true)
  return v, nil
}

func int_lt_imp(g *G, args Vec, env *Env) (Val, E) {
  v := args[0]
  a0 := v.AsInt()
  
  for _, a := range args[1:] {
    if a.AsInt() <= a0 {
      v.Init(g.BoolType, false)
      return v, nil
    }
  }
  
  v.Init(g.BoolType, true)
  return v, nil
}

func int_add_imp(g *G, args Vec, env *Env) (Val, E) {
  if len(args) == 1 {
    v := args[0]
    v.imp = Abs(v.AsInt())
    return v, nil
  }
  
  var v int

  for _, iv := range args {
    v += iv.AsInt()
  }
  
  var out Val
  out.Init(g.IntType, v)
  return out, nil
}

func int_sub_imp(g *G, args Vec, env *Env) (Val, E) {
  var out Val
  v := args[0].AsInt()

  if len(args) == 1 {
    out.Init(g.IntType, -v)
  } else {    
    for _, iv := range args[1:] {
      v -= iv.AsInt()
    }
    
    out.Init(g.IntType, v)
  }
  
  return out, nil
}

func vec_imp(g *G, args Vec, env *Env) (Val, E) {
  var out Val
  out.Init(g.VecType, args)
  return out, nil
}

func vec_len_imp(g *G, args Vec, env *Env) (Val, E) {
  v := args[0]  
  v.Init(g.IntType, len(v.AsVec()))
  return v, nil
}

func vec_push_imp(g *G, args Vec, env *Env) (Val, E) {
  id := args[0].AsSym()
  _, found := env.Find(id)

  if found == nil {
    return g.NIL, g.E("Unknown var: %v", id)
  }

  v := &found.Val
  v.imp = v.AsVec().Push(args[1:]...)
  return *v, nil
}

func vec_peek_imp(g *G, args Vec, env *Env) (Val, E) {
  return args[0].AsVec().Peek(g), nil
}

func vec_pop_imp(g *G, args Vec, env *Env) (Val, E) {
  id := args[0].AsSym()
  _, found := env.Find(id)

  if found == nil {
    return g.NIL, g.E("Unknown var: %v", id)
  }

  v := &found.Val
  var out Val
  out, v.imp = v.AsVec().Pop(g)
  return out, nil
}

func (e *Env) InitAbc(g *G) {
  g.MetaType = e.AddType(g, "Meta", new(MetaType))
  g.BoolType = e.AddType(g, "Bool", new(BoolType))
  g.FunType = e.AddType(g, "Fun", new(FunType))
  g.IntType = e.AddType(g, "Int", new(IntType))
  g.MacroType = e.AddType(g, "Macro", new(MacroType))
  g.NilType = e.AddType(g, "Nil", new(NilType))
  g.OptType = e.AddType(g, "Opt", new(OptType))
  g.PrimType = e.AddType(g, "Prim", new(PrimType))
  g.QuoteType = e.AddType(g, "Quote", new(QuoteType))
  g.SpliceType = e.AddType(g, "Splice", new(SpliceType))
  g.SplatType = e.AddType(g, "Splat", new(SplatType))
  g.SymType = e.AddType(g, "Sym", new(SymType))
  g.VecType = e.AddType(g, "Vec", new(VecType))
  
  e.AddVal(g, "_", g.NilType, nil, &g.NIL)
  e.AddVal(g, "T", g.BoolType, true, &g.T)
  e.AddVal(g, "F", g.BoolType, false, &g.F)
  
  e.AddPrim(g, "do", do_imp, "body..")
  e.AddPrim(g, "fun", fun_imp, "args", "body..")
  e.AddPrim(g, "macro", macro_imp, "args", "body..")
  e.AddPrim(g, "let", let_imp, "args..")
  e.AddPrim(g, "if", if_imp, "cond", "t", "f?")
  e.AddPrim(g, "or", or_imp, "conds..")
  e.AddPrim(g, "and", and_imp, "conds..")
  e.AddPrim(g, "inc", inc_imp, "var", "delta?")
  e.AddPrim(g, "for", for_imp, "nreps", "body..")
  e.AddPrim(g, "test", test_imp, "cases..")
  e.AddPrim(g, "bench", bench_imp, "nreps", "body..")

  e.AddFun(g, "dump", dump_imp, "vals..")
  e.AddFun(g, "eval", eval_imp, "form")
  e.AddFun(g, "recall", recall_imp, "args..")
  e.AddFun(g, "g-sym", g_sym_imp, "prefix?")

  e.AddFun(g, "bool", bool_imp, "val")
  e.AddFun(g, "not", not_imp, "val")
  
  e.AddFun(g, "=", eq_imp, "vals..")
  e.AddFun(g, "==", is_imp, "vals..")
  
  e.AddFun(g, "<", int_lt_imp, "vals..")
  e.AddFun(g, "+", int_add_imp, "vals..")
  e.AddFun(g, "-", int_sub_imp, "vals..")

  e.AddFun(g, "vec", vec_imp, "items..")
  e.AddFun(g, "len", vec_len_imp, "vec")
  e.AddPrim(g, "push", vec_push_imp, "vec val..")
  e.AddFun(g, "peek", vec_peek_imp, "vec")
  e.AddPrim(g, "pop", vec_pop_imp, "vec")
}
