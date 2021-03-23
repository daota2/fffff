package gfu

import (
  "fmt"
)

type Error interface {
  String() string
}

type BasicError struct {
  pos Pos
  msg string
}

func (e *BasicError) Init(pos Pos, msg string) *BasicError {
  e.pos = pos
  e.msg = msg
  return e
}

func (e *BasicError) String() string {
  p := &e.pos
  
  return fmt.Sprintf(
    "Error in '%s' on row %v, col %v:\n%v",
    p.src, p.Row, p.Col, e.msg)
}

func (g *G) E(pos Pos, msg string, args...interface{}) *BasicError {
  msg = fmt.Sprintf(msg, args...)  
  e := new(BasicError).Init(pos, msg)

  if g.Debug {
    panic(e.String())
  }

  return e
}
