(let switch (mac (alts..)
  (tr (reverse alts) _
      (fun (acc alt)
        '(if %(head alt)
           (do %(tail alt)..)
           %acc)))))