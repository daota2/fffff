![Logo](../logo.png)
  
### Intro
[g-fu](https://github.com/codr7/g-fu) is a pragmatic [Lisp](https://xkcd.com/297/) developed and embedded in Go.

This document describes the initial release; which implements an extensible, tree-walking interpreter with quasi-quotation and macros, lambdas, optimized tail-recursion, opt-/varargs, threads and channels; weighing in at 3 kloc.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.14

Press Return twice to evaluate.

  (fun fib (n)
    (if (< n 2)
      n
      (+ (fib (- n 1)) (fib (- n 2)))))
```
```
  (fib 20)

6765
```

### Syntax
g-fu quasi-quotes using `'` and splices using `%`, `_` is used for missing values and `..` to splat sequences.

### Vectors
g-fu uses vectors (or slices in Go lingo), rather than linked lists as basic data structure.

The empty vector is written `()`, which is not the same thing as `_`.

```
  (len ())

0

  (len _)

Error: Len not supported: Nil
```

One consequence of using vectors is that items are pushed/popped last rather than first.

```
  (let (v ())
    (push v 1 2 3)
    (pop v)
    v)

(1 2)
```

### Types
`type` may be used to get the type of any value.

```
  (type 42)

Int

  (type Int)

Meta

  (type Meta)

Meta
```

Calling a type returns the direct parent if the argument is a child.

```
  (Int 42)
  
Int

  (Int T)
_

  (Seq Vec)

Seq

  (Seq IntIter)

Iter
```

### Conditions
```(load "lib/cond.gf")```

Every value has a boolean representation that may be retrieved using `bool`.

```
  (bool 42)

T

  (bool "")

F
```

Values may be combined using `or`/`and`. Only used values are evaluated. Comparisons are performed using boolean representations while preserving the original values. The last evaluated argument is always returned.

```
  (or 0 1)

1

  (or 0 F)
F

  (and 1 2 3)

3

  (and 1 _ 3)

_
```

`if` may be used to branch on a condition.

```
  (if 42 'foo 'bar)

'foo
```

The else-branch is optional.

```

  (if "" 'foo)  
_
```

`switch` may be used to combine multiple branches.

```
  (switch
    (F 'foo)
    (T 'bar)
    (T 'baz))

'bar
```

### Bindings
All identifiers except constants like `_`/`T`/`F` live in the same namespace. New bindings may be created using `let`.

`let` comes in two flavors. When called with arguments, it creates bindings and evaluates its body in a fresh environment.

```
  (let (foo 'outer)
    (let (foo 'inner)
      (say foo))
    (say foo))

inner
outer
```

And when called without arguments, it creates the specified bindings in the current environment.

```
  (let (foo 1)
    (let bar 2 baz (+ bar 1))
    (say foo bar baz))

123
```

Shadowing is not allowed within the same environment.

```
  (let (foo 1)
    (let foo 2))

Error: Dup binding: foo 1
```

`set` may be used to change the value of existing bindings.

```
  (let (foo 1)
    (let _
      (set 'foo 3))
    (say foo))

3
```

### Environments
Environments are first class, `this-env` evaluates to the current environment.

Used bindings (`this-env` in the following example) are automatically captured.

```
  (let (foo 42) this-env)

(foo:42 this-env:(prim this-env))
```

### Functions
Functions are created using `fun`; which accepts an optional name to bind in the current environment, a list of arguments and a body; and returns the function.

```
  (let _
    (fun say-hello (x) (say "Hello " x))
    (dump say-hello)
    (say-hello "World"))

(fun say-hello (x) (say "Hello " x))
Hello World
```

Arguments may be defined as optional by specifying default values.

```
  (let _
    (fun say-hello ((x "World")) (say "Hello " x))
    (say-hello))

Hello World
```

Arguments suffixed with `..` consume any number of remaining values.

```
  (let _
    (fun say-hello (xs..) (say "Hello " xs))
    (say-hello 'Moe 'Larry 'Curly))

Hello Moe Larry Curly
```

The following example implements a simple counter as a closure.

```
  (let (i 0
        c (fun ((d 1)) (inc i d)))
    (say (c 1))
    (say (c 2)))

1
3
```

### Macros
From a distance, macros look much like functions. They are defined the same way, accept arguments and return results. The difference is that macros have to deal with two dimensions, expansion time and evaluation time. As a consequence, macro arguments are not automatically evaluated. What is eventually evaluated is the result of expanding the macro.

The following example defines a macro called `foo` that expands to it's argument.

```
  (let _
    (mac foo (x) x)
    (foo 42))

42
```

Raising the bar one notch, the `call`-macro below expands into code calling the specified target with arguments. `expand` may be used to expand any macro call to the specified depth.

```
  (let _
    (mac call (x args..) '(%x %args..))
    (dump (expand 1 '(call + 35 7)))
    (call + 35 7))

(+ 35 7)
42
```

The next example is taken straight from the [standard library](https://github.com/codr7/g-fu/blob/master/v1/lib/abc.gf), and expands recursively to a nested series of calls using previous result as first argument.

```
(mac @ (f1 fs..)
  '(fun (args..)
     %(tr fs '(call %f1 args..) (fun (acc x) '(call %x %acc)))))
```

### Iterators
```(load "lib/iter.gf")```

All loops support exiting with a result using `(break ...)` and skipping to the start of next iteration using `(continue)`.

The most fundamental loop is called `loop`, and that's exactly what it does until exited using `break` or external means such as `recall` and `fail`.

```
  (say (loop (say 'foo) (break 'bar) (say 'baz)))

foo
bar
```

The `while`-loop keeps iterating until the specified condition turns false.

```
  (let (i 0)
    (while (< i 3)
      (say (inc i))))

1
2
3
```

The `for`-loop accepts any iterable and an optional variable name, and runs one iteration for each value.

```
  (for 3 (say 'foo))

foo
foo
foo
```

```
  (for ('(foo bar baz) v) (say v))

foo
bar
baz
```

### Tasks
Tasks are first class, preemptive green threads (or goroutines) that run in separate environments and interact with the outside world using channels. New tasks are started using `task`, which takes an optional task id and channel or buffer size and returns the new task. `wait` may be used to wait for task completion and get the results.

```
  (let _
    (task t1 () (say 'foo) 'bar)
    (task t2 () (say 'baz) 'qux)
    (say (wait t1 t2)))

baz
foo
bar qux
```

The defining environment is cloned.

```
  (let (v 42)
    (say (wait (task () (inc v))))
    (say v))

43
42
```

#### Channels
Channels are optionally buffered, thread-safe pipes. `chan` may be used to create new channels, and `push`/`pop` to transfer values. `len` returns the current number of buffered values.

```
  (let (c (chan 1))
    (push c 42)
    (say (len c))
    (say (pop c)))

1
42
```

Unbuffered channels are useful for synchronizing tasks. The following example starts a task called `t` and puts the current task in its inbox, `t` then replies `'foo` and returns `'bar`.

```
  (let _
    (task t ()
      (post (fetch) 'foo)
      'bar)
      
    (post t this-task)
    (say (fetch))
    (say (wait t)))

foo
bar
```

### Profiles
CPU profiling may be enabled by passing `-prof` on the command line; results are written to the specified file, `fib_tail.prof` in the following example.

```
$ ./gfu -prof fib_tail.prof bench/fib_tail.gf

$ go tool pprof fib_tail.prof
File: gfu
Type: cpu
Time: Apr 12, 2019 at 12:52am (CEST)
Duration: 16.52s, Total samples = 17.31s (104.79%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 10890ms, 62.91% of 17310ms total
Dropped 99 nodes (cum <= 86.55ms)
Showing top 10 nodes out of 79
      flat  flat%   sum%        cum   cum%
    2130ms 12.31% 12.31%    15980ms 92.32%  _/home/a/Dev/g-fu/v1/src/gfu.Vec.EvalVec
    1580ms  9.13% 21.43%    15980ms 92.32%  _/home/a/Dev/g-fu/v1/src/gfu.(*VecType).Eval
    1500ms  8.67% 30.10%     1680ms  9.71%  runtime.heapBitsSetType
    1180ms  6.82% 36.92%     1520ms  8.78%  _/home/a/Dev/g-fu/v1/src/gfu.(*Env).Find
     970ms  5.60% 42.52%     2290ms 13.23%  _/home/a/Dev/g-fu/v1/src/gfu.(*SymType).Eval
     830ms  4.79% 47.31%    15980ms 92.32%  _/home/a/Dev/g-fu/v1/src/gfu.Val.Eval
     780ms  4.51% 51.82%     4440ms 25.65%  runtime.mallocgc
     770ms  4.45% 56.27%      770ms  4.45%  runtime.memclrNoHeapPointers
     730ms  4.22% 60.49%    15980ms 92.32%  _/home/a/Dev/g-fu/v1/src/gfu.(*FunType).Call
     420ms  2.43% 62.91%      420ms  2.43%  runtime.memmove
```

### License
LGPL3