(test (= (_ foo bar baz) _))

(test (bool 42))
(test (not (vec)))

(test (= (and T 42) 42))
(test (= (or 42 T _) 42))

(let (x 35)
  (test (= (inc x 7) 42))
  (test (= x 42)))

(test (= (- 42) -42))
(test (= (- 10 1 2 3) 4))

(test (= (+ -42) 42))
(test (= (+ 35 7) 42))

(test (== 'foo 'foo))
(test (= ''foo ''foo))
(test (not (= (g-sym) (g-sym))))

(test (= '(1 2 3) (vec 1 2 3)))
(test (= '(1 %(+ 2 3) 4) (vec 1 5 4)))
(test (= (vec 1 2 3) (vec 1 2 3)))
(test (== (vec 1 2 3) (vec 1 2 3)))
(test (= '(1 %(vec 2 3)..) (vec 1 2 3)))
(test (= (+ (vec 1 2 3)..) 6))

(let (v (vec))
  (push v 1)
  (push v 2 3)
  (test (= (len v) 3))
  (test (= (pop v) 3))
  (test (= (peek v) 2)))

(test (= (do 1 2 3) 3))

(test (= ((fun (xs..) xs) 1 2 3) (vec 1 2 3)))
(test (= ((fun (xs..) (+ xs..)) 1 2 3) 6))
(test (= (let (x 35) ((fun (y) (+ x y)) 7)) 42))

(let (foo (fun (x?) x))
  (test (= (foo 42) 42))
  (test (= (foo) _)))

(let (foo (macro () ''bar))
  (test (= (foo) 'bar)))

(let (foo 42 bar (macro () 'foo))
  (test (= (bar) 42)))

(let (foo (macro (x) '(+ %x 7)))
  (test (= (foo 35) 42)))

(let (fib (fun (n)
            (if (< n 2)
              n
              (+ (fib (- n 1)) (fib (- n 2))))))
  (test (= (fib 20) 6765)))

(let (fib (fun (n a b)
            (if n 
              (if (= n 1)
                b
                (recall (- n 1) b (+ a b)))
              a)))
  (test (= (fib 20 0 1) 6765)))

(let (foo 42)
  (test (= (eval 'foo) 42)))

(let (foo 35)
  (test (= (eval '(+ %foo 7)) 42)))

(let (foo (vec 35 7))
  (test (= (eval '(+ %foo..)) 42)))


(let (t (task _ 'foo))
  (test (= (wait t) 'foo)))

(let (t1 (task _ 35)
      t2 (task _ 7))
  (test (= (+ (wait t1 t2)..) 42)))

(let (v 42
      t (task _ (inc v)))
  (test (= (wait t) 43))
  (test (= v 42)))

(let (t (task _
          (post (fetch) 'foo)
          'bar))
  (post t (this-task))
  (test (= (fetch) 'foo))
  (test (= (wait t) 'bar)))


(let loop (macro (body..)
  (let done (g-sym) result (g-sym))
  
  '(let (break (macro (args..) '(recall T %args..)))
     ((fun (%done? %result..)
        (if %done %result.. (do %body.. (recall))))))))

(test (= (loop (break 'foo)) 'foo))

(let while (macro (cond body..)
  '(loop
     (if %cond _ (break))
     %body..)))

(let (i 0)
  (while (< (inc i) 7))
  (test (= i 7)))