package gfu

import (
  "bufio"
  "fmt"
  //"log"
  "strings"
  "time"
)

func do_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  return args.EvalExpr(g, task, env, args_env)
}

func fun_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  i := 0
  id, ok := args[0].(*Sym)

  if ok {
    i++
  }

  as, e := ParseArgs(g, task, env, ParsePrimArgs(g, args[i]), args_env)

  if e != nil {
    return nil, e
  }

  i++
  f, e := NewFun(g, env, id, as)

  if e != nil {
    return nil, e
  }

  f.body = args[i:]

  if e = f.InitEnv(g, env); e != nil {
    return nil, e
  }

  return f, nil
}

func mac_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  i := 0
  id, ok := args[0].(*Sym)

  if ok {
    i++
  }

  as, e := ParseArgs(g, task, env, ParsePrimArgs(g, args[i]), args_env)

  if e != nil {
    return nil, e
  }

  i++
  m, e := NewMac(g, env, id, as)

  if e != nil {
    return nil, e
  }

  m.body = args[i:]

  if e = m.InitEnv(g, env); e != nil {
    return nil, e
  }
  
  return m, nil
}

func call_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  t, e := g.Eval(task, args_env, args[0], args_env)

  if e != nil {
    return nil, e
  }

  return g.Call(task, env, t, args[1:], env)
}

func let_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (v Val, e E) {
  if len(args) == 0 {
    return &g.NIL, nil
  }

  bsf := args[0]
  bs, is_scope := bsf.(Vec)

  if bsf == &g.NIL {
    bs = nil
    is_scope = true
  }

  var le *Env

  if is_scope {
    le = new(Env)

    if e = g.Extenv(env, le, args, false); e != nil {
      return nil, e
    }
  } else {
    bs = args
    le = env
  }

  if e = g.Extenv(args_env, le, args, false); e != nil {
    return nil, e
  }

  v = &g.NIL

  for i := 0; i+1 < len(bs); i += 2 {
    kf, vf := bs[i], bs[i+1]

    if _, ok := kf.(*Sym); !ok {
      return nil, g.E("Invalid let key: %v", kf)
    }

    k := kf.(*Sym)
    v, e = g.Eval(task, le, vf, args_env)

    if e != nil {
      return nil, e
    }

    if e = le.Let(g, k, v); e != nil {
      return nil, e
    }
  }

  if !is_scope {
    return v, nil
  }

  rv, e := args[1:].EvalExpr(g, task, le, args_env)

  if e != nil {
    return nil, e
  }

  return rv, nil
}

func val_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  if v, _, _, e = args[0].(*Sym).Lookup(g, task, env, env, true); e != nil {
    return nil, e
  }

  if v == nil { v = &g.NIL }
  return v, nil
}

func set_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (v Val, e E) {
  for i := 0; i+1 < len(args); i += 2 {
    var k Val
    k, v = args[i], args[i+1]
    
    if v, e = g.Eval(task, env, v, args_env); e != nil {
      return nil, e
    }

    if e = env.Set(g, task, k, v, args_env); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func use_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  prefix := args[0]

  if prefix == &g.NIL {
    return &g.NIL, nil
  }

  var v *Var
  var e E

  for _, k := range args[1:] {
    ks := g.Sym(fmt.Sprintf("%v/%v", prefix.(*Sym), k))

    if v, _, _, _, e = ks.LookupVar(g, args_env, nil, false); e != nil {
      return nil, e
    }

    if i, found := env.Find(v.key); found == nil {
      env.InsertVar(i, v)
    }
  }

  return &g.NIL, nil
}

func env_this_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  return env, nil
}

func if_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  c, e := g.Eval(task, env, args[0], args_env)

  if e != nil {
    return nil, e
  }

  v, e := g.Bool(c)

  if e != nil {
    return nil, e
  }

  if v {
    return g.Eval(task, args_env, args[1], args_env)
  }

  if len(args) > 2 {
    return g.Eval(task, args_env, args[2], args_env)
  }

  return &g.NIL, nil
}

