(mac loop (body..)
  '(let _
     (mac break (args..) '(recall T %args..))
     (mac continue () '(recall))
     
     (call (fun ((%$done? F) %$result..)
             (if %$done?
               (if %$result (if (= 1 (len %$result)) (head %$result) %$result))
               (do %body.. (recall)))))))

(mac while (cond body..)
  '(loop
     (if %cond _ (break))
     %body..))

(mac for (args body..)
  (let v? (= (type args) Vec)
       out (if (and v? (> (len args) 1)) (pop args) $out))
       
  '(let (%$in (iter %(if v? (pop args) args)))
     (loop
       (let %out (pop %$in))
       (if (_? %out) (break))
       %body..)))

(mac t@ (rf fs..)
  '(call (@ %(reverse fs)..) %rf))

(fun t@@ (rf fs..)
  (call (@@ (reverse fs)..) rf))

(mac tr-fun (rf args body..)
  (let f '(fun %args %body..))
  
  '(if (_? %rf)
     (fun (%rf) %f)
     %f))

(fun tmap (f (rf _))
  (tr-fun rf (acc val)
    (rf acc (f val))))

(Vec/fun map (in f)
  (tr in () (tmap f push)))

(fun tcat ((rf _))
  (tr-fun rf (acc val)
    (if (Seq/? val)
      (tr val acc rf)
      (rf acc val))))

(fun tflat ((rf _))
  (tr-fun rf (acc val)
    (if (Seq/? val)
      (tr val acc (tflat rf))
      (rf acc val))))

(fun tfilt (f (rf _))
  (tr-fun rf (acc val)
    (if (f val)
      (rf acc val)
      acc)))

(fun find-if (in pred)
  (for (in v)
    (if (set v (pred v)) (break v))))