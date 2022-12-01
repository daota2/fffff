(let loop (mac (body..)
  (let done? (g-sym) result (g-sym))
  
  '(let (break (mac (args..) '(recall T %args..)))
     ((fun ((%done? F) %result..)
        (if %done? %result.. (do %body.. (recall))))))))

(let while (mac (cond body..)
  '(loop
     (if %cond _ (break))
     %body..)))

(let for (mac (args body..)
  (let v? (= (type args) Vec)
       i (if (and v? (> (len args) 1)) (pop args) (g-sym))
       n (g-sym))
       
  '(let (%i 0 %n %(if v? (pop args) args))
     (while (< %i %n)
       %body..
       (inc %i)))))

(let map (fun (f)
  (fun (rf)
    (fun (acc val)
      (rf acc (f val))))))

(let cat (fun (rf)
  (fun (acc val)
    (fold val rf acc))))

(let keep (fun (f)
  (fun (rf)
    (fun (acc val)
      (if (f val)
        (rf acc val)
        acc)))))

(let @ (fun (fs..)
  (fun (in)
    (fold (reverse fs) (fun (acc x) (x acc)) in))))