func inc_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  d, e := g.Eval(task, args_env, args[1], args_env)

  if e != nil {
    return nil, e
  }

  p := args[0]

  switch p.(type) {
  case *Sym, Vec:
    return args_env.Update(g, task, p, func(v Val) (Val, E) {
      return g.Add(v, d)
    }, args_env)
  }

  if p, e = g.Eval(task, args_env, p, args_env); e != nil {
    return nil, e
  }

  return g.Add(p, d)
}

func test_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  for _, in := range args {
    v, e := g.Eval(task, env, in, args_env)

    if e != nil {
      return nil, e
    }

    bv, e := g.Bool(v)

    if e != nil {
      return nil, e
    }

    if !bv {
      return nil, g.E("Test failed: %v", in)
    }
  }

  return &g.NIL, nil
}

func bench_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  as := ParsePrimArgs(g, args[0])

  if as == nil {
    return nil, g.E("Invalid bench args: %v", as)
  }

  a, e := g.Eval(task, args_env, as[0], args_env)

  if e != nil {
    return nil, e
  }

  n := a.(Int)
  b := args[1:]

  for i := Int(0); i < n; i++ {
    b.EvalExpr(g, task, env, args_env)
  }

  t := time.Now()

  for i := Int(0); i < n; i++ {
    if _, e = b.EvalExpr(g, task, env, args_env); e != nil {
      return nil, e
    }
  }

  return Int(time.Now().Sub(t).Nanoseconds() / 1000000), nil
}

func debug_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  g.Debug = true
  return &g.NIL, nil
}

func fail_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return nil, g.E(string(args[0].(Str)))
}

func load_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Load(task, env, string(args[0].(Str)))
}

func dup_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Dup(args[0])
}

func clone_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Clone(args[0])
}

func type_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].Type(g), nil
}

func eval_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  var e E
  v := args[0]

  if v, e = g.Eval(task, args_env, v, args_env); e != nil {
    return nil, e
  }

  return g.Eval(task, env, v, args_env)
}

func expand_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  return g.Expand(task, env, args[1], args[0].(Int))
}

func recall_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return &g.NIL, NewRecall(args)
}

func new_sym_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.NewSym(string(args[0].(Str))), nil
}

func sym_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var out strings.Builder
  w := bufio.NewWriter(&out)
  
  for _, a := range args {
    g.Print(a, w)
  }

  w.Flush()
  return g.Sym(out.String()), nil
}

func str_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var out strings.Builder
  w := bufio.NewWriter(&out)

  for _, a := range args {
    g.Print(a, w)
  }

  w.Flush()
  return Str(out.String()), nil
}

func bool_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.BoolVal(args[0])
}

func float_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Float(args[0])
}

func int_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Int(args[0])
}

func eq_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]

  for _, iv := range args[1:] {
    ok, e := g.Eq(iv, v)

    if e != nil {
      return nil, e
    }

    if !ok {
      return &g.F, nil
    }
  }

  return &g.T, nil
}

func is_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]

  for _, iv := range args[1:] {
    if !g.Is(iv, v) {
      return &g.F, nil
    }
  }

  return &g.T, nil
}

func int_lt_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  lhs := args[0].(Int)

  for _, a := range args[1:] {
    rhs := a.(Int)

    if rhs <= lhs {
      return &g.F, nil
    }

    lhs = rhs
  }

  return &g.T, nil
}

func int_gt_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  lhs := args[0].(Int)

  for _, a := range args[1:] {
    rhs := a.(Int)

    if rhs >= lhs {
      return &g.F, nil
    }

    lhs = rhs
  }

  return &g.T, nil
}

