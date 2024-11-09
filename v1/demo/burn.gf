(debug)
(load "../lib/all.gf")

(let burn (let (width 25 height 25
                esc (str 0x1b "[")
                data (new-bin (* width height 1))
                out _)
  (fun get-offs (x y)
    (+ (- width x 1) (* (- height y 1) width)))

  (fun xy (x y)
    (# data (get-offs x y)))

  (fun set-xy (f x y)
    (let i (get-offs x y))
    (set (# data i) (f (# data i)))
    _)

  (fun move-to (x y)
    (print out (str esc (if (or x y) (str (+ y 1) ";" (+ x 1))) "H")))

  (fun cr-lf ((n 1))
    (print out (str esc (if (> n 1) n) "E"))) 
  
  (fun pick-color (r g b)
    (print out (str esc "48;2;" r ";" g ";" b "m")))

  (fun init ()
    (for (width x)
      (set (xy x 0) 0xff)))

  (fun render ()
    (for ((- height 1) y)
      (for (width x)
        (let v (xy x y))
        (set (xy x (+ y 1)) (- v (rand (+ v 1))))))

    (move-to 0 0)
  
    (for (height y)
      (if y (cr-lf))

      (for (width x)
        (let g (xy x y) r (if g 0xff 0x00) b (if (= g 0xff) 0xff 0x00))
        (pick-color r g b)
        (print " " ))))

  this-env))

(burn/init)
(burn/render)