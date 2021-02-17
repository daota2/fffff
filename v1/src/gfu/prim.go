package gfu

type PrimImp func (*G, ListForm, *Env, Pos) (Val, Error)

type Prim struct {
  id *Sym
  min_args, max_args int
  imp PrimImp
}

func NewPrim(id *Sym, min_args, max_args int, imp PrimImp) *Prim {
  p := new(Prim)
  p.id = id
  p.min_args, p.max_args = min_args, max_args
  p.imp = imp
  return p
}

func (e *Env) AddPrim(g *G, id *Sym, min_args, max_args int, imp PrimImp) {
  var p Val
  p.Init(g.Prim, NewPrim(id, min_args, max_args, imp))
  e.Put(id, p)
}