func add_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  a0 := args[0]

  if len(args) == 1 {
    return g.Abs(a0)
  }

  v = args[0]

  for _, a := range args[1:] {
    if v, e = g.Add(v, a); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func sub_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  a0 := args[0]

  if len(args) == 1 {
    return g.Neg(a0)
  }

  v = args[0]

  for _, a := range args[1:] {
    if v, e = g.Sub(v, a); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func mul_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  a0 := args[0]

  if len(args) == 1 {
    return g.Mul(a0, a0)
  }

  v = args[0]

  for _, a := range args[1:] {
    if v, e = g.Mul(v, a); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func div_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  a0 := args[0]

  if len(args) == 1 {
    var x, y Float
    x.SetInt(1)
    y.SetInt(a0.(Int))
    x.Div(y)
    return x, nil
  }

  v = args[0]

  for _, a := range args[1:] {
    if v, e = g.Div(v, a); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func iter_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Iter(args[0])
}

func push_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  place := args[0]
  vs, e := args[1:].EvalVec(g, task, args_env)

  if e != nil {
    return nil, e
  }

  switch p := place.(type) {
  case *Sym:
    return args_env.Update(g, task, p, func(v Val) (Val, E) {
      return g.Push(v, vs...)
    }, args_env)
  default:
    if place, e = g.Eval(task, args_env, place, args_env); e != nil {
      return nil, e
    }
  }

  return g.Push(place, vs...)
}

func pop_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  var it, rest Val
  place := args[0]
  var e E

  switch p := place.(type) {
  case *Sym:
    args_env.Update(g, task, p, func(v Val) (Val, E) {
      if it, rest, e = g.Pop(v); e != nil {
        return nil, e
      }

      return rest, nil
    }, args_env)

    return it, nil
  default:
    if place, e = g.Eval(task, args_env, place, args_env); e != nil {
      return nil, e
    }

    if it, rest, e = g.Pop(place); e != nil {
      return nil, e
    }
  }

  return NewSplat(g, Vec{it, rest}), nil
}

func drop_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  place := args[0]
  var e E

  switch p := place.(type) {
  case *Sym:
    return args_env.Update(g, task, p, func(v Val) (Val, E) {
      return g.Drop(v, args[1].(Int))
    }, args_env)
  default:
    if place, e = g.Eval(task, args_env, place, args_env); e != nil {
      return nil, e
    }
  }

  return g.Drop(place, args[1].(Int))
}

func len_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Len(args[0])
}

func index_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Index(args[0], args[1:])
}

func set_index_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.SetIndex(args[1], args[2:], args[0].(Setter))
}

func vec_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args, nil
}

func vec_peek_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].(Vec).Peek(g), nil
}

func find_key_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  in, k := args[0].(Vec), args[1]

  for i := 0; i < len(in)-1; i += 2 {
    if in[i] == k {
      return in[i+1], nil
    }
  }

  return &g.NIL, nil
}

func pop_key_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  in, k := args[0], args[1]
  var e E

  if k, e = g.Eval(task, args_env, k, args_env); e != nil {
    return nil, e
  }

  if id, ok := in.(*Sym); ok {
    var v Val

    if _, e = args_env.Update(g, task, id, func(in Val) (Val, E) {
      var out Val

      if v, out, e = in.(Vec).PopKey(g, k.(*Sym)); e != nil {
        return nil, e
      }

      return out, nil
    }, args_env); e != nil {
      return nil, e
    }

    return v, nil
  }

  if in, e = g.Eval(task, args_env, in, args_env); e != nil {
    return nil, e
  }

  var v, out Val

  if v, out, e = in.(Vec).PopKey(g, k.(*Sym)); e != nil {
    return nil, e
  }

  return NewSplat(g, Vec{v, out}), nil
}

func head_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]

  switch v := v.(type) {
  case Vec:
    if len(v) == 0 {
      return &g.NIL, nil
    }

    return v[0], nil
  case *Nil:
    return &g.NIL, nil
  default:
    break
  }

  return nil, g.E("Invalid head target: %v", v)
}

func tail_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]

  switch v := v.(type) {
  case Vec:
    if len(v) < 2 {
      return Vec(nil), nil
    }

    return v[1:], nil
  default:
    break
  }

  return nil, g.E("Invalid tail target: %v", v)
}

func reverse_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].(Vec).Reverse(), nil
}

func new_bin_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return NewBin(int(args[0].(Int))), nil
}

func bin_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  b := NewBin(len(args))
  
  for i, v := range args {
    bv, ok := v.(Byte)

    if !ok {
      return nil, g.E("Expected Byte: %v", v.Type(g))
    }

    b[i] = byte(bv)
  }
  
  return b, nil
}

