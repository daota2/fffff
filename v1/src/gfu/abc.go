package gfu

import (
  //"log"
  "os"
  "strings"
  "time"
)

func do_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args.EvalExpr(g, task, env)
}

func fun_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  i := 0
  id, ok := args[0].(*Sym)
  
  if ok {
    i++
  }
  
  as, e := ParseArgs(g, task, env, ParsePrimArgs(g, args[i]))

  if e != nil {
    return nil, e
  }

  i++
  f := NewFun(g, env, id, as)
  f.body = args[i:]
  return f, nil
}

func mac_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  i := 0
  id, ok := args[0].(*Sym)

  if ok {
    i++
  }

  as, e := ParseArgs(g, task, env, ParsePrimArgs(g, args[i]))

  if e != nil {
    return nil, e
  }

  i++
  m := NewMac(g, env, id, as)
  m.body = args[i:]
  return m, nil
}

func let_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  bsf := args[0]
  bs, is_scope := bsf.(Vec)

  if bsf == &g.NIL {
    bs = nil
    is_scope = true
  }
  
  var le *Env
  
  if is_scope {
    le = new(Env)

    if e := args.Extenv(g, env, le, false); e != nil {
      return nil, e
    }
  } else {
    bs = args
    le = env
  }

  for i := 0; i+1 < len(bs); i += 2 {
    kf, vf := bs[i], bs[i+1]

    if _, ok := kf.(*Sym); !ok {
      return nil, g.E("Invalid let key: %v", kf)
    }

    k := kf.(*Sym)
    v, e := vf.Eval(g, task, le)
    
    if e != nil {
      return nil, e
    }

    le.Let(k, v)
  }

  if !is_scope {
    return &g.NIL, nil
  }

  rv, e := args[1:].EvalExpr(g, task, le)

  if e != nil {
    return nil, e
  }

  return rv, nil
}

func set_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  for i := 0; i+1 < len(args); i += 2 {
    var k Val
    k, v = args[i], args[i+1]

    if _, ok := k.(*Sym); !ok {
      return nil, g.E("Invalid set key: %v", k)
    }

    if _, e = env.Set(g, k.(*Sym), v); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func if_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  c, e := args[0].Eval(g, task, env)

  if e != nil {
    return nil, e
  }

  if c.Bool(g) {
    return args[1].Eval(g, task, env)
  }

  if len(args) > 2 {
    return args[2].Eval(g, task, env)
  }

  return &g.NIL, nil
}

func inc_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  d, e := args[1].Eval(g, task, env)

  if e != nil {
    return nil, e
  }

  p := args[0]
  
  if id, ok := p.(*Sym); ok {
    return env.Update(g, id, func(v Val) (Val, E) {
      return v.(Int) + d.(Int), nil
    })
  }

  if p, e = p.Eval(g, task, env); e != nil {
    return nil, e
  }
  
  return p.(Int)+d.(Int), nil
}

func test_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  for _, in := range args {
    v, e := in.Eval(g, task, env)
    
    if e != nil {
      return nil, e
    }

    if !v.Bool(g) {
      return nil, g.E("Test failed: %v", in)
    }
  }

  return &g.NIL, nil
}

func bench_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  as := ParsePrimArgs(g, args[0])

  if as == nil {
    return nil, g.E("Invalid bench args: %v", as)
  }

  a, e := as[0].Eval(g, task, env)

  if e != nil {
    return nil, e
  }

  n := a.(Int)
  b := args[1:]

  for i := Int(0); i < n; i++ {
    b.EvalExpr(g, task, env)
  }

  t := time.Now()

  for i := Int(0); i < n; i++ {
    if _, e = b.EvalExpr(g, task, env); e != nil {
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

func dump_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var out strings.Builder

  for _, v := range args {
    v.Dump(&out)
    out.WriteRune('\n')
  }

  os.Stderr.WriteString(out.String())
  return &g.NIL, nil
}

func say_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var out strings.Builder

  for _, v := range args {
    v.Print(&out)
  }

  out.WriteRune('\n')
  os.Stdout.WriteString(out.String())
  return &g.NIL, nil
}

func load_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Load(task, env, string(args[0].(Str)))
}

func dup_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].Dup(g)
}

func clone_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].Clone(g)
}

func type_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].Type(g), nil
}

func eval_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {  
  var e E
  v := args[0]
  
  if v, e = v.Eval(g, task, env); e != nil {
      return nil, e
  }

  return v.Eval(g, task, env)
}

func expand_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  return args[1].Expand(g, task, env, args[0].(Int))
}

func recall_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return &g.NIL, NewRecall(args)
}

