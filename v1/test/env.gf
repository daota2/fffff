(test (= (len (let _ this-env)) 1))

(let (e (let (foo 1)
          (use _ do let)
          this-env))
  (test (= e/foo 1))
  (e/do (let bar 2))
  (test (= e/bar 2))

  (let ee (let _ (use e foo))
    (test (= ee/foo 1))))

(let (foo (let (bar 7)
            (fun resolve (key) 42)
            this-env))
  (test (= foo/bar 7))
  (test (= foo/baz 42)))

(let (foo (let (bar 42)
            this-env)
      baz (let _
            (use foo bar)
            this-env))
  (test (= baz/bar 42)))

(let (super this-env
      Counter (fun ((n 0))
                (fun inc ((d 1)) (super/inc n d))
                this-env)
      c (Counter))
  (for 3 (c/inc))
  (test (= (c/inc -1) 2)))

(let (key _ val _)
  (fun foo (x y))
  (fun set-foo (f k..) (set key k val (f _)))
  (set (foo 'bar 'baz) 42)
  (test (= key '(bar baz)))
  (test (= val 42)))

(let (proxy (fun (d)
              (fun resolve (key)
                (d/val key))
              this-env)
      p (proxy (let (foo 42)
                 (use _ val)
                 this-env)))
  (test (= p/foo 42)))

(let _

(let Widget (let _
  (fun new (args..)
    (let left 0 top 0
         width (or (pop-key args 'width) (fail "Missing width"))
         height (or (pop-key args 'height) (fail "Missing height")))

    (fun move (dx dy)
      (vec (inc left dx)
           (inc top dy)))

    (fun resize (dx dy)
      (vec (inc width dx)
           (inc height dy)))
  
    this-env)

  this-env))

(let Button (let _
  (fun new (args..)
    (let this (this-env)
         w (Widget/new args..)
         click-event ())
         
    (use w move)

    (fun click ()
      (for (click-event f) (f this)))
      
    (fun on-click (f)
      (push click-event f))
    
    (fun resize (dx dy)
      (w/resize (min (+ w/width dx) (- 200 w/width))
                (min (+ w/height dy) (- 100 w/height))))
    
    this)

  this-env))

(let b (Button/new 'width 100 'height 50))

(test (= (b/move 10 10) '(10 10)))
(test (= (b/resize 400 200) '(200 100)))

(let (called F)
  (b/on-click (fun (b) (set called T)))
  (test (not called))
  (b/click)
  (test called))

)