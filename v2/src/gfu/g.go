package gfu

import (
	"fmt"
	"io/ioutil"
	//"log"
	"path/filepath"
	"strings"
	"sync"
)

type G struct {
	Debug    bool
	MainTask Task
	RootEnv  Env

	AbortType   AbortType
	BinType     BinType
	BinIterType BinIterType
	ByteType    ByteType
	ChanType    ChanType
	CharType    CharType
	EType       EType
	EnvType     EnvType
	FalseType   FalseType
	FloatType   FloatType
	FunType     FunType
	IntType     IntType
	IntIterType IntIterType
	IterType    IterType
	MacType     MacType
	MetaType    MetaType
	NilType     NilType
	NanosType   NanosType
	NumType     NumType
	OpsType     OpsType
	PrimType    PrimType
	PunType     PunType
	QuoteType   QuoteType
	RecallType  RecallType
	RestartType RestartType
	RetryType   RetryType
	SetterType  SetterType
	SeqType     SeqType
	SplatType   SplatType
	SpliceType  SpliceType
	StrType     StrType
	SymType     SymType
	TaskType    TaskType
	TimeType    TimeType
	TrueType    TrueType
	VecType     VecType
	VecIterType VecIterType
	WriterType  WriterType

	EImpureType  EImpureType
	EReadType    EReadType
	EUnknownType EUnknownType

	NIL Nil
	T   True
	F   False

	syms,
	consts sync.Map
	nsyms uint64

	let_sym,
	nil_sym,
	nop_sym,
	resolve_sym,
	set_sym *Sym

	abort_fun,
	retry_fun *Fun

	load_path string
}

func NewG() (*G, E) {
	return new(G).Init()
}

func (g *G) Init() (*G, E) {
	g.let_sym = g.Sym("let")
	g.nil_sym = g.Sym("_")
	g.nop_sym = g.Sym("__")
	g.resolve_sym = g.Sym("resolve")
	g.set_sym = g.Sym("set")
	g.MainTask.Init(g, &g.RootEnv, g.Sym("main-task"), NewChan(0), nil)
	return g, nil
}

func (g *G) Version() [3]byte {
	return [3]byte{2, 0, 0}
}

func (g *G) VersionStr() string {
	var s strings.Builder

	for i, v := range g.Version() {
		if i > 0 {
			s.WriteRune('.')
		}

		fmt.Fprintf(&s, "%v", v)
	}

	return s.String()
}

func (g *G) NewEnv() *Env {
	var env Env
	g.RootEnv.Dup(&env)
	return &env
}

func (g *G) EvalString(task *Task, env *Env, pos Pos, s string) (Val, E) {
	in := strings.NewReader(s)
	out, e := g.ReadAll(&pos, in, nil)

	if e != nil {
		return nil, e
	}

	if out, e = out.ExpandVec(g, task, env, -1); e != nil {
		return nil, e
	}

	return out.EvalExpr(g, &g.MainTask, env, env)
}

func (g *G) Load(task *Task, env, args_env *Env, path string, eval bool) (v Val, e E) {
	use_filename := NewFun(g, env, g.Sym("use-filename"), A("new"))
	use_filename.imp = func(g *G, task *Task, env *Env, args Vec) (Val, E) {
		ps, ok := args[0].(Str)

		if !ok {
			return nil, g.E("Invalid filename: \"%v\"", args[0].Type(g))
		}

		path = string(ps)
		return nil, Retry{}
	}

	var bs []byte

	g.Try(task, env, args_env, func() (Val, E) {
		path = filepath.Join(g.load_path, path)
		var re error
		bs, re = ioutil.ReadFile(path)

		if re != nil {
			return nil, g.E("Failed loading file: \"%v\"\n%v", path, re)
		}

		return &g.NIL, nil
	}, use_filename)

	var pos Pos
	pos.Init(path)
	prev_path := g.load_path
	defer func() { g.load_path = prev_path }()
	g.load_path = filepath.Dir(path)

	if v, e = g.ReadAll(&pos, strings.NewReader(string(bs)), nil); e != nil {
		return nil, e
	}

	if !eval {
		return v, nil
	}

	return v.(Vec).EvalExpr(g, &g.MainTask, env, env)
}

func (env *Env) AddConst(g *G, id string, val Val) E {
	if _, dup := g.consts.LoadOrStore(g.Sym(id), val); dup {
		return g.E("Dup const: %v", id)
	}

	imp := func(g *G, task *Task, env *Env, args Vec) (Val, E) {
		ok, e := g.Eq(val, args[0])

		if e != nil {
			return nil, e
		}

		if !ok {
			return &g.F, nil
		}

		return &g.T, nil
	}

	_, e := env.AddPun(g, fmt.Sprintf("%v?", id), imp, A("val"))
	return e
}

func (g *G) FindConst(id *Sym) Val {
	v, ok := g.consts.Load(id)

	if !ok {
		return nil
	}

	return v.(Val)
}