func new_sym_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.NewSym(string(args[0].(Str))), nil
}

func sym_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var out strings.Builder

  for _, a := range args {
    a.Print(&out);
  }
  
  return g.Sym(out.String()), nil
}

func str_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var out strings.Builder

  for _, a := range args {
    a.Print(&out);
  }
  
  return Str(out.String()), nil
}

func bool_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Bool(args[0].Bool(g)), nil
}

func eq_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]

  for _, iv := range args[1:] {
    if !iv.Eq(g, v) {
      return &g.F, nil
    }
  }

  return &g.T, nil
}

func is_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]

  for _, iv := range args[1:] {
    if !iv.Is(g, v) {
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

func int_add_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  if len(args) == 1 {
    return args[0].(Int).Abs(), nil
  }

  var v Int

  for _, iv := range args {
    v += iv.(Int)
  }

  return v, nil
}

func int_sub_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0].(Int)

  if len(args) == 1 {
    return -v, nil
  }

  for _, iv := range args[1:] {
    v -= iv.(Int)
  }

  return v, nil
}

func int_mul_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0].(Int)

  if len(args) == 1 {
    return Int(v * v), nil
  }

  for _, iv := range args[1:] {
    v *= iv.(Int)
  }

  return v, nil
}

func iter_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  if len(args) == 1 {
    return args[0].Iter(g)
  }
  
  return args.Iter(g)
}

func push_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  place := args[0]
  vs, e := args[1:].EvalVec(g, task, env)

  if e != nil {
    return nil, e
  }

  switch p := place.(type) {
  case *Sym:
    return env.Update(g, p, func(v Val) (Val, E) {
      return v.Push(g, vs...)
    })
  default:
    if place, e = place.Eval(g, task, env); e != nil {
      return nil, e
    }
  }

  return place.Push(g, vs...)
}

func pop_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var it, rest Val
  place := args[0]  
  var e E
  
  switch p := place.(type) {
  case *Sym:  
    env.Update(g, p, func(v Val) (Val, E) {
      if it, rest, e = v.Pop(g); e != nil {
        return nil, e
      }
      
      return rest, nil
    })

    return it, nil
  default:
    if place, e = place.Eval(g, task, env); e != nil {
      return nil, e
    }

    if it, rest, e = place.Pop(g); e != nil {
      return nil, e
    }    
  }

  return NewSplat(g, Vec{it, rest}), nil
}

func drop_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  place := args[0]
  var e E

  switch p := place.(type) {
  case *Sym:  
    return env.Update(g, p, func(v Val) (Val, E) {
      return v.Drop(g, args[1].(Int))
    })
  default:
    if place, e = place.Eval(g, task, env); e != nil {
      return nil, e
    }
  }

  return place.Drop(g, args[1].(Int))
}

func len_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].Len(g)
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

func pop_key_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  in, k := args[0], args[1]
  var e E
  
  if k, e = k.Eval(g, task, env); e != nil {
    return nil, e
  }
  
  if id, ok := in.(*Sym); ok {
    var v Val

    if _, e = env.Update(g, id, func(in Val) (Val, E) {
      var out Val
      
      if v, out, e = in.(Vec).PopKey(g, k.(*Sym)); e != nil {
        return nil, e
      }

      return out, nil
    }); e != nil {
      return nil, e
    }

    return v, nil
  }

  if in, e = in.Eval(g, task, env); e != nil {
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
      return &g.NIL, nil
    }
  
    return v[1:], nil
  case *Nil:
    return &g.NIL, nil
  default:
    break
  }

  return nil, g.E("Invalid tail target: %v", v) 
}

func cons_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var tail Vec

  switch a := args[1].(type) {
  case Vec:
    tail = a
  case *Nil:
    break
  default:
    return nil, g.E("Invalid cons target: %v", args[1].Type(g))
  }
  
  return append(Vec{args[0]}, tail...), nil
}

func reverse_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].(Vec).Reverse(), nil
}

func task_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var e E
  as := ParsePrimArgs(g, args[0])
  nargs := len(as)
  var inbox Chan
  safe := true
  
  if as == nil {
    inbox = NewChan(0)
  } else {
    var a Val

    if a, e = as[0].Eval(g, task, env); e != nil {
      return nil, e
    }

    if v, ok := a.(Int); ok {
      inbox = NewChan(v)
    } else if v, ok := a.(Chan); ok {
      inbox = v
    } else {
      return nil, g.E("Invalid task args: %v", as)
    }

    if nargs > 1 {
      if a, e = as[1].Eval(g, task, env); e != nil {
        return nil, e
      }

      safe = a.Bool(g)
    }
  }

  t := NewTask(g, inbox, safe, args[1:])

  if e := t.Start(g, env); e != nil {
    return nil, e
  }
  
  return t, nil
}

