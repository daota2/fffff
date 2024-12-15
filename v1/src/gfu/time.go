package gfu

import (
  "bufio"
  //"log"
  "time"
)

type Time time.Time
var TimeFormat = "2006-01-02 15:04:05"

type TimeType struct {
  NumType
}

func (_ Time) Type(g *G) Type {
  return &g.TimeType
}

func (_ *TimeType) Dump(g *G, val Val, out *bufio.Writer) E {
  t := time.Time(val.(Time))
  out.WriteString(t.Format(TimeFormat))
  return nil
}

func (t *TimeType) Sub(g *G, x, y Val) (Val, E) {
  yt, ok := y.(Time)

  if !ok {
    return nil, g.E("Expected Time: %v", y.Type(g))
  }
  
  return NSecs(time.Time(x.(Time)).Sub(time.Time(yt))), nil
}
