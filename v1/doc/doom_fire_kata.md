## The DOOM Fire Kata

### Intro
Ever since I came across the DOOM fire [trick](https://fabiensanglard.net/doom_fire_psx/), I've been itching to work my way through it using console graphics for use as kata to exercise new languages. This post describes how I would approach it in [g-fu](https://github.com/codr7/g-fu/tree/master/v1), a pragmatic Lisp embedded in Go.

![Fire](fire.gif)
[Source](https://github.com/codr7/g-fu/blob/master/v1/demo/fire.gf)

### Setup
If you're in the mood to light your own fire, the following shell spell will take you where you need to go.

```
$ git clone https://github.com/codr7/g-fu.git
$ cd g-fu/v1
$ go build src/gfu.go
$ rlwrap ./gfu
g-fu v1.15

Press Return twice to evaluate.

  (load "demo/fire.gf")
```

### Syntax
g-fu quasi-quotes using `'` and splices using `%`, `_` is used for missing values and `..` to splat sequences.

### Idea
The idea is to model each "particle" of the fire as a value that decays from white to black along a reddish scale while moving upwards. This is the reason for the white line at the bottom, that's where new particles are born. Add a touch of pseudo-chaos to make it interesting and that's pretty much it.

### Implementation
Particles are implemented using an array of bytes representing the green part of their colors. Red is locked at 255 and blue at 0 to get a gradient of red/yellow colors. 

We start with module variables and a set of utilities for manipulating the console, more info may be found on [Wikipedia](https://en.wikipedia.org/wiki/ANSI_escape_code).

```
(env fire (width 50 height 25
           buf (new-bin (* width height))
           esc (str 0x1b "[")
           out stdout
           max-fade 50
           tot-frames 0 tot-time .0)
  (fun clear ()
    (print out (str esc "2J")))

  (fun home ()
    (print out (str esc "H")))

  (fun pick-color (r g b)
    (print out (str esc "48;2;" (int r) ";" (int g) ";" (int b) "m")))
```

Before we can start rendering, the bottom row needs to be initialized and the screen cleared.

```
  (fun init ()
    (for (width i)
      (set (# buf i) 0xff))

    (clear))
```

Rendering begins with a loop that fades all particles. Particles may rise straight or diagonally, the three cases are handled by the `if`-statement. Next the color is faded if not black and the particle is moved up one row. Note that particles are actually stored bottom-top, since that's the direction they move.

```
  (fun render ()
    (let t0 (now) i -1)

    (for ((- height 1) y)
      (for (width x)
        (let v (# buf (inc i))
             j (+ i width))
        
        (if (and x (< x (- width 1)))
          (inc j (- 1 (rand 3))))
        
        (set (# buf j)
             (if v (- v (rand (min max-fade (+ (int v) 1)))) v))))

        ...
```

Once particles are faded and moved, its time to generate console output. We start by adding the top row to the index and moving the cursor home, then pick the right color and print a blank for each particle last-first. Before exiting, the output is flushed and average frame rate recorded.

```
    ...

    (inc i (+ width 1))
    (home)
    
    (for (height y)
      (for (width x)
        (let g (# buf (dec i))
             r (if g 0xff 0)
             b (if (= g 0xff) 0xff 0))
             
        (pick-color r g b)
        (print out " "))

      (print out \n))

    (flush out)
    (inc tot-time (- (now) t0))
    (inc tot-frames))
```

Since it's rude to mess around with user console settings, we make sure that everything is put back in the right place before leaving; the first line resets the color.

```
  ...
  
  (fun restore ()
    (print out (str esc "0m"))
    (clear)
    (home)))
```

The final few lines run 50 frames and print the average frame rate.

```
(fire/init)
(for 50 (fire/render))
(fire/restore)

(say (/ (* 1000000000.0 fire/tot-frames) fire/tot-time))
```

### Performance
While there's nothing seriously wrong with this implementation from my perspective, it's not going to win any performance prizes yet. [g-fu](https://github.com/codr7/g-fu/tree/master/v1) is still very young and I'm mostly focusing on correctness at this point. More mature languages with comparable features should be able to run this plenty faster. One thing that does come to mind is using a separate buffer for output and dumping that all at once to the console, the code supports switching output stream but g-fu is still missing support for memory streams.

Until next time,<br/>
c7