func task_this_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return task, nil
}

func task_post_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  t := args[0].(*Task)
  var e E
  
  if task.safe || t.safe {
    for _, v := range args[1:] {
      if v, e = v.Clone(g); e != nil {
        return nil, e
      }
      
      t.Inbox.Push(g, v)
    }
  } else {
    t.Inbox.Push(g, args[1:]...)
  }
  
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
  e.AddType(g, &g.ChanType, "Chan")
  e.AddType(g, &g.FalseType, "False")
  e.AddType(g, &g.FunType, "Fun")
  e.AddType(g, &g.IntType, "Int")
  e.AddType(g, &g.IterType, "Iter")
  e.AddType(g, &g.MacType, "Mac")
  e.AddType(g, &g.NilType, "Nil")
  e.AddType(g, &g.PrimType, "Prim")
  e.AddType(g, &g.QuoteType, "Quote")
  e.AddType(g, &g.SpliceType, "Splice")
  e.AddType(g, &g.SplatType, "Splat")
  e.AddType(g, &g.StrType, "Str")
  e.AddType(g, &g.SymType, "Sym")
  e.AddType(g, &g.TaskType, "Task")
  e.AddType(g, &g.TrueType, "True")
  e.AddType(g, &g.VecType, "Vec")

  e.AddConst(g, "_", g.NIL.Init(g))
  e.AddConst(g, "T", g.T.Init(g))
  e.AddConst(g, "F", g.F.Init(g))

  e.AddPrim(g, "do", do_imp, ASplat("body"))
  e.AddPrim(g, "fun", fun_imp, AOpt("id", nil), A("args"), ASplat("body"))
  e.AddPrim(g, "mac", mac_imp, AOpt("id", nil), A("args"), ASplat("body"))
  e.AddPrim(g, "let", let_imp, ASplat("args"))
  e.AddFun(g, "set", set_imp, ASplat("args"))
  e.AddPrim(g, "if", if_imp, A("cond"), A("t"), AOpt("f", nil))
  e.AddPrim(g, "inc", inc_imp, A("var"), AOpt("delta", Int(1)))
  e.AddPrim(g, "test", test_imp, ASplat("cases"))
  e.AddPrim(g, "bench", bench_imp, A("nreps"), ASplat("body"))

  e.AddFun(g, "debug", debug_imp) 
  e.AddFun(g, "fail", fail_imp, A("reason"))
  e.AddFun(g, "dump", dump_imp, ASplat("vals"))
  e.AddFun(g, "say", say_imp, ASplat("vals"))
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

  e.AddFun(g, "=", eq_imp, ASplat("vals"))
  e.AddFun(g, "==", is_imp, ASplat("vals"))

  e.AddFun(g, "<", int_lt_imp, ASplat("vals"))
  e.AddFun(g, ">", int_gt_imp, ASplat("vals"))
  
  e.AddFun(g, "+", int_add_imp, ASplat("vals"))
  e.AddFun(g, "-", int_sub_imp, ASplat("vals"))
  e.AddFun(g, "*", int_mul_imp, ASplat("vals"))

  e.AddFun(g, "iter", iter_imp, ASplat("vals"))
  e.AddPrim(g, "push", push_imp, A("out"), ASplat("vals"))
  e.AddPrim(g, "pop", pop_imp, A("in"))
  e.AddPrim(g, "drop", drop_imp, A("in"), AOpt("n", Int(1)))
  e.AddFun(g, "len", len_imp, A("in"))

  e.AddFun(g, "vec", vec_imp, ASplat("vals"))
  e.AddFun(g, "peek", vec_peek_imp, A("vec"))
  e.AddFun(g, "find-key", find_key_imp, A("in"), A("key"))
  e.AddPrim(g, "pop-key", pop_key_imp, A("in"), A("key"))
  e.AddFun(g, "head", head_imp, A("vec"))
  e.AddFun(g, "tail", tail_imp, A("vec"))
  e.AddFun(g, "cons", cons_imp, A("val"), A("vec"))
  e.AddFun(g, "reverse", reverse_imp, A("vec"))

  e.AddPrim(g, "task", task_imp, A("args"), ASplat("body"))
  e.AddFun(g, "this-task", task_this_imp)
  e.AddFun(g, "post", task_post_imp, A("task"), ASplat("vals"))
  e.AddFun(g, "fetch", task_fetch_imp)
  e.AddFun(g, "wait", task_wait_imp, ASplat("tasks"))
  e.AddFun(g, "chan", chan_imp, AOpt("buf", Int(0)))
}
