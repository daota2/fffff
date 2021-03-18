package gfu

type Vec struct {
  items []Val
}

func NewVec() *Vec {
  return new(Vec)
}

func (v *Vec) Push(item Val) {
  v.items = append(v.items, item)
}

func (v *Vec) Pop() *Val {
  if v.items == nil {
    return nil
  }

  is := v.items
  n := len(is)
  var it Val
  it, v.items = is[n-1], is[:n-1]
  return &it
}