func task_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  id, ok := args[0].(*Sym)
  i := 0

  if ok {
    i++
  }

  as := ParsePrimArgs(g, args[i])
  var inbox Chan
  var e E

  if as == nil {
    inbox = NewChan(0)
  } else {
    var a Val

    if a, e = g.Eval(task, args_env, as[0], args_env); e != nil {
      return nil, e
    }

    if v, ok := a.(Int); ok {
      inbox = NewChan(v)
    } else if v, ok := a.(Chan); ok {
      inbox = v
    } else {
      return nil, g.E("Invalid task args: %v", as)
    }
  }

  i++
  t, e := NewTask(g, env, id, inbox, args[i:])

  if e != nil {
    return nil, e
  }

  if e = t.Start(g, env); e != nil {
    return nil, e
  }

  return t, nil
}

func this_task_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  return task, nil
}

func task_post_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  t := args[0].(*Task)
  t.Inbox.Push(g, args[1:]...)
  return t, nil
}

func task_fetch_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := <-task.Inbox

  if v == nil {
    v = &g.NIL
  }

  return v, nil
}

func task_wait_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  nargs := len(args)

  if nargs == 1 {
    return args[0].(*Task).Wait(), nil
  }

  out := make(Vec, nargs)

  for i, a := range args {
    out[i] = a.(*Task).Wait()
  }

  return out, nil
}

func chan_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return NewChan(args[0].(Int)), nil
}

