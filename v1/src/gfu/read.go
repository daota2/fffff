package gfu

import (
  "io"
  //"log"
  "strconv"
  "strings"
  "unicode"
)

func (g *G) ReadChar(pos *Pos, in *strings.Reader) (rune, E) {
  c, _, e := in.ReadRune()

  if e == io.EOF {
    return 0, nil
  }

  if e != nil {
    return 0, g.ReadE(*pos, "Failed reading char: %v", e)
  }

  if c == '\n' {
    pos.Col = INIT_POS.Col
    pos.Row++
  } else {
    pos.Col++
  }

  return c, nil
}

func (g *G) Unread(pos *Pos, in *strings.Reader, c rune) E {
  if e := in.UnreadRune(); e != nil {
    return g.ReadE(*pos, "Failed unreading char")
  }

  if c == '\n' {
    pos.Row--
  } else {
    pos.Col--
  }

  return nil
}

func (g *G) Read(pos *Pos, in *strings.Reader, out Vec, end rune) (Vec, E) {
  var c rune
  var e E

  for {
    c, e = g.ReadChar(pos, in)

    if e != nil || c == 0 || c == end {
      return nil, e
    }

    switch c {
    case ' ', '\n':
      break
    case '(':
      return g.ReadVec(pos, in, out)
    case '\'':
      return g.ReadQuote(pos, in, out, end)
    case '.':
      return g.ReadSplat(pos, in, out)
    case '%':
      return g.ReadSplice(pos, in, out, end)
    case '"':
      return g.ReadStr(pos, in, out)
    default:
      if unicode.IsDigit(c) {
        if e = g.Unread(pos, in, c); e != nil {
          return nil, e
        }

        return g.ReadNum(pos, in, out, false)
      } else if c == '-' {
        var nc rune
        nc, e = g.ReadChar(pos, in)

        if e != nil {
          return nil, e
        }

        is_num := unicode.IsDigit(nc)

        if e = g.Unread(pos, in, nc); e != nil {
          return nil, e
        }

        if is_num {
          return g.ReadNum(pos, in, out, true)
        }

        return g.ReadId(pos, in, out, "-")
      } else if unicode.IsGraphic(c) {
        if e = g.Unread(pos, in, c); e != nil {
          return nil, e
        }

        return g.ReadId(pos, in, out, "")
      }

      return nil, g.ReadE(*pos, "Unexpected input: %v", c)
    }
  }
}

func (g *G) ReadId(pos *Pos, in *strings.Reader, out Vec, prefix string) (Vec, E) {
  var buf strings.Builder
  buf.WriteString(prefix)

  for {
    c, e := g.ReadChar(pos, in)

    if e != nil {
      return nil, e
    }

    if c == 0 {
      break
    }

    if unicode.IsSpace(c) ||
      c == '.' || c == '?' || c == '%' || c == '(' || c == ')' {
      if e := g.Unread(pos, in, c); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.ReadE(*pos, "Failed writing char: %v", we)
    }
  }

  return append(out, g.Sym(buf.String())), nil
}

func (g *G) ReadNum(pos *Pos, in *strings.Reader, out Vec, is_neg bool) (Vec, E) {
  var buf strings.Builder

  for {
    c, e := g.ReadChar(pos, in)

    if e != nil {
      return nil, e
    }

    if c == 0 {
      break
    }

    if !unicode.IsDigit(c) {
      if e := g.Unread(pos, in, c); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.ReadE(*pos, "Failed writing char: %v", we)
    }
  }

  s := buf.String()
  n, e := strconv.ParseInt(s, 10, 64)

  if e != nil {
    return nil, g.ReadE(*pos, "Invalid num: %v", s)
  }

  if is_neg {
    n = -n
  }

  return append(out, Int(n)), nil
}

func (g *G) ReadQuote(pos *Pos, in *strings.Reader, out Vec, end rune) (Vec, E) {
  vpos := *pos
  vs, e := g.Read(pos, in, nil, end)

  if e != nil {
    return nil, e
  }

  if len(vs) == 0 {
    return nil, g.ReadE(vpos, "Nothing to quote")
  }

  return append(out, NewQuote(vs[0])), nil
}

func (g *G) ReadSplat(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
  vpos := *pos
  vpos.Col--

  var nc rune
  var e E

  nc, e = g.ReadChar(pos, in)

  if e != nil {
    return nil, e
  }

  if nc != '.' {
    return nil, g.ReadE(*pos, "Invalid input: .%v", nc)
  }

  i := len(out)

  if i == 0 {
    return nil, g.ReadE(*pos, "Missing splat value")
  }

  out[i-1] = NewSplat(out[i-1])
  return out, nil
}

func (g *G) ReadSplice(pos *Pos, in *strings.Reader, out Vec, end rune) (Vec, E) {
  vpos := *pos
  vpos.Col--

  vs, e := g.Read(pos, in, nil, end)

  if e != nil {
    return nil, e
  }

  if len(vs) == 0 {
    return nil, g.ReadE(vpos, "Nothing to eval")
  }

  return append(out, NewSplice(vs[0])), nil
}

func (g *G) ReadStr(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
  var buf strings.Builder

  for {
    c, e := g.ReadChar(pos, in)

    if e != nil {
      return nil, e
    }

    if c == 0 || c == '"' {
      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.ReadE(*pos, "Failed writing char: %v", we)
    }
  }

  return append(out, Str(buf.String())), nil
}

func (g *G) ReadVec(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
  var body Vec

  for {
    vs, e := g.Read(pos, in, body, ')')

    if e != nil {
      return nil, e
    }

    if vs == nil {
      break
    }

    body = vs
  }

  return append(out, body), nil
}