func (e *Env) InitAbc(g *G) {
  e.AddType(g, &g.MetaType, "Meta")

  e.AddType(g, &g.NumType, "Num")
  e.AddType(g, &g.SeqType, "Seq")

  e.AddType(g, &g.IterType, "Iter", &g.SeqType)

  e.AddType(g, &g.BinType, "Bin", &g.SeqType)
  e.AddType(g, &g.BinIterType, "BinIter", &g.SeqType)
  e.AddType(g, &g.ByteType, "Byte", &g.NumType)
  e.AddType(g, &g.ChanType, "Chan")
  e.AddType(g, &g.CharType, "Char")
  e.AddType(g, &g.EnvType, "Env")
  e.AddType(g, &g.FalseType, "False")
  e.AddType(g, &g.FloatType, "Float", &g.NumType)
  e.AddType(g, &g.FunType, "Fun")
  e.AddType(g, &g.IntType, "Int", &g.NumType)
  e.AddType(g, &g.IntIterType, "IntIter", &g.IterType)
  e.AddType(g, &g.MacType, "Mac")
  e.AddType(g, &g.NilType, "Nil")
  e.AddType(g, &g.PrimType, "Prim")
  e.AddType(g, &g.QuoteType, "Quote")
  e.AddType(g, &g.SetterType, "Setter")
  e.AddType(g, &g.SpliceType, "Splice")
  e.AddType(g, &g.SplatType, "Splat")
  e.AddType(g, &g.StrType, "Str")
  e.AddType(g, &g.SymType, "Sym")
  e.AddType(g, &g.TaskType, "Task")
  e.AddType(g, &g.TrueType, "True")
  e.AddType(g, &g.VecType, "Vec", &g.SeqType)
  e.AddType(g, &g.VecIterType, "VecIter", &g.IterType)
  e.AddType(g, &g.WriterType, "Writer")

  e.AddConst(g, "_", &g.NIL)
  e.AddConst(g, "T", &g.T)
  e.AddConst(g, "F", &g.F)
  e.AddConst(g, "\\e", Char('\x1b'))
  e.AddConst(g, "\\n", Char('\n'))

  e.AddPrim(g, "do", do_imp, ASplat("body"))
  e.AddPrim(g, "fun", fun_imp, AOpt("id", nil), A("args"), ASplat("body"))
  e.AddPrim(g, "mac", mac_imp, AOpt("id", nil), A("args"), ASplat("body"))
  e.AddPrim(g, "call", call_imp, A("target"), ASplat("args"))
  e.AddPrim(g, "let", let_imp, ASplat("args"))
  e.AddFun(g, "val", val_imp, A("key"))
  e.AddPrim(g, "set", set_imp, ASplat("args"))
  e.AddPrim(g, "use", use_imp, AOpt("prefix", nil), ASplat("ids"))
  g.EnvType.Env().AddPrim(g, "this", env_this_imp)
  e.AddPrim(g, "if", if_imp, A("cond"), A("t"), AOpt("f", nil))
  e.AddPrim(g, "inc", inc_imp, A("var"), AOpt("delta", Int(1)))
  e.AddPrim(g, "test", test_imp, ASplat("cases"))
  e.AddPrim(g, "bench", bench_imp, A("nreps"), ASplat("body"))

  e.AddFun(g, "debug", debug_imp)
  e.AddFun(g, "fail", fail_imp, A("reason"))
  e.AddFun(g, "load", load_imp, A("path"))

  e.AddFun(g, "dup", dup_imp, A("val"))
  e.AddFun(g, "clone", clone_imp, A("val"))
  e.AddFun(g, "type", type_imp, A("val"))
  e.AddPrim(g, "eval", eval_imp, A("expr"))
  e.AddFun(g, "expand", expand_imp, A("n"), A("expr"))
  e.AddFun(g, "recall", recall_imp, ASplat("args"))
  e.AddFun(g, "new-sym", new_sym_imp, AOpt("prefix", Str("")))
  e.AddFun(g, "sym", sym_imp, ASplat("args"))
  e.AddFun(g, "str", str_imp, ASplat("args"))

  e.AddFun(g, "bool", bool_imp, A("val"))
  e.AddFun(g, "float", float_imp, A("val"))
  e.AddFun(g, "int", int_imp, A("val"))

  e.AddFun(g, "=", eq_imp, ASplat("vals"))
  e.AddFun(g, "==", is_imp, ASplat("vals"))

  e.AddFun(g, "<", int_lt_imp, ASplat("vals"))
  e.AddFun(g, ">", int_gt_imp, ASplat("vals"))

  e.AddFun(g, "+", add_imp, A("x"), ASplat("ys"))
  e.AddFun(g, "/", div_imp, A("x"), ASplat("ys"))
  e.AddFun(g, "-", sub_imp, A("x"), ASplat("ys"))
  e.AddFun(g, "*", mul_imp, A("x"), ASplat("ys"))

  e.AddFun(g, "iter", iter_imp, A("val"))
  e.AddPrim(g, "push", push_imp, A("out"), ASplat("vals"))
  e.AddPrim(g, "pop", pop_imp, A("in"))
  e.AddPrim(g, "drop", drop_imp, A("in"), AOpt("n", Int(1)))
  e.AddFun(g, "len", len_imp, A("in"))

  e.AddFun(g, "#", index_imp, A("source"), ASplat("key"))
  e.AddFun(g, "set-#", set_index_imp, A("set"), A("dest"), ASplat("key"))
  
  e.AddFun(g, "vec", vec_imp, ASplat("vals"))
  e.AddFun(g, "peek", vec_peek_imp, A("vec"))
  e.AddFun(g, "find-key", find_key_imp, A("in"), A("key"))
  e.AddPrim(g, "pop-key", pop_key_imp, A("in"), A("key"))
  e.AddFun(g, "head", head_imp, A("vec"))
  e.AddFun(g, "tail", tail_imp, A("vec"))
  e.AddFun(g, "reverse", reverse_imp, A("vec"))

  e.AddFun(g, "new-bin", new_bin_imp, AOpt("len", Int(0)))
  e.AddFun(g, "bin", bin_imp, ASplat("vals"))
  
  e.AddPrim(g, "task", task_imp, A("args"), ASplat("body"))
  e.AddPrim(g, "this-task", this_task_imp)
  g.TaskType.Env().AddFun(g, "post", task_post_imp, A("task"), ASplat("vals"))
  e.AddFun(g, "fetch", task_fetch_imp)
  e.AddFun(g, "wait", task_wait_imp, ASplat("tasks"))
  e.AddFun(g, "chan", chan_imp, AOpt("buf", Int(0)))
}
